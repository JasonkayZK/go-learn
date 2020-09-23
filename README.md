A repository for learning Golang.

# Nil

## 引出问题

最初遇到问题是因为，在某一次调试的时候发现：

当项目启动的时候，即使某个组件初始化失败，返回了nil，也是可以通过这个nil调用一些方法的；

例如下面的这个例子：

```go
package main

import (
	"fmt"
)

type Student struct {
	Name string
	Age  int
}

func (s *Student) Say() {
	s.sayHello()
}

func (s *Student) sayHello() {
	fmt.Println("hello")
}

func (s *Student) SayName() {
	fmt.Println(s.Name)
}

func main() {
	student, err := InitStudent(true)
	if err != nil {
		fmt.Println(err)
	}
	student.Say()
}

func InitStudent(needErr bool) (*Student, error) {
	if needErr {
		return nil, fmt.Errorf("init student err")
	}
	return &Student{}, nil
}

```

尽管在初始化方法InitStudent中由于错误导致了初始化失败，返回了nil；

但是我们的程序还是可以正常的执行Say方法，并同时打印出错误：

```
init student err
hello
```

直接一点，甚至我们可以直接在main中使用nil来调用Say方法：

```go
func main() {
	var nilPointer *Student = nil
	nilPointer.Say()
}
```

这个方法也是可以正常执行的！并不会报空指针错误！

但是如果我们调用SayName去访问student实例中的Name，则就会报空指针错误了：

```go
func main() {
	var nilPointer *Student = nil
	nilPointer.SayName()
}
```

<BR/>

## 关于Nil

下面是Go对于nil的定义：

```go
// nil is a predeclared identifier representing the zero value for a
// pointer, channel, func, interface, map, or slice type.
// Type must be a pointer, channel, func, interface, map, or slice type
var nil Type 

// Type is here for the purposes of documentation only. It is a stand-in
// for any Go type, but represents the same type for any given function
// invocation.
type Type int
```

>   很多人都误以为 golang中的`nil`与Java、PHP等编程语言中的null一样；
>
>   但是实际上Golang的niu复杂得多了；

可以看出，对于Go这种强类型的语言来说即使nil（准确来说是空值）也是有类型区别的！

### nil的零值

按照Go语言规范，任何类型在未初始化时都对应一个零值：布尔类型是false，整型是0，字符串是""，而指针、函数、interface、slice、channel和map的零值都是nil。  

>   **PS：这里没有说结构体struct的零值为nil，因为struct的零值与其属性有关**

`nil`没有默认的类型，尽管它是多个类型的零值，但是必须显式或隐式指定每个nil用法的明确类型；

例如：

```go
func main() {
	// 明确.
	_ = (*struct{})(nil)
	_ = []int(nil)
	_ = map[int]bool(nil)
	_ = chan string(nil)
	_ = (func())(nil)
	_ = interface{}(nil)

	// 隐式.
	var _ *struct{} = nil
	var _ []int = nil
	var _ map[int]bool = nil
	var _ chan string = nil
	var _ func() = nil
	var _ interface{} = nil
}
```

>   如果关注过golang关键字的同学就会发现，里面并没有`nil`，也就是说`nil`并不是关键字；
>
>   那么就可以在代码中定义`nil`，那么`nil`就会被隐藏！
>
>   例如下面的代码也是合法的（虽然这样做是强烈不推荐的！）：
>
>   ```go
>   func main() {
>   	// 123
>   	nil := 123
>   	fmt.Println(nil)
>   	//cannot use nil (type int) as type map[string]int in assignment
>   	//var _ map[string]int = nil
>   }
>   
>   ```

### nil类型的地址和值大小

<font color="#f00">**`nil`类型的所有值的内存布局始终相同，换一句话说就是：不同类型`nil`的内存地址是一样的！**</font>

如下：

```go
func main() {
	var m map[int]string
	var ptr *int
	var sl []int
    var i interface{} = nil
	fmt.Printf("%p\n", m)   //0x0
	fmt.Printf("%p\n", ptr) //0x0
	fmt.Printf("%p\n", sl)  //0x0
    fmt.Printf("%p\n", i)   //%!p(<nil>)
}

```

<font color="#f00">**但是nil值的大小始终与其类型与`nil`值相同的`non-nil`值大小相同；因此, 表示不同零值的nil标识符可能具有不同的大小。**</font>

>   这里有一个例外：
>
>   <font color="#f00">**interface{}类型的变量是真的空的，它是真的不会被分配内存空间！**</font>

例如，虽然下面的各个nil变量的地址是相同的，但是指针的大小是不同的：

```go
func main() {
	var p *struct{} = nil
	fmt.Printf("%p\n", p)  //0x0
	fmt.Println(unsafe.Sizeof(p)) // 8

	var s []int = nil
	fmt.Printf("%p\n", s)  //0x0
	fmt.Println(unsafe.Sizeof(s)) // 24

	var m map[int]bool = nil
	fmt.Printf("%p\n", m)  //0x0
	fmt.Println(unsafe.Sizeof(m)) // 8

	var c chan string = nil
	fmt.Printf("%p\n", c)  //0x0
	fmt.Println(unsafe.Sizeof(c)) // 8

	var f func() = nil
	fmt.Printf("%p\n", f)  //0x0
	fmt.Println(unsafe.Sizeof(f)) // 8

	var i interface{} = nil
	fmt.Printf("%p\n", i)  //%!p(<nil>)
	fmt.Println(unsafe.Sizeof(i)) // 16
}

```

大小是编译器和体系结构所依赖的，以上打印结果为64位体系结构和正式 Go 编译器；

对于32位体系结构, 打印的大小将是一半。  

<font color="#f00">**对于正式 Go 编译器, 同一种类的不同类型的两个nil值的大小始终相同；**</font>

例如，两个不同的切片类型 ( []int和[]string) 的两个nil值始终相同；

### nil值比较

**① 不同类型的`nil`是不能比较的**

由于Go是强类型的语言，所以这个结论是很容易得出的；

例如：

```go
func main() {
	var m map[int]string
	var ptr *int

	//invalid operation: m == ptr (mismatched types map[int]string and *int)
	fmt.Printf(m == ptr)
}

```

在 Go 中，两个不同可比较类型的两个值只能在一个值可以隐式转换为另一种类型的情况下进行比较。

具体来说，有两个案例两个不同的值可以比较：

-   两个值之一的类型是另一个的基础类型；
-   两个值之一的类型实现了另一个值的类型 (必须是接口类型)；

`nil`值比较也没有脱离上述规则，例如：

```go
func main() {
	type IntPtr *int
	fmt.Println(IntPtr(nil) == (*int)(nil))        //true
	fmt.Println((interface{})(nil) == (*int)(nil)) //false
}

```

****

**② 同一类型的两个`nil`值可能无法比较**

 因为golang中，map、slice和函数类型是不可比较类型；

它们有一个别称为`不可比拟的类型`，所以比较它们的`nil`亦是非法的！

例如：

```go
func main() {
	var v1 []int = nil
	var v2 []int = nil
	fmt.Println(v1 == v2)
	fmt.Println((map[string]int)(nil) == (map[string]int)(nil))
	fmt.Println((func())(nil) == (func())(nil))
}

```

<font color="#f00">**但是`不可比拟的类型`的值是可以与`纯nil`进行比较的！**</font>

例如：

```go
func main() {
	fmt.Println((map[string]int)(nil) == nil) //true
	fmt.Println((func())(nil) == nil)         //true
}

```

****

**③ 两`nil`值可能不相等**

 如果两个比较的nil值一个是接口值，另一个不是，假设它们是可比较的, 则比较结果总是 false！

原因在于：进行比较之前，接口值将转换为接口具体类型。转换后的接口值具有具体的动态类型，但其他接口值没有，这就是为什么比较结果总是错误的；

例如：

```go
func main() {
	// false
	fmt.Println((interface{})(nil) == (*int)(nil))
}

```

<BR/>

## 结论

回到最开始我们的问题，nil也能够调用方法的原因其实很简单：

虽然值为nil，但是它的类型是`*Student`，而`*Stundent`类型绑定了Say函数，而且Say并没有访问对象的任何变量，而导致panic的SayName是因为访问了对象中的成员；

使用nil进行函数调用的情况有点类似于Java中的静态方法：在调用静态方法时，不需要创建实例；

最后，再给一个比较有意思的例子：通过强转将一个nil转换为其他类型的nil，然后使用nil调用其方法：

```go
package main

import "fmt"

type Printer interface {
	Print()
}

type Student struct {
	Name string
	Age int
}

func (s *Student) Print() {
	fmt.Println("hello")
}

func main() {
	nilPointer := (*Student)(nil)
	nilPointer.Print()
}

```

 其实也比较好理解，我们把一个空指针，强转成一个*Student类型，必然给它附上了对应的函数指针，它就跟纯nil不太一样了；

