package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/gorilla/securecookie"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/tryy3/webbforum"
	"github.com/tryy3/webbforum/api"
	loghandler "github.com/tryy3/webbforum/log"
	"github.com/tryy3/webbforum/utils"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// handler for logging to file
	file := loghandler.NewFile("log_%date%.log")

	// handler for logging to cli
	cli := loghandler.NewCli(os.Stderr)

	// handler for logging to multiple handlers
	handler := loghandler.NewMulti(cli, file)

	// set default handler
	log.SetHandler(handler)

	// configure the config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")

	// default config values
	viper.SetDefault("database.host", "127.0.0.1")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.name", "webbforum")

	viper.SetDefault("http.host", "127.0.0.1")
	viper.SetDefault("http.port", 80)

	viper.SetDefault("xsrf.name", "csrf_token")

	viper.SetDefault("cookie.key", securecookie.GenerateRandomKey(64))
	viper.SetDefault("cookie.expiry", 24*365)

	viper.SetDefault("session.key", securecookie.GenerateRandomKey(64))
	viper.SetDefault("session.name", "ab_webbforum")

	viper.SetDefault("smtp.host", "smtp.gmail.com:587")
	viper.SetDefault("smtp.username", "example@gmail.com")
	viper.SetDefault("smtp.password", "example_password")
	viper.SetDefault("smtp.identity", "")
	viper.SetDefault("smtp.email", "webbforum@gmail.com")
	viper.SetDefault("smtp.name", "Webbforum Administrator")

	// try to read the config
	err := viper.ReadInConfig()
	if err != nil {
		// write to config if config file does not exists
		err := viper.WriteConfigAs("config.json")
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	// connect to mysql database
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true",
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.name")))
	if err != nil {
		log.Fatal(err.Error())
	}

	// generate cookie key
	cookieKey, err := base64.StdEncoding.DecodeString(viper.GetString("cookie.key"))
	if err != nil {
		log.Fatal(err.Error())
	}

	// generate session key
	sessionKey, err := base64.StdEncoding.DecodeString(viper.GetString("session.key"))
	if err != nil {
		log.Fatal(err.Error())
	}

	// create context Config
	config := &utils.Config{
		HTTPIP:   viper.GetString("http.host"),
		HTTPPort: viper.GetInt("http.port"),

		XSRFName: viper.GetString("xsrf.name"),

		CookieStoreKey: cookieKey,
		CookieExpiry:   time.Duration(viper.GetInt64("cookie.expiry")) * time.Hour,

		SessionStoreKey: sessionKey,
		SessionName:     viper.GetString("session.name"),

		SMTPHost:     viper.GetString("smtp.host"),
		SMTPUsername: viper.GetString("smtp.username"),
		SMTPPassword: viper.GetString("smtp.password"),
		SMTPIdentity: viper.GetString("smtp.identity"),
		SMTPEmail:    viper.GetString("smtp.email"),
		SMTPName:     viper.GetString("smtp.name"),
	}

	// create context object
	context := &utils.Context{
		Config:   config,
		Database: db,
	}

	// create the API
	a, err := api.NewAPI(context.Database)
	if err != nil {
		log.Fatal(err.Error())
	}

	// start the http server
	webbforum.StartServer(context, a)
}
