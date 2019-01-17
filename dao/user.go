package dao

// DAO上のユーザー
type User struct {
	ID     int64 `gorm:"PRIMARY_KEY"`
	Name   string
	EMails []EMail
}
