package caching

import (
	"bcca_crawler/internal/auth/roles"
	"fmt"	
	"sync"
	"time"

	"github.com/google/uuid"
)

var roleCache sync.Map

type CacheEntry struct {	
	Role roles.Role
	Expiration time.Time
}

func SetRoleCache(userID uuid.UUID, role roles.Role, expiration time.Time) {
	roleCache.Store(userID, CacheEntry{role, expiration})
}

func GetRoleCache(userID uuid.UUID) (roles.Role, error) {
	if v, ok := roleCache.Load(userID); ok {
		entry := v.(CacheEntry)
		if entry.Expiration.After(time.Now()) {
			return entry.Role, nil
		}
		roleCache.Delete(userID)
	}
	return -1, fmt.Errorf("cache not found or expired")
}

func DeleteRoleCache(userID uuid.UUID) {
	roleCache.Delete(userID)
}


