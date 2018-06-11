package vault

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/crypto/bcrypt"
)

//Business logic as interface
type Service interface {
	Hash(ctx context.Context, password string) (string, error)
	Validate(ctx context.Context, password string, hash string) (bool, error)
}

//implementation with empty struct (stateless)
type vaultService struct {
}

//constructor - we can later add initialization if needed
func NewService() Service {
	return vaultService{}
}

//implementation
func (vaultService) Hash(ctx context.Context, password string) (string, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := span.Tracer().StartSpan("Hash", opentracing.ChildOf(span.Context()))
		span.SetTag("email", ctx.Value("email"))
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

//implementation
func (vaultService) Validate(ctx context.Context, password string, hash string) (bool, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := span.Tracer().StartSpan("Validate", opentracing.ChildOf(span.Context()))
		span.SetTag("email", ctx.Value("email"))
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
