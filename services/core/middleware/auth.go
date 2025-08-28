package middleware

import (
	"altron/common/dto"
	"altron/config"
	"altron/core/interfaces"
	"altron/pkg/auth"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthMiddleware struct {
	cfg           *config.Config
	authenticator auth.Authenticator
	jwtAuth       auth.JwtAuthenticator
	userRepo      interfaces.UserRepository
}

func NewAuthMiddleware(
	cfg *config.Config,
	authenticator auth.Authenticator,
	jwtAuth auth.JwtAuthenticator,
	userRepo interfaces.UserRepository,
) *AuthMiddleware {
	return &AuthMiddleware{
		cfg:           cfg,
		authenticator: authenticator,
		jwtAuth:       jwtAuth,
		userRepo:      userRepo,
	}
}

func (a *AuthMiddleware) HandlerFunc() func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		secure, err := a.authenticator.IsSecure(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    dto.ErrorUnauthorized,
				Message: err.Error(),
			})
			c.Abort()
			return
		}
		if !secure || a.isServiceRequest(c.Request) {
			user, err := a.userRepo.GetUser(ctx, a.cfg.App.AltronName)
			if err != nil {
				c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Code:    dto.ErrorInternalServerError,
					Message: err.Error(),
				})
				c.Abort()
				return
			}

			c.Request = c.Request.WithContext(context.WithValue(ctx, auth.UserContextKey, auth.AccessClaims{
				UserID: user.ID,
			}))
			c.Next()
			return
		}
		jwt, err := a.authenticator.GetToken(c.Request)
		if err != nil || jwt == nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    dto.ErrorUnauthorized,
				Message: "missing token",
			})
			c.Abort()
			return
		}

		var accessClaims *auth.AccessClaims
		cleanToken := auth.GetBearer(*jwt)
		accessClaims, err = a.jwtAuth.ParseAccessToken(cleanToken)
		if err != nil {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Code:    dto.ErrorForbidden,
				Message: "invalid security token",
			})
			c.Abort()
			return
		}

		user, err := a.userRepo.GetUserByID(ctx, accessClaims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    dto.ErrorUnauthorized,
				Message: err.Error(),
			})
			c.Abort()
			return
		}

		if isAdminHandler(c.Request) && user.Username != "manager" {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Code:    dto.ErrorForbidden,
				Message: "access denied",
			})
			c.Abort()
			return
		}

		c.Request = c.Request.WithContext(accessClaims.WithContext(ctx))
		c.Next()
	}
}

func GetUserIdFromContext(ctx context.Context) uuid.UUID {
	claims, ok := ctx.Value(auth.UserContextKey).(auth.AccessClaims)

	if !ok {
		return uuid.Nil
	}

	return claims.UserID
}

func (a *AuthMiddleware) isServiceRequest(r *http.Request) bool {
	return !strings.Contains(r.Header.Get("User-Agent"), "Mozilla")//r.Header.Get("User-Agent") != "client"
}

func isAdminHandler(r *http.Request) bool {
	return strings.Contains(r.URL.Path, "admin") || strings.Contains(r.URL.Path, "users/info")
}
