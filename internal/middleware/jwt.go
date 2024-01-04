package middleware

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func JWTProtected() func(*fiber.Ctx) error {
	// Create config for JWT authentication middleware.
	config := jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte("secret")},
		ContextKey:   "jwt", // used in private routes
		ErrorHandler: jwtError,
	}

	return jwtware.New(config)
}

func jwtError(c *fiber.Ctx, err error) error {
	// Return status 401 and failed authentication error.
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 401 and failed authentication error.
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": true,
		"msg":   err.Error(),
	})
}

type TokenWithClaims struct {
	Token   string
	Expires int64
	UserId  int
	Role    string
}

func GenerateNewAccessToken(userId int, role string) (*TokenWithClaims, error) {
	expires := time.Now().Add(time.Hour * 48).Unix()

	claims := jwt.MapClaims{
		"user_id": userId,
		"role":    role,
		"exp":     expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}

	return &TokenWithClaims{
		Token:   t,
		Expires: expires,
		UserId:  userId,
		Role:    role,
	}, nil
}
