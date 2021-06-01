package postgres

import (
	"errors"
	"fmt"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

const (
	table      = "hello"
	relevantID = "58db0d969fa44de4ba14631c07622761"
)

var errNotFound = errors.New("not found")

type Client struct {
	Session db.Session
	DSN     DSN
}

type entry struct {
	ID   string `json:"id" db:"id"`
	Hits int    `json:"hits" db:"hits"`
}

func (d DSN) ConnectionURL() *postgresql.ConnectionURL {
	return &postgresql.ConnectionURL{
		User:     d.Username,
		Password: d.Password,
		Host:     d.URI,
		Database: d.DatabaseName,
	}
}

// Write bumps the one and only row's hits attribute by 1
func (c *Client) Write() error {
	result := entry{ID: relevantID}

	current, err := c.Read()
	if err != nil {
		return fmt.Errorf("fetching current hits value: %w", err)
	}

	result.Hits = current + 1

	collection := c.Session.Collection(table)

	err = collection.UpdateReturning(&result)
	if err != nil {
		return fmt.Errorf("updating hits for %s: %w", relevantID, err)
	}

	return nil
}

// Read returns the hits attribute from the one and only row
func (c *Client) Read() (int, error) {
	collection := c.Session.Collection(table)

	condition := db.Cond{"id": relevantID}

	results := collection.Find(condition)

	exists, err := results.Exists()
	if err != nil {
		return -1, fmt.Errorf("checking existence: %w", err)
	}

	if !exists {
		return -1, fmt.Errorf("finding hits for %s: %w", relevantID, errNotFound)
	}

	entryResult := entry{ID: relevantID}

	err = results.One(&entryResult)
	if err != nil {
		return -1, fmt.Errorf("fetching result: %w", err)
	}

	return entryResult.Hits, nil
}

// Close closes a database connection
func (c *Client) Close() error {
	return c.Session.Close()
}

// Open opens a database connection
func (c *Client) Open() error {
	sess, err := postgresql.Open(c.DSN.ConnectionURL())
	if err != nil {
		return fmt.Errorf("connecting to Postgres: %w", err)
	}

	c.Session = sess

	return c.setup()
}

func (c *Client) setup() error {
	sql := c.Session.SQL()

	// Create table if it doesn't exist
	_, err := sql.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id text primary key,
		hits integer
	)`, table))
	if err != nil {
		return fmt.Errorf("creating tables: %w", err)
	}

	// Insert the one and only row if it doesn't exist
	if _, err := c.Read(); errors.Is(err, errNotFound) {
		collection := c.Session.Collection(table)

		initialEntry := entry{
			ID:   relevantID,
			Hits: 0,
		}

		_, err := collection.Insert(&initialEntry)
		if err != nil {
			return fmt.Errorf("creating initial entry: %w", err)
		}
	}

	return nil
}
