package dashboard

import (
	"github.com/learnToCrypto/lakoposlati/internal/messages"
	"github.com/learnToCrypto/lakoposlati/internal/user"
)

type DashboardUser struct {
	ActiveDemands    []user.Demand
	CompletedDemands []user.Demand
	Inbox            []messages.Message
}
