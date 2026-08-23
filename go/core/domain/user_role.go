package domain

import "fmt"

type UserRole string

const (
	UserRoleSuperAdmin   UserRole = "super_admin"
	UserRoleGymOwner     UserRole = "gym_owner"
	UserRoleGymAdmin     UserRole = "gym_admin"
	UserRoleTrainer      UserRole = "trainer"
	UserRoleReceptionist UserRole = "receptionist"
	UserRoleMember       UserRole = "member"
)

func ParseUserRole(value string) (UserRole, error) {
	role := UserRole(value)
	switch role {
	case UserRoleSuperAdmin, UserRoleGymOwner, UserRoleGymAdmin, UserRoleTrainer, UserRoleReceptionist, UserRoleMember:
		return role, nil
	default:
		return "", fmt.Errorf("invalid user role %q", value)
	}
}

// CanManageEmployeeUsers states the current user-management policy. It does
// not create users; future application use cases must enforce this policy.
func (r UserRole) CanManageEmployeeUsers() bool {
	switch r {
	case UserRoleSuperAdmin, UserRoleGymOwner, UserRoleGymAdmin:
		return true
	default:
		return false
	}
}
