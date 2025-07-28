package models

// UserModel user table
type UserModel struct {
	Id    int32  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name  string `gorm:"column:name" json:"name"`
	Email string `gorm:"column:email" json:"email"`
}

func (u *UserModel) TableName() string {
	return "user"
}
