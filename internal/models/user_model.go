package models

// User user table
type User struct {
	Id    int32  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name  string `gorm:"column:name" json:"name"`
	Email string `gorm:"column:email" json:"email"`
}

func (u *User) TableName() string {
	return "user"
}
