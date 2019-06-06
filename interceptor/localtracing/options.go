package localtracing

var (
	defaultOptions = &options{
		idGenerateFunc: genId,
	}
)

type IdGenerateFunc func(seeds ...interface{}) string

type options struct {
	idGenerateFunc IdGenerateFunc
}

func evaluateOptions(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}

	return optCopy
}

type Option func(*options)

func WithIdGenerateFunc(f IdGenerateFunc) Option {
	return func(o *options) {
		o.idGenerateFunc = f
	}
}
