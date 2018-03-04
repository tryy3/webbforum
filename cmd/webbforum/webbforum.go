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

func main() {
	setupLogger()
	configureViper()
	defaultConfig()

	if err := readConfig(); err != nil {
		log.Fatal(fmt.Sprintf("error creating config: %s", err.Error()))
	}

	db, err := connectDatabase()
	if err != nil {
		log.Fatal(fmt.Sprintf("error connecting to database: %s", err.Error()))
	}
	db.LogMode(true)

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

func setupLogger() {
	// handler for logging to file
	file := loghandler.NewFile("log_%date%.log")

	// handler for logging to cli
	cli := loghandler.NewCli(os.Stderr)

	// handler for logging to multiple handlers
	handler := loghandler.NewMulti(cli, file)

	// set default handler
	log.SetHandler(handler)
}

func configureViper() {
	// configure the config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
}

func defaultConfig() {
	// default config values
	viper.SetDefault("database.host", "127.0.0.1")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.name", "webbforum")

	viper.SetDefault("http.host", "127.0.0.1")
	viper.SetDefault("http.port", 80)

	viper.SetDefault("xsrf.name", "csrf_token")

	viper.SetDefault("cookie.key", utils.GenerateRandomKey(64))
	viper.SetDefault("cookie.expiry", 24*365)

	viper.SetDefault("session.key", utils.GenerateRandomKey(64))
	viper.SetDefault("session.name", "ab_webbforum")

	viper.SetDefault("smtp.host", "smtp.gmail.com")
	viper.SetDefault("smtp.port", 587)
	viper.SetDefault("smtp.username", "example@gmail.com")
	viper.SetDefault("smtp.password", "example_password")
	viper.SetDefault("smtp.identity", "")
	viper.SetDefault("smtp.email", "webbforum@gmail.com")
	viper.SetDefault("smtp.name", "Webbforum Administrator")

	viper.SetDefault("views.folder", "views")
	viper.SetDefault("views.partials", "views/partials")
}

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
	return db, err
}

func setupDatabase(db *gorm.DB) error {
	// auto migrate all of the models
	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.Group{})
	db.AutoMigrate(&models.Permission{})
	db.AutoMigrate(&models.Post{})
	db.AutoMigrate(&models.Thread{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Token{})

	// insert default permissions, update if permission has changed
	perms := models.DefaultPermission.Permissions()
	tx := db.Begin()
	for k, v := range perms {
		err := tx.Set("gorm:insert_option", "ON DUPLICATE KEY UPDATE name = VALUES(name)").Create(&models.Permission{Bit: v, Name: k}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}
