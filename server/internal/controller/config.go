package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"strings"

	"server/config"
	"server/internal/model"
	"server/pkg/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ConfigController struct {
	BaseController[model.Config]
}

func NewConfigController() *ConfigController {
	return &ConfigController{
		BaseController: BaseController[model.Config]{
			Model: model.Config{},
		},
	}
}

func (b *ConfigController) List(ctx *gin.Context) {
	model := b.Model
	var list []map[string]interface{}
	// 先查询所有分组，去重，给前端做 tab 标题
	var groups []string
	if err := database.DB.Model(&model).Distinct("`group`").Pluck("`group`", &groups).Error; err != nil {
		b.Error(ctx, err.Error())
		return
	}
	// 列表查询与过滤
	db := database.DB.Model(&model)
	if err := db.Scan(&list).Error; err != nil {
		b.Error(ctx, err.Error())
		return
	}
	// 对富文本类型进行反转义，确保在编辑器能正确显示
	for _, row := range list {
		if t, ok := row["type"].(string); ok && t == "editor" {
			if v, vok := row["value"].(string); vok {
				row["value"] = html.UnescapeString(v)
			}
		}
	}
	b.Success(ctx, "获取列表成功", gin.H{
		"groups": groups,
		"list":   &list,
	})
}

// Getconfig 获取配置：优先从 Redis 读取；若某项值为空或不存在则从数据库读取并回填 Redis，最后返回
func (b *ConfigController) Getconfig(ctx *gin.Context) {
	// 1) 先尝试从 Redis 获取全部配置
	cached, err := RedisGetAllConfig(context.Background())
	if err != nil {
		b.Error(ctx, err.Error())
		return
	}
	if len(cached) == 0 {
		var list []model.Config
		if err := database.DB.Model(&model.Config{}).Find(&list).Error; err != nil {
			b.Error(ctx, err.Error())
			return
		}
		result := make(map[string]string, len(list))
		for _, c := range list {
			result[c.Name] = c.Value
			// 富文本在返回前做反转义
			if c.Type == "editor" {
				result[c.Name] = html.UnescapeString(c.Value)
			}
		}
		cached = result
		if err := RedisSetConfig(context.Background(), result); err != nil {
			b.Error(ctx, err.Error())
			return
		}
	}
	b.Success(ctx, "获取成功", cached)
}

func (b *ConfigController) Update(ctx *gin.Context) {
	type updateReq struct {
		Data map[string]any `json:"data"`
	}
	var req updateReq
	body, err := ctx.GetRawData()
	if err != nil {
		b.Error(ctx, err.Error())
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		b.Error(ctx, err.Error())
		return
	}
	if len(req.Data) == 0 {
		b.Error(ctx, "参数 data 不能为空")
		return
	}

	if err := b.withTransaction(ctx, func(tx *gorm.DB) error {
		// 构造单条 SQL 使用 CASE WHEN 实现批量更新
		names := make([]string, 0, len(req.Data))
		args := make([]any, 0, len(req.Data)*2)
		caseSql := "CASE name "
		for k, v := range req.Data {
			names = append(names, k)
			// 允许空值，nil 统一存为空串，数组拼逗号
			var val string
			if v == nil {
				val = ""
			} else if arr, ok := v.([]any); ok {
				parts := make([]string, 0, len(arr))
				for _, it := range arr {
					if it == nil {
						continue
					}
					parts = append(parts, fmt.Sprintf("%v", it))
				}
				val = strings.Join(parts, ",")
			} else {
				val = fmt.Sprintf("%v", v)
			}
			caseSql += "WHEN ? THEN ? "
			args = append(args, k, val)
		}
		caseSql += "END"

		inPlaceholders := ""
		for i := 0; i < len(names); i++ {
			if i > 0 {
				inPlaceholders += ","
			}
			inPlaceholders += "?"
		}
		// 完整 SQL
		sql := "UPDATE " + config.GetConfig().MySQL.Prefix + "config SET value = " + caseSql + " WHERE name IN (" + inPlaceholders + ")"
		// 追加 IN 子句参数
		inArgs := make([]any, len(names))
		for i := range names {
			inArgs[i] = names[i]
		}
		args = append(args, inArgs...)
		if err := tx.Exec(sql, args...).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return
	}

	// 同步写入 Redis 哈希
	kv := make(map[string]string, len(req.Data))
	for k, v := range req.Data {
		var val string
		if v == nil {
			val = ""
		} else if arr, ok := v.([]any); ok {
			parts := make([]string, 0, len(arr))
			for _, it := range arr {
				if it == nil {
					continue
				}
				parts = append(parts, fmt.Sprintf("%v", it))
			}
			val = strings.Join(parts, ",")
		} else {
			val = fmt.Sprintf("%v", v)
		}
		kv[k] = val
	}
	err = RedisSetConfig(context.Background(), kv)
	if err != nil {
		b.Error(ctx, err.Error())
		return
	}
	b.Success(ctx, "更新成功", nil)
}
