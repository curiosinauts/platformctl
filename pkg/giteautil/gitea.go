package giteautil

import (
	"code.gitea.io/sdk/gitea"
	"github.com/curiosinauts/platformctl/internal/database"
	"github.com/spf13/viper"
)

func NewGiteaClient() (*gitea.Client, error) {
	accessToken := viper.Get("gitea_access_token").(string)
	giteaURL := viper.Get("gitea_url").(string)
	return gitea.NewClient(giteaURL, gitea.SetToken(accessToken))
}

func AddUser(user database.User) error {
	api, err := NewGiteaClient()
	if err != nil {
		return err
	}
	mustChangePassword := false
	option := gitea.CreateUserOption{
		Username:           user.Username,
		Password:           user.Password,
		Email:              user.Email,
		MustChangePassword: &mustChangePassword,
	}
	_, _, err = api.AdminCreateUser(option)
	if err != nil {
		return err
	}
	return nil
}

func RemoveUser(username string) error {
	api, err := NewGiteaClient()
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
	api, err := NewGiteaClient()
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
	api, err := NewGiteaClient()
	if err != nil {
		return err
	}

	_, err = api.DeleteRepo(username, "project")
	if err != nil {
		return err
	}
	return nil
}

func CreateUserPublicKey(user database.User) (int64, error) {
	api, err := NewGiteaClient()
	if err != nil {
		return 0, err
	}

	option := gitea.CreateKeyOption{
		Key:      user.PublicKey,
		ReadOnly: false,
		Title:    user.Username + " public key",
	}

	publicKey, _, err := api.AdminCreateUserPublicKey(user.Username, option)
	if err != nil {
		return 0, err
	}

	return publicKey.ID, nil
}

func DeleteUserPublicKey(user database.User, keyID int64) error {
	api, err := NewGiteaClient()
	if err != nil {
		return err
	}

	_, err = api.AdminDeleteUserPublicKey(user.Username, int(keyID))
	if err != nil {
		return err
	}
	return nil
}
