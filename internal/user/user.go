package user

import (
	"time"

	"github.com/learnToCrypto/lakoposlati/internal/platform/postgres"
)

//type User struct {
//Id        int
//Uuid      string
//Name      string
//Email     string
//Password  string
//CreatedAt time.Time
//}

// User respresents a registered user with email/password authentication  (see netlify/gotrue)
type User struct {
	Id   int    `json:"id" db:"id"`
	Uuid string `json:"Uuid" db:"uuid"`

	Name string
	//NickName          string          username
	FirstName string
	LastName  string

	Aud  string `json:"aud" db:"aud"`
	Role string `json:"role" db:"role"`

	Email       string     `json:"email" db:"email"`
	Password    string     `json:"-" db:"encrypted_password"`
	ConfirmedAt *time.Time `json:"confirmed_at,omitempty" db:"confirmed_at"`
	InvitedAt   *time.Time `json:"invited_at,omitempty" db:"invited_at"`

	ConfirmationToken  string     `json:"-" db:"confirmation_token"`
	ConfirmationSentAt *time.Time `json:"confirmation_sent_at,omitempty" db:"confirmation_sent_at"`

	RecoveryToken  string     `json:"-" db:"recovery_token"`
	RecoverySentAt *time.Time `json:"recovery_sent_at,omitempty" db:"recovery_sent_at"`

	EmailChangeToken  string     `json:"-" db:"email_change_token"`
	EmailChange       string     `json:"new_email,omitempty" db:"email_change"`
	EmailChangeSentAt *time.Time `json:"email_change_sent_at,omitempty" db:"email_change_sent_at"`

	LastSignInAt *time.Time `json:"last_sign_in_at,omitempty" db:"last_sign_in_at"`

	//used to store information (e.g., a user's support plan, security roles, or access control groups)
	//that can impact a user's core functionality, such as how an application functions or what the user can access.
	AppMetaData map[string]interface{} `json:"app_metadata" db:"raw_app_meta_data"`
	//UserMetaData (RawData) used to store user attributes (e.g., user preferences) that do not impact a user's core functionality;
	UserMetaData map[string]interface{} `json:"user_metadata" db:"raw_user_meta_data"`

	IsAdmin bool `json:"-" db:"is_super_admin"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Create a new user, save user info into the database
func (user *User) Create() (err error) {
	// Postgres does not automatically return the last insert id, because it would be wrong to assume
	// you're always using a sequence.You need to use the RETURNING keyword in your insert to get this
	// information from postgres.
	statement := "insert into users (uuid, name, email, password, created_at) values ($1, $2, $3, $4, $5) returning id, uuid, created_at"
	stmt, err := postgres.Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	// use QueryRow to return a row and scan the returned id into the User struct
	err = stmt.QueryRow(postgres.CreateUUID(), user.Name, user.Email, postgres.Encrypt(user.Password), time.Now().UTC()).Scan(&user.Id, &user.Uuid, &user.CreatedAt)
	return
}

// Delete user from database
func (user *User) Delete() (err error) {
	statement := "delete from users where id = $1"
	stmt, err := postgres.Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id)
	return
}

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
