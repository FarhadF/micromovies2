package jwtauth

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"context"
)

const mySigningKey = "Super_Dup3r_S3cret"

type Service interface {
	GenerateToken(ctx context.Context, email string, role string) (string, error)
	ParseToken(ctx context.Context, token string) (map[string]interface{}, error)
}

type jwtService struct {
}

func NewService () Service {
	return jwtService{}
}

func (jwtService) GenerateToken(ctx context.Context, email string, role string) (string, error) {
	// Create the token
	tokenObject := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	tokenObject.Claims = jwt.MapClaims{
		"exp":   time.Now().UTC().Add(time.Hour * time.Duration(1)).Unix(),
		"iat":   time.Now().UTC().Unix(),
		"email": email,
		"role":  role,
	}

	// Sign and get the complete encoded token as a string
	tokenString, err := tokenObject.SignedString([]byte(mySigningKey))
	return tokenString, err
}


func (jwtService) ParseToken (ctx context.Context, myToken string) (map[string]interface{}, error) {
	parsedToken, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySigningKey), nil
	})
	c := parsedToken.Claims.(jwt.MapClaims)
	//fmt.Println(token.Claims)
	if err == nil && parsedToken.Valid {
		return c, nil
	} else {
		return nil, err
	}
}

