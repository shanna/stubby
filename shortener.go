package stubby

// TODO: LRU in memory cache & snapshot visit updates from LRU to take the load off the DB for all the hits I'll be getting.
// TODO: Onion JWT redirects so extra data can be included in redirect. Bundle meta-data into signed JWT redirect. Noice!
// TODO: Add a google tracking URL into all redirects for even more stats I'll never look at!
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

// Database storage interface.
type Store interface {
	Get(id string) (*Record, error)
	Create(url string, meta Meta) (*Record, error)
}

type ID interface {
	Encode(id string) string
	Decode(token string) string
}

type Shortener struct {
	Store Store
	ID    ID
}

func New(store Store, id ID) (*Shortener, error) {
	s := &Shortener{
		Store: store,
		ID:    id,
	}
	return s, nil
}

func (s Shortener) Get(token string) (*Record, error) {
	return s.Store.Get(s.ID.Decode(token))
}

func (s Shortener) Create(url string, meta Meta) (*Record, error) {
	record, err := s.Store.Create(url, meta)
	if err != nil {
		return nil, err
	}
	record.ID = s.ID.Encode(record.ID)
	return record, nil
}
