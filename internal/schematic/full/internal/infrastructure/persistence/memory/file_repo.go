package memory

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"{{.ModulePath}}/internal/domain/entity"
	"{{.ModulePath}}/internal/domain/repository"

	"github.com/linkeunid/ligo"
	ligomemory "github.com/linkeunid/ligo-memory"
)

// FileRepository is an in-memory implementation of repository.FileRepository.
// Metadata is stored in a ligo-memory Store; file content is written to disk.
type FileRepository struct {
	store *ligomemory.Store[int, *entity.File]
	dir   string
	log   ligo.Logger
}

// NewFileRepository creates a new in-memory file repository.
func NewFileRepository(dir string, store *ligomemory.Store[int, *entity.File], log ligo.Logger) repository.FileRepository {
	os.MkdirAll(dir, 0o755)
	return &FileRepository{store: store, dir: dir, log: log}
}

// EnsureUploadDir ensures the upload directory exists.
// This method is registered as a lifecycle hook via Register().
func (r *FileRepository) EnsureUploadDir() error {
	r.log.Info("Ensuring upload directory exists", ligo.LoggerField{Key: "dir", Value: r.dir})
	return os.MkdirAll(r.dir, 0o755)
}

// Register implements ligo.Registerable interface for compile-time safe hook registration.
func (r *FileRepository) Register(registry *ligo.HookRegistry) {
	registry.OnInit(r.EnsureUploadDir)
}

func (r *FileRepository) Save(file io.Reader, filename string) (*entity.File, error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	id := nextID()
	path := filepath.Join(r.dir, fmt.Sprintf("%d_%s", id, filename))

	if err := os.WriteFile(path, content, 0o644); err != nil {
		return nil, err
	}

	fileEntity := &entity.File{
		ID:          id,
		Name:        filename,
		ContentType: detectContentType(filename, content),
		Size:        int64(len(content)),
		Path:        path,
		CreatedAt:   time.Now(),
	}

	r.store.Set(id, fileEntity)
	return fileEntity, nil
}

func (r *FileRepository) FindByID(id int) (*entity.File, bool) {
	return r.store.Get(id)
}

func (r *FileRepository) GetContent(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

func (r *FileRepository) FindAll() []*entity.File {
	return r.store.All()
}

func (r *FileRepository) Delete(id int) error {
	file, found := r.store.Get(id)
	if !found {
		return fmt.Errorf("file not found")
	}
	if err := os.Remove(file.Path); err != nil {
		return err
	}
	r.store.Delete(id)
	return nil
}

func detectContentType(filename string, content []byte) string {
	ct := http.DetectContentType(content)
	if ct != "application/octet-stream" {
		return ct
	}
	switch filepath.Ext(filename) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".txt":
		return "text/plain"
	case ".json":
		return "application/json"
	default:
		return "application/octet-stream"
	}
}
