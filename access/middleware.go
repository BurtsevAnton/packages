package access

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PolicyMW struct {
	service *Service
}

// NewPolicyMW возвращает новый экземпляр мидлваря политики доступа.
func NewPolicyMW(service *Service) *PolicyMW {
	return &PolicyMW{service: service}
}

// CheckApiAccess проверяет доступ пользователя к API.
func (mw *PolicyMW) CheckApiAccess(ctx *gin.Context, roleId int) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := ctx.Request.URL.Path
		method := ctx.Request.Method
		resourceCode := fmt.Sprintf("%s_%s", method, path)

		isAccessAllowed := mw.service.IsAccessAllowed(roleId, resourceCode)
		if !isAccessAllowed {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Message": "Доступ запрещен!"})

			return
		}

		ctx.Next()
	}
}
