package controller

import (
	"html"
	"time"

	"server/internal/model"
	"server/pkg/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleController struct {
	BaseController[model.Article]
}

func NewArticleController() *ArticleController {
	return &ArticleController{
		BaseController: BaseController[model.Article]{
			Model: model.Article{},
		},
	}
}

func (b *ArticleController) List(ctx *gin.Context) {
	model := b.Model
	var article []map[string]interface{}
	var total int64
	var pagedata struct {
		Title      string    `json:"title"`
		Is_hot     *int      `json:"is_hot"`
		Is_top     *int      `json:"is_top"`
		Status     *int      `json:"status"`
		Created_at DateRange `json:"created_at"`
		Updated_at DateRange `json:"updated_at"`
		Pagination
	}
	if err := b.HandleValidationError(ctx, &pagedata); err != nil {
		return
	}
	db := database.DB.Model(&model)
	if pagedata.Title != "" {
		db.Where("title like ?", "%"+pagedata.Title+"%")
	}
	if pagedata.Is_hot != nil {
		db.Where("is_hot = ?", *pagedata.Is_hot)
	}
	if pagedata.Is_top != nil {
		db.Where("is_top = ?", *pagedata.Is_top)
	}
	if pagedata.Status != nil {
		db.Where("status = ?", *pagedata.Status)
	}
	// 处理Created_at日期区间查询
	if !pagedata.Created_at.IsEmpty() {
		startTime, endTime, err := pagedata.Created_at.ParseDateRange()
		if err != nil {
			b.Error(ctx, "创建日期区间格式错误: "+err.Error())
			return
		}
		db.Where("created_at between ? and ?", startTime, endTime)
	}
	// 处理Updated_at日期区间查询
	if !pagedata.Updated_at.IsEmpty() {
		startTime, endTime, err := pagedata.Updated_at.ParseDateRange()
		if err != nil {
			b.Error(ctx, "更新日期区间格式错误: "+err.Error())
			return
		}
		db.Where("updated_at between ? and ?", startTime, endTime)
	}

	db.Count(&total)
	offset := (pagedata.Page - 1) * pagedata.PageSize
	if err := db.Omit("content").Offset(int(offset)).Limit(int(pagedata.PageSize)).Scan(&article).Error; err != nil {
		b.Error(ctx, err.Error())
		return
	}
	for _, v := range article {
		v["created_at"] = v["created_at"].(time.Time).Format("2006-01-02 15:04:05")
		v["updated_at"] = v["updated_at"].(time.Time).Format("2006-01-02 15:04:05")
	}
	b.Success(ctx, "获取列表成功", gin.H{
		"list":      &article,
		"pageTotal": total,
	})
}

func (b *ArticleController) GetByID(ctx *gin.Context) {
	Model := b.Model
	var id model.ID
	if err := b.HandleValidationError(ctx, &id); err != nil {
		return
	}
	if err := database.DB.First(&Model, id.ID).Error; err != nil {
		b.Error(ctx, "记录不存在")
		return
	}
	Model.Content = html.UnescapeString(Model.Content)
	b.Success(ctx, "获取成功", Model)
}

func (b *ArticleController) Create(ctx *gin.Context) {
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

func (b *ArticleController) Update(ctx *gin.Context) {
	oldmodel := b.Model
	model := b.Model
	data := ctx.PostForm("data")
	if err := b.HandleJSONUnmarshalValidation(ctx, data, &model); err != nil {
		return
	}
	if err := database.DB.First(&oldmodel, model.ID).Error; err != nil {
		b.Error(ctx, err.Error())
	}
	// 日志不记录富文本
	model2 := model
	model2.Content = ""
	ctx.Set("requestData", &model2)
	if err := b.withTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Save(&model).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return
	}

	b.Success(ctx, "更新成功", "")
}

func (b *ArticleController) UpdateStatus(ctx *gin.Context) {
	Model := b.Model
	var data struct {
		ID     uint `json:"id" binding:"required"`
		Is_hot int  `gorm:"column:is_hot" json:"is_hot" binding:"oneof=0 1"`
		Is_top int  `gorm:"column:is_top" json:"is_top" binding:"oneof=0 1"`
		Status int  `gorm:"column:status" json:"status" binding:"oneof=0 1"`
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
