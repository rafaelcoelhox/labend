package users

import "github.com/rafaelcoelhox/labbend/pkg/database"

// init - registra automaticamente os modelos do módulo users
func init() {
	database.RegisterModel(&User{})
	database.RegisterModel(&UserXP{})
}
