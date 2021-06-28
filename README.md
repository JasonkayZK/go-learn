# Receiver

在Go中定义方法时可以将方法的接收器声明为值接收器和引用接收器两种；

废话不多说，直接来看下面代码定义的两个方法有什么区别：

```go
type Object struct {
	Elem string
}

func (o *Object) refEqual(o2 *Object) bool {
	fmt.Printf("refEqual str: %v\n", o.Elem == o2.Elem)
	return o == o2
}

func (o Object) copyEqual(o2 *Object) bool {
	fmt.Printf("copyEqual str: %v\n", o.Elem == o2.Elem)
	return &o == o2
}
```

两个Equal方法的接收器声明分别为值接收器和引用接收器；

下面进行测试：

receiver_test.go

```go
package main

import (
	"fmt"
	"testing"
)

type Object struct {
	Elem string
}

func (o *Object) refEqual(o2 *Object) bool {
	fmt.Printf("refEqual str: %v\n", o.Elem == o2.Elem)
	return o == o2
}

func (o Object) copyEqual(o2 *Object) bool {
	fmt.Printf("copyEqual str: %v\n", o.Elem == o2.Elem)
	return &o == o2
}

func TestEqual(t *testing.T) {
	elem := Object{Elem: "haha"}

	fmt.Printf("ref func: %v\n", elem.refEqual(&elem))
	fmt.Printf("copy func: %v\n", elem.copyEqual(&elem))
}
```

执行后输出：

```
refEqual str: true
ref func: true
copyEqual str: true
copy func: false
```

其实从方法名就可以看出区别：

-   <font color="#f00">**对于值接收器而言：在调用方法时会将当前对象完整的Copy一份，然后使用这个对象进行调用；**</font>
-   <font color="#f00">**对于引用接收器而言：在调用方法时会将当前对象的引用Copy一份，然后使用这个引用进行调用，因此，此时两个引用指向的是同一个堆中对象！**</font>

因此不难得出结论：

-   对于对象中的string而言：
    -   对于值接收器而言：**由于将对象完整的Copy了一份，所以对象中的字符串必定是相同的；**
    -   对于引用接收器而言：**由于两个引用指向的是同一个堆中对象，所以对象中的字符串也必定是相同的；**
-   对于对象本身而言：
    -   对于值接收器而言：**由于将对象完整的Copy了一份，因此两个对象是堆中不同的对象！**
    -   对于引用接收器而言：**由于两个引用指向的是同一个堆中对象，所以对象本身也是相同的；**

<br/>

## **接收器和Interface**

由于Go是强类型的语言，你会认为`*Object`和`Object`在实现Interface时会被认为是两个类型！

>   **我们理所当然的会想：**
>
>   -   **`Object`就是Object类型；**
>   -   **`*Object`是指向Object类型的指针类型；**
>
>   **所以对于Interface而言，在实现时`(o *Object)`和`(o Object)`也是不同的！**

<font color="#f00">**然而在Go中，对于Interface而言，由于Go不支持重载，所以实际上对于`(o *Object)`和`(o Object)`，我们只能实现一个！**</font>

来看下面的代码：

```go
type Introduce interface {
	Introduce() string
}

type Item struct {
	Elem string
}

func (i *Item) Introduce() string {
	return "haha from reference"
}

func (i Item) Introduce() string {
	return "haha from object"
}
```

**上面的代码无法编译，因为会被认为`Introduce`方法被重复定义！**

这就奇怪了，明明类型和类型的引用是不同的类型，为什么会被认为是重复定义呢？

其实，这是由Go中方法调用的机制决定的：

<font color="#f00">**Go会自动判断调用方法的是具体类型还是类型引用，并增加`&`或者`*`来帮助你完成方法调用；**</font>

例子如下：

interface_test.go

```go
package main

import (
	"fmt"
	"testing"
)

type Introduce interface {
	Introduce() string
}

type Item struct {
	Elem string
}

//func (i *Item) Introduce() string {
//	return "haha from reference"
//}

func (i Item) Introduce() string {
	return "haha from object"
}

func TestInterface(t *testing.T) {
	elem := &Item{Elem: "haha"}

	fmt.Println(elem.Introduce())
	fmt.Println((*elem).Introduce())
}
```

程序是可以被正常编译并输出的：

```
haha from object
haha from object
```

我们将`(i Item)`替换为`(i *Item)`，重新执行：

```
haha from reference
haha from reference
```

可以看到，也是可以正常执行的！

<font color="#f00">**因此，拜Go的方法调用补全和无重载机制所赐，我们无法对一个类型的值接收器和引用接收器同时实现一个Interface；换句话说，我们只能“二选一”！**</font>

<font color="#f00">**否则，Go将不能推断到底是采用的值调用还是引用调用；**</font>

>   **无法使用同一个方法名同时实现值接收器和引用接收器的问题不是Interface本身的问题，只要是在Go中声明方法，都存在这样的问题；**
>
>   **但是由于Interface需要比对函数签名，所以我们无法同时实现接口！**
>
>   **对于普通方法，我们可以根据方法名区分，如：**
>
>   -   **EqualByValue(o Object)；**
>   -   **EqualByRef(o *Object)；**

<br/>

## **Interface方法调用和Receiver**

经过上面的总结，我们基本了解了Go中的值接收器和引用接收器，以及他们和Interface的关系；

那么，对于实现了Interface的值接收器和引用接收器和普通的方法有什么区别呢？

答案是没有区别！

来看下面的例子：

interface2_test.go

```go
package main

import (
	"fmt"
	"testing"
)

type SelfCompare interface {
	Equal(item *MyElem) bool
}

type MyElem struct {
	Elem string
}

func (m *MyElem) Equal(item *MyElem) bool {
	return m == item
}

//func (m MyElem) Equal(item *MyElem) bool {
//	return &m == item
//}

func TestMyElem(t *testing.T) {
	elem := &MyElem{Elem: "haha"}

	fmt.Printf("ref func: %v\n", elem.Equal(elem))
}
```

结果如下：

-   当采用值接收器调用，返回为false；
-   当采用引用接收器调用，返回true；

<font color="#f00">**方法调用和普通的方法完全相同，唯一的区别在于此时MyElem引用可以被当作SelfCompare的Interface类型！**</font>

<br/>

<font color="#f00">**需要注意的是，对于使用Interface调用来说，是有区别的！**</font>

下面来看下面的例子：

```go
func callEqual(s SelfCompare) {
	fmt.Println(s.Equal(s.(*MyElem)))
}

func TestInterfaceCall(t *testing.T) {
	elem := &MyElem{Elem: "haha"}
	callEqual(elem)
}
```

输出结果如下：

-   当采用值接收器调用，返回为false；
-   当采用引用接收器调用，返回true；

输出结果和上面一致；

这是由于：<font color="#f00">**本质上，`SelfCompare`类型（或者说interface类型）就是一个其他类型的引用！**</font>

如果你不相信，你可以将代码稍作修改，声明方法为值引用，同时向`callEqual`函数传参时，仅传递对象值，而非引用：

```go
package main

import (
	"fmt"
	"testing"
)

type SelfCompare interface {
	Equal(item *MyElem) bool
}

type MyElem struct {
	Elem string
}

func (m MyElem) Equal(item *MyElem) bool {
	return &m == item
}

func callEqual(s SelfCompare) {
	fmt.Println(s.Equal(s.(*MyElem)))
}

func TestInterfaceCall(t *testing.T) {
	elem := MyElem{Elem: "haha"}
    fmt.Println(elem.Equal(&elem))
	callEqual(elem)
}
```

>   <font color="#f00">**注意：在这里我们向函数`callEqual`直接传递的是值！**</font>

尝试执行代码会产生一个Panic：

```
=== RUN   TestInterfaceCall
false

--- FAIL: TestInterfaceCall (0.00s)
panic: interface conversion: main.SelfCompare is main.MyElem, not *main.MyElem [recovered]
	panic: interface conversion: main.SelfCompare is main.MyElem, not *main.MyElem

goroutine 6 [running]:
testing.tRunner.func1.1(0xee4020, 0xc00007a4e0)
	E:/golang/src/testing/testing.go:1057 +0x310
testing.tRunner.func1(0xc000045080)
	E:/golang/src/testing/testing.go:1060 +0x43a
panic(0xee4020, 0xc00007a4e0)
	E:/golang/src/runtime/panic.go:969 +0x176
receiver.callEqual(0xf33a40, 0xc000050520)
	D:/workspace/Go_Learn/interface2_test.go:31 +0xd9
receiver.TestInterfaceCall(0xc000045080)
	D:/workspace/Go_Learn/interface2_test.go:37 +0xd7
testing.tRunner(0xc000045080, 0xf118e0)
	E:/golang/src/testing/testing.go:1108 +0xef
created by testing.(*T).Run
	E:/golang/src/testing/testing.go:1159 +0x397

Process finished with exit code 1
```

首先，我们使用`elem`本身去调用方法，是可以正常输出`false`的！

但是为什么到了函数中，转换为了`SelfCompare`的接口类型就不行了呢？

这是由于：<font color="#f00">**Go中的接口在调用方法时总是希望获取到一个对象的引用类型，而非对象本身；这一点从Panic输出的错误中也可以看出！**</font>

<font color="#f00">**同时，此时Go编译器并不会自动判断Interface是否是引用，并自动添加`*`和`&`！**</font>

因此，当我们传入一个对象时，产生了Panic；

**我想，这可能就是Go不允许类型的值接收器和引用接收器同时实现同一个接口的另一个原因吧？**

>   这里我还想吐槽一下Go的编译器，为了保证编译速度，这么重要的类型检查也都做的这么粗糙；
>
>   和Rust相比差了不是一点半点！

<br/>

## **Interface指针**

最后，再跑题说几句Interface指针吧；

从上面我们知道，对于一个对象方法的直接调用来说Go编译器会自动判断是否是引用，并自动添加`*`和`&`，但是对于Interface而言并不存在这个优化；

因此我们单纯将`callEqual`的入参修改为指针类型，函数将直接报错：

```go
func callEqual(s *SelfCompare) {
	fmt.Println(s.Equal(s.(*MyElem)))
}
```

我们可以通过手动添加解引用来修复错误：

```go
func callEqual(s *SelfCompare) {
	fmt.Println((*s).Equal((*s).(*MyElem)))
}
```

修改完成后，下面的测试代码也会报错：

```go
func TestInterfaceCall(t *testing.T) {
	elem := MyElem{Elem: "haha"}
	callEqual(elem) // Cannot use 'elem' (type MyElem) as type *SelfCompare
}
```

为了调用这段代码，我们需要将elem强制转换为`*SelfCompare`类型；

当你写下下面的代码后，会发现还是报错：

```go
func TestInterfaceCall(t *testing.T) {
	elem := MyElem{Elem: "haha"}
	e := elem.(SelfCompare) // Invalid type assertion: elem.(SelfCompare) (non-interface type MyElem on left)
	callEqual(&e)
}
```

我们不能将一个非Interface类型转换为Interface！

至此，你会发现这段代码写起来无比的别扭：因为Interface本身就已经是一个引用了，你没必要再去声明一个引用的引用；

因此，比较好的实践就是：<font color="#f00">**对于Interface类型入参永远使用对象引用，同时永远不要使用`*Interface`的骚操作！**</font>

<br/>

## **总结**

本文比较深入的探讨了Go中方法实现的值接收器和引用接收器以及他们和Interface的联系；

最后跑题聊了聊Go中的Interface；

文中内容都是本人编写Go代码中的一些思索，如有不对之处，还请批评指出！

<br/>

## **其他**

相关博文：

-   Github Pages：[一文看懂Go方法中的值接收器和引用接收器](https://jasonkayzk.github.io/2021/06/28/一文看懂Go方法中的值接收器和引用接收器/)
-   国内Gitee镜像：[一文看懂Go方法中的值接收器和引用接收器](https://jasonkay.gitee.io/2021/06/28/一文看懂Go方法中的值接收器和引用接收器/)

<br/>