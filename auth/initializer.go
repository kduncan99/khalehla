//	authd.go - Kahlehla project
//	Copyright 2023 by Kurt Duncan
//	Code which initializes the db tables associated with the authenticator

package auth

import "log"

var dropCommands = []string{
	`DROP TABLE auth_accounts_groups;`,
	`DROP TABLE auth_accounts_privileges`,
	`DROP TABLE auth_accounts;`,
	`DROP TABLE auth_groups;`,
	`DROP TABLE auth_privileges;`}

var initCommands = []string{
	`CREATE TABLE auth_accounts (
		accountID uuid DEFAULT gen_random_uuid() PRIMARY KEY,
		loginName VARCHAR UNIQUE NOT NULL,
		passwordHash VARCHAR NOT NULL,
		disabled int NOT NULL,
		userName VARCHAR NOT NULL,
		emailAddress VARCHAR UNIQUE NOT NULL)`,

	`CREATE TABLE auth_groups (
		groupID uuid DEFAULT gen_random_uuid() PRIMARY KEY,
		parentGroupID uuid REFERENCES auth_groups(groupID) ON DELETE CASCADE,
		groupName VARCHAR UNIQUE NOT NULL)`,

	`CREATE TABLE auth_privileges (
		privilegeID uuid DEFAULT gen_random_uuid() PRIMARY KEY,
		privilegeName VARCHAR UNIQUE NOT NULL)`,

	`CREATE TABLE auth_accounts_groups (
    	accountID uuid REFERENCES auth_accounts ON DELETE CASCADE,
    	groupID   uuid REFERENCES auth_groups   ON DELETE CASCADE,
    	PRIMARY KEY (accountID, groupID))`,

    `CREATE TABLE auth_accounts_privileges (
		accountID   uuid REFERENCES auth_accounts   ON DELETE CASCADE,
		privilegeID uuid REFERENCES auth_privileges ON DELETE CASCADE,
		PRIMARY KEY (accountID, privilegeID))`}

func (a *Authenticator) Initialize() error {
	if a.database == nil {
		log.Fatal("Database is not open")
	}

	for _, cmd := range dropCommands {
		_, _ = a.database.SQL.Exec(cmd)
	}

	for _, cmd := range initCommands {
		_, err := a.database.SQL.Exec(cmd)
		if err != nil {
			log.Print(err.Error())
			return err
		}
	}

	gid, err := a.AddBaseGroup("administrators")
	if err != nil {
		return err
	}

	aid, err := a.AddAccount("admin", "admin", "administrator", "n/a")
	if err != nil {
		return err
	}

	err = a.ConnectAccountToGroup(aid, gid)
	if err != nil {
		return err
	}

	return nil
}
