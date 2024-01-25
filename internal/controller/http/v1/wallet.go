package v1

import (
	"errors"

	"github.com/1boombacks1/testTaskInfotecs/internal/service"
	"github.com/gofiber/fiber/v2"
)

type routes struct {
	srvc service.Service
}

func newWalletRoutes(g fiber.Router, srvc service.Service) {
	r := routes{srvc}

	g.Post("/wallet/:walletID/send", r.send)
	g.Get("/wallet/:walletID/history", r.getHistoryOperations)
	g.Get("/wallet/:walletID", r.getWallet)
	g.Post("/wallet", r.createWallet)
}

func (r *routes) createWallet(c *fiber.Ctx) error {
	wallet, err := r.srvc.CreateWallet(c.Context())
	if err != nil {
		return newErrorResponse(c, 400, ErrBadReqStr)
	}

	type response struct {
		ID      string  `json:"id"`
		Balance float32 `json:"balance"`
	}

	return c.Status(fiber.StatusOK).JSON(response{
		ID:      wallet.ID.String(),
		Balance: wallet.Balance,
	})
}

type transferInput struct {
	To     string  `json:"to"`
	Amount float32 `json:"amount"`
}

func (r *routes) send(c *fiber.Ctx) error {
	from := c.Params("walletID")
	var input transferInput
	if err := c.BodyParser(&input); err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, ErrInternalStr)
	}

	if err := r.srvc.Transfer(c.Context(), from, input.To, input.Amount); err != nil {
		if errors.Is(err, service.ErrInvalidID) {
			return newErrorResponse(c, fiber.StatusBadRequest, ErrBadReqStr)

		} else if errors.Is(err, service.ErrFromWalletNotFound) {
			return newErrorResponse(c, fiber.StatusNotFound, ErrFromWalletNotFoundStr)
		} else if errors.Is(err, service.ErrToWalletNotFound) {
			return newErrorResponse(c, fiber.StatusNotFound, ErrToWalletNotFoundStr)
		}
		return newErrorResponse(c, fiber.StatusInternalServerError, ErrInternalStr)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": "The transfer was successfully completed",
	})
}

func (r *routes) getHistoryOperations(c *fiber.Ctx) error {
	walletId := c.Params("walletID")
	operations, err := r.srvc.GetHistoryOperations(c.Context(), walletId)
	if err != nil {
		if errors.Is(err, service.ErrInvalidID) {
			return newErrorResponse(c, fiber.StatusBadRequest, ErrBadReqStr)
		} else if errors.Is(err, service.ErrWalletNotFound) {
			return newErrorResponse(c, fiber.StatusNotFound, ErrWalletNotFoundStr)
		}
	}
	// type response struct {
	// 	Time time.Time `json:"time"`
	// 	From string `json:"from"`
	// 	To string `json:"to"`
	// 	Amount float32 `json:"amount"`
	// }
	return c.Status(fiber.StatusOK).JSON(operations)
}

func (r *routes) getWallet(c *fiber.Ctx) error {
	walletID := c.Params("walletID")
	wallet, err := r.srvc.GetWalletByID(c.Context(), walletID)
	if err != nil {
		return newErrorResponse(c, fiber.StatusNotFound, ErrWalletNotFoundStr)
	}

	type response struct {
		Id      string  `json:"id"`
		Balance float32 `json:"balance"`
	}

	return c.Status(fiber.StatusOK).JSON(response{
		Id:      wallet.ID.String(),
		Balance: wallet.Balance,
	})
}
