package main

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/tryy3/webbforum"
	loghandler "github.com/tryy3/webbforum/log"
	"github.com/tryy3/webbforum/models"
	"github.com/tryy3/webbforum/utils"

	_ "github.com/go-sql-driver/mysql"
)

// main function is the function that will be ran first when the program is ran
func main() {
	// Initialize all configurations and components
	setupLogger()
	configureViper()
	defaultConfig()

	// attempt to read the config
	if err := readConfig(); err != nil {
		log.Fatal(fmt.Sprintf("error creating config: %s", err.Error()))
	}

	// attempt to connect to the database
	db, err := connectDatabase()
	if err != nil {
		log.Fatal(fmt.Sprintf("error connecting to database: %s", err.Error()))
	}
	db.LogMode(true)

	// attempt to initialize setup tasks on the database
	err = setupDatabase(db)
	if err != nil {
		log.Fatal(fmt.Sprintf("error setup database: %s", err.Error()))
	}

	// start the http server
	err = webbforum.StartServer(db)
	if err != nil {
		log.Fatal(fmt.Sprintf("error starting the server: %s", err.Error()))
	}
}

// setupLogger configures default settings on the logger
func setupLogger() {
	// handler for logging to file
	file := loghandler.NewFile("log_%date%.log")

	// handler for logging to cli
	cli := loghandler.NewCli(os.Stderr)

	// handler for logging to multiple handlers
	handler := loghandler.NewMulti(cli, file)

	// set default handler and logging level
	log.SetHandler(handler)
	log.SetLevel(log.DebugLevel)
}

// configureViper configure the config system
func configureViper() {
	// configure the config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
}

// defaultConfig initialize default config values
func defaultConfig() {
	// database settings
	viper.SetDefault("database.host", "127.0.0.1")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.name", "webbforum")

	// http server settings
	viper.SetDefault("http.host", "127.0.0.1")
	viper.SetDefault("http.port", 80)

	// xsrf protection value
	viper.SetDefault("xsrf.name", "csrf_token")

	// cookie settings
	viper.SetDefault("cookie.key", utils.GenerateRandomKey(64))
	viper.SetDefault("cookie.expiry", 24*365)

	// local session settings
	viper.SetDefault("session.key", utils.GenerateRandomKey(64))
	viper.SetDefault("session.name", "ab_webbforum")

	// smtp settings
	viper.SetDefault("smtp.host", "smtp.gmail.com")
	viper.SetDefault("smtp.port", 587)
	viper.SetDefault("smtp.username", "example@gmail.com")
	viper.SetDefault("smtp.password", "example_password")
	viper.SetDefault("smtp.identity", "")
	viper.SetDefault("smtp.email", "webbforum@gmail.com")
	viper.SetDefault("smtp.name", "Webbforum Administrator")

	// views settings
	viper.SetDefault("views.folder", "views")
	viper.SetDefault("views.partials", "views/partials")

	// http content settings
	viper.SetDefault("content.base", "content")
	viper.SetDefault("content.tmp", "tmp")
	viper.SetDefault("content.image.folder", "image")
	viper.SetDefault("content.image.size", 10*1024*1024) // 10 MB
	viper.SetDefault("content.js.folder", "js")
	viper.SetDefault("content.css.folder", "css")
}

// readConfig will attempt to read existing config file, if file doesn't exists it will create one with default settings
func readConfig() error {
	// try to read the config
	err := viper.ReadInConfig()
	if err != nil {
		// write to config if config file does not exists
		err := viper.WriteConfigAs("config.json")
		if err != nil {
			return err
		}
	}
	return nil
}

// connectDatabase attempts to connect to mysql
func connectDatabase() (*gorm.DB, error) {
	// connect to mysql database
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true&charset=utf8",
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.name")))
	if err != nil {
		return nil, err
	}

	db = db.Set("gorm:auto_preload", true)
	return db, err
}

// setupDatabase takes care of the default table migrations
func setupDatabase(db *gorm.DB) error {
	// auto migrate all of the models
	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.Group{})
	db.AutoMigrate(&models.Permission{})
	db.AutoMigrate(&models.Post{})
	db.AutoMigrate(&models.Thread{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Token{})
	db.AutoMigrate(&models.File{})
	db.AutoMigrate(&models.Permission{})
	return db.Error
}
