package auth

import (
	"github.com/jinzhu/gorm"
	"github.com/tryy3/webbforum/models"
	"github.com/volatiletech/authboss"
)

// Storer is the implementation of authboss.Storer and all of the extra store functionalities
type Storer struct {
	db *gorm.DB
}

// New will create a new storer and migrate User and Token with the database
func New(db *gorm.DB) *Storer {
	return &Storer{
		db: db,
	}
}

// Create will create a new user in the database
func (s Storer) Create(key string, attr authboss.Attributes) error {
	var user models.User
	if err := attr.Bind(&user, true); err != nil {
		return err
	}

	return s.db.Create(&user).Error
}

// Put will modify an existing user, if user is not found it will return authboss.ErrUserNotFound
func (s Storer) Put(key string, attr authboss.Attributes) error {
	var user models.User
	if err := attr.Bind(&user, true); err != nil {
		return err
	}

	result := s.db.Model(&user).Where("username = ?", key).Updates(&user)
	if result.RecordNotFound() {
		return authboss.ErrUserNotFound
	}
	return result.Error
}

// Get will return a User from the database, if user is not found it will return authboss.ErrUserNotFound
func (s Storer) Get(key string) (interface{}, error) {
	var user models.User
	result := s.db.Where("username = ?", key).First(&user)
	if result.RecordNotFound() {
		return nil, authboss.ErrUserNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// AddToken will create a new token in the database
func (s Storer) AddToken(key, token string) error {
	var tok = models.Token{
		User:  key,
		Token: token,
	}
	return s.db.Create(&tok).Error
}

// DelTokens will remove all tokens related to a specific user
func (s Storer) DelTokens(key string) error {
	return s.db.Where("user = ?", key).Delete(models.Token{}).Error
}

// UseToken will remove a specific token from the database
func (s Storer) UseToken(key, token string) error {
	result := s.db.Where("user = ? AND token = ?", key, token)
	if result.RowsAffected <= 0 {
		return authboss.ErrTokenNotFound
	}
	return result.Error
}

// ConfirmUser will retrieve a user based on their confirm_token
func (s Storer) ConfirmUser(token string) (interface{}, error) {
	var user models.User
	result := s.db.Where("confirm_token = ?", token).First(&user)
	if result.RecordNotFound() {
		return nil, authboss.ErrTokenNotFound
	}
	return &user, nil
}

// RecoverUser will retrieve a user based on their recover_token
func (s Storer) RecoverUser(token string) (interface{}, error) {
	var user models.User
	result := s.db.Where("recover_token = ?", token).First(&user)
	if result.RecordNotFound() {
		return nil, authboss.ErrTokenNotFound
	}
	return &user, nil
}
