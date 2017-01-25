package stubby

// TODO: LRU in memory cache & snapshot visit updates from LRU to take the load off the DB for all the hits I'll be getting.
// TODO: Tags? Increment all tag counts for a visited URL.

// Record meta data.
type Record struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Visits int    `json:"visits"`
	Meta   Meta   `json:"meta"`
}

// Meta data for a record.
type Meta interface{}

// Store database storage interface.
type Store interface {
	Get(id string) (*Record, error)
	Create(url string, meta Meta) (*Record, error)
}

// ID encoding/decoding interface.
type ID interface {
	Encode(id string) (string, error)
	Decode(token string) (string, error)
}

// Shortener service.
type Shortener struct {
	Store Store
	ID    ID
}

// New shortener service.
func New(store Store, id ID) (*Shortener, error) {
	s := &Shortener{
		Store: store,
		ID:    id,
	}
	return s, nil
}

// Get a URL record from the store by token.
func (s Shortener) Get(token string) (*Record, error) {
	id, err := s.ID.Decode(token)
	if err != nil {
		return nil, err
	}
	return s.Store.Get(id)
}

// Create a URL record in the store.
func (s Shortener) Create(url string, meta Meta) (*Record, error) {
	record, err := s.Store.Create(url, meta)
	if err != nil {
		return nil, err
	}
	record.ID, err = s.ID.Encode(record.ID)
	if err != nil {
		return nil, err
	}
	return record, nil
}
