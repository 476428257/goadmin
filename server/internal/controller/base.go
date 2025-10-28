package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"slices"
	"strings"
	"time"

	"server/config"
	"server/internal/middleware"
	"server/internal/model"
	"server/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var CFG = config.GetConfig()

type Pagination struct {
	Page     uint `json:"page" binding:"required,numeric"`
	PageSize uint `json:"pagesize" binding:"required,numeric"`
}

// HandleValidationError 处理验证错误的公共方法
// 参数: ctx - gin上下文, obj - 要验证的结构体指针
// 返回: error - 如果有验证错误则返回错误，否则返回nil
func (b *BaseController[T]) HandleValidationError(ctx *gin.Context, obj interface{}) error {
	if err := ctx.ShouldBindJSON(obj); err != nil {
		// 使用validator中间件的中文错误提示
		if validationError := middleware.GetValidationError(err); validationError != "参数验证失败" {
			b.Error(ctx, validationError)
		} else {
			b.Error(ctx, err.Error())
		}
		return err
	}
	return nil
}

// HandleJSONUnmarshalValidation 处理json.Unmarshal后的验证
// 参数: ctx - gin上下文, jsonData - JSON字符串, obj - 要验证的结构体指针
// 返回: error - 如果有解析或验证错误则返回错误，否则返回nil
func (b *BaseController[T]) HandleJSONUnmarshalValidation(ctx *gin.Context, jsonData string, obj interface{}) error {
	// 先进行JSON解析
	if err := json.Unmarshal([]byte(jsonData), obj); err != nil {
		b.Error(ctx, "JSON格式错误: "+err.Error())
		return err
	}

	// 使用reflection处理binding标签转换为validate标签进行验证
	if err := b.validateStructWithBindingTags(obj); err != nil {
		// 使用validator中间件的中文错误提示
		if validationError := middleware.GetValidationError(err); validationError != "参数验证失败" {
			b.Error(ctx, validationError)
		} else {
			b.Error(ctx, err.Error())
		}
		return err
	}
	return nil
}

// DateRange 日期区间类型
type DateRange []string

// ParseDateRange 解析日期区间，返回开始时间和结束时间
func (dr DateRange) ParseDateRange() (time.Time, time.Time, error) {
	if len(dr) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("日期区间必须包含开始和结束两个时间")
	}

	// 获取本地时区，避免UTC时区导致的时差问题
	loc := time.Local

	// 解析开始时间
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", dr[0], loc)
	if err != nil {
		// 尝试按日期格式解析，然后设置为当天00:00:00
		if startTime, err = time.ParseInLocation("2006-01-02", dr[0], loc); err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("开始日期格式错误: %v", err)
		}
	}

	// 解析结束时间
	endTime, err := time.ParseInLocation("2006-01-02 15:04:05", dr[1], loc)
	if err != nil {
		// 尝试按日期格式解析，然后设置为当天23:59:59
		if endTime, err = time.ParseInLocation("2006-01-02", dr[1], loc); err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("结束日期格式错误: %v", err)
		} else {
			// 设置为当天的23:59:59
			endTime = endTime.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
	}

	return startTime, endTime, nil
}

// IsEmpty 检查日期区间是否为空
func (dr DateRange) IsEmpty() bool {
	return len(dr) == 0 || (len(dr) == 2 && (dr[0] == "" || dr[1] == ""))
}

// validateStructWithBindingTags 使用binding标签验证结构体
// 这个方法创建一个配置了字段名映射的validator实例来处理binding标签
func (b *BaseController[T]) validateStructWithBindingTags(obj interface{}) error {
	// 获取gin的默认validator引擎，它已经配置了字段名映射
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 使用已配置的validator实例进行验证
		// 这样可以保证字段名映射和自定义验证器都能正常工作
		return v.Struct(obj)
	}

	// 如果无法获取gin的validator引擎，创建新的实例并配置
	validate := validator.New()

	// 注册字段名获取函数，优先使用name标签，然后是json标签，最后是结构体字段名
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		// 首先尝试获取name标签
		nameTag := fld.Tag.Get("name")
		if nameTag != "" && nameTag != "-" {
			return nameTag
		}

		// 如果没有name标签，尝试获取json标签
		jsonTag := fld.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			// 提取json标签中的字段名（去掉omitempty等选项）
			name := strings.SplitN(jsonTag, ",", 2)[0]
			if name != "" {
				return name
			}
		}

		// 如果都没有，返回结构体字段名
		return fld.Name
	})

	return validate.Struct(obj)
}

// BaseController 基础控制器
type BaseController[T interface {
	model.IDGetter
}] struct {
	Model T
}

// 键值对
type K_V struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

func filterEmpty(s []string) []string {
	var out []string
	for _, v := range s {
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

func recordOperationLog(ctx *gin.Context, tx *gorm.DB) error {
	// 指定方法后缀不记录日志
	path := ctx.Request.URL.Path
	parts := filterEmpty(strings.Split(path, "/"))
	if len(parts) == 0 {
		return fmt.Errorf("路由错误")
	}
	// 获取请求数据
	requestData, exists := ctx.Get("requestData")
	var reqStr string
	if exists {
		// 尝试将requestData转换为JSON字符串
		if data, err := json.Marshal(requestData); err == nil {
			reqStr = string(data)
		}
	}
	// 不记录内容
	NoAOperationLogVague := []string{"login"}
	if slices.Contains(NoAOperationLogVague, parts[len(parts)-1]) {
		reqStr = ""
	}

	// 获取管理员信息
	adminID, exists := ctx.Get("userID")
	if !exists {
		return fmt.Errorf("用户信息不存在")
	}

	var admin model.Admin
	if err := tx.Where("id = ?", adminID).First(&admin).Error; err != nil {
		return fmt.Errorf("用户信息出错")
	}

	// 记录日志
	log := model.OperationLog{
		AdminID:     admin.ID,
		Username:    admin.Username,
		Path:        path,
		IP:          ctx.ClientIP(),
		RequestData: reqStr,
		CreatedAt:   time.Now(),
	}

	if err := tx.Create(&log).Error; err != nil {
		return fmt.Errorf("操作日志记录失败: %v", err)
	}

	return nil
}

func (b *BaseController[T]) withTransaction(ctx *gin.Context, fn func(tx *gorm.DB) error) error {
	tx := database.DB.Begin()

	// 执行业务逻辑
	if err := fn(tx); err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		b.Error(ctx, err.Error())
		return err
	}
	if err := recordOperationLog(ctx, tx); err != nil {
		tx.Rollback()
		b.Error(ctx, err.Error())
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		b.Error(ctx, "事务提交失败: "+err.Error())
		return err
	}
	return nil
}

// Create 创建记录
func (b *BaseController[T]) Create(ctx *gin.Context) {
	model := b.Model
	if err := b.HandleValidationError(ctx, &model); err != nil {
		return
	}

	if err := b.withTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Create(&model).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return
	}
	b.Success(ctx, "创建成功", "")
}

// 修改记录
func (b *BaseController[T]) Update(ctx *gin.Context) {
	var data, Model T
	// Model := b.Model
	if err := b.HandleValidationError(ctx, &data); err != nil {
		return
	}

	if err := b.withTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.First(&Model, data.GetID()).Error; err != nil {
			return err
		}
		if err := tx.Save(&data).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return
	}

	b.Success(ctx, "更新成功", "")
}

func (b *BaseController[T]) UpdateStatus(ctx *gin.Context) {
	Model := b.Model
	var data struct {
		ID     uint `json:"id" binding:"required"`
		Status uint `gorm:"column:status;" json:"status" binding:"oneof=0 1"`
	}
	if err := b.HandleValidationError(ctx, &data); err != nil {
		return
	}

	if err := database.DB.First(&Model, data.ID).Error; err != nil {
		b.Error(ctx, err.Error())
		return
	}

	if err := b.withTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Model(&Model).Select("*").Updates(data).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return
	}

	b.Success(ctx, "更新成功", "")
}

// BuildTreeGeneric 通用树构建函数
func BuildTreeGeneric[T any](
	rules []T,
	pid uint,
	getID func(T) uint,
	getPID func(T) uint,
	getTitle func(T) string,
	setChildren func(*T, []T),
	setTitle func(*T, string),
) ([]T, []K_V) {
	tree := make([]T, 0)
	tree2 := make([]K_V, 0)
	for _, rule := range rules {
		if getPID(rule) == pid {
			// 创建节点副本避免修改原始数据
			node := rule

			// 递归构建子节点
			children, children2 := BuildTreeGeneric(
				rules,
				getID(node),
				getID,
				getPID,
				getTitle,
				setChildren,
				setTitle,
			)
			if len(children) > 0 {
				setChildren(&node, children)
			}

			tree = append(tree, node)
			tree2 = append(tree2, K_V{ID: getID(node), Title: getTitle(node)})
			if len(children2) > 0 {
				tree2 = append(tree2, children2...)
			}
		}
	}

	return tree, tree2
}

// 删除记录
func (b *BaseController[T]) Delete(ctx *gin.Context) {
	bmodel := b.Model
	ids := new(model.IDs)
	if err := b.HandleValidationError(ctx, ids); err != nil {
		return
	}
	if err := b.withTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Delete(&bmodel, "id in ?", ids.IDs).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return
	}

	b.Success(ctx, "删除成功", "")
}

// 数据列表
func (b *BaseController[T]) List(ctx *gin.Context) {
	model := b.Model
	var data []map[string]interface{}
	var total int64
	var pagedata Pagination
	if err := b.HandleValidationError(ctx, &pagedata); err != nil {
		return
	}
	db := database.DB.Model(&model)
	db.Count(&total)
	offset := (pagedata.Page - 1) * pagedata.PageSize
	if err := db.Offset(int(offset)).Limit(int(pagedata.PageSize)).Scan(&data).Error; err != nil {
		b.Error(ctx, "获取列表失败")
		return
	}
	b.Success(ctx, "获取列表成功", gin.H{
		"list":      &model,
		"pageTotal": total,
	})
}

func (b *BaseController[T]) GetByID(ctx *gin.Context) {
	Model := b.Model
	var id model.ID
	if err := b.HandleValidationError(ctx, &id); err != nil {
		return
	}
	if err := database.DB.First(&Model, id.ID).Error; err != nil {
		b.Error(ctx, "记录不存在")
		return
	}

	b.Success(ctx, "获取成功", Model)
}

// Success 成功响应
func (b *BaseController[T]) Success(ctx *gin.Context, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  msg,
		"data": data,
	})
}

// Error 错误响应
func (b *BaseController[T]) Error(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 400,
		"msg":  msg,
	})
}
