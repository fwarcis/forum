package models

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
)

type PasswordHasher interface {
	PasswordAndHashComparator
	Generate(password []byte) ([]byte, error)
}

type PasswordAndHashComparator interface {
	Compare(password, hash []byte) error
}

type User struct {
	login    string
	email    string
	passHash []byte
	session  Session

	passAndHashCmpr PasswordAndHashComparator
	idGener         IDGenerator
}

func NewUser(
	login string,
	email string,
	password []byte,
	passHasher PasswordHasher,
	idGener IDGenerator,
) (*User, error) {
	err := errors.Join(
		checkLogin(login),
		checkEmail(email),
		checkPasswLen(string(password)))
	if err != nil {
		return nil, err
	}

	passHash, err := passHasher.Generate(password)
	if err != nil {
		return nil, err
	}

	session, err := newSession(idGener)
	if err != nil {
		return nil, err
	}
	return &User{login, email, passHash, *session, passHasher, idGener}, err
}

func (u User) CheckPassword(password []byte) error {
	return u.passAndHashCmpr.Compare(u.passHash, password)
}

func (u User) Session() Session {
	return u.session
}

func (u *User) GenerateNewSession() error {
	sessn, err := newSession(u.idGener)
	if err != nil {
		return err
	}
	u.session = *sessn
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

type validationError struct {
	prefix string
	msg    string
}

func (e validationError) Error() string {
	return fmt.Sprintf("%s: %s", e.prefix, e.msg)
}

var (
	ErrUserExists   = fmt.Errorf("user already exists")
	ErrUserNotFound = fmt.Errorf("user not found")
)
