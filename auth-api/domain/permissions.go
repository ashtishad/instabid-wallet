package domain

type RolePermissions struct {
	Permissions map[string][]string
}

// ToDo: Load permissions from policy.json
