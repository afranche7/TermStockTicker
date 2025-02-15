package circularbuffer

import "errors"

type CircularBuffer struct {
	data  []int
	size  int
	start int
	count int
}

func CreateBuffer(size int) *CircularBuffer {
	return &CircularBuffer{
		data: make([]int, size),
		size: size,
	}
}

func (cb *CircularBuffer) Add(value int) {
	index := (cb.start + cb.count) % cb.size
	cb.data[index] = value

	if cb.count < cb.size {
		cb.count++
	} else {
		cb.start = (cb.start + 1) % cb.size
	}
}

func (cb *CircularBuffer) GetAll() []int {
	result := make([]int, cb.count)
	for i := 0; i < cb.count; i++ {
		result[i] = cb.data[(cb.start+i)%cb.size]
	}
	return result
}

func (cb *CircularBuffer) Get(index int) (int, error) {
	if index < 0 || index >= cb.count {
		return 0, errors.New("index out of range")
	}
	return cb.data[(cb.start+index)%cb.size], nil
}

func (cb *CircularBuffer) Size() int {
	return cb.count
}
