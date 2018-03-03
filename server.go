package webbforum

import (
	"fmt"
	"net/http"
	"net/smtp"
	"time"

	"github.com/apex/log"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/justinas/nosurf"
	"github.com/tryy3/webbforum/api"
	"github.com/tryy3/webbforum/auth"
	"github.com/tryy3/webbforum/handlers"
	loghandler "github.com/tryy3/webbforum/log"
	"github.com/tryy3/webbforum/utils"
	"github.com/volatiletech/authboss"
)

// StartServer takes care of initializing the http server with all of the routes
func StartServer(context *utils.Context, a *api.API) error {
	// initialize the mux router
	r := mux.NewRouter()

	// create the api routes
	s := r.PathPrefix("/api").Subrouter()
	api.CreateAPIRoutes(a, s)

	// setup authboss
	ab, err := SetupAuthboss(context.Config)
	if err != nil {
		return err
	}
	r.Handle("/auth", ab.NewRouter())

	// create routes for regular handlers/paths
	r.HandleFunc("/", handlers.HomeHandler)

	// create the http server
	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("%s:%d", context.Config.HTTPIP, context.Config.HTTPPort),

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.WithField("Host", fmt.Sprintf("%s:%d", context.Config.HTTPIP, context.Config.HTTPPort)).Info("http server started")
	return srv.ListenAndServe()
}

// SetupAuthboss creates a new Authboss instance with everything configured
func SetupAuthboss(config *utils.Config) (*authboss.Authboss, error) {
	// initialize authboss
	ab := authboss.New()

	// set general configurations
	ab.MountPath = "/auth"
	ab.ViewsPath = "ab_views" // TODO: Try this out
	ab.LogWriter = loghandler.NewAuthbossLogger()
	ab.RootURL = fmt.Sprintf("%s:%d", config.HTTPIP, config.HTTPPort)

	// TODO: Implement
	/*ab.LayoutDataMaker = layoutData

	b, err := ioutil.ReadFile(filepath.Join("views", "layout.html.tpl"))
	if err != nil {
		panic(err)
	}
	ab.Layout = template.Must(template.New("layout").Funcs(funcs).Parse(string(b))) */

	// xsrf protection
	ab.XSRFName = config.XSRFName
	ab.XSRFMaker = func(_ http.ResponseWriter, r *http.Request) string {
		return nosurf.Token(r)
	}

	// cookie configurations
	// TODO: Test block key
	cookieStore := securecookie.New(config.CookieStoreKey, nil)
	cookieManager := auth.NewCookieManager(cookieStore, config.CookieExpiry)
	ab.CookieStoreMaker = cookieManager.NewCookieStorer

	// session configurations
	sessionStore := sessions.NewCookieStore(config.SessionStoreKey)
	sessionManager := auth.NewSessionManager(sessionStore, config.SessionName)
	ab.SessionStoreMaker = sessionManager.NewSessionStorer

	// email configurations
	ab.Mailer = authboss.SMTPMailer(config.SMTPHost, smtp.PlainAuth(
		config.SMTPIdentity,
		config.SMTPUsername,
		config.SMTPPassword,
		config.SMTPHost,
	))
	ab.EmailFrom = config.SMTPEmail
	ab.EmailFromName = config.SMTPName

	// policies/validators
	ab.Policies = []authboss.Validator{
		authboss.Rules{
			FieldName:       "email",
			Required:        true,
			AllowWhitespace: false,
		},
		authboss.Rules{
			FieldName:       "password",
			Required:        true,
			MinLength:       4,
			MaxLength:       32,
			AllowWhitespace: false,
		},
	}

	// initialize the authboss
	if err := ab.Init(); err != nil {
		return nil, err
	}
	return ab, nil
}
