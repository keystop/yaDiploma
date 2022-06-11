package models

import "context"

type ServerDB interface {
	CheckDBConnection(context.Context)
	NewDBUserRepo() UsersRepo
	NewDBOrdersRepo() OrdersRepo
	NewDBBalanceRepo() BalanceRepo
	Close()
}
