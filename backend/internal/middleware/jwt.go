package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Authorizationヘッダーを取得
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authorization header is required",
				})
			}

			// "Bearer "プレフィックスをチェック
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid authorization header format",
				})
			}

			// JWTトークンを抽出
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Token is required",
				})
			}

			// JWTトークンを検証
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// 署名方法を確認
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				// シークレットキーを取得
				secret := os.Getenv("JWT_SECRET")
				if secret == "" {
					return nil, jwt.ErrTokenMalformed
				}
				return []byte(secret), nil
			})

			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token: " + err.Error(),
				})
			}

			// トークンが有効かチェック
			if !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Token is not valid",
				})
			}

			// クレームからユーザーIDを取得
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token claims",
				})
			}

			userID, ok := claims["sub"].(float64)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid user ID in token",
				})
			}

			// コンテキストにユーザーIDを保存
			c.Set("user_id", uint(userID))

			// 次のハンドラーを実行
			return next(c)
		}
	}
}

func GetUserID(c echo.Context) uint {
	return c.Get("user_id").(uint)
}
