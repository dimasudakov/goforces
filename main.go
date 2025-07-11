package main

import (
	"math"
)

// Пример решения задачи https://codeforces.com/contest/2123/problem/C

func main() {
	t := scanT[int]()

	for ; t > 0; t-- {
		solve()
		_ = out.Flush()
	}
}

func solve() {
	var n int
	scan(&n)

	arr := scanArrT[int](n)

	ans := make([]bool, n)

	sufMx := make([]int, n+1)
	sufMx[n] = math.MinInt
	for i := n - 1; i >= 0; i-- {
		sufMx[i] = max(sufMx[i+1], arr[i])
	}

	prefMin := math.MaxInt
	for i := 0; i < n; i++ {
		if prefMin > arr[i] || sufMx[i+1] < arr[i] {
			ans[i] = true
		}
		prefMin = min(prefMin, arr[i])
	}

	for _, val := range ans {
		if val {
			gout("1")
		} else {
			gout("0")
		}
	}
	nl()
}
