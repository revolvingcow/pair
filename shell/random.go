package shell

import (
	"math/rand"
	"time"
)

var (
	portRandomizer = rand.New(rand.NewSource(time.Now().UnixNano()))
)
