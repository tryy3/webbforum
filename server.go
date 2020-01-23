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

	// initialize the storer
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

	// static folders
	imageFolder := filepath.Join(viper.GetString("content.base"), viper.GetString("content.image.folder"))
	cssFolder := filepath.Join(viper.GetString("content.base"), viper.GetString("content.css.folder"))
	jsFolder := filepath.Join(viper.GetString("content.base"), viper.GetString("content.js.folder"))

	// serve the static folders
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images", http.FileServer(http.Dir(imageFolder))))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js", http.FileServer(http.Dir(jsFolder))))
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css", http.FileServer(http.Dir(cssFolder))))

	// user routes
	r.Methods("GET").
		PathPrefix("/user/{username}").
		Handler(handlers.NewMemberHandler(storer))

	// profile routes
	r.Methods("GET").
		PathPrefix("/profile").
		Handler(middleware.LoggedInProtect(handlers.NewProfileHandler(ab), ab))
	r.Methods("POST").
		PathPrefix("/profile/upload").
		Handler(middleware.LoggedInProtect(handlers.NewProfileUploadHandler(db, ab), ab))
	r.Methods("POST").
		PathPrefix("/profile").
		Handler(middleware.LoggedInProtect(handlers.NewProfileEditHandler(storer, ab), ab))

	// admin routes
	r.Methods("GET").
		PathPrefix("/admin/category/remove/{category}").
		Handler(handlers.CategoryDeleteHandler{db})
	r.Methods("POST").
		PathPrefix("/admin/category/create").
		Handler(handlers.CategoryCreateHandler{db})
	r.Methods("POST").
		PathPrefix("/admin/category/modify").
		Handler(handlers.CategoryEditHandler{db})
	r.Methods("GET").
		PathPrefix("/admin/group/remove/{group}").
		Handler(handlers.GroupDeleteHandler{db})
	r.Methods("POST").
		PathPrefix("/admin/group/create").
		Handler(handlers.GroupCreateHandler{db})
	r.Methods("POST").
		PathPrefix("/admin/group/modify").
		Handler(handlers.GroupEditHandler{db})
	r.Methods("GET").
		PathPrefix("/admin").
		Handler(handlers.AdminHandler{db})

	// Posts routes
	r.Methods("POST").
		PathPrefix("/forums/thread/{thread}/new_post").
		Handler(handlers.PostCreateHandler{db, ab})
	r.Methods("GET", "POST").
		PathPrefix("/forums/thread/{thread}/edit/{post}").
		Handler(handlers.PostEditHandler{db, ab})
	r.Methods("GET").
		PathPrefix("/forums/thread/{thread}").
		Handler(handlers.PostsShowHandler{db})

	// Thread routes
	r.Methods("GET", "POST").
		PathPrefix("/forums/{category}/create_thread").
		Handler(handlers.ThreadCreateHandler{db, ab})
	r.Methods("POST").
		PathPrefix("/thread/modify").
		Handler(handlers.ThreadEditHandler{db, ab})
	r.Methods("POST").
		PathPrefix("/thread/remove").
		Handler(handlers.ThreadDeleteHandler{db, ab})
	r.Methods("GET").
		PathPrefix("/forums/{category}").
		Handler(handlers.ThreadShowHandler{db})

	// main route
	r.Methods("GET").
		PathPrefix("/").
		Handler(handlers.HomeHandler{Database: db})

	// NotFoundHandler
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "Not found")
	})

	// set up the middleware chain
	//stack := alice.New(middleware.Logging, middleware.NoSurfingMiddleware, ab.ExpireMiddleware).Then(r)
	stack := alice.New(middleware.Logging, ab.ExpireMiddleware).Then(r)

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
		return errors.Wrap(err, "error when compiling templates")
	}
	handlers.Templates = t
	return nil
}
