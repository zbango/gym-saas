package domain

import "testing"

func TestParseUserRole(t *testing.T) {
	for _, role := range []UserRole{
		UserRoleSuperAdmin, UserRoleGymOwner, UserRoleGymAdmin,
		UserRoleTrainer, UserRoleReceptionist, UserRoleMember,
	} {
		got, err := ParseUserRole(string(role))
		if err != nil || got != role {
			t.Fatalf("ParseUserRole(%q) = %q, %v", role, got, err)
		}
	}

	if _, err := ParseUserRole("owner"); err == nil {
		t.Fatal("ParseUserRole accepted unknown role")
	}
}

func TestUserRoleCanManageEmployeeUsers(t *testing.T) {
	for _, role := range []UserRole{UserRoleSuperAdmin, UserRoleGymOwner, UserRoleGymAdmin} {
		if !role.CanManageEmployeeUsers() {
			t.Fatalf("%q cannot manage employee users", role)
		}
	}

	for _, role := range []UserRole{UserRoleTrainer, UserRoleReceptionist, UserRoleMember} {
		if role.CanManageEmployeeUsers() {
			t.Fatalf("%q can manage employee users", role)
		}
	}
}
