package handlers

import (
	"fmt"
	"html"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/CalebQ42/bbConvert"
	"github.com/apex/log"
	"github.com/volatiletech/authboss"
)

func badRequest(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest)
	log.WithError(err).Error("bad request")

	return true
}

func getUser(w http.ResponseWriter, r *http.Request, ab *authboss.Authboss) (interface{}, bool) {
	u, err := ab.CurrentUser(w, r)
	if err != nil {
		log.WithError(err).Error("error fetching current user")
		w.WriteHeader(http.StatusInternalServerError)
		return nil, false
	} else if u == nil {
		log.Errorf("redirecting unauthorized user from: %s", r.URL.Path)
		http.Redirect(w, r, "/", http.StatusFound)
		return nil, false
	}
	return u, true
}

func getCategoryID(attr authboss.Attributes) (uint, string) {
	idStr, ok := attr.String("category_id")
	if !ok {
		return 0, "missing category ID"
	}

	id, err := strconv.ParseUint(idStr, 32, 10)
	if err != nil {
		return 0, "category ID is a not valid number"
	}

	return uint(id), ""
}

func getGroupID(attr authboss.Attributes) (uint, string) {
	idStr, ok := attr.String("group_id")
	if !ok {
		return 0, "missing group ID"
	}

	id, err := strconv.ParseUint(idStr, 32, 10)
	if err != nil {
		return 0, "group ID is a not valid number"
	}

	return uint(id), ""
}

func getThreadID(attr authboss.Attributes) (uint, string) {
	idStr, ok := attr.String("thread_id")
	if !ok {
		return 0, "missing thread ID"
	}

	id, err := strconv.ParseUint(idStr, 32, 10)
	if err != nil {
		return 0, "thread ID is a not valid number"
	}

	return uint(id), ""
}

func urlEncoded(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}


type BBCodeConverter struct {
	conv      *bbConvert.Converter
	emoticons map[string]string
}

func (b BBCodeConverter) Convert(str string) template.HTML {
	c := b.conv.Convert(str)
	c = strings.Replace(c, "\n", "<br />", -1)

	rep := make([]string, len(b.emoticons)*2)
	i := 0
	for k, v := range b.emoticons {
		rep[i] = k
		rep[i+1] = v
		i += 2
	}

	rep = append(rep, "\n", "<br />", "[hr]", "<hr />")

	c = strings.NewReplacer(rep...).Replace(c)

	return template.HTML(c)
}

func NewBBCodeConverter() *BBCodeConverter {
	var htmlConv bbConvert.HTMLConverter
	htmlConv.ImplementDefaults()

	conv := htmlConv.Converter()
	conv.AddCustom("sub", func(_ bbConvert.Tag, meat string) string {
		return "<sub>" + meat + "</sub>"
	})
	conv.AddCustom("sup", func(_ bbConvert.Tag, meat string) string {
		return "<sup>" + meat + "</sup>"
	})
	conv.AddCustom("left", func(_ bbConvert.Tag, meat string) string {
		return "<div align=\"left\">" + meat + "</div>"
	})
	conv.AddCustom("center", func(_ bbConvert.Tag, meat string) string {
		return "<div align=\"center\">" + meat + "</div>"
	})
	conv.AddCustom("right", func(_ bbConvert.Tag, meat string) string {
		return "<div align=\"right\">" + meat + "</div>"
	})
	conv.AddCustom("justify", func(_ bbConvert.Tag, meat string) string {
		return "<div align=\"justify\">" + meat + "</div>"
	})
	conv.AddCustom("table", func(_ bbConvert.Tag, meat string) string {
		return "<table>" + meat + "</table>"
	})
	conv.AddCustom("ul", func(_ bbConvert.Tag, meat string) string {
		return "<ul>" + meat + "</ul>"
	})
	conv.AddCustom("li", func(_ bbConvert.Tag, meat string) string {
		return "<li>" + meat + "</li>"
	})
	conv.AddCustom("td", func(_ bbConvert.Tag, meat string) string {
		return "<td>" + meat + "</td>"
	})
	conv.AddCustom("tr", func(_ bbConvert.Tag, meat string) string {
		return "<tr>" + meat + "</tr>"
	})
	conv.AddCustom("code", func(_ bbConvert.Tag, meat string) string {
		return "<code>" + meat + "</code>"
	})
	conv.AddCustom("email", func(t bbConvert.Tag, meat string) string {
		email := t.Value("starting")
		if email != "" {
			email = meat
		}
		return "<a href=\"mailto:" + html.EscapeString(email) + "\">" + meat + "</a>"
	})
	conv.AddCustom("quote", func(_ bbConvert.Tag, meat string) string {
		return "<blockquote>" + meat + "</blockquote>"
	})
	conv.AddCustom("rtl", func(_ bbConvert.Tag, meat string) string {
		return "<div style=\"direction: rtl\">" + meat + "</div>"
	})
	conv.AddCustom("ltr", func(_ bbConvert.Tag, meat string) string {
		return "<div style=\"direction: ltr\">" + meat + "</div>"
	})

	emoticons := map[string]string{
		":)":          "emoticons/smile.png",
		":angel:":     "emoticons/angel.png",
		":angry:":     "emoticons/angry.png",
		"8-)":         "emoticons/cool.png",
		":'(":         "emoticons/cwy.png",
		":ermm:":      "emoticons/ermm.png",
		":D":          "emoticons/grin.png",
		"<3":          "emoticons/heart.png",
		":(":          "emoticons/sad.png",
		":O":          "emoticons/shocked.png",
		":P":          "emoticons/tongue.png",
		";)":          "emoticons/wink.png",
		":alien:":     "emoticons/alien.png",
		":blink:":     "emoticons/blink.png",
		":blush:":     "emoticons/blush.png",
		":cheerful:":  "emoticons/cheerful.png",
		":devil:":     "emoticons/devil.png",
		":dizzy:":     "emoticons/dizzy.png",
		":getlost:":   "emoticons/getlost.png",
		":happy:":     "emoticons/happy.png",
		":kissing:":   "emoticons/kissing.png",
		":ninja:":     "emoticons/ninja.png",
		":pinch:":     "emoticons/pinch.png",
		":pouty:":     "emoticons/pouty.png",
		":sick:":      "emoticons/sick.png",
		":sideways:":  "emoticons/sideways.png",
		":silly:":     "emoticons/silly.png",
		":sleeping:":  "emoticons/sleeping.png",
		":unsure:":    "emoticons/unsure.png",
		":woot:":      "emoticons/w00t.png",
		":wassat:":    "emoticons/wassat.png",
		":whistling:": "emoticons/whistling.png",
		":love:":      "emoticons/wub.png",
	}

	emoticonsRoot := "/images/"
	for key, value := range emoticons {
		emoticons[key] = fmt.Sprintf("<img src=\"%s\" data-sceditor-emoticon=\"%s\" alt=\"%s\" title=\"%s\" />",
			emoticonsRoot+value,
			key,
			key,
			key,
		)
	}

	return &BBCodeConverter{conv, emoticons}
}
