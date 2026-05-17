//go:build wireinject

package simple

type Config struct{ Addr string }
type Logger struct{ Prefix string }
type Database struct {
	Cfg *Config
	Log *Logger
}
type Service struct {
	DB  *Database
	Log *Logger
}
type App struct{ Svc *Service }

func NewConfig() *Config         { return &Config{Addr: ":8080"} }
func NewLogger() *Logger         { return &Logger{Prefix: "app"} }
func NewDatabase(cfg *Config, log *Logger) *Database {
	return &Database{Cfg: cfg, Log: log}
}
func NewService(db *Database, log *Logger) (*Service, error) {
	return &Service{DB: db, Log: log}, nil
}
func NewApp(svc *Service) *App { return &App{Svc: svc} }

func Build() (*App, error) {
	return nil, wire(
		NewConfig,
		NewLogger,
		NewDatabase,
		NewService,
		NewApp,
	)
}

// wire is the user-defined stub. ligo-cli replaces Build's body in the
// generated counterpart; this function is never called at runtime.
func wire(_ ...any) error { return nil }
