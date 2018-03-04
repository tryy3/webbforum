package models

import (
	"fmt"
)

// Permission is the database model for the permission list
type Permission struct {
	Bit  uint64 `gorm:"primary_key"`
	Name string
}

// DefaultPermission is a list of default permissions
var DefaultPermission = PermissionHandler{
	map[string]uint64{
		"testy-permission": 1,
		"foo-bar":          2,
		"nana-ko":          4,
		"cho-cho":          8,
	},
}

// PermissionHandler is a read only representation of the permissions
type PermissionHandler struct {
	permissions map[string]uint64
}

// Permission will retrieve the permission bit of a specific permission
func (p PermissionHandler) Permission(key string) uint64 {
	perm, ok := p.permissions[key]
	if !ok {
		return 0
	}
	return perm
}

// Permissions returns the default permission list
func (p PermissionHandler) Permissions() map[string]uint64 {
	return p.permissions
}

// PermissionList will retrieve a list of all the permissions that the permission bits owns
func (p PermissionHandler) PermissionList(perms uint64) []string {
	var permList []string
	for k, v := range p.permissions {
		if perms&v != 0 {
			permList = append(permList, k)
		}
	}
	return permList
}

// ParsePermissions will take a list of permissions and return the permission bit
func (p PermissionHandler) ParsePermissions(perms []string) (uint64, error) {
	var perm uint64
	for _, permName := range perms {
		permBit := p.Permission(permName)
		if permBit == 0 {
			return 0, fmt.Errorf("unknown permission name: %s", permName)
		}
		if perm&permBit != 0 {
			return 0, fmt.Errorf("duplicated permission name: %s", permName)
		}
		perm += permBit
	}
	return perm, nil
}
