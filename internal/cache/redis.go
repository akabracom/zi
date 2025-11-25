package cache

import (
    "context"
    "fmt"
    "os"
    "time"

    "github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// InitRedis инициализирует подключение к Redis
func InitRedis() error {
    host := os.Getenv("REDIS_HOST")
    if host == "" {
        host = "localhost"
    }
    port := os.Getenv("REDIS_PORT")
    if port == "" {
        port = "6280"
    }

    RedisClient = redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%s", host, port),
        Password: "", // без пароля
        DB:       0,  // БД по умолчанию
    })

    ctx := context.Background()
    _, err := RedisClient.Ping(ctx).Result()
    if err != nil {
        return fmt.Errorf("failed to connect to redis: %v", err)
    }

    fmt.Println("Successfully connected to Redis")
    return nil
}

// SaveToken сохраняет JWT токен в Redis с TTL 24 часа
func SaveToken(email string, token string) error {
    ctx := context.Background()
    key := fmt.Sprintf("token:%s", email)
    return RedisClient.Set(ctx, key, token, 24*time.Hour).Err()
}

// GetToken получает JWT токен из Redis
func GetToken(email string) (string, error) {
    ctx := context.Background()
    key := fmt.Sprintf("token:%s", email)
    return RedisClient.Get(ctx, key).Result()
}

// SaveUserSession сохраняет сессию пользователя
func SaveUserSession(userID int64, email string, role string) error {
    ctx := context.Background()
    key := fmt.Sprintf("session:user:%d", userID)
    data := map[string]interface{}{
        "email":    email,
        "role":     role,
        "login_at": time.Now().Format(time.RFC3339),
    }
    return RedisClient.HSet(ctx, key, data).Err()
}

// GetAllSessions получает все активные сессии (для демонстрации)
func GetAllSessions() ([]string, error) {
    ctx := context.Background()
    keys, err := RedisClient.Keys(ctx, "session:user:*").Result()
    if err != nil {
        return nil, err
    }
    return keys, nil
}
