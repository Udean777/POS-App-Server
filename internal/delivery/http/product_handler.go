package http

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sajudin/pos-app-server/internal/domain"
)

type ProductHandler struct {
	ProductUsecase domain.ProductUsecase
	StorageService domain.StorageService
}

func NewProductHandler(pu domain.ProductUsecase, ss domain.StorageService) *ProductHandler {
	return &ProductHandler{ProductUsecase: pu, StorageService: ss}
}

func (h *ProductHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File tidak ditemukan"})
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuka file"})
		return
	}
	defer openedFile.Close()

	fileSize := file.Size
	fileBuffer := make([]byte, fileSize)
	_, err = openedFile.Read(fileBuffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca file"})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	contentType := file.Header.Get("Content-Type")

	url, err := h.StorageService.Upload(c.Request.Context(), fileBuffer, fileName, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
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

func (h *ProductHandler) GetByID(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	businessID, _ := uuid.Parse(c.GetString("business_id"))

	product, err := h.ProductUsecase.GetProductByID(c.Request.Context(), id, businessID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	businessID, _ := uuid.Parse(c.GetString("business_id"))

	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	product.ID = id
	product.BusinessID = businessID
	for i := range product.Variants {
		product.Variants[i].ProductID = id
		product.Variants[i].BusinessID = businessID
	}

	if err := h.ProductUsecase.UpdateProduct(c.Request.Context(), &product); err != nil {
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "Gagal merubah data. Varian produk yang dihapus sedang digunakan dalam riwayat transaksi."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	businessID, _ := uuid.Parse(c.GetString("business_id"))

	if err := h.ProductUsecase.DeleteProduct(c.Request.Context(), id, businessID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus produk"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil dihapus"})
}

func (h *ProductHandler) Restock(c *gin.Context) {
	variantIDStr := c.Param("variantId")
	variantID, err := uuid.Parse(variantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "variantId tidak valid"})
		return
	}

	businessIDStr := c.GetString("business_id")
	businessID, err := uuid.Parse(businessIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id tidak valid"})
		return
	}

	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	if err := h.ProductUsecase.RestockVariant(c.Request.Context(), variantID, businessID, req.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Restock berhasil"})
}
