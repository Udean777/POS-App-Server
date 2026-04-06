package http

import (
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sajudin/pos-app-server/internal/domain"
)

type ProductHandler struct {
	ProductUsecase domain.ProductUsecase
}

func NewProductHandler(pu domain.ProductUsecase) *ProductHandler {
	return &ProductHandler{ProductUsecase: pu}
}

func (h *ProductHandler) Create(c *gin.Context) {
	businessIDStr := c.GetString("business_id")
	businessID, err := uuid.Parse(businessIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id tidak valid"})
		return
	}

	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data tidak valid"})
		return
	}

	product.BusinessID = businessID

	// Set BusinessID for each variant to ensure SKU uniqueness per business
	for i := range product.Variants {
		product.Variants[i].BusinessID = businessID
	}

	if err := h.ProductUsecase.AddProduct(c.Request.Context(), &product); err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "idx_business_sku") {
			c.JSON(http.StatusConflict, gin.H{"error": "SKU sudah digunakan di bisnis ini"})
			return
		}
		// Validation errors from usecase - we can assume these are safe for client
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	businessIDStr := c.GetString("business_id")
	businessID, err := uuid.Parse(businessIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id tidak valid"})
		return
	}

	products, err := h.ProductUsecase.GetAllProducts(c.Request.Context(), businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data produk"})
		return
	}

	c.JSON(http.StatusOK, products)
}
