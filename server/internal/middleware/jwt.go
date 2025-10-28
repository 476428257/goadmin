package middleware

import (
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"

	"server/config"
	"server/internal/model"
	"server/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Version  uint   `json:"version"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, username string, version uint) (string, error) {
	cfg := config.GetConfig()
	expireTime := time.Now().Add(time.Hour * time.Duration(cfg.JWT.ExpireHours))

	claims := Claims{
		UserID:   userID,
		Username: username,
		Version:  version,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// 验证token中的version是否与数据库中的一致
func validateTokenVersion(userID uint, tokenVersion uint) error {
	var admin model.Admin
	if err := database.DB.First(&admin, userID).Error; err != nil {
		return err
	}
	cfg := config.GetConfig()
	if admin.Version != tokenVersion && !cfg.JWT.Multi_login {
		return errors.New("token已刷新,帐号已在其他地方登录,请重新登录")
	}
	return nil
}

func Error(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 400,
		"msg":  msg,
	})
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		NoJWT := []string{"/api/v1/login", "/api/v1/auth/rule/getmenu", "/api/v1/config/getconfig"}
		if slices.Contains(NoJWT, c.Request.URL.Path) {
			c.Next()
			return
		}
		// 上传图片文件夹
		NoJWTSTRING := []string{"/uploads/"}
		for _, v := range NoJWTSTRING {
			if strings.Contains(c.Request.URL.Path, v) {
				c.Next()
				return
			}
		}
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			Error(c, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			Error(c, "Authorization header format must be Bearer {token}")
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetConfig().JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			Error(c, "token验证失败,请重新登录")
			c.Abort()
			return
		}
		// 验证token版本
		if err := validateTokenVersion(claims.UserID, claims.Version); err != nil {
			Error(c, err.Error())
			c.Abort()
			return
		}
		// logger.Logger.Info("claims", zap.Any("claims", claims))
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
