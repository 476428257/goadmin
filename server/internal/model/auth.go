package model

// AuthRole 权限角色模型
type AuthRole struct {
	ID       uint       `gorm:"column:id;primarykey" json:"id"`
	PID      uint       `gorm:"column:pid;default:0;" json:"pid"`
	Title    string     `gorm:"column:title;size:50;not null;uniqueIndex;" json:"title" binding:"required,min=2,max=50" name:"角色名称"`
	Status   uint       `gorm:"column:status;not null;" json:"status" binding:"oneof=0 1" name:"状态"`
	Rule     string     `gorm:"column:rule;type:text;" json:"rule" binding:"required" name:"规则字符串"`
	Children []AuthRole `gorm:"-" json:"children,omitempty"`
}

// GetID 获取ID
func (a AuthRole) GetID() uint {
	return a.ID
}

// AuthRule 权限路由模型
type AuthRule struct {
	ID       uint       `gorm:"column:id;primarykey" json:"id"`
	PID      uint       `gorm:"column:pid;default:0;" json:"pid"`
	Title    string     `gorm:"column:title;size:100;not null;" json:"title" binding:"required,min=2,max=100" name:"规则标题"`
	Route    string     `gorm:"column:route;size:200;not null;uniqueIndex;" json:"route" binding:"required,min=1,max=200" name:"路由路径"`
	Pagepath string     `gorm:"column:pagepath;size:200;" json:"pagepath" name:"页面路径"`
	IsMenu   uint       `gorm:"column:is_menu;not null;" json:"is_menu" binding:"oneof=0 1" name:"是否菜单"`
	Status   uint       `gorm:"column:status;not null;" json:"status" binding:"oneof=0 1" name:"状态"`
	Icon     string     `gorm:"column:icon;" json:"icon" `
	Weigh    uint       `gorm:"column:weigh;default:0;" json:"weigh" `
	Children []AuthRule `gorm:"-" json:"children,omitempty"`
}

// GetID 获取ID
func (a AuthRule) GetID() uint {
	return a.ID
}

// AuthUser 用户角色关联模型
type AuthUser struct {
	ID      uint     `gorm:"column:id;primarykey" json:"id"`
	RoleID  uint     `gorm:"column:role_id;not null;comment:角色ID" json:"role_id" binding:"required"`
	AdminID uint     `gorm:"column:admin_id;not null;uniqueIndex;comment:管理员ID" json:"admin_id" binding:"required"`
	Role    AuthRole `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// GetID 获取ID
func (a AuthUser) GetID() uint {
	return a.ID
}
