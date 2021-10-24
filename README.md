# Goroutine limit method

Some method to limit goroutines spawn numbers, in order to avoid resource leakage!



## The Job Abstration

All jobs in the demos is doing kinda like this:

```go
func job(str string, jobIdx int, res *[]string) {
	fmt.Printf("str: %s, jobIdx: %d\n", str, jobIdx)
	(*res)[jobIdx] = prefix + str
}
```

Add prefix for str and generate new result array: res.



## Content

There are some demos:

-   No-Limit
-   Limit with Channel & WaitGroup (Standard Golang lib)
-   Semaphore (golang.org/x)
-   Ants: https://github.com/panjf2000/ants
-   Go-Playground: https://github.com/go-playground/pool
-   Tunny: https://github.com/Jeffail/tunny



## Linked Blog

-   [《控制Goroutine数量的方法》](https://jasonkayzk.github.io/2021/10/22/控制Goroutine数量的方法/)
