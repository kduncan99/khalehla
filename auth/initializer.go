//	authd.go - Kahlehla project
//	Copyright 2023 by Kurt Duncan
//	Code which initializes the db tables associated with the authenticator

package auth

import (
	"fmt"
	"log"
)

const AdminAccountName = "admin"
const AdminGroupName = "administrators"
const AddUserPrivilege = "AUTH:ADD_USER"
const RemoveUserPrivilege = "AUTH:REMOVE_USER"
const UpdateUserPrivilege = "AUTH:UPDATE_USER"
const ViewUserPrivilege = "AUTH:VIEW_USER"

var dropCommands = []string{
	`DROP TABLE auth_accounts_groups;`,
	`DROP TABLE auth_groups_privileges`,
	`DROP TABLE auth_accounts;`,
	`DROP TABLE auth_groups;`,
	`DROP TABLE auth_privileges;`}

var initCommands = []string{
	`CREATE TABLE auth_accounts (
		accountID uuid DEFAULT gen_random_uuid() PRIMARY KEY,
		accountName VARCHAR UNIQUE NOT NULL,
		passwordHash VARCHAR NOT NULL,
		disabled int NOT NULL,
		userName VARCHAR NOT NULL,
		emailAddress VARCHAR UNIQUE NOT NULL)`,

	`CREATE TABLE auth_groups (
		groupID uuid DEFAULT gen_random_uuid() PRIMARY KEY,
		groupName VARCHAR UNIQUE NOT NULL)`,

	`CREATE TABLE auth_privileges (
		privilegeID uuid DEFAULT gen_random_uuid() PRIMARY KEY,
		privilegeName VARCHAR UNIQUE NOT NULL)`,

	`CREATE TABLE auth_accounts_groups (
    	accountID uuid REFERENCES auth_accounts ON DELETE CASCADE,
    	groupID   uuid REFERENCES auth_groups   ON DELETE CASCADE,
    	PRIMARY KEY (accountID, groupID))`,

    `CREATE TABLE auth_groups_privileges (
		groupID     uuid REFERENCES auth_groups     ON DELETE CASCADE,
		privilegeID uuid REFERENCES auth_privileges ON DELETE CASCADE,
		PRIMARY KEY (groupID, privilegeID))`}

func (a *Authenticator) Initialize() error {
	if a.database == nil {
		log.Fatal("Database is not open")
	}

	for _, cmd := range dropCommands {
		fmt.Println(cmd)
		_, err := a.database.SQL.Exec(cmd)
		if err != nil {
			fmt.Println("***" + err.Error())
		}
	}

	for _, cmd := range initCommands {
		fmt.Println(cmd)
		_, err := a.database.SQL.Exec(cmd)
		if err != nil {
			log.Print(err.Error())
			return err
		}
	}

	_, _ = a.AddPrivilege(AddUserPrivilege)
	_, _ = a.AddPrivilege(RemoveUserPrivilege)
	_, _ = a.AddPrivilege(UpdateUserPrivilege)
	_, _ = a.AddPrivilege(ViewUserPrivilege)
	_, _ = a.AddGroup(AdminGroupName)
	_, _ = a.AddAccount(AdminAccountName, "admin", "administrator", "n/a")

	_ = a.ConnectAccountToGroup(AdminAccountName, AdminGroupName)
	_ = a.ConnectGroupToPrivilege(AdminGroupName, AddUserPrivilege)
	_ = a.ConnectGroupToPrivilege(AdminGroupName, RemoveUserPrivilege)
	_ = a.ConnectGroupToPrivilege(AdminGroupName, UpdateUserPrivilege)
	_ = a.ConnectGroupToPrivilege(AdminGroupName, ViewUserPrivilege)

	return nil
}
