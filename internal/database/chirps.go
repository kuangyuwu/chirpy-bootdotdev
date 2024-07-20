package database

import "errors"

type Chirp struct {
	Id       int    `json:"id"`
	AuthorId int    `json:"author_id"`
	Body     string `json:"body"`
}

var ErrChirpNotFound = errors.New("chirp not found")
var ErrUnauthorized = errors.New("user not authorized to perform the action")

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, authorId int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	newChirp := Chirp{
		Id:       id,
		AuthorId: authorId,
		Body:     body,
	}
	dbStructure.Chirps[id] = newChirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	n := len(dbStructure.Chirps)
	chirps := make([]Chirp, len(dbStructure.Chirps))
	for i := 0; i < n; i++ {
		chirps[i] = dbStructure.Chirps[i+1]
	}

	return chirps, nil
}

func (db *DB) DeleteChirp(userId, chirpId int) error {

	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	chirp, ok := dbStructure.Chirps[chirpId]
	if !ok {
		return ErrChirpNotFound
	}

	if chirp.AuthorId != userId {
		return ErrUnauthorized
	}

	delete(dbStructure.Chirps, chirpId)

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}
