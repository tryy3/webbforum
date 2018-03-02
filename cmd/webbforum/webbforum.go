package main

import (
	"github.com/tryy3/webbforum"
	"github.com/apex/log"
	loghandler "github.com/tryy3/webbforum/log"
	"os"
	"github.com/tryy3/webbforum/api"
	"github.com/spf13/viper"
	"github.com/jinzhu/gorm"
	"fmt"

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
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@%s:%d/%s",
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.name")))
	if err != nil {
		log.Fatal(err.Error())
	}

	// create the API
	a, err := api.NewAPI(db)
	if err != nil {
		log.Fatal(err.Error())
	}

	// start the http server
	webbforum.StartServer(
		viper.GetString("http.host"),
		viper.GetInt("http.port"),
		a)
}