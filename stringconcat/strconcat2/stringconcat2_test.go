package strconcat2

import "testing"

//+号拼接不再具有优势，因为string是不可变的，每次拼接都会生成一个新的string,也就是会进行一次内存分配，
//我们现在是10个大小的切片，每次操作要进行9次进行分配，占用内存，所以每次操作时间都比较长，自然性能就低下。

//+号和我们上面分析得一样，这次是99次内存分配，性能体验越来越差，在后面的测试中，会排除掉。
//
//fmt和bufrer已经的性能也没有提升，继续走低。剩下比较坚挺的是Join和Builder。

//1000个字符串整体和100个字符串的时候差不多，表现好的还是Join和Builder。这两个方法的使用侧重点有些不一样，
//如果有现成的数组、切片那么可以直接使用Join,但是如果没有，并且追求灵活性拼接，还是选择Builder。
//Join还是定位于有现成切片、数组的（毕竟拼接成数组也要时间），并且使用固定方式进行分解的，比如逗号、空格等，局限比较大。

//+ 连接适用于短小的、常量字符串（明确的，非变量），因为编译器会给我们优化。
//Join是比较统一的拼接，不太灵活
//fmt和buffer基本上不推荐
//builder从性能和灵活性上，都是上佳的选择。
func BenchmarkStringPlus(b *testing.B) {
	p := initStrings(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringPlus(p)
	}
}

func BenchmarkStringFmt(b *testing.B) {
	p := initStringi(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringFmt(p)
	}
}

func BenchmarkStringJoin(b *testing.B) {
	p := initStrings(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringJoin(p)
	}
}

func BenchmarkStringBuffer(b *testing.B) {
	p := initStrings(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringBuffer(p)
	}
}

func BenchmarkStringBuilder(b *testing.B) {
	p := initStrings(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringBuilder(p)
	}
}

func BenchmarkStringBuilder1(b *testing.B) {
	p := initStrings(1000)
	b.ResetTimer()
	cap := 1000 * len(BLOG)
	for i := 0; i < b.N; i++ {
		StringBuilder1(p, cap)
	}
}
