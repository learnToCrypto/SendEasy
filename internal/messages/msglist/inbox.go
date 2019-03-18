package msglist

import (
	"fmt"

	"github.com/learnToCrypto/lakoposlati/internal/demands"
	"github.com/learnToCrypto/lakoposlati/internal/messages"
	"github.com/learnToCrypto/lakoposlati/internal/platform/postgres"
	"github.com/learnToCrypto/lakoposlati/internal/user"
)

type Inbox struct {
	Owner    user.User
	Messages []messages.Message
}

// get messages to a demand
func Messages(demand *demands.Demand) (msgs []messages.Message, err error) {
	rows, err := postgres.Db.Query("SELECT id, uuid, body, user_id, demand_id, created_at FROM messages where demand_id = $1", demand.Id)
	fmt.Println(rows, "err = ", err)
	if err != nil {
		return
	}
	for rows.Next() {
		msg := messages.Message{}
		if err = rows.Scan(&msg.Id, &msg.Uuid, &msg.Body, &msg.UserId, &msg.DemandId, &msg.CreatedAt); err != nil {
			return
		}
		msgs = append(msgs, msg)
	}
	rows.Close()
	return
}
