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
		Expires:  time.Now().Add(7 * (24 * time.Hour)),
		Path:     "/",
		Secure:   false,
	}
}

func EmptyCookie() *fiber.Cookie {
	return &fiber.Cookie{
		Name:     "accessToken",
		Value:    "",
		HTTPOnly: true,
		Expires:  time.Unix(0, 0),
		Path:     "/",
		Secure:   false,
	}
}
