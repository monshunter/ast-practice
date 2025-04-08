package testdata

import (
	"fmt"
	"time"
)

func Demo() {
	fmt.Println("start")
	if true {
		fmt.Println("true")
	} else {
		fmt.Println("false")
		if 1+2 > 3 {
			fmt.Println("1+2>3")
		} else if 1+2 == 3 {
			fmt.Println("1+2==3")
		} else {
			fmt.Println("1+2<=3")
		}
	}
	fmt.Println("middle")
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}
	fmt.Println("middle2")
	switch 1 + 2 {
	case 3:
		fmt.Println("1+2=3")
	case 4:
		fmt.Println("1+2=4")
	}

	fmt.Println("end")
}

func Demo2() {
	fmt.Println("start2")
	Demo()
	select {}
	fmt.Println("middle3")
	switch 1 + 2 {
	case 3:
		fmt.Println("1+2=3")
	case 4:
		fmt.Println("1+2=4")
	}
	fmt.Println("middle4")
	var x any = 1
	switch x.(type) {
	case int:
		fmt.Println("int")
	case string:
		fmt.Println("string")
	}

	fmt.Println("middle5")
	type T struct {
		x int
		y int
	}

	var t T = T{x: 1, y: 2}

	fmt.Println("middle6")
	const a = 1
	const b = 2

	const (
		c = 3
		d = 4
	)

	handle := func() {
		fmt.Println("handle")
	}

	handle()

	var dfs func(int)
	dfs = func(n int) {
		if n == 0 {
			return
		}
		dfs(n - 1)
	}
	dfs(10)

	fmt.Println("start2")
	defer func() {
		fmt.Println("defer")
	}()
	fmt.Println("end2")

	fmt.Println("start3")
	go func() {
		fmt.Println("go")
	}()
	fmt.Println("end3")

	select {}

	returnDemo()
	deferDemo()
	goDemo()
	switchDemo()
	typeSwitchDemo()
	closureDemo()
	handleDemo()
	handleDemo2()
	returnDemo2()
	returnDemo3()
	returnDemo4()
	returnDemo5()
	returnDemo6()
}

func returnDemo() int {
	return 1
}

func deferDemo() {
	defer func() {
		fmt.Println("defer")
	}()
	fmt.Println("end")
	defer fmt.Println("defer2")
}

func goDemo() {
	go func() {
		fmt.Println("go")
	}()
	go func() {
		fmt.Println("go2")
	}()
	go fmt.Println("go3")
}

func switchDemo() {
	switch 1 + 2 {
	case 3:
		fmt.Println("1+2=3")
	default:
		fmt.Println("1+2!=3")
	}
}

func selectDemo() {
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("1 second")
	default:
		fmt.Println("default")
	}
}

func typeSwitchDemo() {
	var x any = 1
	switch x.(type) {
	case int:
		fmt.Println("int")
	case string:
		fmt.Println("string")
	}
}

func closureDemo() {
	func() {
		fmt.Println("closure")
	}()
}

func elseDemo() {
	if 1+1 > 2 {
		fmt.Println("true")
	} else if true {
		fmt.Println("false")
	} else {
		fmt.Println("else")
	}

}

type handle func(i int) int

func handleDemo() handle {
	return func(i int) int {
		return i + 1
	}
}

func handleDemo2() handle {
	h := func(i int) int {
		return i * 2
	}
	return h
}

func returnDemo2() (int, int) {
	return 1, 2
}

func returnDemo3() (string, int) {
	return "hello", 1
}

func returnDemo4() (handle, int, handle) {
	return handleDemo(), 1, handleDemo2()
}

func returnDemo5() (handle, int, handle) {
	h1 := handleDemo()
	h2 := handleDemo2()
	return h1, 1, h2
}

func returnDemo6() (handle, int, handle) {
	return func(i int) int {
			return i + 1
		}, 1, func(i int) int {
			return i * 2
		}
}

func typeDeclDemo() int {
	s := struct {
		x int
		y int
	}{
		x: 1,
		y: 2,
	}
	return s.x + s.y
}

func forRangeDemo() {
	for i, v := range []int{1, 2, 3} {
		fmt.Println(i, v)
	}
}

func forRangeDemo2() {
	for i, v := range map[string]int{"a": 1, "b": 2, "c": 3} {
		fmt.Println(i, v)
	}
}
