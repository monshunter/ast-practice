package main

import (
	"fmt"
	"time"

	"github.com/monshunter/ast-practice/cmd/testfirstclassfunc/wait"
)

func main() {
	fmt.Println("test first class function")
}

var stop chan struct{}

type agentManager struct{}
type nodeManager struct{}
type controller struct {
}

func (a agentManager) Ready() bool {
	return true
}

func (n nodeManager) Ready() bool {
	return true
}

func (c controller) Run(stop <-chan struct{}) {
	fmt.Println("Run controller")
}

func (c controller) Stopped() bool {
	return true
}

// type

type RegistryController struct {
	controllers  []controller
	agentManager agentManager
	nodeManager  nodeManager
	Internal     Internal
}

type Internal struct{}

func (i Internal) Run(stop <-chan struct{}) {
	fmt.Println("Run Internal")
}

func (r *RegistryController) waitAgentsReady(stop <-chan struct{}) error {
	return wait.PollUntil(time.Second, func() (done bool, err error) {
		return r.agentManager.Ready(), nil
	}, stop)
}

func (r *RegistryController) waitNodesReady(stop <-chan struct{}) error {
	return wait.PollUntil(time.Second, func() (done bool, err error) {
		return r.nodeManager.Ready(), nil
	}, stop)
}

func (r *RegistryController) run(stop <-chan struct{}) {
	for _, c := range r.controllers {
		go c.Run(stop)
	}
	fmt.Println("go r.Internal.Run(stop)")
	go r.Internal.Run(stop)
	_ = wait.PollUntil(time.Second, func() (done bool, err error) {
		for _, c := range r.controllers {
			if !c.Stopped() {
				return false, nil
			}
		}
		return true, nil
	}, stop)
}

func switchCase1(a int) {
	switch a {
	case 1:
		fmt.Println("a is 1")
	case 2:
		fmt.Println("a is 2")
	default:
		fmt.Println("a is not 1 or 2")
	}
}

func switchCase2(a int) {
	switch a {
	case 1:
		fmt.Println("a is 1")
		fallthrough
	case 2:
		fmt.Println("a is 2")
		fallthrough
	case 3:
		fmt.Println("a is 3")
	default:
		fmt.Println("a is not 1 or 2")
	}
}

func typeswitch(a interface{}) {
	switch a.(type) {
	case int:
		fmt.Println("a is int")
	case string:
		fmt.Println("a is string")
		fmt.Println(a)
	default:
		fmt.Println("a is not int or string")
	}
}

func selectCase() {
	select {
	case <-stop:
		fmt.Println("stop")
	case <-time.After(time.Second):
		fmt.Println("timeout")
	}
}

func selectCase2() {
	fmt.Println("selectCase2")
	stop = make(chan struct{})

	select {
	case <-stop:
		fmt.Println("stop")
	case <-time.After(time.Second):
		fmt.Println("timeout")
	}
}

func selectCase3() {
	fmt.Println("selectCase3")
	select {}
}

func selectCase4() {
	fmt.Println("selectCase4")
	stop = make(chan struct{})
	select {
	case <-stop:
		fmt.Println("stop")
	default:
		fmt.Println("default")
		fmt.Println("default")
	}
}
