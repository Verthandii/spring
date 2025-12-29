package cors

type Option func(c *CORS)

func WithOrigin(origin []string) Option {
	return func(c *CORS) {
		c.origin = origin
	}
}

func WithCredentials(credentials string) Option {
	return func(c *CORS) {
		c.credentials = credentials
	}
}

func WithHeaders(headers string) Option {
	return func(c *CORS) {
		c.headers = headers
	}
}

func WithMethods(methods string) Option {
	return func(c *CORS) {
		c.methods = methods
	}
}

func WithOriginReflection(originReflection bool) Option {
	return func(c *CORS) {
		c.originReflection = originReflection
	}
}

type CORS struct {
	origin           []string
	credentials      string
	headers          string
	methods          string
	originReflection bool
}

func New(opts ...Option) *CORS {
	c := &CORS{
		origin:           []string{"*"},
		credentials:      "true",
		headers:          "*",
		methods:          "*",
		originReflection: false,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}
