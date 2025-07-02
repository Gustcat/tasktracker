package model

type Role int32

const (
	UNSPECIFIED Role = 0
	USER        Role = 1
	ADMIN       Role = 2
)

type User struct {
	ID    int64
	Name  string
	Email string
	Role  Role
}
