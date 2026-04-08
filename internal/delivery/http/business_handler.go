package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sajudin/pos-app-server/internal/domain"
)

type BusinessHandler struct {
	BusinessUsecase domain.BusinessUsecase
}

func NewBusinessHandler(bu domain.BusinessUsecase) *BusinessHandler {
	return &BusinessHandler{BusinessUsecase: bu}
}

func (h *BusinessHandler) UpdateBusiness(c *gin.Context) {
	var req domain.UpdateBusinessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data tidak valid"})
		return
	}

	bizIDStr := c.GetString("business_id")
	businessID, _ := uuid.Parse(bizIDStr)

	err := h.BusinessUsecase.UpdateBusiness(c.Request.Context(), businessID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "konfigurasi toko berhasil diperbarui"})
}
