package model

type Config struct {
	ID      uint   `gorm:"column:id;primaryKey;type:int(12) unsigned;comment:ID" json:"id"`
	Name    string `gorm:"column:name;type:varchar(100);comment:变量名称;not null" json:"name" `
	Title   string `gorm:"column:title;type:varchar(100);comment:标题;not null" json:"title"`
	Value   string `gorm:"column:value;type:text;comment:变量值" json:"value"`
	Group   string `gorm:"column:group;type:varchar(100);comment:分组;not null" json:"group"`
	Type    string `gorm:"column:type;type:varchar(100);comment:类型;not null" json:"type"`
	Content string `gorm:"column:content;type:text;comment:扩展内容" json:"content"`
	Extend  string `gorm:"column:extend;type:varchar(255);comment:扩展属性" json:"extend"`
}

func (c Config) GetID() uint {
	return c.ID
}
