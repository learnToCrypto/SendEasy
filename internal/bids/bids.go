package bid

import (
	"time"

	"github.com/learnToCrypto/lakoposlati/internal/platform/postgres"
)

type Bid struct {
	Id              int
	Uuid            string
	StartingBid     float64
	LowestBid       float64
	Currency        string
	ProviderId      int
	DemandId        int
	PickUpTime      string
	DropOffTime     string
	QuoteExpiration string
	Note            string
	Payment         string
	CreatedAt       time.Time
	Status          int // active, canceled
}

func CreateBid(demand_id, provider_id, status int, starting_bid float64, currency, pick_up_time, drop_off_time, quote_expiration, note, payment string) (bid Bid, err error) {
	statement := "insert into bids (uuid, starting_bid, currency, provider_id, demand_id, pick_up_time, drop_off_time, quote_expiration, note, payment, created_at, status) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) returning id,uuid, starting_bid, currency, provider_id, demand_id, pick_up_time, drop_off_time, quote_expiration, note, payment, created_at, status"
	stmt, err := postgres.Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
	err = stmt.QueryRow(postgres.CreateUUID(), starting_bid, currency, provider_id, demand_id, pick_up_time, drop_off_time, quote_expiration, note, payment, time.Now().UTC(), status).Scan(&bid.Id, &bid.Uuid, &bid.StartingBid, &bid.Currency, &bid.ProviderId, &bid.DemandId, &bid.PickUpTime, &bid.DropOffTime, &bid.QuoteExpiration, &bid.Note, &bid.Payment, &bid.CreatedAt, &bid.Status)
	return
}
