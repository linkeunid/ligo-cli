package simple

type (
	Config   struct{ Addr string }
	Logger   struct{ Prefix string }
	Database struct {
		Cfg *Config
		Log *Logger
	}
	Service struct {
		DB  *Database
		Log *Logger
	}
	App struct{ Svc *Service }
)

func NewConfig() *Config { return &Config{Addr: ":8080"} }
func NewLogger() *Logger { return &Logger{Prefix: "app"} }
func NewDatabase(cfg *Config, log *Logger) *Database {
	return &Database{Cfg: cfg, Log: log}
}
func NewService(db *Database, log *Logger) (*Service, error) {
	return &Service{DB: db, Log: log}, nil
}
func NewApp(svc *Service) *App { return &App{Svc: svc} }
