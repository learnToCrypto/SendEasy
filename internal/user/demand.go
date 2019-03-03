package user

import (
	"fmt"
	"time"

	"github.com/learnToCrypto/lakoposlati/internal/platform/postgres"
)

type Demand struct {
	Id         int
	Uuid       string
	Object     string
	Collection string
	Delivery   string
	Timeframe  string
	UserId     int
	CreatedAt  time.Time
	Status     int // if 0 active, 1 success, 2 cancelled

	Length     string //m
	Width      string
	Height     string
	Weight     string //kg
	Desciption string
	DateListed time.Time
	Expires    time.Time
}

type Message struct {
	Id        int
	Uuid      string
	Body      string
	UserId    int
	DemandId  int
	CreatedAt time.Time
}

// Create a new demand
func (user *User) CreateDemand(object, collection, delivery, timeframe string, status int) (conv Demand, err error) {
	statement := "insert into demands (uuid, object, collection, delivery, timeframe, user_id, created_at, status) values ($1, $2, $3, $4, $5, $6, $7, $8) returning id, uuid, object, collection, delivery, timeframe, user_id, created_at, status"
	stmt, err := postgres.Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
	err = stmt.QueryRow(postgres.CreateUUID(), object, collection, delivery, timeframe, user.Id, time.Now(), status).Scan(&conv.Id, &conv.Uuid, &conv.Object, &conv.Collection, &conv.Delivery, &conv.Timeframe, &conv.UserId, &conv.CreatedAt, &conv.Status)
	return
}

// Get the user who made the demand
func (demand *Demand) User() (user User) {
	user = User{}
	postgres.Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", demand.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}

// Create a new message to a demand
func (user *User) CreateMessage(conv Demand, body string) (msg Message, err error) {
	statement := "insert into messages (uuid, body, user_id, demand_id, created_at) values ($1, $2, $3, $4, $5) returning id, uuid, body, user_id, demand_id, created_at"
	stmt, err := postgres.Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
	err = stmt.QueryRow(postgres.CreateUUID(), body, user.Id, conv.Id, time.Now()).Scan(&msg.Id, &msg.Uuid, &msg.Body, &msg.UserId, &msg.DemandId, &msg.CreatedAt)
	return
}

// Get the user who wrote the message
func (msg *Message) User() (user User) {
	user = User{}
	postgres.Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", msg.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}

// CreatedAtDate
func (demand *Demand) CreatedAtDate() string {
	return demand.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

func (message *Message) CreatedAtDate() string {
	return message.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

// get the number of messages for a demand
func (demand *Demand) NumReplies() (count int) {
	rows, err := postgres.Db.Query("SELECT count(*) FROM messages where demand_id = $1", demand.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return
		}
	}
	rows.Close()
	return
}

// get messages to a demand
func (demand *Demand) Messages() (messages []Message, err error) {
	rows, err := postgres.Db.Query("SELECT id, uuid, body, user_id, demand_id, created_at FROM messages where demand_id = $1", demand.Id)
	fmt.Println(rows, "err = ", err)
	if err != nil {
		return
	}
	for rows.Next() {
		msg := Message{}
		if err = rows.Scan(&msg.Id, &msg.Uuid, &msg.Body, &msg.UserId, &msg.DemandId, &msg.CreatedAt); err != nil {
			return
		}
		messages = append(messages, msg)
	}
	rows.Close()
	return
}

// Get all demands in the database and returns it
func Demands() (demands []Demand, err error) {
	rows, err := postgres.Db.Query("SELECT id, uuid, object, collection, delivery, timeframe, user_id, created_at, status FROM demands ORDER BY created_at DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		conv := Demand{}
		if err = rows.Scan(&conv.Id, &conv.Uuid, &conv.Object, &conv.Collection, &conv.Delivery, &conv.Timeframe, &conv.UserId, &conv.CreatedAt, &conv.Status); err != nil {
			return
		}
		demands = append(demands, conv)
	}
	rows.Close()
	return
}

// Get a demand by the UUID
func DemandByUUID(uuid string) (conv Demand, err error) {
	conv = Demand{}
	err = postgres.Db.QueryRow("SELECT id, uuid, object, collection, delivery, timeframe, user_id, created_at, status FROM demands WHERE uuid = $1", uuid).
		Scan(&conv.Id, &conv.Uuid, &conv.Object, &conv.Collection, &conv.Delivery, &conv.Timeframe, &conv.UserId, &conv.CreatedAt, &conv.Status)
	return
}
