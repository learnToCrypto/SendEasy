package demands

import (
	"strconv"
	"time"

	"github.com/learnToCrypto/lakoposlati/internal/platform/postgres"
	"github.com/learnToCrypto/lakoposlati/internal/user"
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

// Create a new demand
func CreateDemand(user *user.User, object, collection, delivery, timeframe string, status int) (conv Demand, err error) {
	statement := "insert into demands (uuid, object, collection, delivery, timeframe, user_id, created_at, status) values ($1, $2, $3, $4, $5, $6, $7, $8) returning id, uuid, object, collection, delivery, timeframe, user_id, created_at, status"
	stmt, err := postgres.Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
	err = stmt.QueryRow(postgres.CreateUUID(), object, collection, delivery, timeframe, user.Id, time.Now().UTC(), status).Scan(&conv.Id, &conv.Uuid, &conv.Object, &conv.Collection, &conv.Delivery, &conv.Timeframe, &conv.UserId, &conv.CreatedAt, &conv.Status)
	return
}

// Get the user who made the demand
func (demand *Demand) User() (userI user.User) {
	userI = user.User{}
	postgres.Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", demand.UserId).
		Scan(&userI.Id, &userI.Uuid, &userI.Name, &userI.Email, &userI.CreatedAt)
	return
}

// CreatedAtDate
func (demand *Demand) CreatedAtDate() string {
	return demand.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
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

// Get all demands in the database and returns it
func Demands(limit int, offset int) (demands []Demand, err error) {
	//var stm string
	//stm =
	//	fmt.Println(stm)
	rows, err := postgres.Db.Query("SELECT id, uuid, object, collection, delivery, timeframe, user_id, created_at, status FROM demands ORDER BY created_at DESC LIMIT " + strconv.Itoa(limit) + "OFFSET " + strconv.Itoa(offset))
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

func DemandsNum() (count int, err error) {
	rows, err := postgres.Db.Query("SELECT count(*) FROM demands")
	if err != nil {
		return
	}
	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return
		}
	}
	rows.Close()
	return count, err
}

// Get a demand by the UUID
func DemandByUUID(uuid string) (conv Demand, err error) {
	conv = Demand{}
	err = postgres.Db.QueryRow("SELECT id, uuid, object, collection, delivery, timeframe, user_id, created_at, status FROM demands WHERE uuid = $1", uuid).
		Scan(&conv.Id, &conv.Uuid, &conv.Object, &conv.Collection, &conv.Delivery, &conv.Timeframe, &conv.UserId, &conv.CreatedAt, &conv.Status)
	return
}
