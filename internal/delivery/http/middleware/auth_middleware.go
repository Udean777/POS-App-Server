package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sajudin/pos-app-server/pkg/utils"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "butuh login"})
			ctx.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		claims, err := utils.ValidateToken(tokenString, secret)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token tidak valid atau expired"})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("business_id", claims.BusinessID)
		ctx.Set("role", claims.Role) // Simpan role ke context

		ctx.Next()
	}
}

// RoleMiddleware membatasi akses endpoint berdasarkan role tertentu
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role := ctx.GetString("role")
		if role != requiredRole {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "akses ditolak, butuh role " + requiredRole})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
