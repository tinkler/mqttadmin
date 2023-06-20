package qm

import (
	"math"
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"
)

var currentValue int32 = int32(rand.Intn(math.MaxInt32) / 3)

var currentWokerID = strconv.FormatInt(time.Now().Unix(), 10)

func GetNextVal() string {
	atomic.AddInt32(&currentValue, 1)
	return strconv.Itoa(int(currentValue))
}
