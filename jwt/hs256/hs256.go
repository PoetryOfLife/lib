package hs256

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"time"
)

type Claims struct {
	UserID   int
	UserName string
	jwt.RegisteredClaims
}

const SignKey = "KEY"

var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStr(strLen int) string {
	randBytes := make([]rune, strLen)
	for i := range randBytes {
		randBytes[i] = letters[rand.Intn(len(letters))]
	}
	return string(randBytes)
}

func GenerateTokenUsingHs256(userID int, userName string) (string, error) {
	claim := Claims{
		UserID:   userID,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth_server",
			Subject:   userName,
			Audience:  jwt.ClaimStrings{"app"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        randStr(10),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(SignKey))
	return token, err
}

func ParseTokenHS256(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SignKey), nil //返回签名密钥
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("claim invalid")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid claim type")
	}
	return claims, nil
}
