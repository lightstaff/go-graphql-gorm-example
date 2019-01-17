package dao

import "github.com/lightstaff/go-graphql-gorm-example/dao/types"

// DAO上のメール
type EMail struct {
	ID      int64 `gorm:"PRIMARY_KEY"`
	Address string
	Remarks types.NullString
	UserID  int64
	User    *User
}
