package domain

import (
	"context"

	"github.com/ashtishad/instabid-wallet/lib"
)

type UserRepository interface {
	Insert(ctx context.Context, u User) (*User, lib.APIError)
	InsertProfile(ctx context.Context, uuid string, up Profile) (*Profile, lib.APIError)

	findByUUID(ctx context.Context, uuid string) (*User, lib.APIError)
	findProfile(ctx context.Context, id int64) (*Profile, lib.APIError)
	checkExists(ctx context.Context, email, username string) lib.APIError
}
