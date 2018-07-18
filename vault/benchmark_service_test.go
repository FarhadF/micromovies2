package vault

import (
	"testing"
	"context"
	"math/rand"
	"time"
)

func BenchmarkHash(b *testing.B) {
	svc := NewService()
	ctx := context.Background()

	var p []string
	for i := 0; i< 200000; i++ {
		password := RandStringBytesMaskImprSrc(8)
		p = append(p, password)

	}
	rand.Seed(time.Now().Unix())
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := svc.Hash(ctx, p[rand.Intn(len(p))])
		if err != nil {
			b.Errorf("Hash: %s", err)
		}
	}
}

func BenchmarkValidate(b *testing.B) {
	svc := NewService()
	ctx := context.Background()

	var p []string
	for i := 0; i< 200000; i++ {
		password := RandStringBytesMaskImprSrc(8)
		p = append(p, password)

	}
	rand.Seed(time.Now().Unix())
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := svc.Validate(ctx, p[rand.Intn(len(p))], "$2a$10$3U0jF7G.7BbFHoP.t2cI2uvoPzQtfOIZ7AXwTFo/Kh44VlV1ALWaS")
		if err != nil{
		if err.Error() != "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			b.Errorf("Hash: %s", err)
		}
		}
	}
}


var src = rand.NewSource(time.Now().UnixNano())
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}