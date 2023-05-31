package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func NewHTTPOnlyCookie(token string) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     "accessToken",
		Value:    token,
		HTTPOnly: true,
		Expires:  time.Now().Add(2 * 24 * time.Hour),
		Path:     "/",
		Secure:   false,
	}
}
