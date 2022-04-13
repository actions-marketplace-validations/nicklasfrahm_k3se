package k3se

// Options contains the configuration for an operation.
type Options struct {
	ConfigPath string
}

// Option applies a configuration option
// for the execution of an operation.
type Option func(options *Options) error

// Apply applies the option functions to the current set of options.
func (o *Options) Apply(options ...Option) (*Options, error) {
	for _, option := range options {
		if err := option(o); err != nil {
			return nil, err
		}
	}
	return o, nil
}

// GetDefaultOptions returns the default options
// for all operations of this library.
func GetDefaultOptions() *Options {
	return &Options{
		ConfigPath: Program + ".yml",
	}
}

// WithConfigPath overrides the default configuration path.
func WithConfigPath(configPath string) Option {
	return func(options *Options) error {
		options.ConfigPath = configPath
		return nil
	}
}