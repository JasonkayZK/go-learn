package prime

var filter []int

func SerialPrintPrimeN(n int) {
	filter = make([]int, 0, n+1)

	//fmt.Println("0: 2")

	for prime, i := 3, 1; i < n; prime += 2 {
		if FilterPrime(prime) {
			filter = append(filter, prime)
			//fmt.Printf("%d: %d\n", i, prime)
			i++
		}
	}
}

func FilterPrime(x int) bool {
	for _, f := range filter {
		if x < f * f {
			return true
		}
		if x % f == 0 {
			return false
		}
	}
	return true
}
