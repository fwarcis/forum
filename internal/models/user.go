package models

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Session struct {
	id        uuid.UUID
	expiresAt time.Time
}

func newSession() (Session, error) {
	uuid_, err := uuid.NewRandom()
	if err != nil {
		return Session{}, err
	}

	return Session{
		id:        uuid_,
		expiresAt: time.Now().Add(1 * time.Hour),
	}, nil
}

func (s Session) ID() uuid.UUID {
	return s.id
}

func (s Session) IsExpired() bool {
	return s.expiresAt.Before(time.Now())
}

func (s Session) ExpiresAt() time.Time {
	return s.expiresAt
}

type User struct {
	login     string
	email     string
	passwHash []byte

	session Session
}

func NewUser(
	login string,
	email string,
	password []byte,
) (User, error) {
	err := errors.Join(
		checkLogin(login),
		checkEmail(email),
		checkPasswLen(string(password)))
	if err != nil {
		return User{}, err
	}

	passwHash, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	sessn, err := newSession()
	if err != nil {
		return User{}, err
	}
	return User{login, email, passwHash, sessn}, err
}

func (u User) CheckPassword(passw []byte) error {
	return checkPasswWithHash(u.passwHash, passw)
}

func (u User) Session() Session {
	return u.session
}

func (u *User) GenerateNewSession() error {
	sessn, err := newSession()
	if err != nil {
		return err
	}
	u.session = sessn
	return nil
}

var loginCharsLimitPatterm = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]*$")

func checkLogin(login string) error {
	length := len([]rune(login))
	switch {
	case length < 3:
		return &validationError{"login", "min length = 3"}
	case length > 20:
		return &validationError{"login", "max length = 20"}
	}
	return nil
}

// RFC 5322
func checkEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return validationError{prefix: "email", msg: err.Error()}
	}
	return nil
}

// NIST SP 800-63B
func checkPasswLen(passw string) error {
	length := len([]rune(passw))
	switch {
	case length < 16:
		return validationError{"password", "min length = 16"}
	case length > 128:
		return validationError{"password", "max length = 127"}
	}
	return nil
}

func checkPasswWithHash(hash, passw []byte) error {
	return bcrypt.CompareHashAndPassword(hash, passw)
}

type validationError struct {
	prefix string
	msg    string
}

func (e validationError) Error() string {
	return fmt.Sprintf("%s: %s", e.prefix, e.msg)
}
