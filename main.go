package main

import (
	"fmt"
	"math/rand"

	"github.com/zinx/utils"
)

func main() {
	source := rand.NewSource(int64(10))
	randNumGenetor := rand.New(source)

	for i := 0; i < 10; i++ {
		randomSendInterval := randNumGenetor.Intn(utils.GlobalObject.MaxSendInterval-utils.GlobalObject.MinSendInterval) + utils.GlobalObject.MinSendInterval
		fmt.Println(randomSendInterval)
	}
}
