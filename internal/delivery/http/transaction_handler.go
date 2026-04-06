package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sajudin/pos-app-server/internal/domain"
)

type TransactionHandler struct {
	txUsecase domain.TransactionUsecase
}

func NewTransactionHandler(txUsecase domain.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{txUsecase: txUsecase}
}

func (h *TransactionHandler) Checkout(c *gin.Context) {
	var req domain.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data tidak valid", "details": err.Error()})
		return
	}

	businessIDStr := c.GetString("business_id")
	businessID, _ := uuid.Parse(businessIDStr)

	staffIDStr := c.GetString("user_id")
	staffID, _ := uuid.Parse(staffIDStr)

	tx, err := h.txUsecase.ProcessCheckout(c.Request.Context(), req, businessID, staffID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Transaksi berhasil",
		"data":    tx,
	})
}

func (h *TransactionHandler) GetAll(c *gin.Context) {
	businessIDStr := c.GetString("business_id")
	businessID, _ := uuid.Parse(businessIDStr)

	txs, err := h.txUsecase.GetTransactions(c.Request.Context(), businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil riwayat transaksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": txs})
}
