package main

type lifecycle[T lifecycler] struct {
	v T
}

type lifecycler interface {
	comparable
	begin()
	finish()
}

func (l *lifecycle[T]) begin(v T) {
	var zero T
	if l.v != zero {
		panic("lifecycle.begin: already begun")
	}
	l.v = v
	l.v.begin()
}

func (l *lifecycle[T]) finish() {
	var zero T
	if l.v == zero {
		panic("lifecycle.finish: not begun yet")
	}
	l.v.finish()
	l.v = zero
}

func (l *lifecycle[T]) get() (v T, ok bool) {
	var zero T
	if l.v == zero {
		return zero, false
	}
	return l.v, true
}
