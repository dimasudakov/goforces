package main

var (
	MOD int64
	inv []int64
)

func setMOD(mod int64) {
	MOD = mod
}

func calcInverses(n int64) {
	inv = make([]int64, n+1)
	inv[1] = 1
	for i := int64(2); i <= n; i++ {
		inv[i] = MOD - (MOD/i)*inv[MOD%i]%MOD
	}
}

type Mint int64

func NewMint(val int64) Mint {
	val = val % MOD
	if val < 0 {
		val += MOD
	}
	return Mint(val)
}

func (m Mint) mul(a any) Mint {
	aInt64 := getInt64(a) % MOD
	return Mint((int64(m) * aInt64) % MOD)
}

func (m Mint) div(a any) Mint {
	aInt64 := getInt64(a)
	if aInt64 <= 0 || aInt64 >= int64(len(inv)) {
		panic("invalid argument a")
	}
	return m.mul(Mint(inv[aInt64]))
}

func (m Mint) add(a any) Mint {
	aInt64 := getInt64(a) % MOD
	res := int64(m) + aInt64
	if res < 0 {
		res += MOD
	}
	return Mint(res % MOD)
}

func (m Mint) pow(p any) Mint {
	pn := getInt64(p)
	a := m
	res := NewMint(1)
	for pn > 0 {
		if pn&1 != 0 {
			res = res.mul(a)
		}
		a = a.mul(a)
		pn >>= 1
	}
	return m
}

func (m Mint) int64() int64 {
	return int64(m)
}

func (m Mint) int() int {
	return int(m)
}

func getInt64(a any) int64 {
	var val int64
	switch v := a.(type) {
	case Mint:
		val = int64(v)
	case int64:
		val = v
	case int:
		val = int64(v)
	case int32:
		val = int64(v)
	default:
		panic("unsupported type")
	}
	return val
}
