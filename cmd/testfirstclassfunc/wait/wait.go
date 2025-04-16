package wait

import (
	"fmt"
	"time"
)

type ConditionFunc func() (done bool, err error)

func PollUntil(interval time.Duration, condition ConditionFunc, stopCh <-chan struct{}) error {
	fmt.Println("run PollUntil")
	return nil
}
