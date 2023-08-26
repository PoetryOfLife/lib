package rs256

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"time"
)

const pri_key = `-----BEGIN RSA PRIVATE KEY-----
MIIBOwIBAAJBAL00QsML/ovZle3Lq3C7QBo9s00ivsLhG2xlamhHOZDrjTGJX4OA
H27qQbDREcYXpUt5JqOt+KzB4MA/vUKCbT0CAwEAAQJBAINbkS5RWXxGqCzcRj6S
AkM1qxJWmRI7rwpmrqWPLYxKiS1i/i3bwSA3H+NODWIk1p2BWtycWzx5s3cNLn4b
gIECIQD6WuNzXxZHRIxRJQDRyEeWLsrRv9nkZJXHde78DoIZuQIhAMF4ZOgQX2hV
+y9YZmca2tW7etwGPmVjFWQd6JFtjyGlAiBFR9GZo76uijGqYusPIrVswhYuZUEP
CybHw8MWzY0DQQIgc4DDDWCo9QtP+MYX7Lo1p6BUCwOXQMRUwv6wGBKGfxkCIQDn
EKF3Ee6bnLT5DMfrnGY20RNg1Yes+14KkEyYsx0++Q==
-----END RSA PRIVATE KEY-----
`

const pub_key = `-----BEGIN RSA PUBLIC KEY-----
MEgCQQC9NELDC/6L2ZXty6twu0AaPbNNIr7C4RtsZWpoRzmQ640xiV+DgB9u6kGw
0RHGF6VLeSajrfisweDAP71Cgm09AgMBAAE=
-----END RSA PUBLIC KEY-----
`

type Claims struct {
	UserID   int
	UserName string
	jwt.RegisteredClaims
}

var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStr(strLen int) string {
	randBytes := make([]rune, strLen)
	for i := range randBytes {
		randBytes[i] = letters[rand.Intn(len(letters))]
	}
	return string(randBytes)
}

func parsePriKeyBytes(buf []byte) (*rsa.PrivateKey, error) {
	p := &pem.Block{}
	p, buf = pem.Decode(buf)
	if p == nil {
		return nil, errors.New("parse key error")
	}
	return x509.ParsePKCS1PrivateKey(p.Bytes)
}

func GenerateTokenUsingRS256(userID int, userName string) (string, error) {
	claim := Claims{
		UserID:   userID,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth_server",
			Subject:   userName,
			Audience:  jwt.ClaimStrings{"app"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Second)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        randStr(10),
		},
	}
	rsaPriKey, err := parsePriKeyBytes([]byte(pri_key))
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claim).SignedString(rsaPriKey)
	return token, err
}

func parsePubKeyBytes(pubKey []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return nil, errors.New("block nil")
	}
	pubRet, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pubRet, nil
}

func ParseTokenRs256(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		pub, err := parsePubKeyBytes([]byte(pub_key))
		if err != nil {
			return nil, err
		}
		return pub, nil
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
