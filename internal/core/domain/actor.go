package domain

type Actor struct {
	UserID int
	Role   UserRole
}

func (a Actor) IsAdmin() bool {
	return a.Role == UserRoleAdmin
}
