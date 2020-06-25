package prime

// The max value of the prime filter
// Update only when new prime filter added
var MaxFilter chan int

func NewGenerateNaturalSeq() chan int {
	ch := make(chan int)
	go func() {
		for i := 3; ; i += 2 {
			ch <- i
		}
	}()
	return ch
}

func NewAddPrimeFilterChan(in <-chan int, prime int) chan int {
	out := make(chan int)
	MaxFilter = out

	go func() {
		for {
			if i := <-in; i%prime != 0 {
				// Boundary condition:
				// Pass the i to the last filter
				if i < prime*prime {
					MaxFilter <- i
				} else {
					out <- i
				}
			}
		}
	}()

	return out
}

func NewPrintPrimeN(n int) {
	ch := NewGenerateNaturalSeq()
	//fmt.Println("0: 2")
	for i := 1; i < n; i++ {
		prime := <-ch
		//fmt.Printf("%d: %d\n", i, prime)
		ch = NewAddPrimeFilterChan(ch, prime)
	}
}
