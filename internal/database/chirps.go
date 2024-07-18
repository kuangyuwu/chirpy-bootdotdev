package database

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	newChirp := Chirp{
		Id:   id,
		Body: body,
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
