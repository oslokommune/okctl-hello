package postgres

import (
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

var dsnRegex = regexp.MustCompile(`(\w+)://(\w+):(.+)@([\w-\.]+):(\d+)/(\w+)`)

type DSN struct {
	Scheme       string
	Username     string
	Password     string
	URI          string
	Port         string
	DatabaseName string
}

// ParseDSN interprets a Data Source Name in the following format
// databaseType://username:password@host:port/databaseName
func ParseDSN(rawDSN string) DSN {
	matches := dsnRegex.FindStringSubmatch(rawDSN)

	if len(matches) != 7 {
		return DSN{}
	}

	return DSN{
		Scheme:       matches[1],
		Username:     matches[2],
		Password:     matches[3],
		URI:          matches[4],
		Port:         matches[5],
		DatabaseName: matches[6],
	}
}

// Validate ensures the DSN has the necessary information
func (d DSN) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Scheme, validation.Required),
		validation.Field(&d.Username, validation.Required),
		validation.Field(&d.Password, validation.Required),
		validation.Field(&d.URI, validation.Required, is.DNSName),
		validation.Field(&d.Port, validation.Required, is.UTFNumeric),
		validation.Field(&d.DatabaseName, validation.Required),
	)
}

// String returns the DSN as a string
func (d DSN) String() string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s",
		d.Scheme,
		d.Username,
		d.Password,
		d.URI,
		d.Port,
		d.DatabaseName,
	)
}
