package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	arr := [...]int{1: 2, 3, 4, 9}
	modify(arr)
	fmt.Println(arr)
	fmt.Println("-------------------------------------------")
	fmt.Println("-------------------------------------------")
	fmt.Println("-------------------------------------------")
	fmt.Println("-------------------------------------------")
	slice := make([]int, 5, 10)
	fmt.Println(len(slice))
	slice1 := []float64{1: 2}
	fmt.Println(slice1)
	fmt.Println(len(slice1))
	//nil切片（不存在的切片）
	var nilSlice []int

	fmt.Println(nilSlice)
	//空切片（空集合）
	slice2 := []int{}
	fmt.Println(slice2)
	slice = append(slice, 1, 2, 3, 4, 5)
	fmt.Println(slice)
	fmt.Println(slice[:])
	fmt.Println(slice[1:6])
	//新切片的值包含原切片的i索引，但是不包含j索引
	fmt.Println(slice[0:10])
	fmt.Println(slice[1:5])
	newSlice := slice[0:7]
	fmt.Println(newSlice)
	newSlice[0] = 10
	//新的切片和原切片共用的是一个底层数组，所以当修改的时候，
	//底层数组的值就会被改变，所以原切片的值也改变了。
	//当然对于基于数组的切片也一样的。
	fmt.Println(newSlice)
	fmt.Println(slice)

	slice3 := []int{1, 2, 3, 4, 5}
	newSlice1 := slice3[1:3]
	//对于底层数组容量是k的切片slice[i:j]来说
	//长度：j-i
	//容量:k-i
	fmt.Printf("newSlice长度:%d,容量:%d", len(newSlice1), cap(newSlice1))
	newSlice2 := slice3[1:2:3]
	fmt.Println(newSlice2)
	fmt.Println(cap(newSlice2))
	//如果切片的底层数组，没有足够的容量时，就会新建一个底层数组，
	// 把原来数组的值复制到新底层数组里，
	// 再追加新值，这时候就不会影响原来的底层数组了。

	fmt.Println(slice3)
	// 这两个切片的地址不一样，
	// 所以可以确认切片在函数间传递是复制的。
	// 而我们修改一个索引的值后，
	// 发现原切片的值也被修改了，
	// 说明它们共用一个底层数组。

	// 在函数间传递切片非常高效，
	// 而且不需要传递指针和处理复杂的语法，
	// 只需要复制切片，然后根据自己的业务修改，
	// 最后传递回一个新的切片副本即可，
	// 这也是为什么函数间传递参数，使用切片，而不是数组的原因。
	for i, v := range slice3 {
		fmt.Println(i, v)
	}
	fmt.Println(&slice3)
	modifySlice(slice3)
	fmt.Println(slice3)

	//这个是因为s1容量不够，
	// 所以append调用的时候会产生一个新的切片，
	// 所以append对原来的s1没有影响
	s1 := make([]int, 3)
	fmt.Println(cap(s1)) //容量为3 append会容量不够，导致s1 没加上append 的值，对原有的类容的修改仍然有效
	fmt.Println(s1)
	modify1(s1)
	fmt.Println(s1)

	//key： string  value: int
	dict := make(map[string]int)
	fmt.Println(dict == nil)
	dict["张三"] = 43
	fmt.Println(dict)
	dict1 := map[string]int{"zhangsan": 43, "lisi": 32}
	fmt.Println(dict1)
	age, exists := dict["李四"]
	fmt.Println(age)
	fmt.Println(exists)
	//delete函数删除不存在的键也是可以的，只是没有任何作用。
	delete(dict1, "lisi")
	dict1["wada"] = 22
	dict1["wada1"] = 22
	dict1["wada2"] = 22
	for _, v := range dict1 {
		fmt.Println(v)
	}
	names := make([]string, 0)
	for k := range dict1 {
		names = append(names, k)
	}
	sort.Strings(names)
	for _, key := range names {
		fmt.Println(key, dict1[key])
	}
	//未初始化 nil 的map

	// 键可以使用==运算符进行比较，
	// 所以像切片、函数以及含有切片的结构类型就不能用于Map的键了，
	// 因为他们具有引用的语义，不可比较
	var dict2 map[string]int
	fmt.Println(dict2 == nil)
	// 本来name的值并没有被改变,
	// 也就是说，我们传递的时一个副本，
	// 并且返回一个新创建的字符串。
	// 基本类型因为是拷贝的值，
	// 并且在对他进行操作的时候，
	// 生成的也是新创建的值，
	// 所以这些类型在多线程里是安全的，
	// 我们不用担心一个线程的修改影响了另外一个线程的数据。
	name := "张三"
	fmt.Println(modifyName(name))
	fmt.Println(name)

	// 用类型和原始的基本类型恰恰相反，它的修改可以影响到任何引用到它的变量。
	// 在Go语言中，引用类型有切片、map、接口、函数类型以及chan。

	// 引用类型之所以可以引用，是因为我们创建引用类型的变量，其实是一个标头值，
	// 标头值里包含一个指针，指向底层的数据结构，当我们在函数中传递引用类型时，
	// 其实传递的是这个标头值的副本，
	// 它所指向的底层结构并没有被复制传递，这也是引用类型传递高效的原因。

	// 本质上，我们可以理解函数的传递都是值传递，只不过引用类型传递的是一个指向底层数据的指针，
	// 所以我们在操作的时候， 可以修改共享的底层数据的值，
	// 进而影响到所有引用到这个共享底层数据的变量。
	ages := map[string]int{"zhangsan": 22}
	fmt.Println(ages)
	modifyAges(ages)
	fmt.Println(ages)

	var p person
	fmt.Println(p.age)
	jim := person{10, "Jim"}
	jack := person{name: "Jack", age: 12}
	fmt.Println(jim, jack)
	//函数传参是值传递，所以对于结构体来说也不例外，结构体传递的是其本身以及里面的值的拷贝。
	//modifyJim(jim)
	//要修改age的值可以通过传递结构体的指针
	modifyJim(&jim)
	fmt.Println(jim)

	//结构体里引用类型
	sliceName := make([]string, 5, 10)
	sliceName[0] = "dasdad"
	sliceName[1] = "dasdad"
	sliceName[3] = "dasdad"
	fmt.Println(sliceName)
	names1 := Names{slice: sliceName}
	fmt.Println(names1)
	modifySl(names1)
	fmt.Println(names1)
	fmt.Println(sliceName)

	// 这就是Go灵活的地方，我们可以使用自定义的类型做很多事情，
	// 比如添加方法，比如可以更明确的表示业务的含义等等
	//var i Duration = 100
	//var j int64 = 100
	var dur Duration
	//dur = int64(dur)
	fmt.Println(dur)

	// 这个函数名称是小写开头的add，所以它的作用域只属于所声明的包内使用，不能被其他包使用，
	// 如果我们把函数名以大写字母开头，该函数的作用域就大了，可以被其他包调用。这也是Go语言中大小写的用处，
	// 比如Java中，就有专门的关键字来声明作用域private、protect、public等。
	sum := add(1, 2)
	fmt.Println(sum)

	//Go语言里有两种类型的接收者：值接收者和指针接收者。

	//使用值类型接收者定义的方法，在调用的时候，使用的其实是值接收者的一个副本，
	// 所以对该值的任何操作，不会影响原来的类型变量。
	p1 := person1{name: "张三"}
	fmt.Println(p1.String())
	fmt.Println(p1.name)

	// 如果我们使用一个指针作为接收者，那么就会其作用了，因为指针接收者传递的是一个指向原值指针的副本，
	// 指针的副本，指向的还是原来类型的值，
	// 所以修改时，同时也会影响原来类型变量的值。
	p1.modifyPerson1()
	fmt.Println(p1.name)

	//Go的编译器自动会帮我们取指针，以满足接收者的要求。
	(&p1).modifyPerson1()
	fmt.Println(p1.name)
	//如果是一个值接收者的方法，使用指针也是可以调用的，Go编译器自动会解引用，以满足接收者的要求，
	// 比如例子中定义的String()方法，也可以这么调用
	fmt.Println((&p1).String())

	//如果返回的值，我们不想使用，可以使用_进行忽略。
	file, _ := os.Open("/usr/tmp")
	fmt.Println(file)

	// 可以变参数，可以是任意多个。
	// 我们自己也可以定义可以变参数，可变参数的定义，在类型前加上省略号…即可
	print("1", 3, 5, 7)

	// 函数方法还有其他一些知识点，比如painc异常处理，递归等，
	// 这些在《Go语言实战》书里也没有介绍，
	// 这些基础知识，可以参考Go语言的那本圣经。

	// 抽象就是接口的优势，它不用和具体的实现细节绑定在一起，我们只需定义接口，告诉编码人员它可以做什么，
	// 这样我们可以把具体实现分开，这样编码就会更加灵活方面，适应能力也会非常强。
	var b bytes.Buffer
	_, _ = fmt.Fprint(&b, "Hello World")
	fmt.Println(b.String())

	///因为bytes.Buffer实现了接口io.Writer,所以我们可以通过w = &b赋值，
	// 这个赋值的操作会把定义类型的值存入接口类型的值。
	var w io.Writer
	w = &b
	fmt.Println(w)

	// 赋值操作执行后，如果我们对接口方法执行调用，其实是调用存储的用户定义类型的对应方法，
	// 这里我们可以把用户定义的类型称之为实体类型。

	//我们可以定义很多类型，让它们实现一个接口，那么这些类型都可以赋值给这个接口，
	// 这时候接口方法的调用，其实就是对应实体类型对应方法的调用，这就是多态。
	var a animal
	//var c cat

	//a = c
	//a.printInfo() //使用另外一个类型赋值

	var d dog
	a = d
	a.printInfo()

	// 如果要实现一个接口，必须实现这个接口提供的所有方法，但是实现方法的时候，我们可以使用指针接收者实现，
	// 也可以使用值接收者实现，这两者是有区别的，下面我们就好好分析下这两者的区别。
	//var e cat
	//invoke(e)
	//实体类型以值接收者实现接口的时候，不管是实体类型的值，还是实体类型值的指针，都实现了该接口。
	//实体类型以指针接收者实现接口的时候，只有指向这个类型的指针才被认为实现了该接口
	var f cat
	invoke(&f)
	//Methods Receivers	Values
	//(t T)	T and *T
	//(t *T)	*T
	//上面的表格可以解读为：如果是值接收者，实体类型的值和指针都可以实现对应的接口；
	// 如果是指针接收者，那么只有类型的指针能够实现对应的接口。

	// 嵌入类型，或者嵌套类型，这是一种可以把已有的类型声明在新的类型里的一种方式，这种功能对代码复用非常重要。
	// 在其他语言中，有继承可以做同样的事情，但是在Go语言中，没有继承的概念，Go提倡的代码复用的方式是组合，
	// 所以这也是嵌入类型的意义所在，组合而不是继承，所以Go才会更灵活
	ad := admin{user{"张三", "zhangshan@flysnow.com"}, "guanliyuan"}
	ad.sayHello()
	ad.user.sayHello() //内部类型user有一个sayHello方法，外部类型对其进行了覆盖，同名重写sayHello，
	// 然后我们在main方法里分别访问这两个类型的方法

	//这里就可以说明admin实现了接口Hello,但是我们又没有显示的声明类型admin实现，所以这个实现是通过内部类型user实现的，
	// 因为admin包含了user所有的方法函数，所以也就实现了接口Hello。
	sayHello(ad)
	sayHello(ad.user)

	// Go语言提供的是以大小写的方式进行区分的，如果一个类型的名字是以大写开头，那么其他包就可以访问；
	// 如果以小写开头，其他包就不能访问。
	// Go 可以推断变量的类型。
	// 这是一种非常好的能力，试想，我们在和其他人进行函数方法通信的时候，只需约定好接口，就可以了，
	// 至于内部实现，使用方是看不到的，隐藏了实现。
	l := NewLoginer()
	l.Login()
	// 以上例子，我们对于函数间的通信，通过Loginer接口即可，在main函数中，使用者只需要返回一个Loginer接口，
	// 至于这个接口的实现，使用者是不关心的，所以接口的设计者可以把defaultLogin类型设计为不可见，
	// 并让它实现接口Loginer，这样我们就隐藏了具体的实现。如果以后重构这个defaultLogin类型的具体实现时，
	// 也不会影响外部的使用者，极为方便，这也就是面向接口的编程。

	// .操作符前面的部分导出了，.操作符后面的部分才有可能被访问；
	// 如果.前面的部分都没有导出，那么即使.后面的部分是导出的，也无法访问。
	// 例子	可否访问
	// Admin.User.Name	是
	// Admin.User.name	否
	// Admin.user.Name	否
	// Admin.user.name	否
	// 以上表格中Admin 为外部类型,User(user)为内部类型,Name(name)为字段，
	// 以此来更好的理解最后的总结，当然方法也适用这个表格。

	// go语言中并发指的是让某个函数独立于其他函数运行的能力，一个goroutine就是一个独立的工作单元，
	// Go的runtime（运行时）会在逻辑处理器上调度这些goroutine来运行，一个逻辑处理器绑定一个操作系统线程，
	// 所以说goroutine不是线程，它是一个协程，也是这个原因，它是由Go语言运行时本身的算法实现的。
	// 概念	说明
	// 进程	一个程序对应一个独立程序空间 是一个容器，
	// 是属于这个程序的工作空间，比如它里面有内存空间、文件句柄、设备和线程等等。
	// 线程	一个执行空间，一个进程可以有多个线程
	// 逻辑处理器	执行创建的goroutine，绑定一个线程
	// 调度器	Go运行时中的，分配goroutine给不同的逻辑处理器
	// 全局运行队列	所有刚创建的goroutine都会放到这里
	// 本地运行队列	逻辑处理器的goroutine队列
	// 当我们创建一个goroutine的后，会先存放在全局运行队列中，等待Go运行时的调度器进行调度，
	// 把他们分配给其中的一个逻辑处理器，并放到这个逻辑处理器对应的本地运行队列中，最终等着被逻辑处理器执行即可。

	// 并发的概念和并行不一样，并行指的是在不同的物理处理器上同时执行不同的代码片段， 并行可以同时做很多事情，
	// 而并发是同时管理很多事情，因为操作系统和硬件的总资源比较少，
	// 所以并发的效果要比并行好的多，使用较少的资源做更多的事情，也是Go语言提倡的。

	// Go的并行:
	// 多创建一个逻辑处理器就好了，
	// 这样调度器就可以同时分配全局运行队列中的goroutine到不同的逻辑处理器上并行执行。

	// Go并发
	// 这里的sync.WaitGroup其实是一个计数的信号量，使用它的目的是要main函数等待两个goroutine执行完成后再结束，
	// 不然这两个goroutine还在运行的时候，程序就结束了，看不到想要的结果。
	// sync.WaitGroup的使用也非常简单，先是使用Add 方法设设置计算器为2，每一个goroutine的函数执行完之后，
	// 就调用Done方法减1。Wait方法的意思是如果计数器大于0，就会阻塞，
	// 所以main 函数会一直等待2个goroutine完成后，再结束。
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 1; i < 100; i++ {
			fmt.Println("A:", i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 1; i < 100; i++ {
			fmt.Println("B:", i)
		}
	}()
	wg.Wait()

	// Go默认是给每个可用的物理处理器都分配一个逻辑处理器，因为我的电脑是4核的，
	// 所以上面的例子默认创建了4个逻辑处理器，
	// 所以这个例子中同时也有并行的调度，如果我们强制只使用一个逻辑处理器，我们再看看结果。
	// 对于并发来说，就是Go语言本身自己实现的调度，对于并行来说，是和运行的电脑的物理处理器的核数有关的，
	// 多核就可以并行并发，单核只能并发了。
	for i := 1; i < 10; i++ {
		fmt.Println()
	}

	runtime.GOMAXPROCS(1)
	var wg1 sync.WaitGroup
	wg1.Add(2)
	go func() {
		defer wg1.Done()
		for i := 1; i < 10000; i++ {
			fmt.Println("A:", i)
		}
	}()
	fmt.Println(runtime.NumCPU())
	go func() {
		defer wg1.Done()
		for i := 1; i < 10000; i++ {
			fmt.Println("B:", i)
		}
	}()
	wg1.Wait()

	// 并发本身并不复杂，但是因为有了资源竞争的问题，
	// 就使得我们开发出好的并发程序变得复杂起来，因为会引起很多莫名其妙的问题。

	//这是一个资源竞争的例子，我们可以多运行几次这个程序，会发现结果可能是2，也可以是3，也可能是4。
	//因为共享资源count变量没有任何同步保护，所以两个goroutine都会对其进行读写，会导致对已经计算好的结果覆盖，
	//以至于产生错误结果，这里我们演示一种可能，两个goroutine我们暂时称之为g1和g2。
	//g1读取到count为0。
	//然后g1暂停了，切换到g2运行，g2读取到count也为0。
	//g2暂停，切换到g1，g1对count+1，count变为1。
	//g1暂停，切换到g2，g2刚刚已经获取到值0，对其+1，最后赋值给count还是1
	//有没有注意到，刚刚g1对count+1的结果被g2给覆盖了，两个goroutine都+1还是1
	wg2.Add(2)
	go incCount()
	go incCount()
	wg2.Wait()
	fmt.Println(count)
	//所以我们对于同一个资源的读写必须是原子化的，也就是说，同一时间只能有一个goroutine对共享资源进行读写操作。
	//共享资源竞争的问题，非常复杂，并且难以察觉，好在Go为我们提供了一个工具帮助我们检查，这个就是go build -race命令

	//传统解决资源竞争的办法—对资源加锁。
	wg2.Add(2)
	go incCountLock()
	go incCountLock()
	wg2.Wait()
	fmt.Println(count)

	// Go语言还提供了一个sync包，这个sync包里提供了一种互斥型的锁，可以让我们自己灵活的控制哪些代码，
	// 同时只能有一个goroutine访问，被sync互斥锁控制的这段代码范围，被称之为临界区，临界区的代码，
	// 同一时间 ，只能又一个goroutine访问。 sync.Mutex
	wg2.Add(2)
	go incCountMutex()
	go incCountMutex()
	wg2.Wait()
	fmt.Println(count)

	//除了原子函数和互斥锁，Go还为我们提供了更容易在多个goroutine同步的功能，这就是通道chan

	// 在多个goroutine并发中，我们不仅可以通过原子函数和互斥锁保证对共享资源的安全访问，消除竞争的状态，
	// 还可以通过使用通道，在多个goroutine发送和接受共享的数据，达到数据同步的目的。
	// make函数初始化的时候，只有一个参数，其实make还可以有第二个参数，用于指定通道的大小。
	// 默认没有第二个参数的时候，通道的大小为0，这种通道也被成为无缓冲通道。
	ch := make(chan int, 3)
	ch <- 2
	x := <-ch
	ch <- 4
	println(x)
	<-ch
	close(ch)

	// 无缓冲的通道定义来看，发送goroutine和接收gouroutine必须是同步的，同时准备后，如果没有同时准备好的话，
	// 先执行的操作就会阻塞等待，直到另一个相对应的操作准备好为止。这种无缓冲的通道我们也称之为同步通道
	// 在计算sum和的goroutine没有执行完，把值赋给ch通道之前，fmt.Println(<-ch)会一直等待，所以main主goroutine就不会终止，
	// 只有当计算和的goroutine完成后，并且发送到ch通道的操作准备好后，同时<-ch就会接收计算好的值，然后打印出来
	ch1 := make(chan int)
	go func() {
		var sum int = 0
		for i := 0; i < 10; i++ {
			sum += i
		}
		ch1 <- sum
	}()
	fmt.Println("--------------------------------------")
	fmt.Println(<-ch1)

	//上一个操作的输出，当成下一个操作的输入，连起来，做一连串的处理操作。
	one := make(chan int)
	two := make(chan int)
	go func() {
		one <- 100
	}()
	go func() {
		v := <-one
		two <- v
	}()

	fmt.Println(<-two)

	// 有缓冲通道，其实是一个队列，这个队列的最大容量就是我们使用make函数创建通道时，通过第二个参数指定的。
	//
	ch3 := make(chan int, 3)
	// 这里创建容量为3的，有缓冲的通道。对于有缓冲的通道，向其发送操作就是向队列的尾部插入元素，
	// 接收操作则是从队列的头部删除元素，并返回这个刚刚删除的元素。
	// 当队列满的时候，发送操作会阻塞；当队列空的时候，接受操作会阻塞。有缓冲的通道，不要求发送和接收操作时同步的，
	// 相反可以解耦发送和接收操作。
	// 想知道通道的容量以及里面有几个元素数据怎么办？其实和map一样，使用cap和len函数就可以了。
	fmt.Println(cap(ch3))
	fmt.Println(len(ch3))

	//定义单向通道也很简单，只需要在定义的时候，带上<-即可。
	//
	//var send chan<- int //只能发送
	//var receive <-chan int //只能接收

	log.Println("...开始执行任务...")

	timeout := 3 * time.Second
	r := New(timeout)

	r.Add(createTask(), createTask(), createTask())
	if err := r.Start(); err != nil {
		switch err {
		case ErrTimeOut:
			log.Println(err)
			os.Exit(1)
		case ErrInterrupt:
			log.Println(err)
			os.Exit(2)
		}
	}
	log.Println("...任务执行结束...")

}

//一个安全的资源池，被管理的资源必须都实现io.Close接口
type Pool struct {
	//m是一个互斥锁，这主要是用来保证在多个goroutine访问资源时，池内的值是安全的。
	m sync.Mutex
	// res字段是一个有缓冲的通道，用来保存共享的资源，这个通道的大小，
	// 在初始化Pool的时候就指定的。注意这个通道的类型是io.Closer接口，
	// 所以实现了这个io.Closer接口的类型都可以作为资源，交给我们的资源池管理。
	res chan io.Closer
	// 一个函数类型，它的作用就是当需要一个新的资源时，可以通过这个函数创建，
	// 也就是说它是生成新资源的，至于如何生成、生成什么资源，是由使用者决定的，
	// 所以这也是这个资源池灵活的设计的地方
	factory func() (io.Closer, error)
	//表示资源池是否被关闭，如果被关闭的话，再访问是会有错误的
	closed bool
}

var ErrPoolClosed = errors.New("资源池已经关闭")

//创建一个资源池
func NewPool(fn func() (io.Closer, error), size uint) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("size的值太小了。")
	}
	return &Pool{
		factory: fn,
		res:     make(chan io.Closer, size),
	}, nil
}

//从资源池里获取一个资源
//Acquire方法可以从资源池获取资源，如果没有资源，则调用factory方法生成一个并返回。
//这里同样使用了select的多路复用，因为这个函数不能阻塞，可以获取到就获取，不能就生成一个。
func (p *Pool) Acquire() (io.Closer, error) {
	select {
	case r, ok := <-p.res:
		log.Println("Acquire:共享资源")
		if !ok {
			return nil, ErrPoolClosed
		}
		return r, nil
	default:
		log.Println("Acquire:新生成资源")
		return p.factory()
	}
}

//关闭资源池，释放资源
func (p *Pool) Close() {
	p.m.Lock()
	defer p.m.Unlock()

	if p.closed {
		return
	}

	//关闭通道，不让写入了
	p.closed = true

	//关闭通道里的资源
	close(p.res)
	for r := range p.res {
		_ = r.Close()
	}
}

func (p *Pool) Release(r io.Closer) {
	//保证该操作和Close方法的操作是安全的
	p.m.Lock()
	defer p.m.Unlock()
	//资源池都关闭了，就省这一个没有释放的资源了，释放即可
	if p.closed {
		_ = r.Close()
		return
	}

	select {
	case p.res <- r:
		log.Println("资源释放到池子了")
	default:
		log.Println("资源池满了，释放这个资源吧")
		_ = r.Close()
	}
}

func createTask() func(int) {
	return func(id int) {
		log.Printf("正在执行任务%d", id)
		time.Sleep(time.Duration(id) * time.Second)
	}
}

// 一个执行者，可以执行任何任务，但是这些任务是限制完成的，//该执行者可以通过发送终止信号终止它
type Runner struct {
	tasks     []func(int)      //要执行的任务
	complete  chan error       //用于通知任务全部完成
	timeout   <-chan time.Time //这些任务在多久内完成
	interrupt chan os.Signal   //可以控制强制终止的信号
}

// 工厂函数New,用于返回我们需要的Runner
// 很快的初始化一个Runnner，它只有一个参数，用来设置这个执行者的超时时间。这个超时时间被我们传递给了time.After函数，
// 这个函数可以在tm时间后，会同伙一个time.Time类型的只能接收的单向通道，来告诉我们已经到时间了。
func New(tm time.Duration) *Runner {
	return &Runner{
		// complete是一个无缓冲通道，也就是同步通道，因为我们要使用它来控制我们整个程序是否终止，所以它必须是同步通道，
		// 要让main goroutine等待，一致要任务完成或者被强制终止。
		complete: make(chan error),
		timeout:  time.After(tm),
		// interrupt是一个有缓冲的通道，这样做是因为，我们可以至少接收到一个操作系统的中断信息，
		// 这样Go runtime在发送这个信号的时候不会被阻塞，如果是无缓冲的通道就会阻塞了
		interrupt: make(chan os.Signal, 1),
	}
}

//将需要执行的任务，添加到Runner里
func (r *Runner) Add(tasks ...func(int)) {
	r.tasks = append(r.tasks, tasks...)
}

//两种错误类型，一个表示因为超时错误，一个表示因为被中断错误。
var ErrTimeOut = errors.New("执行者执行超时")
var ErrInterrupt = errors.New("执行者被中断")

//执行任务，执行的过程中接收到中断信号时，返回中断错误//如果任务全部执行完，还没有接收到中断信号，则返回nil
func (r *Runner) run() error {
	for id, task := range r.tasks {
		if r.isInterrupt() {
			return ErrInterrupt
		}
		task(id)
	}
	return nil
}

//检查是否接收到了中断信号
func (r *Runner) isInterrupt() bool {
	// 基于select的多路复用，select和switch很像，只不过它的每个case都是一个通信操作。那么到底选择哪个case块执行呢？
	// 原则就是哪个case的通信操作可以执行就执行哪个，如果同时有多个可以执行的case，那么就随机选择一个执行。
	select {
	case <-r.interrupt:
		signal.Stop(r.interrupt)
		return true
	default:
		return false
	}
}

//开始执行所有任务，并且监视通道事件
func (r *Runner) Start() error {
	//希望接收哪些系统信号
	//signal.Notify(r.interrupt, os.Interrupt)，这个是表示，如果有系统中断的信号，发给r.interrupt即可。
	signal.Notify(r.interrupt, os.Interrupt)
	go func() {
		r.complete <- r.run()
	}()

	select {
	case err := <-r.complete:
		return err
	case <-r.timeout:
		return ErrTimeOut
	}
}

// 这是Go语言圣经里比较有意义的一个例子，例子是想获取服务端的一个数据，不过这个数据在三个镜像站点上都存在，
// 这三个镜像分散在不同的地理位置，而我们的目的又是想最快的获取到数据。
// 所以这里，我们定义了一个容量为3的通道responses，然后同时发起3个并发goroutine向这三个镜像获取数据，
// 获取到的数据发送到通道responses中，最后我们使用return <-responses返回获取到的第一个数据，
// 也就是最快返回的那个镜像的数据。
//func mirroredQuery() string {
//	responses := make(chan string, 3)
//	go func() { responses <- request("asia.gopl.io") }()
//	go func() { responses <- request("europe.gopl.io") }()
//	go func() { responses <- request("americas.gopl.io") }()
//	return <-responses // return the quickest response
//}
//func request(hostname string) (response string) { /* ... */ }
func incCount() {
	defer wg2.Done()
	for i := 0; i < 2; i++ {
		value := count
		// runtime.Gosched()是让当前goroutine暂停的意思，退回执行队列，让其他等待的goroutine运行，
		// 目的是让我们演示资源竞争的结果更明显.
		// 注意，这里还会牵涉到CPU问题，多核会并行，那么资源竞争的效果更明显。
		runtime.Gosched()
		value++
		count = value
	}
}
func incCountLock() {
	defer wg2.Done()
	for i := 0; i < 2; i++ {
		//这里atomic.LoadInt32和atomic.StoreInt32两个函数，一个读取int32类型变量的值，
		// 一个是修改int32类型变量的值，
		// 这两个都是原子性的操作，Go已经帮助我们在底层使用加锁机制，
		// 保证了共享资源的同步和安全，所以我们可以得到正确的结果，
		// 这时候我们再使用资源竞争检测工具go build -race检查，也不会提示有问题了。
		//atomic包里还有很多原子化的函数可以保证并发下资源同步访问修改的问题，
		// 比如函数atomic.AddInt32可以直接对一个int32类型的变量进行修改，在原值的基础上再增加多少的功能，也是原子性的
		value := atomic.LoadInt32(&count)
		runtime.Gosched()
		value++
		atomic.StoreInt32(&count, value)
	}
}
func incCountMutex() {
	defer wg2.Done()
	for i := 0; i < 2; i++ {
		//实例中，新声明了一个互斥锁mutex sync.Mutex，这个互斥锁有两个方法，
		// 一个是mutex.Lock(),
		// 一个是mutex.Unlock(),这两个之间的区域就是临界区，临界区的代码是安全的。

		//调用mutex.Lock()对有竞争资源的代码加锁，这样当一个goroutine进入这个区域的时候，其他goroutine就进不来了，
		// 只能等待，一直到调用mutex.Unlock() 释放这个锁为止。
		mutex.Lock()
		value := count
		runtime.Gosched()
		value++
		count = value
		mutex.Unlock()
	}
}

var (
	count int32
	wg2   sync.WaitGroup
	mutex sync.Mutex
)

func NewLoginer() Loginer {
	return defaultLogin(0)
}

type Loginer interface {
	Login()
}
type defaultLogin int

func (d defaultLogin) Login() {
	fmt.Println("login in...")
}

type user struct {
	name  string
	email string
}
type admin struct {
	user
	level string
}

func (u user) sayHello() {
	fmt.Println("Hello，i am a user")
}
func (a admin) sayHello() {
	fmt.Println("Hello，i am a admin")
}

type Hello interface {
	hello()
}

func (u user) hello() {
	fmt.Println("Hello，i am a user")
}
func sayHello(h Hello) {
	h.hello()
	fmt.Println("Hello，i am a admin")
}

//需要一个animal接口作为参数
func invoke(a animal) {
	a.printInfo()
}

type animal interface {
	printInfo()
}

type cat int
type dog int

//指针接收者实现animal接口
func (c *cat) printInfo() {
	fmt.Println("a cat")
}

////值接收者实现animal接口
//func (c cat) printInfo() {
//	fmt.Println("a cat")
//}

func (d dog) printInfo() {
	fmt.Println("a dog")
}

func print(a ...interface{}) {
	for _, v := range a {
		fmt.Print(v)
	}
	fmt.Println()
}

// 函数方法声明定义的时候，采用逗号分割，因为时多个返回，还要用括号括起来。
// 返回的值还是使用return 关键字，以逗号分割，和返回的声明的顺序一致。
func add1(a, b int) (int, error) {
	return a + b, nil
}
func (p *person1) modifyPerson1() {
	p.name = "LISI"
}

type person1 struct {
	name string
}

func (p person1) String() string {
	p.name = "lisi"
	return "the person name is " + p.name
}

type person struct {
	age  int
	name string
}
type Names struct {
	slice []string
}

type Duration int64

func add(a, b int) int {
	return a + b
}

func modifySl(names Names) {
	names.slice[2] = "21"

}

//func modifyJim(person person) {
//	person.age = 45
//}
func modifyJim(person *person) {
	person.age = 45
}

func modifyAges(ages map[string]int) {
	ages["zhangsan"] = 10
}
func modifyName(name string) string {
	name = name + name
	return name
}
func modify1(s []int) {
	s[0] = 1
	s = append(s, 999)
	fmt.Println(s)
}
func modify(a [5]int) {
	//传递数组的指针会导致原数组变化
	fmt.Println(a)
	fmt.Println(a[0])
	a[1] = 3
	fmt.Println(a)
}

func modifySlice(slice []int) {
	fmt.Println(&slice)
	slice[1] = 10
}
