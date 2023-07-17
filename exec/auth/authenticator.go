package auth

import (
	"golang.org/x/crypto/bcrypt"
	"khalehla/exec/db"
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

func (a *Authenticator) AddGroup(subsystemName string, groupName string) (string, error) {
	query := "INSERT INTO auth_groups " +
		"(subsystemID, groupName) " +
		"SELECT a.subsystemID, $2 " +
		"FROM auth_subsystems a " +
		"WHERE subsystemName = $1 " +
		"RETURNING groupID"

	var id string
	err := a.database.SQL.QueryRow(query, subsystemName, groupName).Scan(&id)
	if err != nil {
		log.Print(err.Error())
		return "", err
	}

	return id, nil
}

func (a *Authenticator) AddPrivilege(subsystemName string, privilegeName string) (string, error) {
	query := "INSERT INTO auth_privileges " +
		"(subsystemID, privilegeName) " +
		"SELECT a.subsystemID, $2 " +
		"FROM auth_subsystems a " +
		"WHERE subsystemName = $1 " +
		"RETURNING privilegeID"

	var id string
	err := a.database.SQL.QueryRow(query, subsystemName, privilegeName).Scan(&id)
	if err != nil {
		log.Print(err.Error())
		return "", err
	}

	return id, nil
}

func (a *Authenticator) AddSubsystem(subsystemName string) (string, error) {
	query := "INSERT INTO subsystems " +
		"(subsystemName) " +
		"VALUES ($1) " +
		"RETURNING subsystemID"

	var id string
	err := a.database.SQL.QueryRow(query, subsystemName).Scan(&id)
	if err != nil {
		log.Print(err.Error())
		return "", err
	}

	return id, nil
}

func (a *Authenticator) ConnectAccountToGroup(accountName string, subsystemName string, groupName string) error {
	query := "INSERT INTO auth_accounts_groups " +
		"(accountID, groupID) " +
		"SELECT a.accountID, b.groupID " +
		"FROM auth_accounts a, auth_groups b " +
		"WHERE accountName = $1 AND subsystemName = $2 AND groupName = $3"
	_, err := a.database.SQL.Exec(query, accountName, subsystemName, groupName)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	return nil
}

func (a *Authenticator) ConnectGroupToPrivilege(subsystemName string, groupName string, privilegeName string) error {
	query := "INSERT INTO auth_groups_privileges " +
		"(groupID, privilegeID) " +
		"SELECT a.groupID, b.privilegeID " +
		"FROM auth_groups a, auth_privileges b " +
		"WHERE subsystemName = $2 groupName = $1 AND privilegeName = $3"
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
