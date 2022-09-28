package actions

import "math/rand"

// helpers
const charsAll = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const charsAbc = "abcdefghijklmnopqrstuvwxyz"

// returns random number in range [low,up)
func RandomRange(rnd *rand.Rand, bottom, up int64) int64 {
	return rnd.Int63n(up-bottom) + bottom
}

// returns random string length n
func RandomString(rnd *rand.Rand, n int64, source string) string {
	var letters = []rune(source)
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rnd.Intn(len(letters))]
	}
	return string(s)
}

func RandomChoice[T any](rnd *rand.Rand, list []T) T {
	return list[rnd.Intn(len(list))]
}

func RandomChoiceMap[K comparable, V any](rnd *rand.Rand, m map[K]V) (K, V) {
	var firstK K
	var firstV V
	n := rand.Intn(len(m))
	i := 0
	for k, v := range m {
		if i == 0 {
			firstK, firstV = k, v
		}
		if i == n {
			return k, v
		}
		i++
	}
	// this line only for compilation
	// it will never executed
	return firstK, firstV
}

func RandomSublist[T any](rnd *rand.Rand, list []T) []T {
	if len(list) == 0 {
		return []T{}
	}
	if len(list) == 1 {
		return []T{list[0]}
	}
	// random indexes to choose
	ids := make([]int, len(list))
	for i := range list {
		ids[i] = i
	}
	rnd.Shuffle(len(ids), func(i, j int) {
		ids[i], ids[j] = ids[j], ids[i]
	})
	n := int(RandomRange(rnd, 1, int64(len(list)+1)))
	result := make([]T, n)
	for i := 0; i < n; i++ {
		result[i] = list[ids[i]]
	}
	return result
}
