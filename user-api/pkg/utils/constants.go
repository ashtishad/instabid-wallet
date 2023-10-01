package utils

import "time"

const (
	UUIDRegex   = `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[1-5][a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`
	StatusRegex = `^(active|inactive|deleted)$`
	RoleRegex   = `^(user|admin|moderator|merchant)$`

	DefaultPageSize = 20

	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
	UserStatusDeleted  = "deleted"

	UserRoleAdmin     = "admin"
	UserRoleUser      = "user"
	UserRoleModerator = "moderator"
	UserRoleMerchant  = "merchant"

	TimeoutCreateUser        = 200 * time.Millisecond
	TimeoutCreateUserProfile = 200 * time.Millisecond
)
