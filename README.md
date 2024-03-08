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

## Инициализация мидлваря

```
    // Создаем мидлварь доступа
	accessMW := access.NewPolicyMW(serviceMW)
```


## Использование мидлваря

```
...
import "github.com/BurtsevAnton/packages/access"
...
    // Инициализируем сервис политики доступа
	service := NewService()
	
	// Получаем политики доступа из базы данных
	policyFromDb, _ := getPolicyFromYourDbAsJson()
	
	// Устанавливаем политики доступа
	err := service.SetAccessPolicies(policyFromDb)
	if err != nil {
		t.Error("Ошибка установки политики доступа")
	}

	// Создаем мидлварь проверки доступа
	mw := NewPolicyMW(service)
	
	r := gin.Default()
	
	group := r.Group("/api")
	{
		// Применяем мидлварь для группы роутов
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


```

