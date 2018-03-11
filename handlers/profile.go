package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/tryy3/webbforum/fileutils"
	"github.com/tryy3/webbforum/models"
	"github.com/volatiletech/authboss"
)

// ProfileHandler handles the /profil page with GET requests
type ProfileHandler struct {
	ab *authboss.Authboss
}

func NewProfileHandler(ab *authboss.Authboss) http.Handler {
	return ProfileHandler{ab}
}

func (p ProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u, ok := getUser(w, r, p.ab)
	if !ok {
		return
	}
	data := LayoutData(w, r).MergeKV("user", u)
	mustRender(w, r, "profil", data)
}

// profileEditValidators are the required fields when editing a profile page
var profileEditValidators = []authboss.Validator{
	authboss.Rules{
		FieldName:       "first_name",
		Required:        true,
		AllowWhitespace: true,
	},
	authboss.Rules{
		FieldName:       "last_name",
		Required:        true,
		AllowWhitespace: true,
	},
	authboss.Rules{
		FieldName:       "email",
		Required:        true,
		AllowWhitespace: false,
	},
}

// ProfileEditHandler takes care of editing a profile page
type ProfileEditHandler struct {
	storer models.UserStorer
	ab     *authboss.Authboss
}

func NewProfileEditHandler(storer models.UserStorer, ab *authboss.Authboss) ProfileEditHandler {
	return ProfileEditHandler{storer, ab}
}

func (p ProfileEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u, ok := getUser(w, r, p.ab)
	if !ok {
		return
	}

	data := LayoutData(w, r)

	// attempt to get authboss attributes from the request
	ab, err := authboss.AttributesFromRequest(r)
	if err != nil {
		data = data.MergeKV("error", "unable to prase request", "user", u)
		mustRender(w, r, "profil", data)
		return
	}

	// get a slice of all the keys from attributes
	keys := make([]string, len(ab))
	i := 0
	for k := range ab {
		keys[i] = k
		i++
	}

	// attempt to validate the request
	errs := authboss.Validate(r, profileEditValidators)
	spew.Dump(errs)
	if len(errs) > 0 {
		data.MergeKV("errs", errs.Map(), "user", u)
		mustRender(w, r, "profil", data)
		return
	}

	// update the user fields
	err = p.storer.Put(u.(*models.User).Username, ab)
	if badRequest(w, err) {
		data = data.MergeKV("error", "internal error", "user", u)
		mustRender(w, r, "profil", data)
		return
	}

	http.Redirect(w, r, r.RequestURI, http.StatusFound)
}

// ProfileUploadHandler handles uploading profile image
type ProfileUploadHandler struct {
	db *gorm.DB
	ab *authboss.Authboss
}

func NewProfileUploadHandler(db *gorm.DB, ab *authboss.Authboss) ProfileUploadHandler {
	return ProfileUploadHandler{db, ab}
}

func (p ProfileUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u, ok := getUser(w, r, p.ab)
	if !ok {
		return
	}

	// convert from interface{} to models.User
	user, ok := u.(*models.User)
	if !ok {
		return
	}

	data := LayoutData(w, r)

	fileUpload := models.File{}
	fileUpload.UserID = user.ID
	fileUpload.ContentType = r.Header.Get("content-type")

	// save 32MB into memory, if the file is large it will be saved to a temporary file
	r.ParseMultipartForm(32 << 20)

	// get the uploaded file
	file, handle, err := r.FormFile("uploadfile")
	if err != nil {
		log.WithError(err).Error("internal error when uploading file")
		data = data.MergeKV("error", "internal error", "user", u)
		mustRender(w, r, "profil", data)
		return
	}
	defer file.Close()

	fileUpload.UploadName = handle.Filename

	maxFileSize := viper.GetInt64("content.image.size")
	tmpFolder := filepath.Join(viper.GetString("content.base"), viper.GetString("content.tmp"))

	// attempt to write the uploaded file to a temporary file
	hash, bytesWritten, tmpDir, err := fileutils.WriteTempFile(file, maxFileSize, tmpFolder)
	if err != nil {
		log.WithError(err).Error("internal erorr when saving file")
		data = data.MergeKV("error", "internal error", "user", u)
		mustRender(w, r, "profil", data)
		return
	}

	fileUpload.Base64Hash = hash
	fileUpload.FileSizeBytes = bytesWritten

	// check if the file has already been uploaded before
	media := models.File{}
	sel := p.db.Where("base64_hash = ?", fileUpload.Base64Hash).First(&media)
	if sel.Error != nil && !sel.RecordNotFound() {
		log.WithError(sel.Error).Error("internal error when retrieving existing data from database")
		data = data.MergeKV("error", "internal error", "user", u)
		mustRender(w, r, "profil", data)
		return
	}

	// if the file haven't been uploaded before, we will need to move it to the correct folder
	// otherwise simply remove the temporary file
	if sel.RecordNotFound() {
		imageFoler := filepath.Join(viper.GetString("content.base"), viper.GetString("content.image.folder"))

		// attempt to move the temp file
		_, duplicate, err := fileutils.MoveFile(tmpDir, imageFoler, fileUpload)
		os.RemoveAll(tmpDir)
		if err != nil {
			log.WithError(err).Error("failed to move file")
			data = data.MergeKV("error", "internal error", "user", u)
			mustRender(w, r, "profil", data)
			return
		}

		if duplicate {
			log.Error("stored file already exists")
			data = data.MergeKV("error", "internal error", "user", u)
			mustRender(w, r, "profil", data)
			return
		}
	} else {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			log.WithError(err).Error("internal error when removing temp folder")
			data = data.MergeKV("error", "internal error", "user", u)
			mustRender(w, r, "profil", data)
			return
		}
	}

	// save the file metadata to the database
	err = p.db.Create(&fileUpload).Error
	if err != nil {
		log.WithError(err).Error("internal error when inserting to database")
		data = data.MergeKV("error", "internal error", "user", u)
		mustRender(w, r, "profil", data)
		return
	}

	// update the user with the new profile data
	user.ProfileImageID = fileUpload.ID
	user.ProfileImage = &fileUpload
	user.Attachments = append(user.Attachments, &fileUpload)

	err = p.db.Model(&user).Where("username = ?", user.Username).Updates(&user).Error
	if err != nil {
		log.WithError(err).Error("internal error when updating user")
		data = data.MergeKV("error", "internal error", "user", u)
		mustRender(w, r, "profil", data)
		return
	}

	http.Redirect(w, r, "/profil", http.StatusFound)
}
