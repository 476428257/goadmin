package model

type Article struct {
	Commonmodel
	Title   string `gorm:"column:title" json:"title" binding:"required" name:"标题"`
	Image   string `gorm:"column:image" json:"image" binding:"required" name:"图片"`
	Is_hot  int    `gorm:"column:is_hot" json:"is_hot" binding:"min=0,oneof=0 1" name:"是否热门"`
	Is_top  int    `gorm:"column:is_top" json:"is_top" binding:"min=0,oneof=0 1" name:"是否置顶"`
	Status  int    `gorm:"column:status" json:"status" binding:"required,oneof=0 1" name:"状态"`
	Content string `gorm:"column:content" json:"content" binding:"required" name:"内容"`
}
