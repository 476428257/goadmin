package controller

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"os"
	"path/filepath"
	"time"

	"server/config"
	"server/internal/model"
	"server/pkg/database"

	"github.com/gin-gonic/gin"
)

type SystemController struct {
	BaseController[model.OperationLog]
}

func NewSystemController() *SystemController {
	return &SystemController{
		BaseController: BaseController[model.OperationLog]{
			Model: model.OperationLog{},
		},
	}
}

func randomString(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "randerr"
	}
	return hex.EncodeToString(b)
}

// ProcessUpload 处理文件上传的公共方法，返回上传后的文件路径列表
func (b *SystemController) ProcessUpload(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		b.Error(ctx, err.Error())
		return
	}
	dateDir := time.Now().Format("20060102")
	uploadPath := filepath.Join(config.GetConfig().Upload.UploadDir, dateDir)
	files := form.File["files"]
	if len(files) == 0 {
		b.Error(ctx, "请选择文件")
		return
	}
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		b.Error(ctx, "创建目录失败: "+err.Error())
		return
	}
	var filePaths []string
	host := ctx.Request.Host
	scheme := "http"
	if ctx.GetHeader("X-Forwarded-Proto") == "https" || ctx.Request.TLS != nil {
		scheme = "https"
	}
	for _, file := range files {
		ext := filepath.Ext(file.Filename)
		newName := randomString(8) + ext
		savePath := filepath.Join(uploadPath, newName)
		// 保存文件
		if err := ctx.SaveUploadedFile(file, savePath); err != nil {
			b.Error(ctx, "上传失败: "+err.Error())
			return
		}
		filePaths = append(filePaths, scheme+"://"+host+"/"+filepath.ToSlash(savePath))
	}
	b.Success(ctx, "上传成功", gin.H{
		"filePaths": filePaths,
	})
}

func (b *SystemController) Aclog(ctx *gin.Context) {
	var OperationLog []map[string]interface{}
	var total int64

	var pagedata struct {
		Username    string    `json:"username"`
		RequestData string    `json:"request_data"`
		Path        string    `json:"path"`
		Date        DateRange `json:"created_at"`
		Pagination
	}
	if err := b.HandleValidationError(ctx, &pagedata); err != nil {
		return
	}
	db := database.DB.Model(&model.OperationLog{})
	if pagedata.Username != "" {
		db = db.Where("username like ?", "%"+pagedata.Username+"%")
	}
	if pagedata.RequestData != "" {
		db = db.Where("request_data like ?", "%"+pagedata.RequestData+"%")
	}
	if pagedata.Path != "" {
		db = db.Where("path like ?", "%"+pagedata.Path+"%")
	}
	if !pagedata.Date.IsEmpty() {
		startTime, endTime, err := pagedata.Date.ParseDateRange()
		if err != nil {
			b.Error(ctx, "创建日期区间格式错误: "+err.Error())
			return
		}
		db.Where("created_at between ? and ?", startTime, endTime)
	}
	userID, _ := ctx.Get("userID")
	// 查询当前管理员的角色ID
	var roleID uint
	if err := database.DB.Model(&model.AuthUser{}).
		Where("admin_id = ?", userID).
		Select("role_id").
		Scan(&roleID).Error; err != nil {
		b.Error(ctx, "查询用户角色失败: "+err.Error())
		return
	}
	if roleID != 1 {
		db = db.Where("admin_id = ?", userID)
	}
	db.Count(&total)
	offset := (pagedata.Page - 1) * pagedata.PageSize
	if err := db.Offset(int(offset)).Order("created_at desc").Limit(int(pagedata.PageSize)).Scan(&OperationLog).Error; err != nil {
		b.Error(ctx, err.Error())
		return
	}
	for _, v := range OperationLog {
		v["created_at"] = v["created_at"].(time.Time).Format("2006-01-02 15:04:05")
	}
	b.Success(ctx, "获取列表成功", gin.H{
		"list":      &OperationLog,
		"pageTotal": total,
	})
}

func secureRandom(max int) int {
	b := make([]byte, 4)
	rand.Read(b)
	return int(binary.BigEndian.Uint32(b)) % max
}

func (b *SystemController) Dashboard(ctx *gin.Context) {
	b.Success(ctx, "获取数据成功", gin.H{
		"n1": secureRandom(10000),
		"n2": secureRandom(10000),
		"n3": secureRandom(10000),
		"n4": secureRandom(10000),
	})
}
