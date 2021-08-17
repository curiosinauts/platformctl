package giteautil

import (
	"code.gitea.io/sdk/gitea"
	"fmt"
	"github.com/curiosinauts/platformctl/internal/database"
)

func AddUser(user database.User) error {
	api, err := gitea.NewClient("https://git-web.curiosityworks.org", gitea.SetToken(""))
	if err != nil {
		return err
	}
	option := gitea.CreateUserOption{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	}
	_, _, err = api.AdminCreateUser(option)
	if err != nil {
		return err
	}
	return nil
}

func RemoveUser(username string) error {
	api, err := gitea.NewClient("https://git-web.curiosityworks.org", gitea.SetToken(""))
	if err != nil {
		return err
	}
	_, err = api.AdminDeleteUser(username)
	if err != nil {
		return err
	}
	return nil
}

func CreateUserRepo(username string) error {
	api, err := gitea.NewClient("https://git-web.curiosityworks.org", gitea.SetToken(""))
	if err != nil {
		return err
	}

	option := gitea.CreateRepoOption{
		Name:     "project",
		AutoInit: true,
	}

	_, _, err = api.AdminCreateRepo(username, option)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUserRepo(username string) error {
	api, err := gitea.NewClient("https://git-web.curiosityworks.org", gitea.SetToken(""))
	if err != nil {
		return err
	}

	_, err = api.DeleteRepo(username, "project")
	if err != nil {
		return err
	}
	return nil
}

func CreateUserPublicKey(user database.User) error {
	api, err := gitea.NewClient("https://git-web.curiosityworks.org", gitea.SetToken(""))
	if err != nil {
		return err
	}

	option := gitea.CreateKeyOption{
		Key:      user.PublicKey,
		ReadOnly: false,
		Title:    user.Username + " public key",
	}

	publicKey, _, err := api.AdminCreateUserPublicKey(user.Username, option)
	if err != nil {
		return err
	}

	fmt.Println(publicKey.ID)
	return nil
}

func DeleteUserPublicKey(user database.User, keyID int) error {
	api, err := gitea.NewClient("https://git-web.curiosityworks.org", gitea.SetToken(""))
	if err != nil {
		return err
	}

	_, err = api.AdminDeleteUserPublicKey(user.Username, keyID)
	if err != nil {
		return err
	}
	return nil
}
