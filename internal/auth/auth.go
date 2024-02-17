package auth

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/zaher1307/IDEANEST-project-assignment/internal/types"
)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var (
	accessSecret  string
	refreshSecret string
	redisHost     string
	redisClient   *redis.Client
)

func init() {
	accessSecret = os.Getenv("ACCESS_SECRET")
	refreshSecret = os.Getenv("REFRESH_SECRET")
	redisHost = os.Getenv("REDIS_HOST")
	redisClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":6379",
	})
}

func GenerateAccessToken(refreshToken string) (string, error) {
	return generateToken(refreshToken, accessSecret, time.Minute*15)
}

func GenerateRefreshToken(user types.User) (string, error) {
	refreshToken := uuid.New().String()
	err := redisClient.Set(refreshToken, user.Email, 0).Err()
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

func GetRefreshTokenUserEmail(refreshToken string) (string, error) {
	return redisClient.Get(refreshToken).Result()
}

func ValidateAccessToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessSecret), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token claims")
	}

	return claims.Email, nil
}

func RevokeRefreshToken(token string) error {
	err := redisClient.Del(token).Err()
	if err != nil {
		return err
	}

	return nil
}

// ======================== helper util function ======================== //

func generateToken(refreshToken string, secretKey string, expiration time.Duration) (string, error) {
	email, err := redisClient.Get(refreshToken).Result()
	if err != nil {
		return "", err
	}

	claims := Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
