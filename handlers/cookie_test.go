package handlers

import (
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
)

func BenchmarkCookie(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	b.Run("generate", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer() // останавливаем таймер
			userID := uuid.New()
			b.StartTimer() // возобновляем таймер
			GenerateCookie(userID)
		}
	})

	b.Run("parse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer() // останавливаем таймер
			userID := uuid.New()
			cookie, _ := GenerateCookie(userID)
			cookieStr := cookie.String()
			b.StartTimer() // возобновляем таймер
			ParseCookie(cookieStr)
		}
	})
}
