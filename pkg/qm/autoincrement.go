package qm

import (
	"math"
	"math/rand"
	"strconv"
	"sync/atomic"
)

var currentValue int32 = int32(rand.Intn(math.MaxInt32) / 3)

func GetNextVal() string {
	atomic.AddInt32(&currentValue, 1)
	return strconv.Itoa(int(currentValue))
}
