package provider

import (
	"github.com/learnToCrypto/lakoposlati/internal/platform/postgres"
)

//Provider respresents a registered user with email/password authentication  (see netlify/gotrue)
type Provider struct {
	Id     int    `json:"id" db:"id"`
	Uuid   string `json:"uuid" db:"uuid"`
	UserId int    `json:"user_id" db:"user_id"`

	MobilePhone string `json:"mobile_phone" db:"mobile_phone"`

	CompanyName    string `json:"company_name" db:"company_name"`
	CompanyAddr    string `json:"company_addr" db:"company_addr"`
	CompanyCity    string `json:"company_city" db:"company_city"`
	CompanyZip     string `json:"company_zip" db:"company_zip"`
	CompanyCountry string `json:"company_country" db:"company_country"`

	Equipment     map[string]bool `json:"equipment" db:"equipment"`
	EligibleItems map[string]bool `json:"eligible_items" db:"eligible_items"`

	OperatingCountries []string

	//	Licence
}

// Create a new user, save user info into the database
func (provider *Provider) Create() (err error) {
	// Postgres does not automatically return the last insert id, because it would be wrong to assume
	// you're always using a sequence.You need to use the RETURNING keyword in your insert to get this
	// information from postgres.

	//	username      varchar(255),
	//	 licence bytea
	statement1 := "insert into providers (uuid, user_id, mobile_phone, company_name, company_addr, company_city, company_zip, company_country) values ($1, $2, $3, $4, $5, $6, $7, $8) returning id, uuid"
	stmt1, err := postgres.Db.Prepare(statement1)
	if err != nil {
		return
	}
	defer stmt1.Close()

	// use QueryRow to return a row and scan the returned id into the User struct
	err = stmt1.QueryRow(postgres.CreateUUID(), provider.UserId, provider.MobilePhone, provider.CompanyName, provider.CompanyAddr, provider.CompanyCity, provider.CompanyZip, provider.CompanyCountry).Scan(&provider.Id, &provider.Uuid)
	return
}

// Delete user from database
func (provider *Provider) Delete() (err error) {
	statement := "delete from users where id = $1"
	stmt, err := postgres.Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(provider.Id)

	statement1 := "delete from providers where provider_id = $1"
	stmt1, err := postgres.Db.Prepare(statement1)
	if err != nil {
		return
	}
	defer stmt1.Close()

	_, err = stmt1.Exec(provider.Id)
	return
}

/*

// Update user information in the database
func (user *User) Update() (err error) {
	statement := "update users set name = $2, email = $3 where id = $1"
	stmt, err := postgres.Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id, user.Name, user.Email)
	return
}

// Delete all users from database
func UserDeleteAll() (err error) {
	statement := "delete from users"
	_, err = postgres.Db.Exec(statement)
	return
}

// Get all users in the database and returns it
func Users() (users []User, err error) {
	rows, err := postgres.Db.Query("SELECT id, uuid, name, email, password, created_at FROM users")
	if err != nil {
		return
	}
	for rows.Next() {
		user := User{}
		if err = rows.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt); err != nil {
			return
		}
		users = append(users, user)
	}
	rows.Close()
	return
}

func UsernamebySession(userId int) (username string, err error) {
	username = ""
	err = postgres.Db.QueryRow("SELECT name FROM users WHERE id = $1", userId).Scan(&username)
	return

}

// Get a single user given the email
func UserByEmail(email string) (user User, err error) {
	user = User{}
	err = postgres.Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE email = $1", email).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}

// Get a single user given the UUID
func UserByUUID(uuid string) (user User, err error) {
	user = User{}
	err = postgres.Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE uuid = $1", uuid).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}
*/
