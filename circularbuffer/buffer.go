package circularbuffer

type CircularBuffer struct {
	data  []float64
	size  int
	start int
	count int
}

func CreateBuffer(size int) *CircularBuffer {
	return &CircularBuffer{
		data: make([]float64, size),
		size: size,
	}
}

func (cb *CircularBuffer) Add(value float64) {
	index := (cb.start + cb.count) % cb.size
	cb.data[index] = value

	if cb.count < cb.size {
		cb.count++
	} else {
		cb.start = (cb.start + 1) % cb.size
	}
}

func (cb *CircularBuffer) GetAll() []float64 {
	result := make([]float64, cb.count)
	for i := 0; i < cb.count; i++ {
		result[i] = cb.data[(cb.start+i)%cb.size]
	}
	return result
}

func (cb *CircularBuffer) GetLastN(n int) []float64 {
	if n > cb.count {
		n = cb.count
	}

	result := make([]float64, n)

	for i := 0; i < n; i++ {
		result[i] = cb.data[(cb.start+cb.count-n+i)%cb.size]
	}

	return result
}
