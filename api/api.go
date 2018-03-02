package api

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// CreateAPI takes care of registering all of the API routes
func CreateAPI(api *API, r *mux.Router) {

}

// API contains all of the functions for APIs
type API struct {
	db *gorm.DB
}

func NewAPI(db *gorm.DB) (*API, error) {
	return &API{db}, nil
}
