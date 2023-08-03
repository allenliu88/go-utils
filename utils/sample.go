package utils

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func DoJob(ctx context.Context, req Request) (*Response, error) {
	wg := ctx.Value(ContextValueKey_WaitGroup)
	wg.(*sync.WaitGroup).Add(1)
	defer wg.(*sync.WaitGroup).Done()

	for i := 0; i < req.WorkSize; i++ {
		select {
		case <-ctx.Done():
			logrus.Errorf("cancel occurred for the job %d of request [%s]\n", i, req.RequestID)
			return nil, nil
		default:
			logrus.Warnf("doing the job %d of request [%s]\n", i, req.RequestID)
			time.Sleep(time.Second * time.Duration(rand.Intn(5))) // Simulating some work
			if i%2 == 1 {
				return nil, fmt.Errorf("error occurred for doing the job %d of request [%s]", i, req.RequestID)
			}

			logrus.Infof("successfully done for the job %d of request [%s]\n", i, req.RequestID)
		}
	}

	logrus.Infof("successfully returned for the request [%s]\n", req.RequestID)
	return &Response{"Success"}, nil
}
