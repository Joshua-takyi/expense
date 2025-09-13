package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UseClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	*jwt.StandardClaims
}

func HashPassword(password string) (string, error) {
	bcryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bcryptedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(claims *UseClaims, secretKey string) (string, error) {
	claims.StandardClaims = &jwt.StandardClaims{
		Issuer:    "expensetracker",
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(24 * 7 * time.Hour).Unix(), // Token expires in 7 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateToken(tokenStr, secretKey string) (*UseClaims, error) {
	claims := &UseClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}

func GenerateCsrfToken() (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(tokenBytes), nil
}

func ValidateCsfrToken(token, expectedToken string) bool {
	return token == expectedToken
}

func IsStrongPassword(password string) bool {
	var (
		hasMinLen  = len(password) >= 8
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()-_=+[]{}|;:',.<>?/", char):
			hasSpecial = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func ParseInt(n string) int {
	x, err := strconv.Atoi(n)
	if err != nil {
		return 0
	}
	return x
}
