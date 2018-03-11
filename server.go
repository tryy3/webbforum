package webbforum

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/aarondl/tpl"
	"github.com/apex/log"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/justinas/alice"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/tryy3/webbforum/auth"
	"github.com/tryy3/webbforum/handlers"
	"github.com/tryy3/webbforum/middleware"
)

// StartServer takes care of initializing the http server with all of the routes
func StartServer(db *gorm.DB) error {
	// load the initial templates
	if err := loadTemplates(); err != nil {
		return errors.Wrap(err, "error when loading templates")
	}

	storer := auth.New(db)

	// setup authboss
	ab, err := auth.SetupAuthboss(storer)
	if err != nil {
		return errors.Wrap(err, "error setting up authboss")
	}

	// initialize the mux router
	r := mux.NewRouter()

	// setup the routes
	r.PathPrefix("/auth").Handler(ab.NewRouter())

	r.Methods("GET").
		PathPrefix("/anvandare/{username}").
		Handler(handlers.NewMemberHandler(storer))

	r.Methods("GET").
		PathPrefix("/profil").
		Handler(middleware.LoggedInProtect(handlers.NewProfileHandler(ab), ab))
	r.Methods("POST").
		PathPrefix("/profil/upload").
		Handler(middleware.LoggedInProtect(handlers.NewProfileUploadHandler(db, ab), ab))
	r.Methods("POST").
		PathPrefix("/profil").
		Handler(middleware.LoggedInProtect(handlers.NewProfileEditHandler(storer, ab), ab))

	r.Methods("GET").
		PathPrefix("/admin").
		Handler(handlers.AdminHandler{db})
	r.Methods("POST").
		PathPrefix("/admin/kategori/skapa").
		Handler(handlers.CategoryCreateHandler{db})
	r.Methods("POST").
		PathPrefix("/admin/kategori/uppdatera").
		Handler(handlers.CategoryEditHandler{db})
	r.Methods("POST").
		PathPrefix("/admin/kategori/ta_bort").
		Handler(handlers.CategoryDeleteHandler{db})

	r.Methods("GET").
		PathPrefix("/").
		HandlerFunc(handlers.HomeHandler)

	// serve image folder
	imageFolder := filepath.Join(viper.GetString("content.base"), viper.GetString("content.image.folder"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images", http.FileServer(http.Dir(imageFolder))))

	// NotFoundHandler
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "Not found")
	})

	// set up the middleware chain
	stack := alice.New(middleware.Logging, middleware.NoSurfingMiddleware, ab.ExpireMiddleware).Then(r)

	// create the http server
	srv := &http.Server{
		Handler: stack,
		Addr:    fmt.Sprintf("%s:%d", viper.GetString("http.host"), viper.GetInt("http.port")),

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// start the http server
	log.WithField("Host",
		fmt.Sprintf("%s:%d",
			viper.GetString("http.host"),
			viper.GetInt("http.port"),
		),
	).Info("http server started")
	return errors.Wrap(srv.ListenAndServe(), "error when starting http server")
}

// loadTemplates will try to load all the templates that will be used
func loadTemplates() error {
	t, err := tpl.Load(
		viper.GetString("views.folder"),
		viper.GetString("views.partials"),
		"layout.html.tpl",
		handlers.LayoutFuncs)
	if err != nil {
		return err
	}
	handlers.Templates = t
	return nil
}
