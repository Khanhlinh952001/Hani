package users

const (
	RoleUser  = 0
	RoleAdmin = 1
)

func IsAdmin(role int) bool {
	return role == RoleAdmin
}
