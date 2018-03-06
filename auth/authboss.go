package auth

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"path/filepath"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/justinas/nosurf"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/tryy3/webbforum/handlers"
	"github.com/tryy3/webbforum/log"
	"github.com/volatiletech/authboss"
	_ "github.com/volatiletech/authboss/auth"
	_ "github.com/volatiletech/authboss/confirm"
	_ "github.com/volatiletech/authboss/lock"
	_ "github.com/volatiletech/authboss/recover"
	_ "github.com/volatiletech/authboss/register"
	_ "github.com/volatiletech/authboss/remember"
)

// SetupAuthboss creates a new Authboss instance with everything configured
func SetupAuthboss(db *gorm.DB) (*authboss.Authboss, error) {
	// initialize authboss
	ab := authboss.New()

	// set general configurations
	ab.Storer = Storer{db: db}
	ab.MountPath = "/auth"
	ab.ViewsPath = "ab_views"
	ab.LogWriter = log.NewAuthbossLogger()
	ab.RootURL = fmt.Sprintf("http://%s:%d", viper.GetString("http.host"), viper.GetInt("http.port"))

	// lock settings
	ab.LockAfter = 3                // Login attempts
	ab.LockWindow = 5 * time.Minute // Wait before reseting lock count
	ab.LockDuration = 1 * time.Hour // How long the account should be locked for

	// layouts
	handlers.Authboss = ab
	ab.LayoutDataMaker = handlers.LayoutData

	b, err := ioutil.ReadFile(filepath.Join(viper.GetString("views.folder"), "layout.html.tpl"))
	if err != nil {
		return nil, errors.Wrap(err, "error trying to read layout template")
	}
	ab.Layout, err = template.New("layout").Funcs(handlers.LayoutFuncs).Parse(string(b))
	if err != nil {
		return nil, errors.Wrap(err, "error trying to parse layout template")
	}

	// xsrf protection
	ab.XSRFName = viper.GetString("xsrf.name")
	ab.XSRFMaker = func(_ http.ResponseWriter, r *http.Request) string {
		return nosurf.Token(r)
	}

	// generate cookie key
	cookieKey, err := base64.StdEncoding.DecodeString(viper.GetString("cookie.key"))
	if err != nil {
		return nil, errors.Wrap(err, "error generating cookie key")
	}

	// generate session key
	sessionKey, err := base64.StdEncoding.DecodeString(viper.GetString("session.key"))
	if err != nil {
		return nil, errors.Wrap(err, "error generating session key")
	}

	// cookie configurations
	// TODO: Test block key
	cookieStore := securecookie.New(cookieKey, nil)
	cookieManager := NewCookieManager(cookieStore, time.Duration(viper.GetInt64("cookie.expiry"))*time.Hour)
	ab.CookieStoreMaker = cookieManager.NewCookieStorer

	// session configurations
	sessionStore := sessions.NewCookieStore(sessionKey)
	sessionManager := NewSessionManager(sessionStore, viper.GetString("session.name"))
	ab.SessionStoreMaker = sessionManager.NewSessionStorer

	// email configurations
	ab.Mailer = authboss.SMTPMailer(fmt.Sprintf("%s:%d", viper.GetString("smtp.host"), viper.GetInt("smtp.port")), smtp.PlainAuth(
		viper.GetString("smtp.identity"),
		viper.GetString("smtp.username"),
		viper.GetString("smtp.password"),
		"smtp.gmail.com",
	))
	ab.EmailFrom = viper.GetString("smtp.email")
	ab.EmailFromName = viper.GetString("smtp.name")

	// policies/validators
	ab.Policies = generatePolicies()
	ab.PreserveFields = []string{"email", "first_name", "last_name"}
	ab.PrimaryID = authboss.StoreUsername

	// initialize the authboss
	if err := ab.Init(); err != nil {
		return nil, errors.Wrap(err, "error initialize authboss")
	}
	return ab, nil
}

// generatePolicies will return with the expected authboss.Validators
func generatePolicies() []authboss.Validator {
	return []authboss.Validator{
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
			FieldName:       "username",
			Required:        true,
			AllowWhitespace: false,
		},
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
}
