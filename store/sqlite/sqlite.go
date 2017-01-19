package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	stubby "github.com/shanna/stubby"
)

type Conn struct {
	db *sql.DB
}

func Open(path string) (*Conn, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
    create table if not exists stubby_urls (
      id integer primary key autoincrement,
      url text not null,
      visits integer not null default 0,
      meta text
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
func (c *Conn) Create(url string, meta stubby.Meta) (*stubby.Record, error) {
	_, err := c.db.Exec(`insert into stubby_urls (url) values ($1)`, url)
	if err != nil {
		return nil, err
	}

	row := c.db.QueryRow(`select last_insert_rowid() from stubby_urls`)
	var id string
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return &stubby.Record{ID: id, URL: url, Visits: 0}, nil
}

func (c *Conn) Get(id string) (*stubby.Record, error) {
	c.db.Exec(`update stubby_urls set visits = visits + 1 where id = $1`, id)
	row := c.db.QueryRow(`select url, visits from stubby_urls where id = $1`, id)

	record := &stubby.Record{}
	if err := row.Scan(&record.ID, &record.URL, &record.Visits); err != nil {
		return nil, err
	}
	return record, nil
}
