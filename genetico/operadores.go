package genetico

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func Inversion(a []uint8, wg *sync.WaitGroup) {
	rand.Seed(time.Now().UnixNano())
	start := rand.Intn(len(a))
	size := rand.Intn(len(a) - 2)
	i := 0
	var tmp uint8
	var j, k int
	for i <= (size / 2) {
		j = start + i
		if j > len(a)-1 {
			j = j - len(a)
		}
		k = start + size - i
		if k > len(a)-1 {
			k = k - len(a)
		}
		tmp = a[k]
		a[k] = a[j]
		a[j] = tmp
		i++
	}
	wg.Done()
}

func Intercambio(a []uint8, wg *sync.WaitGroup) {
	rand.Seed(time.Now().UnixNano())
	size := rand.Intn((len(a) / 2)) + 1

	rand.Seed(time.Now().UnixNano())
	pos1 := rand.Intn((len(a)) - (size * 2) + 1)

	rand.Seed(time.Now().UnixNano())
	pos2 := rand.Intn(len(a)-pos1-(size*2)+1) + pos1 + size

	var tmp uint8
	for i := 0; i < size; i++ {
		tmp = a[pos1+i]
		a[pos1+i] = a[pos2+i]
		a[pos2+i] = tmp
	}
	wg.Done()
}

func Cruza(c1, c2 []uint8) ([]uint8, []uint8) {
	rand.Seed(time.Now().UnixNano())
	slicePoint := rand.Intn(23) + 1
	var s1, s2 string
	for _, v := range c1 {
		s1 += fmt.Sprintf("%08b", v)
	}
	for _, v := range c2 {
		s2 += fmt.Sprintf("%08b", v)
	}
	t1 := fmt.Sprintf("%s%s", s1[:slicePoint], s2[slicePoint:])
	t2 := fmt.Sprintf("%s%s", s2[:slicePoint], s1[slicePoint:])
	r1 := make([]uint8, len(c1))
	r2 := make([]uint8, len(c2))
	for i := 0; i < len(t1)/8; i++ {
		if v, err := strconv.ParseUint(t1[i*8:(i+1)*8], 2, 8); err == nil {
			r1[i] = uint8(v)
		}
		if v, err := strconv.ParseUint(t2[i*8:(i+1)*8], 2, 8); err == nil {
			r2[i] = uint8(v)
		}
	}
	return r1, r2
}
