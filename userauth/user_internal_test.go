package userauth

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/tomatorpg/tomatorpg/models"
)

type NopLogwriter int

func (logger NopLogwriter) Println(v ...interface{}) {
	// don't give a damn
}

func TestLoadOrCreateUser(t *testing.T) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	defer db.Close()
	db.AutoMigrate(
		models.User{},
		models.UserEmail{},
	)

	// attempt to create user on first login
	u1, err := loadOrCreateUser(
		db,
		models.User{
			Name:         "dummy user",
			PrimaryEmail: "dummy@foobar.com",
		},
		[]string{},
	)

	if err != nil {
		t.Errorf("unexpected error: %#v", err)
	}
	if u1 == nil {
		t.Errorf("expected user, got nil")
	}

	u1db := models.User{}
	db.First(&u1db, "id = ?", u1.ID)
	if want, have := u1.ID, u1db.ID; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := u1.Name, u1db.Name; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := u1.PrimaryEmail, u1db.PrimaryEmail; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	// should retrieve the same user on second login
	// regardless of the name
	u2, err := loadOrCreateUser(
		db,
		models.User{
			Name:         "dummy user another time",
			PrimaryEmail: "dummy@foobar.com",
		},
		[]string{},
	)
	if want, have := u1.ID, u2.ID; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := u1.Name, u2.Name; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := u1.PrimaryEmail, u2.PrimaryEmail; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	// try to login with no email
	u3, err := loadOrCreateUser(
		db,
		models.User{
			Name:         "dummy user",
			PrimaryEmail: "",
		},
		[]string{},
	)
	if u3 != nil {
		t.Errorf("expected u3 to be nil, got %#v", u3)
	}
	if want, have := ErrNoEmail, err.(*LoginError).Type; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestLoadOrCreateUser_UserEmail(t *testing.T) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	defer db.Close()
	db.AutoMigrate(
		models.User{},
		models.UserEmail{},
	)

	// attempt to create user on first login
	u1, err := loadOrCreateUser(
		db,
		models.User{
			Name:         "dummy user",
			PrimaryEmail: "dummy1@foobar.com",
		},
		[]string{"dummy2@foobar.com"},
	)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if u1 == nil {
		t.Errorf("expected user, got nil")
	}

	u2, err := loadOrCreateUser(
		db,
		models.User{
			Name:         "dummy user",
			PrimaryEmail: "dummy2@foobar.com",
		},
		[]string{},
	)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if u2 == nil {
		t.Errorf("expected user, got nil")
	}

	// both attempts should result the same user
	if want, have := u1.ID, u2.ID; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := u1.Name, u2.Name; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := u1.PrimaryEmail, u2.PrimaryEmail; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

}

func TestLoadOrCreateUser_DatabaseError(t *testing.T) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	defer db.Close()
	db.SetLogger(gorm.Logger{LogWriter: NopLogwriter(0)})

	// test inserting with a missing table
	db.AutoMigrate(
		models.User{},
		models.UserEmail{},
	)
	db.Exec("DROP TABLE user_emails;")
	u1, err := loadOrCreateUser(
		db,
		models.User{
			Name:         "dummy user",
			PrimaryEmail: "dummy@foobar.com",
		},
		[]string{},
	)
	if u1 != nil {
		t.Errorf("expected u1 to be nil, got %#v", u1)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	u1db := models.User{}
	db.First(&u1db, "email = ?", "dummy@foobar.com")
	if u1db.ID != 0 {
		t.Errorf("expected u1db to have id=0, got %#v", u1db)
	}

	// test inserting with a duplicated ValidatedEmail
	db.AutoMigrate(
		models.User{},
		models.UserEmail{},
	)
	u2, err := loadOrCreateUser(
		db,
		models.User{
			Name:         "dummy user",
			PrimaryEmail: "dummy2@foobar.com",
		},
		[]string{"dummy2@foobar.com"},
	)
	if u2 != nil {
		t.Errorf("expected u1 to be nil, got %#v", u1)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	u2db := models.User{}
	db.First(&u2db, "email = ?", "dummy@foobar.com")
	if u2db.ID != 0 {
		t.Errorf("expected u2db to have id=0, got %#v", u2db)
	}

	db.Exec("DROP TABLE users;")
	u3, err := loadOrCreateUser(
		db,
		models.User{
			Name:         "dummy user",
			PrimaryEmail: "dummy3@foobar.com",
		},
		[]string{},
	)
	if u3 != nil {
		t.Errorf("expected u1 to be nil, got %#v", u1)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	u3db := models.User{}
	db.First(&u3db, "email = ?", "dummy@foobar.com")
	if u3db.ID != 0 {
		t.Errorf("expected u3db to have id=0, got %#v", u3db)
	}

}
