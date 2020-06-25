package prime

func GenerateNaturalSeq() chan int {
	ch := make(chan int)
	go func() {
		for i := 2; ; i++ {
			ch <- i
		}
	}()
	return ch
}

func AddPrimeFilterChan(in <-chan int, prime int) chan int {
	out := make(chan int)
	go func() {
		for {
			if i := <-in; i%prime != 0 {
				out <- i
			}
		}
	}()
	return out
}

func PrintPrimeN(n int) {
	ch := GenerateNaturalSeq()
	for i := 0; i < n; i++ {
		prime := <-ch
		//fmt.Printf("%d: %d\n", i, prime)
		ch = AddPrimeFilterChan(ch, prime)
	}
}
