package database

import "errors"

type User struct {
	Id     int    `json:"id"`
	Email  string `json:"email"`
	Hashed string `json:"hashed"`
}

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUser(email, hashed string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	_, err = db.GetUserByEmail(email)
	if err == nil {
		return User{}, errors.New("user already exists")
	}

	id := len(dbStructure.Users) + 1
	newUser := User{
		Id:     id,
		Email:  email,
		Hashed: hashed,
	}
	dbStructure.Users[id] = newUser

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, errors.New("user not found")
}
