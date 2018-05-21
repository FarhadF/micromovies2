package vault

import (
	"testing"
	"context"
)

func TestHashService(t *testing.T) {
	svc := NewService()
	ctx := context.Background()
	h, err := svc.Hash(ctx, "password")
	if err != nil {
		t.Errorf("Hash: %s", err)
	}
	ok, err := svc.Validate(ctx, "password", h)
	if err != nil {
		t.Errorf("Validate: %s", err)
	}
	if !ok {
		t.Error("expected valid but its not!")
	}
	ok, err = svc.Validate(ctx, "wrongpass", h)
	if err != nil {
		t.Logf("Validate: %s", err)
	}
	if ok {
		t.Error("expected invalid but its not")
	}
}
