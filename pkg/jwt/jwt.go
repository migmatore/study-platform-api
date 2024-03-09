package jwt

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/migmatore/study-platform-api/internal/core"
)

//func ParseToken(token string) core.TokenMetadata {
//	jwt.Parse(token)
//}

// ExtractTokenMetadata func to extract metadata from JWT.
func ExtractTokenMetadata(c *fiber.Ctx) core.TokenMetadata {
	jwtCtx := c.Locals("jwt").(*jwt.Token)
	claims := jwtCtx.Claims.(jwt.MapClaims)

	return core.TokenMetadata{
		UserId:  int(claims["user_id"].(float64)),
		Role:    claims["role"].(string),
		Expires: int64(claims["exp"].(float64)),
	}
}

func JwtError(c *fiber.Ctx, err error) error {
	// Return status 401 and failed authentication error.
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Return status 401 and failed authentication error.
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": err.Error(),
	})
}
