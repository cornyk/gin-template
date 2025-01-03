package models

// User 定义模型
type User struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// TableName 返回模型对应的表名
func (User) TableName() string {
	return "user" // 自定义表名
}
