package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sajudin/pos-app-server/internal/domain"
)

type AuthHandler struct {
	AuthUsecase domain.AuthUsecase
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data tidak valid"})
		return
	}

	resp, err := h.AuthUsecase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token wajib diisi"})
		return
	}

	resp, err := h.AuthUsecase.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email        string `json:"email" binding:"required,email"`
		Password     string `json:"password" binding:"required,min=6"`
		BusinessName string `json:"business_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data tidak valid atau password kurang panjang"})
		return
	}

	err := h.AuthUsecase.Register(c.Request.Context(), req.Email, req.Password, req.BusinessName)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{"error": "Email sudah terdaftar"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal melakukan registrasi"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "registrasi berhasil"})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIDStr)

	user, err := h.AuthUsecase.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profil tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, user)
}
