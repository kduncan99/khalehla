package auth

import (
	"golang.org/x/crypto/bcrypt"
	"kalehla/db"
	"log"
)

type Authenticator struct {
	database *db.Database
}

func (a *Authenticator) Close() error {
	return a.database.Close()
}

func (a *Authenticator) AddAccount(accountName string,
									accountPassword string,
									userName string,
									emailAddress string) (string, error) {
	query := "INSERT INTO auth_accounts " +
		"(accountName, passwordHash, disabled, userName, emailAddress) " +
		"VALUES ($1, $2, 0, $3, $4) " +
		"RETURNING accountID"

	hash, _ := a.hashPassword(accountPassword)
	var id string
	err := a.database.SQL.QueryRow(query, accountName, hash, userName, emailAddress).Scan(&id)
	if err != nil {
		log.Print(err.Error())
		return "", err
	}

	return id, nil
}

func (a *Authenticator) AddGroup(groupName string) (string, error) {
	query := "INSERT INTO auth_groups " +
		"(groupName) " +
		"VALUES ($1) " +
		"RETURNING groupID"

	var id string
	err := a.database.SQL.QueryRow(query, groupName).Scan(&id)
	if err != nil {
		log.Print(err.Error())
		return "", err
	}

	return id, nil
}

func (a *Authenticator) AddPrivilege(privilegeName string) (string, error) {
	query := "INSERT INTO auth_privileges " +
		"(privilegeName) " +
		"VALUES ($1) " +
		"RETURNING privilegeID"

	var id string
	err := a.database.SQL.QueryRow(query, privilegeName).Scan(&id)
	if err != nil {
		log.Print(err.Error())
		return "", err
	}

	return id, nil
}

func (a *Authenticator) ConnectAccountToGroup(accountName string, groupName string) error {
	query := "INSERT INTO auth_accounts_groups " +
		"(accountID, groupID) " +
		"SELECT a.accountID, b.groupID " +
		"FROM auth_accounts a, auth_groups b " +
		"WHERE accountName = $1 AND groupName = $2"
	_, err := a.database.SQL.Exec(query, accountName, groupName)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	return nil
}

func (a *Authenticator) ConnectGroupToPrivilege(groupName string, privilegeName string) error {
	query := "INSERT INTO auth_groups_privileges " +
		"(groupID, privilegeID) " +
		"SELECT a.groupID, b.privilegeID " +
		"FROM auth_groups a, auth_privileges b " +
		"WHERE groupName = $1 AND privilegeName = $2"
	_, err := a.database.SQL.Exec(query, groupName, privilegeName)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	return nil
}

func (a *Authenticator) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (a *Authenticator) Open() error {
	cfg, err := db.LoadConfiguration()
	if err != nil {
		return err
	}

	a.database, err = db.New(cfg)
	if err != nil {
		return err
	}

	return nil
}
