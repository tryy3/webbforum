package webbforum

import (
	"github.com/gorilla/mux"
	"github.com/tryy3/webbforum/handlers"
	"net/http"
	"fmt"
	"time"
	"github.com/apex/log"
	"strconv"
	"github.com/tryy3/webbforum/api"
)

// StartServer takes care of initializing the http server with all of the routes
func StartServer(ip string, port int, a *api.API) {
	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()

	api.CreateAPI(a, s)

	r.HandleFunc("/", handlers.HomeHandler)

	srv := &http.Server{
		Handler: r,
		Addr: fmt.Sprintf("%s:%d", ip, port),

		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}

	log.WithFields(log.Fields{"IP": ip, "Port": strconv.Itoa(port)}).Info("http server started")
	log.Fatal(srv.ListenAndServe().Error())
}