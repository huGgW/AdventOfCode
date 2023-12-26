package main

type empty struct{}

type Set[T comparable] struct {
    m map[T]empty
}

func NewSet[T comparable]() *Set[T] {
    return &Set[T] {make(map[T]empty)}
}

func (s *Set[T])Add(elem T) {
    s.m[elem] = empty{}
}

func (s *Set[T])Remove(elem T) {
    delete(s.m, elem)
}

func (s *Set[T])Exists(elem T) bool {
    _, ok := s.m[elem]
    return ok
}

func (s *Set[T])Length() int {
    return len(s.m)
}

func (s *Set[T])ToSlice() []T {
    sl := []T{}

    for k, _ := range s.m {
        sl = append(sl, k)
    }

    return sl
}
