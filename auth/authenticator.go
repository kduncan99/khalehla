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

func (a *Authenticator) AddAccount(loginName string,
									accountPassword string,
									userName string,
									emailAddress string) (string, error) {
	query := "INSERT INTO auth_accounts " +
		"(loginName, hashedPassword, disabled, userName, emailAddress) " +
		"VALUES ($1, $2, false, $3, $4) " +
		"RETURNING accountId"

	hash, _ := a.hashPassword(accountPassword)
	var id string
	err := a.database.SQL.QueryRow(query, loginName, hash, userName, emailAddress).Scan(&id)
	if err != nil {
		log.Print(err.Error())
		return "", err
	}

	return id, nil
}

func (a *Authenticator) AddBaseGroup(groupName string) (string, error) {
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

func (a *Authenticator) AddPrivilege(privName string) (string, error) {
	query := "INSERT INTO auth_privileges " +
		"(privilegeName) " +
		"VALUES ($1) " +
		"RETURNING privilegeID"

	var id string
	err := a.database.SQL.QueryRow(query, privName).Scan(&id)
	if err != nil {
		log.Print(err.Error())
		return "", err
	}

	return id, nil
}

func (a *Authenticator) AddSubGroup(groupName string, parentGroupID string) (string, error) {
	query := "INSERT INTO auth_groups " +
		"(groupName, parentGroupID) " +
		"VALUES ($1, $2) " +
		"RETURNING groupID"

	var id string
	err := a.database.SQL.QueryRow(query, groupName, parentGroupID).Scan(&id)
	if err != nil {
		log.Print(err.Error())
		return "", err
	}

	return id, nil
}

func (a *Authenticator) ConnectAccountToGroup(accountId string, groupId string) error {
	query := "INSERT INTO auth_accounts_groups " +
		"(accountID, groupID) " +
		"VALUES ($1, $2) "
	_, err := a.database.SQL.Exec(query, accountId, groupId)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	return nil
}

func (a *Authenticator) ConnectAccountToPrivilege(accountId string, privilegeId string) error {
	query := "INSERT INTO auth_accounts_privileges " +
		"(accountID, privilegeID) " +
		"VALUES ($1, $2)"
	_, err := a.database.SQL.Exec(query, accountId, privilegeId)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	return nil
}

func (a *Authenticator) DisconnectAccountFromGroup(accountId string, groupId string) error {
	query := "DELETE FROM auth_accounts_groups " +
		"WHERE accountID=$1 AND groupID=$2 "
	_, err := a.database.SQL.Exec(query, accountId, groupId)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	return nil
}

func (a *Authenticator) DisconnectAccountFromPrivilege(accountId string, privilegeId string) error {
	query := "DELETE FROM auth_accounts_privileges " +
		"WHERE accountID=$1 AND privilegeID=$2"
	_, err := a.database.SQL.Exec(query, accountId, privilegeId)
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
