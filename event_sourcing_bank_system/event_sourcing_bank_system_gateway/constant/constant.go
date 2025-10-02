package constant

import "time"

const (
	UserPermissionKeyPrefix = "user::permission::%v"
	PermissionExpiredTime   = time.Hour * time.Duration(30*24)

	UserStorePermissionKeyPrefix = "user::store_permission::%v"
	StorePermissionExpiredTime   = time.Hour * time.Duration(30*24)

	KeyAuthToken      = "auth-token"
	KeyUserLoginID    = "user-login-id"
	KeyUserLoginEmail = "user-login-email"
)
