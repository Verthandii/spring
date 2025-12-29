package iorm

type Config struct {
	Host           string
	Port           int
	UserName       string
	Password       string
	DBName         string
	LoggerLevel    int
	Idle           int
	Open           int
	IdleTime       int
	DisableMetrics bool
	LowerCase      bool
}
