package source

type Source interface {
	Load() ([]byte, error)
}
