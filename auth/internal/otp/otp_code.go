package otp

import (
	"math/rand"
	"time"
)

var Store = make(map[string]int32)

func StoreOTP(email string, otp int32) {
	Store[email] = otp
	go func() {
		time.Sleep(2 * time.Minute)
		delete(Store, email)
	}()
}

func VerifyOTP(email string, otp int32) bool {
	storedCode, exists := Store[email]
	return exists && storedCode == otp
}

func Generate() int32 {
	return int32(rand.Intn(1000000))
}
