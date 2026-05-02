package repository

import (
	"io"

	"{{.ModulePath}}/internal/domain/entity"
)

// FileRepository defines the interface for file data access.
type FileRepository interface {
	// Save stores a file and returns the created file entity.
	Save(file io.Reader, filename string) (*entity.File, error)

	// FindByID retrieves file metadata by ID.
	FindByID(id int) (*entity.File, bool)

	// GetContent returns a reader for the file content.
	GetContent(path string) (io.ReadCloser, error)

	// FindAll returns all file metadata.
	FindAll() []*entity.File

	// Delete removes a file by ID.
	Delete(id int) error
}
