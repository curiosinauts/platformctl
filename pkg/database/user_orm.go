package database

import (
	"github.com/curiosinauts/platformctl/pkg/crypto"
)

// UserObject user object
type UserObject struct {
	User
	UserService UserService
}

// NewUserObject returns new user object
func NewUserObject(userService UserService, email string) (*UserObject, *DBError) {
	user := User{}
	dberr := userService.FindBy(&user, "hashed_email=$1", crypto.Hashed(email))
	if dberr != nil {
		return nil, dberr
	}
	userObject := &UserObject{}
	userObject.User = user
	userObject.UserService = userService
	return userObject, nil
}

// IDEs returns user ides
func (uo *UserObject) IDEs() ([]IDE, *DBError) {
	var ides []IDE
	userIDEs := []UserIDE{}
	dberr := uo.UserService.ListBy("user_ide", &userIDEs, "user_id=$1", uo.User.ID)
	if dberr != nil {
		return ides, dberr
	}
	for _, userIDE := range userIDEs {
		ide := IDE{}
		dberr := uo.UserService.FindBy(&ide, "id=$1", userIDE.IDEID)
		if dberr != nil {
			return ides, dberr
		}
		ides = append(ides, ide)
	}
	return ides, nil
}

// DoesUserHaveIDE checks if the user has given ide or not
func (uo *UserObject) DoesUserHaveIDE(ideName string) (bool, *DBError) {
	ides, dberr := uo.IDEs()
	if dberr != nil {
		return false, dberr
	}

	for _, ide := range ides {
		if ide.Name == ideName {
			return true, nil
		}
	}

	return false, nil
}

// GetIDE returns IDE object for given ide name
func (uo *UserObject) GetIDE(ideName string) (*IDE, *DBError) {
	ides, dberr := uo.IDEs()
	if dberr != nil {
		return nil, dberr
	}
	for _, ide := range ides {
		if ide.Name == ideName {
			return &ide, nil
		}
	}
	return nil, nil
}

// RuntimeInstallsFor returns runtime install for given ide
func (uo *UserObject) RuntimeInstallsFor(ide IDE) ([]RuntimeInstall, *DBError) {
	var runtimeInstalls []RuntimeInstall

	userIDEs := []UserIDE{}
	dberr := uo.UserService.ListBy("user_ide", &userIDEs, "user_id=$1", uo.User.ID)
	if dberr != nil {
		return runtimeInstalls, dberr
	}

	for _, userIDE := range userIDEs {
		if userIDE.IDEID == ide.ID {
			ideRuntimeInstalls := []IDERuntimeInstall{}
			dberr := uo.UserService.ListBy("ide_runtime_install", &ideRuntimeInstalls, "user_ide_id=$1", userIDE.ID)
			if dberr != nil {
				return runtimeInstalls, dberr
			}
			for _, ideRuntimeInstall := range ideRuntimeInstalls {
				runtimeInstall := RuntimeInstall{}
				dberr := uo.UserService.FindBy(&runtimeInstall, "id=$1", ideRuntimeInstall.RuntimeInstallID)
				if dberr != nil {
					// most likely not found
					continue
				}
				runtimeInstalls = append(runtimeInstalls, runtimeInstall)
			}
		}
	}

	return runtimeInstalls, nil
}

// UserIDE returns user ide for user's given ide
func (uo *UserObject) UserIDE(ide IDE) (UserIDE, *DBError) {
	var userIDE UserIDE

	userIDEs := []UserIDE{}
	dberr := uo.UserService.ListBy("user_ide", &userIDEs, "user_id=$1", uo.User.ID)
	if dberr != nil {
		return userIDE, dberr
	}

	for _, userIDE := range userIDEs {
		if userIDE.IDEID == ide.ID {
			return userIDE, nil
		}
	}

	return userIDE, nil
}

// DoesUserHaveRuntimeInstallFor checks if user have runtime install for given ide and runtime install name
func (uo *UserObject) DoesUserHaveRuntimeInstallFor(ide IDE, runtimeInstallName string) (bool, *DBError) {
	runtimeInstalls, dberr := uo.RuntimeInstallsFor(ide)
	if dberr != nil {
		return false, dberr
	}
	for _, runtimeInstall := range runtimeInstalls {
		if runtimeInstall.Name == runtimeInstallName {
			return true, nil
		}
	}
	return false, nil
}
