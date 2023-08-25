package compat

type BoolFieldGetSet struct {
	getter func() bool
	setter func(bool)
}

func (i BoolFieldGetSet) Get() bool {
	return i.getter()
}

func (i BoolFieldGetSet) Set(value bool) {
	i.setter(value)
}

func NewBoolFieldGetSet(getter func() bool, setter func(bool)) BoolFieldGetSet {
	if getter == nil {
		panic("getter is nil")
	}
	if setter == nil {
		panic("setter is nil")
	}
	return BoolFieldGetSet{
		getter: getter,
		setter: setter,
	}
}
