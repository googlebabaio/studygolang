package zujian

import (
	"encoding/binary"
	"io"
	"math/rand"
	"sort"
)

func ArraySource(a ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, v := range a {
			out <- v
		}
		close(out)
	}()
	return out
}

func InMemSort(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		//读取到内存中
		a := []int{}
		for v := range in {
			a = append(a, v)
		}

		//内存中排序
		sort.Ints(a)

		//将排好序的值，再扔回给 channel
		for _, v := range a {
			out <- v
		}
		close(out)
	}()

	return out
}

/**
合并两个输入并输出
 */
func Merge(in1, in2 <-chan int) <-chan int {

	out := make(chan int)

	go func() {

		v1, ok1 := <-in1
		v2, ok2 := <-in2

		for ok1 || ok2 {
			if !ok2 || (ok1 && v1 <= v2) {
				out <- v1
				v1, ok1 = <-in1
			} else {
				out <- v2
				v2, ok2 = <-in2
			}
		}
		close(out)

	}()

	return out
}


func ReadSource(reader io.Reader) <-chan int {

	out := make(chan int)

	go func() {
		buffer := make([]byte, 8)
		for {
			n, err := reader.Read(buffer)

			if err != nil {
				break
			}

			if n > 0 {
				v := int(binary.BigEndian.Uint64(buffer))
				out <- v
			}
		}
		close(out)
	}()

	return out
}

func WirteSink(write io.Writer, in <-chan int) {

	for v := range in {
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, uint64(v))
		write.Write(buffer)
	}
}

func RandomSource(count int) <-chan int {
	out := make(chan int)

	go func() {
		for i := 0; i < count; i++ {
			out <- rand.Int()
		}
		close(out)
	}()

	return out
}

func MergeN(inputs...chan int)<-chan int  {

	if len(inputs)==1{
		return inputs[0]
	}

	m:=len(inputs)/2

	return Merge(MergeN(inputs[:m]...),MergeN(inputs[m:]...))

}