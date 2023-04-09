package utils

import (
	"bytes"
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type JwtToker struct {
	privKey *rsa.PrivateKey
	pubKey  *rsa.PublicKey
}

type TokenClaims struct {
	UserID  string `json:"sub"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

func NewJwtToker(privRsa []byte, pubRsa []byte) (JwtToker, error) {
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(privRsa)
	if err != nil {
		return JwtToker{}, err
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(pubRsa)
	if err != nil {
		return JwtToker{}, err
	}

	return JwtToker{privKey: signKey, pubKey: verifyKey}, nil
}

func (tkr JwtToker) CreateToken(userID string, isAdmin bool) (string, error) {
	claims := &TokenClaims{
		userID,
		isAdmin,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)
	return token.SignedString(tkr.privKey)
}

func HashPassword(password string) (string, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return bytes.NewBuffer(passwordBytes).String(), nil
}

func ValidatePassword(password string, storedHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
}
