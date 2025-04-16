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
