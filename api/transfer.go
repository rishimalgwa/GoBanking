package api

import (
	"database/sql"
	"fmt"
	db "go-banking/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountId int64  `json:"from_account_id" binding:"required"`
	ToAccountId   int64  `json:"to_account_id" binding:"required"`
	Currency      string `json:"currency" binding:"required,oneof=USD EUR INR"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
}

func (server *Server) createTranfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if !server.validateAccount(ctx, req.FromAccountId, req.Currency) {
		return

	}
	if !server.validateAccount(ctx, req.ToAccountId, req.Currency) {
		return

	}
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountId,
		ToAccountID:   req.ToAccountId,
		Amount:        req.Amount,
	}
	resuts, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, resuts)
}

func (server *Server) validateAccount(ctx *gin.Context, accountId int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {

			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}
	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismach: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}
