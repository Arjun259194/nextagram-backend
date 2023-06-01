package utils

import (
	"fmt"
	"time"

	"github.com/Arjun259194/nextagram-backend/types"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const SECRET_TOKEN string = "$2b$10$mMRmJQosiGMXaT7tsst31u"

func GenerateToken(id primitive.ObjectID) (string, error) {
	// Create the claims
	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Set the expiration time
	}

	// Sign the token with a secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SECRET_TOKEN))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(SECRET_TOKEN), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

func JWTMiddleware(c *fiber.Ctx) error {
	tokenCookie := c.Cookies("accessToken")
	if tokenCookie == "" {
		errRes := types.NewErrorResponse(fiber.StatusUnauthorized, fmt.Errorf("There is not access token saved on you device"), "User not authorized")
		return c.Status(fiber.StatusUnauthorized).JSON(errRes)
	}

	claim, err := VerifyToken(tokenCookie)
	if err != nil {
		errRes := types.NewErrorResponse(fiber.StatusUnauthorized, err, "User not authorized")
		return c.Status(fiber.StatusUnauthorized).JSON(errRes)
	}

	id := claim["id"].(string)

	// Here we are converting string Id into mongoDB objectID because to make sure that it's a valid objectID and because our Go driver only uses objectID type
	userObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		errRes := types.NewErrorResponse(fiber.StatusBadRequest, err, "User Id not valid")
		return c.Status(fiber.StatusBadRequest).JSON(errRes)
	}

	c.Locals("id", userObjectID)

	return c.Next()
}
