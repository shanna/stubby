package stubby

// TODO: Configurable storage driver interface so you can use postgres, bolt, whatever to Open.
// TODO: Configurable encoder interface passed to Open.
// TODO: LRU cache & snapshot visit updates from LRU?
// TODO: Identical URLs should use old ID if URL and metadata is identical?
// TODO: Replace bjf implementation? JWT style stalted + signed but reversable ID.
// TODO: Onion JWT redirects so extra data can be included in redirect. Bundle meta-data into signed JWT redirect. Noice!
// TODO: Add a google tracking URL into all redirects for stats.
// TODO: Tags? Increment all tag counts for a visited URL.
// TODO: Include google tracking in URL object?

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xor-gate/bjf"
	"github.com/y0ssar1an/q"
)

type URL struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Visits uint32 `json:"visits"`
	// TODO: Metadata.
}

type Conn struct {
	db *sql.DB
}

func Open(connection string) (*Conn, error) {
	db, err := sql.Open("sqlite3", connection)
	if err != nil {
		return nil, err
	}

	// TODO: Migrations.
	_, err = db.Exec(`
    create table if not exists stubby_urls (
      id integer primary key autoincrement,
      url text not null,
      visits integer not null default 0
    );
  `)
	if err != nil {
		return nil, err
	}

	return &Conn{db}, nil
}

func (c *Conn) Close() error {
	return c.db.Close()
}

// Create a new short URL.
func (c *Conn) Create(original string) (*URL, error) {
	_, err := c.db.Exec(`insert into stubby_urls (url) values ($1)`, original)
	if err != nil {
		return nil, err
	}

	row := c.db.QueryRow(`select last_insert_rowid() from stubby_urls`)
	var id string
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return &URL{ID: bjf.Encode(id), URL: original, Visits: 0}, nil
}

func (c *Conn) Get(key string) (*URL, error) {
	id := bjf.Decode(key)
	c.db.Exec(`update stubby_urls set visits = visits + 1 where id = $1`, id)
	row := c.db.QueryRow(`select url, visits from stubby_urls where id = $1`, id)

	u := &URL{ID: key}
	if err := row.Scan(&u.URL, &u.Visits); err != nil {
		q.Q(err)
		return nil, err
	}
	return u, nil
}
