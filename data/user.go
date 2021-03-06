package data

import (
	"code.google.com/p/go.crypto/bcrypt"
)

// Constants representing the various roles a user may possess
const (
	RoleGuest = iota
	RoleUser
	RoleAdmin
)

// User represents an user registered to wavepipe
type User struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"-"`
	RoleID      int    `db:"role_id" json:"roleId"`
	LastFMToken string `db:"lastfm_token" json:"-"`
}

// NewUser generates and saves a new user, while also hashing the input password
func NewUser(username string, password string, roleID int) (*User, error) {
	// Generate user
	user := new(User)
	user.Username = username
	user.RoleID = roleID

	// Hash input password
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	// Save user
	if err := user.Save(); err != nil {
		return nil, err
	}

	return user, nil
}

// CreateSession generates a new API session for this user
func (u User) CreateSession(client string) (*Session, error) {
	return NewSession(u.ID, u.Password, client)
}

// SetPassword hashes a password using bcrypt, and stores it in the User struct
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	if err != nil {
		return err
	}
	u.Password = string(hash)

	return nil
}

// Delete removes an existing User from the database
func (u *User) Delete() error {
	return DB.DeleteUser(u)
}

// Load pulls an existing User from the database
func (u *User) Load() error {
	return DB.LoadUser(u)
}

// Save creates a new User in the database
func (u *User) Save() error {
	return DB.SaveUser(u)
}

// Update updates an existing User in the database
func (u *User) Update() error {
	return DB.UpdateUser(u)
}
