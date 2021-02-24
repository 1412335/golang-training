package functionaloptions

type Celsius int

type Foobar struct {
	mutable     bool
	temperature Celsius
}

type OptionFoobar func(*Foobar) error

func NewFoobar(options ...OptionFoobar) (*Foobar, error) {
	fb := &Foobar{}
	// Default values...
	fb.mutable = true
	fb.temperature = 37
	// options
	for _, opt := range options {
		if err := opt(fb); err != nil {
			return nil, err
		}
	}
	return fb, nil
}

func OptionReadOnlyFlag(fb *Foobar) error {
	fb.mutable = false
	return nil
}

func OptionTemperature(t Celsius) OptionFoobar {
	return func(fb *Foobar) error {
		fb.temperature = t
		return nil
	}
}
