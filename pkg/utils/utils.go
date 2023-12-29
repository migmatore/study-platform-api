package utils

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgconn"
	"time"
)

func DoWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--

			continue
		}

		return nil
	}

	return
}

func ParsePgError(err error) error {
	var pgErr *pgconn.PgError

	if errors.Is(err, pgErr) {
		pgErr = err.(*pgconn.PgError)
		return fmt.Errorf("database error. message: %s, details: %s, where: %s, sqlstate: %s",
			pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.SQLState())
	}

	return err
}

func FiberError(ctx *fiber.Ctx, status int, err error) error {
	return ctx.Status(status).JSON(fiber.Map{
		"message": err.Error(),
	})
}
