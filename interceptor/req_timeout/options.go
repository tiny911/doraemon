package req_timeout

var (
	defaultOptions = &options{
		timeout: 3000, //缺省超时等待是3000ms
	}
)

type options struct {
	timeout int
}

func evaluateOptions(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}

	return optCopy
}

// Option 类型
type Option func(*options)

// WithTimeout 设置超时
func WithTimeout(timeout int) Option {
	return func(o *options) {
		o.timeout = timeout
	}
}
