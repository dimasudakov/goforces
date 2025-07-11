package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	in  *bufio.Reader
	out *bufio.Writer
)

func init() {
	in = bufio.NewReader(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
}

func scan(a ...any) {
	_, err := fmt.Fscan(in, a...)
	if err != nil {
		panic(err)
	}
}

func scanT[T any]() T {
	var res T
	scan(&res)
	return res
}

func scanArrT[T any](n int) []T {
	res := make([]T, n)
	for i := 0; i < n; i++ {
		scan(&res[i])
	}
	return res
}

func gout(val ...any) {
	_, err := fmt.Fprint(out, val...)
	if err != nil {
		panic(err)
	}
}

func nl() {
	gout("\n")
}

func goutArrT[T any](arr []T) {
	for i := 0; i < len(arr); i++ {
		gout(arr[i])
		if i < len(arr)-1 {
			gout(" ")
		}
	}
	gout("\n")
}
