package middleware

import (
    "net/http"
    "strings"

    "backend/internal/auth"

    "github.com/gin-gonic/gin"
)

// AuthMiddleware проверяет JWT токен (опционально - для гостей)
func AuthMiddleware(optional bool) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        
        // Если токен не передан
        if authHeader == "" {
            if optional {
                // Гость - продолжаем с ролью "guest"
                c.Set("role", "guest")
                c.Next()
                return
            }
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
            c.Abort()
            return
        }

        // Формат: "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
            c.Abort()
            return
        }

        tokenString := parts[1]
        claims, err := auth.ValidateToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
            c.Abort()
            return
        }

        // Сохраняем данные пользователя в контексте
        c.Set("user_id", claims.UserID)
        c.Set("email", claims.Email)
        c.Set("role", claims.Role)
        c.Next()
    }
}

// RequireRole проверяет наличие определенной роли
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("role")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{"error": "role not found"})
            c.Abort()
            return
        }

        userRole := role.(string)
        for _, allowedRole := range allowedRoles {
            if userRole == allowedRole {
                c.Next()
                return
            }
        }

        c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
        c.Abort()
    }
}
