package access

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_NewPolicyMW(t *testing.T) {
	service := &Service{}

	policyMW := NewPolicyMW(service)

	assert.NotNil(t, policyMW)
	assert.Equal(t, service, policyMW.service)
}

func Test_CheckApiAccess(t *testing.T) {
	service := NewService()
	err := service.SetAccessPolicies(policyFromDb)
	if err != nil {
		t.Error("Ошибка установки политики доступа")
	}

	mw := NewPolicyMW(service)

	type args struct {
		roleId       int
		route        string
		method       string
		resourceCode string
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
	}{
		{
			name: "Test_CheckApiAccess_1_POST_/api/test_create",
			args: args{
				roleId:       1, // Developer
				route:        "/api/test_create",
				method:       "POST",
				resourceCode: "POST_/api/test_create",
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Test_CheckApiAccess_1_POST_/api/test_get",
			args: args{
				roleId:       1,
				route:        "/api/test_get",
				method:       "GET",
				resourceCode: "GET_/api/test_get",
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Test_CheckApiAccess_1_DELETE_/api/test_delete",
			args: args{
				roleId:       1,
				route:        "/api/test_delete",
				method:       "DELETE",
				resourceCode: "DELETE_/api/test_delete",
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Test_CheckApiAccess_1_RUN_/api/test_run",
			args: args{
				roleId:       1,
				route:        "/api/test_run",
				method:       "POST",
				resourceCode: "POST_/api/test_run",
			},
			wantCode: http.StatusNotFound,
		},
		{
			name: "Test_CheckApiAccess_4_POST_/api/test_create",
			args: args{
				roleId:       4, // Worker
				route:        "/api/test_create",
				method:       "POST",
				resourceCode: "POST_/api/test_create",
			},
			wantCode: http.StatusForbidden,
		},
		{
			name: "Test_CheckApiAccess_4_POST_/api/test_get",
			args: args{
				roleId:       4,
				route:        "/api/test_get",
				method:       "GET",
				resourceCode: "GET_/api/test_get",
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Test_CheckApiAccess_4_DELETE_/api/test_delete",
			args: args{
				roleId:       4,
				route:        "/api/test_delete",
				method:       "DELETE",
				resourceCode: "DELETE_/api/test_delete",
			},
			wantCode: http.StatusForbidden,
		},
		{
			name: "Test_CheckApiAccess_4_RUN_/api/test_run",
			args: args{
				roleId:       4,
				route:        "/api/test_run",
				method:       "POST",
				resourceCode: "POST_/api/test_run",
			},
			wantCode: http.StatusNotFound,
		},
		{
			name: "Test_CheckApiAccess_3_POST_/api/test_create",
			args: args{
				roleId:       3, // Lead
				route:        "/api/test_create",
				method:       "POST",
				resourceCode: "POST_/api/test_create",
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Test_CheckApiAccess_3_POST_/api/test_get",
			args: args{
				roleId:       3,
				route:        "/api/test_get",
				method:       "GET",
				resourceCode: "GET_/api/test_get",
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Test_CheckApiAccess_3_DELETE_/api/test_delete",
			args: args{
				roleId:       3,
				route:        "/api/test_delete",
				method:       "DELETE",
				resourceCode: "DELETE_/api/test_delete",
			},
			wantCode: http.StatusForbidden,
		},
		{
			name: "Test_CheckApiAccess_3_RUN_/api/test_run",
			args: args{
				roleId:       3,
				route:        "/api/test_run",
				method:       "POST",
				resourceCode: "POST_/api/test_run",
			},
			wantCode: http.StatusNotFound,
		},
		{
			name: "Test_CheckApiAccess_0_POST_/api/test_create",
			args: args{
				roleId:       0, // Default role without access
				route:        "/api/test_create",
				method:       "POST",
				resourceCode: "POST_/api/test_create",
			},
			wantCode: http.StatusForbidden,
		},
		{
			name: "Test_CheckApiAccess_0_POST_/api/test_get",
			args: args{
				roleId:       0,
				route:        "/api/test_get",
				method:       "GET",
				resourceCode: "GET_/api/test_get",
			},
			wantCode: http.StatusForbidden,
		},
		{
			name: "Test_CheckApiAccess_0_DELETE_/api/test_delete",
			args: args{
				roleId:       0,
				route:        "/api/test_delete",
				method:       "DELETE",
				resourceCode: "DELETE_/api/test_delete",
			},
			wantCode: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.ReleaseMode)
			r := gin.Default()
			group := r.Group("/api")
			{
				group.Use(func(c *gin.Context) {
					authHandler := mw.CheckApiAccess(c, tt.args.roleId)
					fmt.Printf("%T", authHandler)
					authHandler(c)
					c.Next()
				})

				group.POST("/test_create", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"Message": "Доступ на чтение разрешен"})
				})
				group.GET("/test_get", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"Message": "Доступ на изменение разрешен"})
				})
				group.DELETE("/test_delete", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"Message": "Доступ на удаление разрешен"})
				})
			}

			// Делаем GET запрос на /test
			req, _ := http.NewRequest(tt.args.method, tt.args.route, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("Ожидается статус %d для пользователя %d и роута %s, но получен статус %d", tt.wantCode, tt.args.roleId, tt.args.route, w.Code)
			}
		})
	}
}

func Test_CheckApiAccess_WithoutInternalMW(t *testing.T) {
	service := NewService()
	err := service.SetAccessPolicies(policyFromDb)
	if err != nil {
		t.Error("Ошибка установки политики доступа")
	}

	type args struct {
		roleId int
		route  string
		method string
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
	}{
		{
			name: "Test_CheckApiAccess_1_POST_/api/test_create",
			args: args{
				roleId: 1, // Developer
				route:  "/api/test_create",
				method: "POST",
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Test_CheckApiAccess_1_POST_/api/test_get",
			args: args{
				roleId: 1,
				route:  "/api/test_get",
				method: "GET",
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Test_CheckApiAccess_1_DELETE_/api/test_delete",
			args: args{
				roleId: 1,
				route:  "/api/test_delete",
				method: "DELETE",
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Test_CheckApiAccess_1_RUN_/api/test_run",
			args: args{
				roleId: 1,
				route:  "/api/test_run",
				method: "POST",
			},
			wantCode: http.StatusNotFound,
		},
		{
			name: "Test_CheckApiAccess_4_POST_/api/test_create",
			args: args{
				roleId: 4, // Worker
				route:  "/api/test_create",
				method: "POST",
			},
			wantCode: http.StatusForbidden,
		},
		{
			name: "Test_CheckApiAccess_4_POST_/api/test_get",
			args: args{
				roleId: 4,
				route:  "/api/test_get",
				method: "GET",
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Test_CheckApiAccess_4_DELETE_/api/test_delete",
			args: args{
				roleId: 4,
				route:  "/api/test_delete",
				method: "DELETE",
			},
			wantCode: http.StatusForbidden,
		},
		{
			name: "Test_CheckApiAccess_4_RUN_/api/test_run",
			args: args{
				roleId: 4,
				route:  "/api/test_run",
				method: "POST",
			},
			wantCode: http.StatusNotFound,
		},
		{
			name: "Test_CheckApiAccess_3_POST_/api/test_create",
			args: args{
				roleId: 3, // Lead
				route:  "/api/test_create",
				method: "POST",
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Test_CheckApiAccess_3_POST_/api/test_get",
			args: args{
				roleId: 3,
				route:  "/api/test_get",
				method: "GET",
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Test_CheckApiAccess_3_DELETE_/api/test_delete",
			args: args{
				roleId: 3,
				route:  "/api/test_delete",
				method: "DELETE",
			},
			wantCode: http.StatusForbidden,
		},
		{
			name: "Test_CheckApiAccess_3_RUN_/api/test_run",
			args: args{
				roleId: 3,
				route:  "/api/test_run",
				method: "POST",
			},
			wantCode: http.StatusNotFound,
		},
		{
			name: "Test_CheckApiAccess_0_POST_/api/test_create",
			args: args{
				roleId: 0, // Default role without access
				route:  "/api/test_create",
				method: "POST",
			},
			wantCode: http.StatusForbidden,
		},
		{
			name: "Test_CheckApiAccess_0_POST_/api/test_get",
			args: args{
				roleId: 0,
				route:  "/api/test_get",
				method: "GET",
			},
			wantCode: http.StatusForbidden,
		},
		{
			name: "Test_CheckApiAccess_0_DELETE_/api/test_delete",
			args: args{
				roleId: 0,
				route:  "/api/test_delete",
				method: "DELETE",
			},
			wantCode: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.ReleaseMode)
			r := gin.Default()
			group := r.Group("/api")
			{
				// Middleware для проверки прав доступа
				group.Use(func(c *gin.Context) {
					path := c.Request.URL.Path
					method := c.Request.Method
					resourceCode := fmt.Sprintf("%s_%s", method, path)

					isAccessAllowed := service.IsAccessAllowed(tt.args.roleId, resourceCode)
					if !isAccessAllowed {
						c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Message": "Доступ запрещен!"})

						return
					}
					c.Next()
				})

				group.POST("/test_create", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"Message": "Доступ на чтение разрешен"})
				})
				group.GET("/test_get", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"Message": "Доступ на изменение разрешен"})
				})
				group.DELETE("/test_delete", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"Message": "Доступ на удаление разрешен"})
				})
			}

			// Делаем GET запрос на /test
			req, _ := http.NewRequest(tt.args.method, tt.args.route, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("Ожидается статус %d для пользователя %d и роута %s, но получен статус %d", tt.wantCode, tt.args.roleId, tt.args.route, w.Code)
			}
		})
	}
}
