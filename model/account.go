package model

import (
	"time"
)

// Account - table
type Account struct {
	ID       string    `db:"id"`
	Account  string    `db:"account"`
	Password string    `db:"password"`
	Ctime    time.Time `db:"ctime"`
	Mtime    time.Time `db:"mtime"`
}

// GetAllAccount -
func GetAllAccount() (accs []Account, err error) {
	err = x.Select(&accs, "select * from public.account order by ctime asc")
	if err != nil {
		return nil, err
	}
	return
}

// GetAccount -
func GetAccount(account string) (acc *Account, err error) {
	err = x.Select(acc, `select * from "public"."account" where "account" = $1`, account)
	if err != nil {
		return nil, err
	}
	return
}
