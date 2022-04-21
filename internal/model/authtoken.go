package model

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

type AuthToken struct {
	Uuid        string `json:"uuid"`
	TokenString string `json:"token"`
	Expires     int64  `json:"exp"`
}

func NewAccessToken(user *User) (*AuthToken, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenUuid := uuid.NewV4().String()
	tokenExpires := time.Now().Add(15 * time.Minute).Unix()

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user_id"] = user.ID
	claims["access_uuid"] = tokenUuid
	claims["admin"] = user.Admin
	claims["exp"] = tokenExpires

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	at := &AuthToken{
		Uuid:        tokenUuid,
		TokenString: tokenString,
		Expires:     tokenExpires,
	}

	return at, nil
}

func NewRefreshToken(user *User) (*AuthToken, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	tokenUuid := uuid.NewV4().String()
	tokenExpires := time.Now().Add(15 * time.Minute).Unix()

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["refresh_uuid"] = tokenUuid
	claims["admin"] = user.Admin
	claims["exp"] = tokenExpires

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	rt := &AuthToken{
		Uuid:        tokenUuid,
		TokenString: tokenString,
		Expires:     tokenExpires,
	}

	return rt, nil
}
