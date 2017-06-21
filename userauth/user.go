package userauth

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/tomatorpg/tomatorpg/models"
)

// LoginErrorType represents the type of LoginError
type LoginErrorType int

// Error implements error interface
func (err LoginErrorType) Error() string {
	return "login error: " + err.String()
}

// String implements Stringer interface
func (err LoginErrorType) String() string {
	switch err {
	case ErrNoEmail:
		return "no email"
	case ErrDatabase:
		return "database error"
	}
	return "unknown error"
}

// GoString implements GoStringer interface
func (err LoginErrorType) GoString() string {
	return "LoginError(\"" + err.String() + "\")"
}

const (
	// ErrUnknown represents all unknown errors
	ErrUnknown LoginErrorType = iota

	// ErrNoEmail happens if the login user did not provide email
	ErrNoEmail

	// ErrDatabase represents database server errors
	ErrDatabase
)

// LoginError is a class of errors occurs in login
type LoginError struct {
	Type LoginErrorType

	// Action triggering of the error
	Action string

	// Err stores inner error type, if any
	Err error
}

// Error implements error interface
func (err LoginError) Error() string {
	return "login error: " + err.String()
}

// String implements error interface
func (err LoginError) String() string {
	if err.Err == nil {
		return err.Type.String()
	}
	if err.Action == "" {
		return fmt.Sprintf(
			"error=%#v",
			err.Err.Error(),
		)
	}
	return fmt.Sprintf(
		"action=%#v error=%#v",
		err.Action,
		err.Err.Error(),
	)
}

// GoString implements GoStringer interface
func (err LoginError) GoString() string {
	if err.Err == nil {
		return "LoginError(\"" + err.String() + "\")"
	}
	return "LoginError(" + err.String() + ")"
}

func loadOrCreateUser(db *gorm.DB, authUser models.User, verifiedEmails []string) (confirmedUser *models.User, err error) {

	// search existing user with the email
	var userEmail models.UserEmail
	var prevUser models.User

	if db.First(&prevUser, "primary_email = ?", authUser.PrimaryEmail); prevUser.PrimaryEmail != "" {
		// TODO: log this?
		authUser = prevUser
	} else if db.First(&userEmail, "email = ?", authUser.PrimaryEmail); userEmail.Email != "" {
		// TODO: log this?
		db.First(&authUser, "id = ?", userEmail.UserID)
	} else if authUser.PrimaryEmail == "" {
		err = &LoginError{Type: ErrNoEmail}
	} else {

		tx := db.Begin()

		// create user
		if res := tx.Create(&authUser); res.Error != nil {
			// append authUser to error info
			err = &LoginError{
				Type:   ErrDatabase,
				Action: "create user",
				Err:    res.Error,
			}
			tx.Rollback()
			return
		}

		// create user-email relation
		newUserEmail := models.UserEmail{
			UserID: authUser.ID,
			Email:  authUser.PrimaryEmail,
		}
		if res := tx.Create(&newUserEmail); res.Error != nil {
			// append newUserEmail to error info
			err = &LoginError{
				Type:   ErrDatabase,
				Action: "create user-email relation " + newUserEmail.Email,
				Err:    res.Error,
			}
			tx.Rollback()
			return
		}

		// also input UserEmail from verifiedEmails, if len not 0
		for _, email := range verifiedEmails {
			newUserEmail := models.UserEmail{
				UserID: authUser.ID,
				Email:  email,
			}
			if res := tx.Create(&newUserEmail); res.Error != nil {
				// append newUserEmail to error info
				err = &LoginError{
					Type:   ErrDatabase,
					Action: "create user-email relation " + newUserEmail.Email,
					Err:    res.Error,
				}
				tx.Rollback()
				return
			}
		}

		tx.Commit()
	}

	if err == nil {
		confirmedUser = &authUser
	}
	return
}
