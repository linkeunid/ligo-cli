package entity

import "time"

// File represents a file entity in the domain layer.
type File struct {
	ID          string
	Name        string
	ContentType string
	Size        int64
	Path        string
	CreatedAt   time.Time
}
