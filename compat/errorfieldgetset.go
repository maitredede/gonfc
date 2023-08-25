package compat

type ErrorFieldGetSet struct {
	getter func() error
	setter func(error)
}

func (i ErrorFieldGetSet) Get() error {
	return i.getter()
}

func (i ErrorFieldGetSet) Set(value error) {
	i.setter(value)
}

func NewErrorFieldGetSet(getter func() error, setter func(error)) ErrorFieldGetSet {
	if getter == nil {
		panic("getter is nil")
	}
	if setter == nil {
		panic("setter is nil")
	}
	return ErrorFieldGetSet{
		getter: getter,
		setter: setter,
	}
}
