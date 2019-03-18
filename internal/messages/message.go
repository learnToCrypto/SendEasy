package messages

import (
	"time"

	"github.com/learnToCrypto/lakoposlati/internal/demands"
	"github.com/learnToCrypto/lakoposlati/internal/platform/postgres"
	"github.com/learnToCrypto/lakoposlati/internal/user"
)

type Message struct {
	Id        int
	Uuid      string
	Body      string
	UserId    int
	DemandId  int
	CreatedAt time.Time
}

// Create a new message to a demand
func CreateMessage(user *user.User, conv demands.Demand, body string) (msg Message, err error) {
	statement := "insert into messages (uuid, body, user_id, demand_id, created_at) values ($1, $2, $3, $4, $5) returning id, uuid, body, user_id, demand_id, created_at"
	stmt, err := postgres.Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
	err = stmt.QueryRow(postgres.CreateUUID(), body, user.Id, conv.Id, time.Now().UTC()).Scan(&msg.Id, &msg.Uuid, &msg.Body, &msg.UserId, &msg.DemandId, &msg.CreatedAt)
	return
}

// Get the user who wrote the message
func (msg *Message) User() (userI user.User) {
	userI = user.User{}
	postgres.Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", msg.UserId).
		Scan(&userI.Id, &userI.Uuid, &userI.Name, &userI.Email, &userI.CreatedAt)
	return
}

func (message *Message) CreatedAtDate() string {
	return message.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}
