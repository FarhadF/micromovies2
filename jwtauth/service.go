package jwtauth

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/opentracing/opentracing-go"
	"time"
)

const mySigningKey = "Super_Dup3r_S3cret"

type Service interface {
	GenerateToken(ctx context.Context, email string, role string) (string, error)
	ParseToken(ctx context.Context, token string) (Claims, error)
}

type jwtService struct {
}

func NewService() Service {
	return jwtService{}
}

func (jwtService) GenerateToken(ctx context.Context, email string, role string) (string, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := span.Tracer().StartSpan("GenerateToken", opentracing.ChildOf(span.Context()))
		span.SetTag("email", email)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
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

type Claims struct {
	Exp   int64  `json:"exp"`
	Iat   int64  `json:"iat"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (jwtService) ParseToken(ctx context.Context, myToken string) (Claims, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := span.Tracer().StartSpan("ParseToken", opentracing.ChildOf(span.Context()))
		span.SetTag("email", ctx.Value("email"))
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	parsedToken, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySigningKey), nil
	})
	if err == nil && parsedToken.Valid {
		c := parsedToken.Claims.(jwt.MapClaims)
		claims := Claims{
			//todo: why is it float64 when I defiend int64 and unix returns int64?!
			Exp:   int64(c["exp"].(float64)),
			Iat:   int64(c["iat"].(float64)),
			Email: c["email"].(string),
			Role:  c["role"].(string),
		}
		return claims, nil
	} else {
		return Claims{}, err
	}
}
