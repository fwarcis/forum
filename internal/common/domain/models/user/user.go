package user

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"regexp"

	"forum/internal/common/domain/infras/passhash"
	"forum/internal/common/domain/models/session"
)

var (
	ErrExists   = fmt.Errorf("user already exists")
	ErrNotFound = fmt.Errorf("user not found")
)

type PasswordHasher interface {
	Generate(pass string) (hash string, err error)
	Compare(pass, hash string) error
}

type User struct {
	login    string
	email    string
	passHash string
	session  session.Session

	passHasher PasswordHasher
}

func New(
	login string,
	email string,
	password string,
	passHasher PasswordHasher,
) (*User, error) {
	err := errors.Join(
		checkLogin(login),
		checkEmail(email),
		checkPassLen(string(password)))
	if err != nil {
		return nil, err
	}

	passHash, err := passHasher.Generate(password)
	if err != nil {
		return nil, err
	}

	session, err := session.New()
	if err != nil {
		return nil, err
	}
	return &User{login, email, passHash, *session, passHasher}, err
}

func (u User) VerifyPassword(pass string) error {
	return u.passHasher.Compare(u.passHash, pass)
}

func (u User) Session() session.Session {
	return u.session
}

func (u *User) GenerateNewSession() error {
	sessn, err := session.New()
	if err != nil {
		return err
	}
	u.session = *sessn
	return nil
}

var (
	loginCharsPattern     = "[a-zA-Z][a-zA-Z0-9_]"
	loginCharsLimitRegExp = regexp.MustCompile(
		fmt.Sprintf("^%s*$", loginCharsPattern))
)

func checkLogin(login string) error {
	const minLen = 3
	const maxLen = 20
	length := len([]rune(login))
	switch {
	case length < minLen:
		return fmt.Errorf("login min length is %d characters", minLen)
	case length > maxLen:
		return fmt.Errorf("login max length is %d characters", maxLen)
	case loginCharsLimitRegExp.Match([]byte(login)):
		return errors.New("follow the login pattern: " + loginCharsPattern)
	}
	return nil
}

// RFC 5322
func checkEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("email is invalid")
	}
	return nil
}

func checkPassLen(pass string) error {
	const minLen = 8 // 16
	const maxLen = passhash.MaxPassBytes
	length := len([]rune(pass))
	switch {
	case length < minLen:
		return fmt.Errorf("password min length is %d characters", minLen)
	case length > maxLen:
		return fmt.Errorf("password max length is %d characters", maxLen)
	}
	return nil
}

type key int

var userKey key

func NewContext(ctx context.Context, usr User) context.Context {
	return context.WithValue(ctx, userKey, usr)
}

func FromContext(ctx context.Context) (*User, bool) {
	usr, ok := ctx.Value(userKey).(*User)
	return usr, ok
}
