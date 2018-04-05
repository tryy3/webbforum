package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"time"

	"github.com/aarondl/tpl"
	"github.com/dustin/go-humanize"
	"github.com/justinas/nosurf"
	"github.com/spf13/viper"
	"github.com/volatiletech/authboss"
	"github.com/tryy3/webbforum/models"
)

type translate struct {
	Key *regexp.Regexp
	Value string
}

var translations = []translate {
	{regexp.MustCompile(`invalid username and/or password`), "Ogiltigt användarnamn och/eller lösenord."},
	{regexp.MustCompile(`Does not match password`), "Matchar inte lösenordet."},
	{regexp.MustCompile(`Cannot be blank`), "Får inte vara tom."},
	{regexp.MustCompile(`No whitespace permitted`), "Inga mellanrum är tillåtet."},
	{regexp.MustCompile(`Internal server error`), "Internt serverfel."},
	{regexp.MustCompile(`Your account has been locked\.`), "Ditt konto är låst."},
	{regexp.MustCompile(`Your account has not been confirmed\.`), "Ditt konto har inte aktiverats än."},
	{regexp.MustCompile(`You have logged out`), "Du har loggat ut."},
	{regexp.MustCompile(`You have successfully confirmed your account\.`), "Du har nu aktiverat ditt konto."},
	{regexp.MustCompile(`An email has been sent with further instructions on how to reset your password`), "Ett e-post har skickats med ytterligare instruktioner om hur du återställer ditt lösenord."},
	{regexp.MustCompile(`Account recovery request has expired\. Please try again\.`), "Förfrågan om lösenordsåterställning har löpt ut. Var god försök igen."},
	{regexp.MustCompile(`Account recovery has failed\. Please contact tech support\.`), "Förfrågan om lösenordsåterställning misslyckades. Vänligen kontakta teknisk support."},
	{regexp.MustCompile(`Already in use`), "Används redan"},
	{regexp.MustCompile(`Account successfully created, please verify your e-mail address\.`), "Ditt konto har nu skapats, var god och verifiera din e-postadress."},
	{regexp.MustCompile(`Account successfully created, you are now logged in\.`), "Ditt konto har nu skapats, du är nu inloggad."},
	{regexp.MustCompile(`Must be between (\d+) and (\d+) characters`), "Måste vara mellan $1 och $2 tecken."},
	{regexp.MustCompile(`Must be at least (\d+) characters`), "Måste vara minst $1 tecken."},
	{regexp.MustCompile(`Must be at most (\d+) characters`), "Får vara högst $1 tecken."},
	{regexp.MustCompile(`Must contain at least (\d+) letters`), "Måste innehålla minst $1 bokstäver."},
	{regexp.MustCompile(`Must contain at least (\d+) uppercase letters`), "Måste innehålla minst $1 stora bokstäver."},
	{regexp.MustCompile(`Must contain at least (\d+) lowercase letters`), "Måste innehålla minst $1 små bokstäver."},
	{regexp.MustCompile(`Must contain at least (\d+) numbers`), "Måste innehålla minst $1 siffror."},
	{regexp.MustCompile(`Must contain at least (\d+) symbols`), "Måste innehålla minst $1 symboler."},
	{regexp.MustCompile(`internal error`), "Internt serverfel."},
	{regexp.MustCompile(`Missing category name`), "Kategorins namn saknas."},
	{regexp.MustCompile(`missing category ID`), "Kategorins ID saknas."},
	{regexp.MustCompile(`category ID is a not valid number`), "Kategorins ID är inte ett giltigt nummer."},
	{regexp.MustCompile(`Missing group name`), "Gruppens namn saknas."},
	{regexp.MustCompile(`missing group ID`), "Gruppens ID saknas."},
	{regexp.MustCompile(`group ID is a not valid number`), "Gruppens ID är inte ett giltigt nummer."},
	{regexp.MustCompile(`You do not have permission to do this.`), "Du har inte tillräckligt med rättigheter för att göra detta."},
}

func isSameDay(a, b time.Time) bool {
	aYear, aMonth, aDay:= a.Date()
	bYear, bMonth, bDay:= b.Date()

	if aYear == bYear && aMonth == bMonth && aDay == bDay {
		return true
	}
	return false
}

// LayoutFuncs contains a list of functions that will be used in templates
var LayoutFuncs = template.FuncMap{
	"formatDate": func(date time.Time) string {
		return date.Format("2006/01/02 15:04")
	},
	"humanDate": func(date time.Time) string {
		now := time.Now()
		if isSameDay(date, now) {
			return date.Format("Idag 15:04")
		} else if isSameDay(date, now.AddDate(0, 0, -1)) {
			return date.Format("Igår 15:04")
		}
		return date.Format("2006/01/02 15:04")
	},
	"formatNumber": func(amount int64) string {
		return humanize.Comma(amount)
	},
	"yield": func() string { return "" },
	"map": func(values ...interface{}) (map[string]interface{}, error) {
		if len(values)%2 != 0 {
			return nil, errors.New("invalid map call")
		}

		m := make(map[string]interface{}, len(values)/2)
		for i := 0; i < len(values); i += 2 {
			key, ok := values[i].(string)
			if !ok {
				return nil, errors.New("map keys must be strings")
			}
			m[key] = values[i+1]
		}
		return m, nil
	},
	"translate": func(v string) string {
		for _, t := range translations {
			if t.Key.MatchString(v) {
				return t.Key.ReplaceAllString(v, t.Value)
			}
		}
		return v
	},
}

var (
	Authboss  *authboss.Authboss
	Templates tpl.Templates
)

// LayoutData will retrieve user information for most routes
func LayoutData(w http.ResponseWriter, r *http.Request) authboss.HTMLData {
	currentUserName := ""
	userInter, err := Authboss.CurrentUser(w, r)
	if userInter != nil && err == nil {
		currentUserName = userInter.(*models.User).Username
	}

	return authboss.HTMLData{
		"loggedin":               userInter != nil,
		"username":               "",
		authboss.FlashSuccessKey: Authboss.FlashSuccess(w, r),
		authboss.FlashErrorKey:   Authboss.FlashError(w, r),
		"current_user_name":      currentUserName,
	}
}

// mustrender renders a template
func mustRender(w http.ResponseWriter, r *http.Request, name string, data authboss.HTMLData) {
	data.MergeKV(viper.GetString("xsrf.name"), nosurf.Token(r))
	err := Templates.Render(w, name, data)
	if err == nil {
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, "Error occurred rendering template:", err)
}
