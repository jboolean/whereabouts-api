package model

type Role int

const (
	RoleRead             Role = 1
	RoleMaintainLocation      = 2
	RoleAdmin                 = 3
)

type Roles []Role

func (roles *Roles) Contains(role Role) bool {
	for _, cur := range []Role(*roles) {
		if cur == role {
			return true
		}
	}
	return false
}
