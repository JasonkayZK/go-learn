# Context

## 为什么需要 Context

在CSP并发模型中，WaitGroup 和信道(channel)是常见的 2 种并发控制的方式；

如果并发启动了多个子协程，需要等待所有的子协程完成任务，WaitGroup 非常适合于这类场景，例如下面的例子：

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func doTask(n int) {
	time.Sleep(time.Duration(n) * time.Second)
	fmt.Printf("Task %d Done\n", n)
	wg.Done()
}

func main() {
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go doTask(i + 1)
	}
	wg.Wait()
	fmt.Println("All Task Done")
}

```

 `wg.Wait()` 会等待所有的子协程任务全部完成，所有子协程结束后，才会执行 `wg.Wait()` 后面的代码。 

```
Task 3 Done
Task 1 Done
Task 2 Done
All Task Done
```

 WaitGroup 只是傻傻地等待子协程结束，但是**并不能主动通知子协程退出；**

假如开启了一个定时轮询的子协程，有没有什么办法，通知该子协程退出呢？

这种场景下，可以使用 `select+chan` 的机制：

```go
package main

import (
	"fmt"
	"time"
)

var stop chan interface{}

func reqTask(name string) {
	for {
		select {
		case <-stop:
			fmt.Println("stop", name)
			return
		default:
			fmt.Println(name, "send request")
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	stop = make(chan interface{})
	go reqTask("worker1")
	time.Sleep(3 * time.Second)
	stop <- struct{}{}
	time.Sleep(3 * time.Second)
}

```

子协程使用 for 循环定时轮询，如果 `stop` 信道有值，则退出，否则继续轮询。

```
worker1 send request
worker1 send request
worker1 send request
stop worker1
```

更复杂的场景如何做并发控制呢？比如子协程中开启了新的子协程，或者需要同时控制多个子协程。这种场景下，`select+chan`的方式就显得力不从心了。

Go 语言提供了 Context 标准库可以解决这类场景的问题，Context 的作用和它的名字很像，上下文，即子协程的下上文。

Context 有两个主要的功能：

-   通知子协程退出（正常退出，超时退出等）；
-   传递必要的参数。

<br/>

##  context.Backgroud()

`context.Backgroud()` 创建根 Context；

通常在 main 函数、初始化和测试代码中创建，作为顶层 Context。

<br/>

## context.WithCancel

`context.WithCancel()` 创建可取消的 Context 对象，即可以主动通知子协程退出。

### 控制单个协程

使用 Context 改写上述的例子，效果与 `select+chan` 相同。

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func reqTask(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("stop", name)
			return
		default:
			fmt.Println(name, "send request")
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go reqTask(ctx, "worker1")
	time.Sleep(3 * time.Second)
	cancel()
	time.Sleep(3 * time.Second)
}

```

-   `context.Backgroud()` 创建根 Context，通常在 main 函数、初始化和测试代码中创建，作为顶层 Context。
-   `context.WithCancel(parent)` 创建可取消的子 Context，同时返回函数 `cancel`。
-   在子协程中，使用 select 调用 `<-ctx.Done()` 判断是否需要退出。
-   主协程中，调用 `cancel()` 函数通知子协程退出。

<br/>

### 控制多个协程

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func reqTask(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("stop", name)
			return
		default:
			fmt.Println(name, "send request")
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go reqTask(ctx, "worker1")
	go reqTask(ctx, "worker2")
	go reqTask(ctx, "worker3")

	time.Sleep(3 * time.Second)
	cancel()
	time.Sleep(3 * time.Second)
}

```

为每个子协程传递相同的上下文 `ctx` 即可，调用 `cancel()` 函数后该 Context 控制的所有子协程都会退出。

```
worker3 send request
worker2 send request
worker1 send request
worker1 send request
worker2 send request
worker3 send request
worker3 send request
worker2 send request
worker1 send request
stop worker2
stop worker3
stop worker1
```

>此次，可能有些人会有疑问：
>
>**在reqTask中传递的是ctx而非&ctx，为什么在cancel()的时候也会对复制的ctx产生作用？**
>
>点开 Context 的源码看一下就能明白了：
>
>首先，**context.Context 是一个接口，并不是真正的结构体类型**。定义如下：
>
>```go
>type Context interface {
>	Deadline() (deadline time.Time, ok bool)
>	Done() <-chan struct{}
>	Err() error
>	Value(key interface{}) interface{}
>}
>```
>
> WithCancel 返回的是 context.Context 接口，如果想知道返回的真正类型是什么，点开 WithCancel 源码，如下： 
>
>```go
>func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
>	if parent == nil {
>		panic("cannot create context from nil parent")
>	}
>	c := newCancelCtx(parent)
>	propagateCancel(parent, &c)
>	return &c, func() { c.cancel(true, Canceled) }
>}
>
>// newCancelCtx returns an initialized cancelCtx.
>func newCancelCtx(parent Context) cancelCtx {
>	return cancelCtx{Context: parent}
>}
>```
>
> 它返回的 context.Context 类型，但 **ctx 本身的类型是 *cancelCtx**，的确是地址传递，而非值传递；
>
> cancelCtx 的定义如下： 
>
>```go
>type cancelCtx struct {
>	Context
>
>	mu       sync.Mutex            // protects following fields
>	done     chan struct{}         // created lazily, closed by first cancel call
>	children map[canceler]struct{} // set to nil by the first cancel call
>	err      error                 // set to non-nil by the first cancel call
>}
>```
>
> cancelCtx 并没有暴露出来，而是 Go 的内部实现，返回的底层类型是 *cancelContext，同时也满足 context.Context 接口；

<br/>

## context.WithValue

如果需要往子协程中传递参数，可以使用 `context.WithValue()`。

```go
package main

import (
	"context"
	"fmt"
	"time"
)

type Options struct{
	Interval time.Duration
}

func reqTask(ctx context.Context, name string) {
	for {
		select {
			case <- ctx.Done():
				fmt.Println("stop", name)
				return
		default:
			fmt.Println(name, "send request")
			op := ctx.Value("options").(*Options)
			time.Sleep(op.Interval * time.Second)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	vCtx := context.WithValue(ctx, "options", &Options{1})

	go reqTask(vCtx, "worker1")
	go reqTask(vCtx, "worker2")

	time.Sleep(3 * time.Second)
	cancel()
	time.Sleep(3 * time.Second)
}

```

-   `context.WithValue()` 创建了一个基于 `ctx` 的子 Context，并携带了值 `options`。
-   在子协程中，使用 `ctx.Value("options")` 获取到传递的值，读取/修改该值。

<br/>

## context.WithTimeout

如果需要控制**子协程的执行时间**，可以使用 `context.WithTimeout` 创建具有超时通知机制的 Context 对象。

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func reqTask(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("stop", name)
			return
		default:
			fmt.Println(name, "send request")
			// 0.5s
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	go reqTask(ctx, "worker1")
	go reqTask(ctx, "worker2")

	time.Sleep(3 * time.Second)
	fmt.Println("before cancel")
	cancel()
	time.Sleep(3 * time.Second)
}

```

`WithTimeout()`的使用与 `WithCancel()` 类似，多了一个参数，用于设置超时时间。执行结果如下：

```
worker1 send request
worker2 send request
worker2 send request
worker1 send request
worker1 send request
worker2 send request
stop worker1
stop worker2
before cancel
```

因为超时时间设置为 2s，但是 main 函数中，3s 后才会调用 `cancel()`，因此，在调用 `cancel()` 函数前，子协程因为超时已经退出了。

<br/>

## context.WithDeadline

超时退出可以控制子协程的最长执行时间，那 `context.WithDeadline()` 则可以控制子协程最迟退出的时间点。

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func reqTask(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("stop", name, ctx.Err())
			return
		default:
			fmt.Println(name, "send request")
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
	go reqTask(ctx, "worker1")
	go reqTask(ctx, "worker2")

	time.Sleep(3 * time.Second)
	fmt.Println("before cancel")
	cancel()
	time.Sleep(3 * time.Second)
}
```

-   `WithDeadline` 用于设置截止时间。在这个例子中，将截止时间设置为1s后，`cancel()` 函数在 3s 后调用，因此子协程将在调用 `cancel()` 函数前结束。
-   在子协程中，可以通过 `ctx.Err()` 获取到子协程退出的错误原因。

运行结果如下：

```
worker2 send request
worker1 send request
stop worker2 context deadline exceeded
stop worker1 context deadline exceeded
before cancel
```

可以看到，子协程 `worker1` 和 `worker2` 均是因为截止时间到了而退出。

<br/>

## Context 使用原则和技巧

-   **不要把Context放在结构体中，要以参数的方式传递，parent Context一般为Background**
-   应该要**把Context作为第一个参数传递给入口请求和出口请求链路上的每一个函数，放在第一位，变量名建议都统一，如ctx**
-   **给一个函数方法传递Context的时候，不要传递nil，否则在tarce追踪的时候，就会断了连接**
-   **Context的Value相关方法应该传递必须的数据，不要什么数据都使用这个传递**
-   **Context是线程安全的，可以放心的在多个goroutine中传递**
-   **可以把一个 Context 对象传递给任意个数的 gorotuine，对它执行 取消 操作时，所有 goroutine 都会接收到取消信号**

<br/>

## 附录

参考文章：

-   [Go Context 并发编程简明教程](https://geektutu.com/post/quick-go-context.html)
-   [Golang Context深入理解](https://juejin.im/post/6844903555145400334)



