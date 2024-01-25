package v1

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrBadReqStr             = "bad request"
	ErrBadRequest            = fmt.Errorf("bad request")
	ErrInternalStr           = "intrenal error"
	ErrInternal              = fmt.Errorf("internal error")
	ErrFromWalletNotFoundStr = "outgoing wallet not found"
	ErrFromWalletNotFound    = fmt.Errorf("outgoing wallet not found")
	ErrToWalletNotFoundStr   = "recipient's wallet not found"
	ErrToWalletNotFound      = fmt.Errorf("recipient's wallet not found")
	ErrWalletNotFoundStr     = "wallet not found"
	ErrWalletNotFound        = fmt.Errorf("wallet not found")
)

type errorResponse struct {
	StatusCode   int
	ErrorMessage string
}

func newErrorResponse(c *fiber.Ctx, errStatus int, message string) error {
	return c.Status(errStatus).JSON(errorResponse{
		StatusCode:   errStatus,
		ErrorMessage: message,
	})
}
