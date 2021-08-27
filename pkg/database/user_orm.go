package database

import "github.com/curiosinauts/platformctl/pkg/crypto"

type UserObject struct {
	User
	UserService UserService
}

func NewUserObject(userService UserService, email string) (*UserObject, *DBError) {
	user, dberr := userService.FindUserByHashedEmail(crypto.Hashed(email))
	if dberr != nil {
		return nil, dberr
	}
	userObject := &UserObject{}
	userObject.User = user
	userObject.UserService = userService
	return userObject, nil
}

func (uo *UserObject) IDEs() ([]IDE, *DBError) {
	var ides []IDE
	userIDEs, dberr := uo.UserService.FindUserIDEsByUser(uo.User)
	if dberr != nil {
		return ides, dberr
	}
	for _, userIDE := range userIDEs {
		ide, dberr := uo.UserService.FindIDEByID(userIDE.IDEID)
		if dberr != nil {
			return ides, dberr
		}
		ides = append(ides, ide)
	}
	return ides, nil
}

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

func (uo *UserObject) RuntimeInstallsFor(ide IDE) ([]RuntimeInstall, *DBError) {
	var runtimeInstalls []RuntimeInstall

	userIDEs, dberr := uo.UserService.FindUserIDEsByUser(uo.User)
	if dberr != nil {
		return runtimeInstalls, dberr
	}

	for _, userIDE := range userIDEs {
		if userIDE.IDEID == ide.ID {
			ideRuntimeInstalls, dberr := uo.UserService.FindIDERuntimeInstallsByUserIDE(userIDE)
			if dberr != nil {
				return runtimeInstalls, dberr
			}
			for _, ideRuntimeInstall := range ideRuntimeInstalls {
				runtimeInstall, dberr := uo.UserService.FindRuntimeInstallByID(ideRuntimeInstall.ID)
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

func (uo *UserObject) UserIDE(ide IDE) (UserIDE, *DBError) {
	var userIDE UserIDE

	userIDEs, dberr := uo.UserService.FindUserIDEsByUser(uo.User)
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
