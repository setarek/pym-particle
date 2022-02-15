package utils

import (
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"
)

var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func CheckUrlValidation(link string) bool {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return false
	}
	return true
}

func GenerateRandString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func ParseString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case nil:
		return ""
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func ParseInt64(value interface{}) int64 {
	switch v := value.(type) {
	case string:
		val, _ := strconv.Atoi(v)
		if val >= 0 {
			return int64(val)
		}
		return 0
	case float64:
		if v >= 0 {
			return int64(v)
		}
		return 0
	case uint:
		return int64(v)
	case int:
		if v >= 0 {
			return int64(v)
		}
		return 0
	case int64:
		return v
	default:
		return 0
	}
}
