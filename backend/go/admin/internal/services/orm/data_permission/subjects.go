package datapermission

import "admin/internal/fiberc/handler"

// Subject identifies one permission principal. The permission lookup combines
// the current user, assigned roles, and wildcard user/role subjects.
type Subject struct {
	Type string
	ID   uint64
}

// BuildPermissionSubjects returns all subjects that can grant permissions to
// a user: the user itself, all roles, and wildcard user/role subjects.
func BuildPermissionSubjects(userID uint64, roleIDs []uint64) []Subject {
	subjects := []Subject{
		{Type: "USER", ID: userID},
		{Type: "ANY_USER", ID: 0},
		{Type: "ANY_ROLE", ID: 0},
	}
	for _, roleID := range roleIDs {
		subjects = append(subjects, Subject{Type: "ROLE", ID: roleID})
	}
	return subjects
}

// BuildPermissionSubjectsFromCtx derives permission subjects from the current session.
func BuildPermissionSubjectsFromCtx(ctx *handler.Ctx) []Subject {
	return BuildPermissionSubjects(ctx.SessionInfo.Id, ctx.SessionInfo.RoleIDs)
}
