package data

import (
	"time"

	"github.com/6ill/greenlight/internal/validator"
)

type Movie struct {
	ID        int64		`json:"id"`     	// Unique integer ID for the movie
	CreatedAt time.Time `json:"-"`	// Timestamp for when the movie is added to our database
	Title     string    `json:"title"`		// Movie title
	Year      int32     `json:"year,omitempty"`		// Movie release year
	Runtime   Runtime     `json:"runtime,omitempty,string"`	// Movie runtime (in minutes)
	Genres    []string  `json:"genres,omitempty"`		// Slice of genres for the movie (romance, comedy, etc.)
	Version   int32     `json:"version"`	// The version number starts at 1 and will be incremented each
						// time the movie information is updated
}

func ValidateMovie(v *validator.Validator, m *Movie) {
	v.Check(m.Title != "", "title", "must be provided")
	v.Check(len(m.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(m.Year != 0, "year", "must be provided")
	v.Check(m.Year >= 1888, "year", "must be greater than 1888")
	v.Check(m.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(m.Runtime != 0, "runtime", "must be provided")
	v.Check(m.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(m.Genres != nil, "genres", "must be provided")
	v.Check(len(m.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(m.Genres) <= 5, "genres", "must not contain more than 5 genres")
	// Note that we're using the Unique helper in the line below to check that all
	// values in the m.Genres slice are unique.
	v.Check(validator.Unique(m.Genres), "genres", "must not contain duplicate values")
}