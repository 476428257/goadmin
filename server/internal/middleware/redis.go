package middleware

import (
	"context"

	"server/internal/model"
	"server/pkg/database"
	"server/pkg/redisx"

	"github.com/gin-gonic/gin"
)

func RedisMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		// 1) 先从 Redis 读取
		cfgMap, err := redisx.HGetAll(ctx, redisx.ConfigHashKey)
		if err != nil {
			// 读 Redis 出错不直接中断，继续尝试 DB 回源
			cfgMap = map[string]string{}
		}

		// 2) 如果 Redis 为空，从数据库加载并写回 Redis
		if len(cfgMap) == 0 {
			var rows []model.Config
			if err := database.DB.Model(&model.Config{}).
				Select("name", "value").
				Find(&rows).Error; err == nil {
				m := make(map[string]string, len(rows))
				for _, r := range rows {
					m[r.Name] = r.Value
				}
				// 忽略写 Redis 的错误，不中断请求
				_ = redisx.HSet(ctx, redisx.ConfigHashKey, m)
				cfgMap = m
			}
		}

		// 3) 注入到 Context，供后续处理/页面使用
		c.Set("config", cfgMap)
		c.Next()
	}
}
