package pring

import (
	"fmt"
	"testing"
)

func reader(t *testing.T, buffer []byte, expected int, r *PRing) {
	t.Parallel()
	count := 0
	for {
		n, err := r.Read(buffer)
		if err != nil {
			t.Fatalf("Read error: %v", err)
		}
		fmt.Printf("Read %v\n", buffer[:n])
		count += n
		if count == expected {
			break
		}
	}
}

func writer(t *testing.T, buffer []byte, pattern []int, r *PRing) {
	t.Parallel()
	for _, v := range pattern {
		r.Write(buffer[:v])
		fmt.Printf("Wrote %v\n", buffer[:v])
	}
}

func TestWrite(t *testing.T) {
	buffer := make([]byte, 256)
	for i := 0; i < 256; i++ {
		buffer[i] = byte(i)
	}

	tests := []struct {
		Name     string
		RingSize int
		BufSize  int
		Pattern  []int
		Expected int
	}{
		{"Test 1", 10, 10, []int{1}, 1},
		{"Test 2", 10, 10, []int{4, 4, 4}, 12},
		{"Test 3", 10, 10, []int{1, 1, 1, 1, 1, 1, 1, 1}, 8},
		{"Test 4", 10, 10, []int{8, 7, 6, 5, 4, 3, 2, 1}, 36},
		{"Test 5", 10, 20, []int{10, 10, 10, 10}, 40},
		{"Test 6", 10, 20, []int{3, 3, 3, 3, 3, 3, 3, 3, 3, 3}, 30},
		{"Test 7", 10, 20, []int{5, 5, 5, 5, 5, 5, 5, 5}, 40},
		{"Test 8", 10, 20, []int{0, 0, 3, 10, 0, 7}, 20},
		{"Test 9", 100, 100, []int{50, 50, 50, 60, 40}, 250},
	}

	for _, test := range tests {
		t.Run(test.Name,
			func(t *testing.T) {
				pr := NewPRing(test.RingSize)
				b := make([]byte, test.BufSize)
				t.Run("writer",
					func(t *testing.T) {
						writer(t, buffer, test.Pattern, pr)
					})
				t.Run("reader",
					func(t *testing.T) {
						reader(t, b, test.Expected, pr)
					})
			})
	}

	/*
		b := NewPRing(10)
		n, err := b.Write([]byte{1, 2, 3, 4, 5})
		if n != 5 {
			t.Errorf("Write: Expected %d, have %d.", 5, n)
		}
		if err != nil {
			t.Errorf("Write: Expected %v, have %v.", nil, err)
		}

		s := make([]byte, 10)
		n, err = b.Read(s)
		if n != 5 {
			t.Errorf("Read: Expected %d, have %d", 5, n)
		}
		if err != nil {
			t.Errorf("Read: Expected %v, have %v", nil, err)
		}
	*/

}

// Benchmarks

// Copy 100MB between two goroutines with 1500byte packets

func spoolReceiver(b *testing.B, sp *SPool, length int) {
	read := 0
	dest := make([]byte, 1500)
	for read < length {
		buffer := <-sp.c
		n := copy(dest, buffer)
		read += n
		sp.Put(buffer)
	}

}

func BenchmarkSPool(b *testing.B) {
	// Generate 100MB
	length := 100 * (1 << 20)
	data := make([]byte, length)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sp := NewSPool(1024)
		go spoolReceiver(b, sp, length)
		written := 0
		for written < length {
			buffer := sp.Get()
			n := copy(buffer, data[written:])
			written += n
			sp.c <- buffer
		}
	}
}

func min(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

func pringReceiver(b *testing.B, pr *PRing, length int) {
	read := 0
	//dest := make([]byte, 1500)
	// This simulates the extra copy needed when extracting data from the PRing
	// in normal applications
	//extra := make([]byte, 1<<18)
	for read < length {
		/*
			n, err := pr.Read(extra)
			if err != nil {
				fmt.Println("Read error", err)
			}
		*/

		// We copied a big chunk, now packetize it and copy it to destination
		/*
			copied := 0
			for copied < n {
				n2 := copy(dest, extra[copied:])
				copied += n2
			}
		*/

		//read += n
		read += 1500
	}
}

func BenchmarkPRing(b *testing.B) {
	// Generate 100MB
	length := 100 * (1 << 20)
	//data := make([]byte, length)
	pr := NewPRing(1 << 18)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go pringReceiver(b, pr, length)
		written := 0
		for written < length {
			//maxWrite := min(1500, len(data)-written)
			//n, err := pr.Write(data[written : written+maxWrite])
			/*
				if err != nil {
					fmt.Println("Write error", err)
				}
			*/
			//written += n
			written += 1500
		}
	}
}
