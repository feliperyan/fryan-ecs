package ecs_cpp_style

import "fmt"

const MAX_ENTITIES int = 5000
const MAX_COMPONENTS uint8 = 32

// -------

type DumbSignature []bool

func (s DumbSignature) Set(pos int, val bool) error {
	if pos < 0 || pos >= len(s) {
		return fmt.Errorf("out of bounds")
	}
	s[pos] = val
	return nil
}

func (s DumbSignature) Get(pos int) (error, bool) {
	if pos < 0 || pos >= len(s) {
		return fmt.Errorf("out of bounds"), false
	}
	return nil, s[pos]
}

func (s DumbSignature) Contains(s2 DumbSignature) (bool, error) {
	if len(s) != len(s2) {
		return false, fmt.Errorf("unequal lengths")
	}

	for i, _ := range s {
		if s[i] == true && s2[i] == false {
			return false, nil
		}
	}

	return true, nil
}

func (s DumbSignature) Reset() {

	for i, _ := range s {
		s[i] = false
	}

}

// -------

// TODO: Queue is inneficient = a circular list would require less allocations, re-sizing etc.

type Queue[T any] []T

func (q *Queue[T]) enqueue(v T) {
	*q = append(*q, v)
}

func (q *Queue[T]) dequeue() (T, bool) {
	if len(*q) == 0 {
		var vv T
		return vv, false
	}
	v := (*q)[0]
	*q = (*q)[1:]
	return v, true
}
