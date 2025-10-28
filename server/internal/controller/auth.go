package controller

import (
	"regexp"
	"slices"
	"strings"

	"server/internal/model"
	"server/pkg/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthController 权限控制器
type AuthController struct {
	BaseController[model.AuthRule]
}

// AuthUserController 权限用户控制器
type AuthRoleController struct {
	BaseController[model.AuthRole]
}

// NewAuthController 创建权限控制器实例
func NewAuthController() *AuthController {
	return &AuthController{
		BaseController: BaseController[model.AuthRule]{
			Model: model.AuthRule{},
		},
	}
}

// NewAuthRoleController 创建权限角色控制器实例
func NewAuthRoleController() *AuthRoleController {
	return &AuthRoleController{
		BaseController: BaseController[model.AuthRole]{
			Model: model.AuthRole{},
		},
	}
}

func (b *AuthController) UpdateStatus(ctx *gin.Context) {
	Model := b.Model
	var data struct {
		ID      uint `json:"id" binding:"required"`
		Is_menu uint `gorm:"column:is_menu;" json:"is_menu" binding:"oneof=0 1" name:"是否菜单"`
		Status  uint `gorm:"column:status;" json:"status" binding:"oneof=0 1" name:"状态"`
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

func (b *AuthRoleController) UpdateStatus(ctx *gin.Context) {
	Model := b.Model
	var data struct {
		ID     uint `json:"id" binding:"required"`
		Status uint `gorm:"column:status;" json:"status" binding:"oneof=0 1" name:"状态"`
	}
	if err := b.HandleValidationError(ctx, &data); err != nil {
		return
	}
	if data.ID == 1 {
		b.Error(ctx, "超级管理员角色不能修改状态")
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

func (b *AuthRoleController) Update(ctx *gin.Context) {
	rolemodel := b.Model
	var data model.AuthRole
	if err := b.HandleValidationError(ctx, &data); err != nil {
		return
	}
	if data.ID == 1 {
		b.Error(ctx, "超级管理员不能修改")
		return
	}
	if err := database.DB.First(&rolemodel, data.ID).Error; err != nil {
		b.Error(ctx, "记录不存在")
		return
	}

	var rule []string
	if data.Rule != "*" {
		rule = strings.Split(data.Rule, ",")
		re := regexp.MustCompile(`^[1-9]\d*$`)
		for _, v := range rule {
			if !re.MatchString(v) {
				b.Error(ctx, "规则格式错误")
				return
			}
		}
	}
	var ids []string
	if err := database.DB.Model(&model.AuthRule{}).Select("id").Scan(&ids).Error; err != nil {
		b.Error(ctx, err.Error())
		return
	}
	newids := intersect(rule, ids)
	if len(ids) == len(newids) {
		data.Rule = "*"
	} else {
		data.Rule = strings.Join(newids, ",")
	}

	if err := b.withTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Save(&data).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return
	}

	b.Success(ctx, "更新成功", "")
}

func intersect(slice1, slice2 []string) []string {
	result := []string{}
	for _, v1 := range slice1 {
		for _, v2 := range slice2 {
			if v1 == v2 {
				result = append(result, v1)
				break // 避免重复添加
			}
		}
	}
	return result
}

func (b *AuthRoleController) Create(ctx *gin.Context) {
	rolemodel := b.Model
	if err := b.HandleValidationError(ctx, &rolemodel); err != nil {
		return
	}
	var rule []string
	if rolemodel.Rule != "*" {
		rule = strings.Split(rolemodel.Rule, ",")
		re := regexp.MustCompile(`^[1-9]\d*$`)
		for _, v := range rule {
			if !re.MatchString(v) {
				b.Error(ctx, "规则格式错误")
				return
			}
		}
	}
	var ids []string
	if err := database.DB.Model(&model.AuthRule{}).Select("id").Scan(&ids).Error; err != nil {
		b.Error(ctx, err.Error())
		return
	}
	newids := intersect(rule, ids)
	if len(ids) == len(newids) {
		rolemodel.Rule = "*"
	} else {
		rolemodel.Rule = strings.Join(newids, ",")
	}
	if err := b.withTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Create(&rolemodel).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return
	}

	b.Success(ctx, "创建成功", "")
}

func (b *AuthRoleController) Delete(ctx *gin.Context) {
	rule := model.AuthRole{}
	ids := new(model.IDs)
	if err := b.HandleValidationError(ctx, &ids); err != nil {
		return
	}
	var count int64
	database.DB.Model(&rule).Where("pid in ?", ids.IDs).Count(&count)
	if count > 0 {
		b.Error(ctx, "无法删除含有子角色的角色")
		return
	}
	if slices.Contains(ids.IDs, 1) {
		b.Error(ctx, "无法删除超级管理员角色")
		return
	}
	database.DB.Model(&model.AuthUser{}).Where("role_id in ?", ids.IDs).Count(&count)
	if count > 0 {
		b.Error(ctx, "存在管理员已绑定该角色,无法删除")
		return
	}
	if err := b.withTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Delete(&rule, "id in ?", ids.IDs).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return
	}

	b.Success(ctx, "删除成功", "")
}

func (b *AuthController) Delete(ctx *gin.Context) {
	rule := model.AuthRule{}
	ids := new(model.IDs)
	if err := b.HandleValidationError(ctx, &ids); err != nil {
		return
	}
	if ids.IDs == nil {
		b.Error(ctx, "无效的请求参数")
		return
	}
	var count int64
	database.DB.Model(&rule).Where("pid in ?", ids.IDs).Count(&count)

	if count > 0 {
		b.Error(ctx, "无法删除含有子规则的规则")
		return
	}
	if err := b.withTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Delete(&rule, "id in ?", ids.IDs).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return
	}
	b.Success(ctx, "删除成功", "")
}

func (b *AuthRoleController) Getkv(ctx *gin.Context) {
	var K_V []K_V
	var tree []model.AuthRole
	if err := database.DB.Order("pid asc").Find(&tree).Error; err != nil {
		b.Error(ctx, err.Error())
		return
	}
	_, K_V = b.buildTree(tree, 0)
	b.Success(ctx, "获取成功", K_V)
}

func (b *AuthController) Getkv(ctx *gin.Context) {
	var K_V []K_V
	var tree []model.AuthRule
	if err := database.DB.Order("pid asc,weigh desc").Find(&tree).Error; err != nil {
		b.Error(ctx, err.Error())
		return
	}
	_, K_V = b.buildTree(tree, 0)
	b.Success(ctx, "获取成功", K_V)
}

func (b *AuthController) GetByID(ctx *gin.Context) {
	var K_V []K_V
	var tree []model.AuthRule
	var id model.ID
	if err := b.HandleValidationError(ctx, &id); err != nil {
		return
	}
	model := b.Model
	if err := database.DB.First(&model, id.ID).Error; err != nil {
		b.Error(ctx, err.Error())
		return
	}
	if err := database.DB.Order("pid asc,weigh desc").Find(&tree).Error; err != nil {
		b.Error(ctx, err.Error())
		return
	}
	tree, K_V = b.buildTree(tree, 0)
	b.Success(ctx, "获取成功", gin.H{
		"model": model,
		"tree":  K_V,
	})
}

// 权限列表
func (b *AuthRoleController) List(ctx *gin.Context) {
	var rules []model.AuthRole
	if err := database.DB.Find(&rules).Error; err != nil {
		b.Error(ctx, "获取规则树失败")
		return
	}
	tree, _ := b.buildTree(rules, 0)
	b.Success(ctx, "获取成功", tree)
}

// 权限列表
func (b *AuthController) List(ctx *gin.Context) {
	var rules []model.AuthRule
	var allrules []model.AuthRule
	var data struct {
		Is_menu uint `json:"Is_menu" binding:"oneof=0 1" name:"是否菜单"`
	}
	if err := b.HandleValidationError(ctx, &data); err != nil {
		return
	}
	pid := 0
	query := database.DB
	tree2 := []model.AuthRule{}
	if data.Is_menu == 1 {
		if err := query.Order("weigh desc").Find(&allrules).Error; err != nil {
			b.Error(ctx, "获取规则树数据失败")
			return
		}
		tree2, _ = b.buildTree(allrules, uint(pid))
		query = query.Where("is_menu = 1 and status = 1")
	}
	if err := query.Order("weigh desc").Find(&rules).Error; err != nil {
		b.Error(ctx, "获取规则树失败")
		return
	}
	tree, _ := b.buildTree(rules, uint(pid))
	b.Success(ctx, "获取成功", gin.H{
		"rules":    tree,
		"allrules": tree2,
	})
}

// 权限列表
func (b *AuthController) GetMenu(ctx *gin.Context) {
	var rules []model.AuthRule
	if err := database.DB.Where("is_menu = 1 and pagepath != '' and status=1").Order("weigh desc").Find(&rules).Error; err != nil {
		b.Error(ctx, err.Error())
		return
	}
	b.Success(ctx, "获取成功", rules)
}

// RefreshAuth 刷新权限
func (b *AuthController) RefreshAuth(ctx *gin.Context) {
	// TODO: 实现权限刷新逻辑
	b.Success(ctx, "刷新成功", gin.H{"message": "权限刷新成功"})
}

// buildTree 构建规则树
func (b *AuthController) buildTree(rules []model.AuthRule, pid uint) ([]model.AuthRule, []K_V) {
	return BuildTreeGeneric(
		rules,
		pid,
		func(r model.AuthRule) uint { return r.ID },
		func(r model.AuthRule) uint { return r.PID },
		func(r model.AuthRule) string { return r.Title },
		func(node *model.AuthRule, children []model.AuthRule) { node.Children = children },
		func(node *model.AuthRule, title string) { node.Title = title },
	)
}

// buildTree 构建角色树
func (b *AuthRoleController) buildTree(rules []model.AuthRole, pid uint) ([]model.AuthRole, []K_V) {
	return BuildTreeGeneric(
		rules,
		pid,
		func(r model.AuthRole) uint { return r.ID },
		func(r model.AuthRole) uint { return r.PID },
		func(r model.AuthRole) string { return r.Title },
		func(node *model.AuthRole, children []model.AuthRole) { node.Children = children },
		func(node *model.AuthRole, title string) { node.Title = title },
	)
}
