# packages

## Описание

В этом пакете будут находиться все пакеты, которые могут быть использоваться в разных проектах.

## Состав

- [ ] Пакет access

***

# Пакет access
Предназначен для реализации доступа пользователей к ресурсам.
Представляет собой middleware или набор методов сервиса для реализации доступа к api.

## Установка

Загружаем пакет

```
go get github.com/BurtsevAnton/packages
```

## Инициализация

```
    // Инициализируем сервис доступа
    serviceMW := access.NewService()

    // Получаем политики доступа из базы данных
    policyFromDb, _ := getPolicyFromYourDB()
    
    // Устанавливаем политики доступа
    if err := serviceMW.SetAccessPolicies(policyFromDb); err != nil {
        return nil, err
    }
```

## Использование
```
    // Роль пользователя нужно получить из токена 
    roleId := 1
    
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

```

## Инициализация встроенного мидлваря

```
    // Создаем мидлварь доступа
    accessMW := access.NewPolicyMW(serviceMW)
```


## Использование встроенного мидлваря

```
...
import "github.com/BurtsevAnton/packages/access"
...
    // Ролль пользователя нужно получить из токена 
    roleId := 1
    
    // Инициализируем сервис политики доступа
    serviceMW := NewService()
	
    // Получаем политики доступа из базы данных
    policyFromDb, _ := getPolicyFromYourDbAsJson()
	
    // Устанавливаем политики доступа
    err := serviceMW.SetAccessPolicies(policyFromDb)
    if err != nil {
        t.Error("Ошибка установки политики доступа")
    }

    // Создаем мидлварь проверки доступа
    mw := NewPolicyMW(serviceMW)
	
    r := gin.Default()
	
    group := r.Group("/api")
    {
        // Применяем внутренний мидлварь для группы роутов
        group.Use(func(c *gin.Context) {
            authHandler := mw.CheckApiAccess(c, roleId)
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


```

