package domain

import "strings"

type RolePermissions map[string]map[string]bool

func (p RolePermissions) IsAuthorizedFor(role string, routeName string) bool {
	perms, ok := p[role]
	if !ok {
		return false
	}

	return perms[strings.TrimSpace(routeName)]
}

var Permissions = RolePermissions{
	"admin": {
		"POST:/users":          true,
		"POST:/users/:user_id": true,
	},
	"user": {
		"POST:/users/:user_id": true,
	},
}
