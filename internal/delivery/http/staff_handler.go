package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sajudin/pos-app-server/internal/domain"
)

type StaffHandler struct {
	StaffUsecase domain.StaffUsecase
}

func NewStaffHandler(su domain.StaffUsecase) *StaffHandler {
	return &StaffHandler{StaffUsecase: su}
}

func (h *StaffHandler) CreateStaff(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data tidak valid"})
		return
	}

	bizIDStr := c.GetString("business_id")
	businessID, _ := uuid.Parse(bizIDStr)

	err := h.StaffUsecase.CreateStaff(c.Request.Context(), req.Email, req.Password, req.Role, businessID)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) || errors.Is(err, domain.ErrEmailRegisteredByOtherBusiness) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan staf"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "staf berhasil didaftarkan"})
}

func (h *StaffHandler) GetStaff(c *gin.Context) {
	bizIDStr := c.GetString("business_id")
	businessID, _ := uuid.Parse(bizIDStr)

	staff, err := h.StaffUsecase.GetStaff(c.Request.Context(), businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data staf"})
		return
	}

	c.JSON(http.StatusOK, staff)
}
