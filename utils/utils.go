package utils

import (
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

func ValidateURL(rawText string) bool {
	var re = regexp.MustCompile(`(\b(https?):\/\/)?[-A-Za-z0-9+&@#\/%?=~_|!:,.;]+\.[-A-Za-z0-9+&@#\/%=~_|]+`)
	return re.Match([]byte(rawText))
}

func ValidateShortURL(rawText string) bool {
	var re = regexp.MustCompile(`http:\/\/localhost:8080\/[a-zA-Z]{5}`)
	return re.Match([]byte(rawText))
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func ShortURLGenerator(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

var digits = []rune("1234567890")

func UserIDGenerator(n int) int {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = digits[rand.Intn(len(digits))]
	}
	res, err := strconv.Atoi(string(b))
	if err != nil {
		log.Println(err)
	}
	return res
}
