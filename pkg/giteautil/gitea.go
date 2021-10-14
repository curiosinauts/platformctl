package giteautil

import (
	"code.gitea.io/sdk/gitea"
	"github.com/curiosinauts/platformctl/pkg/database"
	"github.com/spf13/viper"
)

// GitClient gitea api client
type GitClient struct {
	api *gitea.Client
}

// NewGitClient returns a new gitea api client
func NewGitClient() (*GitClient, error) {
	accessToken := viper.Get("gitea_access_token").(string)
	giteaURL := viper.Get("gitea_url").(string)
	api, err := gitea.NewClient(giteaURL, gitea.SetToken(accessToken))
	if err != nil {
		return nil, err
	}
	return &GitClient{api: api}, nil
}

// AddUser adds new user
func (gc *GitClient) AddUser(user database.User) error {
	mustChangePassword := false
	option := gitea.CreateUserOption{
		Username:           user.Username,
		Password:           user.Password,
		Email:              user.Email,
		MustChangePassword: &mustChangePassword,
	}
	_, _, err := gc.api.AdminCreateUser(option)
	if err != nil {
		return err
	}
	return nil
}

// RemoveUser removes user
func (gc *GitClient) RemoveUser(username string) error {
	_, err := gc.api.AdminDeleteUser(username)
	if err != nil {
		return err
	}
	return nil
}

// CreateUserRepo creates user repo
func (gc *GitClient) CreateUserRepo(username string) error {

	option := gitea.CreateRepoOption{
		Name:     "project",
		AutoInit: true,
	}

	_, _, err := gc.api.AdminCreateRepo(username, option)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUserRepo deletes user repo
func (gc *GitClient) DeleteUserRepo(username string) error {
	_, err := gc.api.DeleteRepo(username, "project")
	if err != nil {
		return err
	}
	return nil
}

// CreateUserPublicKey creates user public key
func (gc *GitClient) CreateUserPublicKey(user database.User) (int64, error) {
	option := gitea.CreateKeyOption{
		Key:      user.PublicKey,
		ReadOnly: false,
		Title:    user.Username + " public key",
	}

	publicKey, _, err := gc.api.AdminCreateUserPublicKey(user.Username, option)
	if err != nil {
		return 0, err
	}

	return publicKey.ID, nil
}

// DeleteUserPublicKey deletes user public key
func (gc *GitClient) DeleteUserPublicKey(user database.User, keyID int64) error {
	_, err := gc.api.AdminDeleteUserPublicKey(user.Username, int(keyID))
	if err != nil {
		return err
	}
	return nil
}
