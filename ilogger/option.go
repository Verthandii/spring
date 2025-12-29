package ilogger

type Option interface {
	apply(logger Logger)
}

type option func(logger Logger)

func (fn option) apply(logger Logger) {
	fn(logger)
}

func WithValuer(key string, v Valuer) Option {
	return option(func(logger Logger) {
		logger.addValuer(key, v)
	})
}
