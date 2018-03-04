package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/aarondl/tpl"
	"github.com/justinas/nosurf"
	"github.com/spf13/viper"
	"github.com/tryy3/webbforum/models"
	"github.com/volatiletech/authboss"
)

// LayoutFuncs contains a list of functions that will be used in templates
var LayoutFuncs = template.FuncMap{
	"formatDate": func(date time.Time) string {
		return date.Format("2006/01/02 03:04pm")
	},
	"yield": func() string { return "" },
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
