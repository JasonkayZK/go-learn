# **在Golang发生Panic后打印出堆栈信息**

对于实际的项目来说，框架都会提供`recover`来做业务发生`panic`时的拦截，保证整个服务不会因为一个业务的`panic`而导致整个服务直接挂掉；

同时，通常情况下框架都会记录并打出`panic`的堆栈信息，但是在框架之外，我们该怎么打印出来堆栈信息呢？

其实很简单通过`runtime.Stack`函数即可！

下面的三行代码就能返回当前Goroutine的堆栈信息：

```go
// getCurrentGoroutineStack 获取当前Goroutine的调用栈，便于排查panic异常
func getCurrentGoroutineStack() string {
	var buf [defaultStackSize]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}
```

下面看一个实际项目抽象出的例子：

```go
package main

import (
	"fmt"
	"runtime"
	"sync"
)

const (
	defaultStackSize = 4096
)

func callPanic() {
	panic("test panic")
}

// getCurrentGoroutineStack 获取当前Goroutine的调用栈，便于排查panic异常
func getCurrentGoroutineStack() string {
	var buf [defaultStackSize]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}

func task(arr *[]int, i int, wg *sync.WaitGroup, lock *sync.Mutex) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("[panic] err: %v\nstack: %s\n", err, getCurrentGoroutineStack())
		}
		wg.Done()
	}()

	if i == 500 {
		callPanic()
	}

	lock.Lock()
	defer lock.Unlock()
	*arr = append(*arr, i)
}

func main() {
	wg := sync.WaitGroup{}
	lock := sync.Mutex{}

	arr := make([]int, 0)
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go task(&arr, i, &wg, &lock)
	}
	wg.Wait()

	fmt.Println(len(arr))
}
```

在main函数中，会并发的创建10000个`task`任务；

在每个`task`任务中，会向arr数组的末尾添加一个 i 值；

>   <font color="#f00">**注：Golang中内置的`append`函数是非线程安全的！**</font>

同时，当 i 为500时，代码模拟了业务panic的场景；

并且，为了防止单个 task 的 panic 影响到其他任务，我们在每一个 task 任务的开头都声明了defer函数，在其中使用`recover`对panic进行了拦截；

执行代码后输出：

```
[panic] err: test panic
stack: goroutine 507 [running]:
main.getCurrentGoroutineStack(...)
	D:/workspace/Go_Learn/app.go:20
main.task.func1(0xc000010090)
	D:/workspace/Go_Learn/app.go:27 +0xc5
panic(0x963180, 0x99cfa0)
	E:/golang/src/runtime/panic.go:969 +0x176
main.callPanic(...)
	D:/workspace/Go_Learn/app.go:14
main.task(0xc000004480, 0x1f4, 0xc000010090, 0xc0000100a0)
	D:/workspace/Go_Learn/app.go:33 +0x197
created by main.main
	D:/workspace/Go_Learn/app.go:48 +0x10f

9999
```

可以看到单个 task 的 panic 并不会影响到其他 task：对于添加10000个数的任务，单个任务panic后，其他的9999个任务仍然正常的执行了！

同时，我们可以很容易的定位到，Panic 来源于 `D:/workspace/Go_Learn/app.go:14`，即代码的第14行！

<br/>

## **总结**

对于并发的情况，对于 task 的抽象是非常重要的；

同时，对于每一个单独的并发 task，都推荐采用下面的代码来对 panic 进行拦截，防止一个 task 的 panic 影响到其他所有的 task；

并且，为每一个 task 在 panic 时打印出堆栈来直接定位问题，并保证 WaitGroup 能够正常退出；

```go
defer func() {
    if err := recover(); err != nil {
        fmt.Printf("[panic] err: %v\nstack: %s\n", err, getCurrentGoroutineStack())
    }
    wg.Done()
}()
```

<br/>

# **附录**

源代码：

-   https://github.com/JasonkayZK/Go_Learn/tree/add-panic-log

<br/>
