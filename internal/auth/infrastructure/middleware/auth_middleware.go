package middleware

import (
	"net/http"
	"strings"

	auth "github.com/internalWizzard/idp-server/internal/auth/port"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	Provider auth.AuthProvider
}

func NewAuthMiddleware(provider auth.AuthProvider) *AuthMiddleware {
	return &AuthMiddleware{
		Provider: provider,
	}
}

// Handle es el middleware real que se usa en Echo.
func (m *AuthMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		// 1. Obtener el token del header Authorization
		authHeader := c.Request().Header.Get("Authorization")

		if authHeader == "" {
			// Opcional: buscar cookie si usás sesiones
			cookie, err := c.Cookie("session")
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authentication token")
			}
			authHeader = "cookie " + cookie.Value
		}

		// 2. Normalizar: si es "Bearer <token>"
		var tokenOrCookie string
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			tokenOrCookie = strings.TrimSpace(authHeader[7:])
		} else {
			tokenOrCookie = strings.TrimSpace(authHeader)
		}

		if tokenOrCookie == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid authentication token")
		}

		// 3. Validar usando el proveedor inyectado
		if err := m.Provider.Validate(tokenOrCookie); err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		// 4. Todo bien → continuar
		return next(c)
	}
}
