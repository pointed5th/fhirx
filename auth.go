package fhird

import (
	"strings"
	"time"

	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/avatar"
	"github.com/go-pkgz/auth/token"
)

type Auth struct {
	*auth.Service
	Options auth.Opts
}

func NewAuth() *Auth {
	opts := auth.Opts{
		SecretReader:   token.SecretFunc(TokenSecret), // secret reader function
		TokenDuration:  time.Minute * 5,               // token expires in 5 minutes
		CookieDuration: time.Hour * 24,                // cookie expires in 1 day and will enforce re-login
		Issuer:         "fhird-dev",
		URL:            "http://127.0.0.1:9292",
		AvatarStore:    avatar.NewLocalFS("./avatars"),
		Validator:      token.ValidatorFunc(Validator),
	}

	return &Auth{
		Options: opts,
		Service: auth.NewService(opts),
	}
}

func TokenSecret(id string) (string, error) {
	return "secret", nil
}

func Validator(r string, claims token.Claims) bool {
	return claims.User != nil && strings.HasPrefix(claims.User.Name, "dev_")
}
