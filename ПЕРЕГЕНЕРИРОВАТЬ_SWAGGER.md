# 🔄 ПЕРЕГЕНЕРАЦИЯ SWAGGER ДОКУМЕНТАЦИИ

## ❌ ПРОБЛЕМА:

В Swagger нет кнопки авторизации и не все endpoints видны.

**Причина:** Swagger документация не перегенерирована после добавления `@securityDefinitions`.

---

## ✅ РЕШЕНИЕ: Перегенерировать Swagger

### Вариант 1: Через swag (РЕКОМЕНДУЕТСЯ)

**Установите swag (если еще не установлен):**

```cmd
go install github.com/swaggo/swag/cmd/swag@latest
```

**Перегенерируйте документацию:**

```cmd
cd C:\spa-frontend\SPA_Dep
swag init -g cmd/server/main.go
```

**Пересоберите контейнер:**

```cmd
docker-compose build --no-cache deposits-app
docker-compose up -d deposits-app
```

---

### Вариант 2: Добавить securityDefinitions вручную

Если swag не установлен, можно добавить вручную в `docs/swagger.json`:

1. Откройте файл `SPA_Dep/docs/swagger.json`
2. Добавьте после `"basePath": "/api",` (после строки 10):

```json
"securityDefinitions": {
    "BearerAuth": {
        "type": "apiKey",
        "name": "Authorization",
        "in": "header",
        "description": "Type \"Bearer\" followed by a space and JWT token."
    }
},
```

3. Пересоберите контейнер:

```cmd
docker-compose build --no-cache deposits-app
docker-compose up -d deposits-app
```

---

### Вариант 3: Исправить в коде и пересобрать

Комментарии уже есть в `cmd/server/main.go`, нужно только перегенерировать.

---

## 🔍 ПРОВЕРКА:

После пересборки:

1. Откройте: http://localhost:3001/swagger/index.html
2. Должна появиться кнопка **"Authorize"** вверху справа
3. Все endpoints должны быть видны

---

## ⚠️ ВАЖНО:

Если используете Docker, нужно пересобрать контейнер после изменения swagger.json!


