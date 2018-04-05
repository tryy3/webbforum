package models

import (
	"sort"

	"github.com/jinzhu/gorm"
)

// Permission is the database model for the permission list
type Permission struct {
	// General information
	ID         uint `gorm:"primary_key"`
	UserID     uint
	GroupID    uint
	Permission string
}

// ParsedPermission is used when you need to parse permissions programmatically
type ParsedPermission struct {
	Permission string
	Title      string
	Has        bool
}

// byTitle is a list of ParsedPermission when you need to sort the permission list
type byTitle []ParsedPermission

func (a byTitle) Len() int           { return len(a) }
func (a byTitle) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byTitle) Less(i, j int) bool { return a[i].Permission < a[j].Permission }

// PermissionList is a struct of all the valid permissions
type PermissionList struct {
	CREATEPOST string
	EDITPOST   string
	DELETEPOST string

	EDITSELFPOST   string
	DELETESELFPOST string

	CREATETHREAD string
	EDITTHREAD   string
	DELETETHREAD string

	EDITSELFTHREAD   string
	DELETESELFTHREAD string

	CREATECATEGORY string
	EDITCATEGORY   string
	DELETECATEGORY string

	CREATEGROUP string
	EDITGROUP   string
	DELETEGROUP string

	EDITUSER    string
	EDITPROFILE string
}

// Permissions is a list of the names of all the permissions
var Permissions = PermissionList{
	CREATEPOST: "create_post",
	EDITPOST:   "edit_post",
	DELETEPOST: "delete_post",

	EDITSELFPOST:   "edit_self_post",
	DELETESELFPOST: "delete_self_post",

	CREATETHREAD: "create_thread",
	EDITTHREAD:   "edit_thread",
	DELETETHREAD: "delete_thread",

	EDITSELFTHREAD:   "edit_self_thread",
	DELETESELFTHREAD: "delete_self_thread",

	CREATECATEGORY: "create_category",
	EDITCATEGORY:   "edit_category",
	DELETECATEGORY: "delete_category",

	CREATEGROUP: "create_group",
	EDITGROUP:   "edit_group",
	DELETEGROUP: "delete_group",

	EDITUSER:    "edit_user",
	EDITPROFILE: "edit_profile",
}

// TitledPermissions is a list of all the titles for the permissions
var TitledPermissions = map[string]string{
	Permissions.CREATEPOST: "Skapa ett inlägg",
	Permissions.EDITPOST:   "Modifiera ett inlägg",
	Permissions.DELETEPOST: "Ta bort ett inlägg",

	Permissions.EDITSELFPOST:   "Modifiera ditt egna inlägg",
	Permissions.DELETESELFPOST: "Ta bort ditt egna inlägg",

	Permissions.CREATETHREAD: "Skapa en tråd",
	Permissions.EDITTHREAD:   "Modifiera en tråd",
	Permissions.DELETETHREAD: "Ta bort en tråd",

	Permissions.EDITSELFTHREAD:   "Modifiera din egna tråd",
	Permissions.DELETESELFTHREAD: "Ta bort din egna tråd",

	Permissions.CREATECATEGORY: "Skapa en kategori",
	Permissions.EDITCATEGORY:   "Modifiera en kategori",
	Permissions.DELETECATEGORY: "Ta bort en kategori",

	Permissions.CREATEGROUP: "Skapa en grupp",
	Permissions.EDITGROUP:   "Modifiera en grupp",
	Permissions.DELETEGROUP: "Ta bort en grupp",

	Permissions.EDITUSER:    "Modifiera en användare",
	Permissions.EDITPROFILE: "Modifiera din egna profil",
}

// HasPermission will check if a specific user has access to a permission either directly or through a group
func HasPermission(db *gorm.DB, user *User, perm string) (bool, error) {
	result := db.Where("user_id = ? AND permission = ?", user.ID, perm).First(&Permission{})
	if result.Error != nil && !result.RecordNotFound() {
		return false, result.Error
	}
	if result.Error == nil {
		return true, nil
	}

	var group = user.Group

	if group == nil {
		group = &Group{ID: 1}
	}

	for {
		if group == nil {
			return false, nil
		}

		result := db.Where("group_id = ? AND permission = ?", group.ID, perm).First(&Permission{})
		if result.Error != nil && !result.RecordNotFound() {
			return false, result.Error
		}
		if result.Error == nil {
			return true, nil
		}

		group = group.Parent
	}
}

// GetGroupPermission will retrieve a list of all the permissions a group has
func GetGroupPermission(db *gorm.DB, group *Group) ([]ParsedPermission, error) {
	var permissions []Permission
	stmt := db.Where("group_id = ?", group.ID).Find(&permissions)
	if stmt.Error != nil && !stmt.RecordNotFound() {
		return nil, stmt.Error
	}

	var parsed []ParsedPermission
	for perm, title := range TitledPermissions {
		parsed = append(parsed, ParsedPermission{perm, title, false})
	}

	for i, perm := range parsed {
		for _, p := range permissions {
			if perm.Permission == p.Permission {
				parsed[i].Has = true
				break
			}
		}
	}

	sort.Sort(byTitle(parsed))
	return parsed, nil
}
