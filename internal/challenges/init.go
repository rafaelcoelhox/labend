package challenges

import "github.com/rafaelcoelhox/labbend/pkg/database"

// init - registra automaticamente os modelos do módulo challenges
func init() {
	database.RegisterModel(&Challenge{})
	database.RegisterModel(&ChallengeVote{})
}
