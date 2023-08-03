package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/allenliu88/go-utils/utils"
	"github.com/google/uuid"
)

func main() {
	reqs := make([]utils.Request, 0)
	reqs = append(reqs, utils.Request{RequestID: uuid.NewString(), WorkSize: 5})
	reqs = append(reqs, utils.Request{RequestID: uuid.NewString(), WorkSize: 5})
	reqs = append(reqs, utils.Request{RequestID: uuid.NewString(), WorkSize: 5})

	rootCtx := context.Background()
	ctx := context.WithValue(rootCtx, utils.ContextValueKey_WaitGroup, &sync.WaitGroup{})
	ret, err := utils.DoRequests(ctx, reqs, utils.DoJob)

	fmt.Printf("Result: %v, Error: %v\n", ret, err)
	ctx.Value(utils.ContextValueKey_WaitGroup).(*sync.WaitGroup).Wait()
}
