package main

import (
	"fmt"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Args struct {
	Input   int
	Service string
}

type Return int

var totalResponseTime time.Duration
var nReq int

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("No args passed in\n")
		os.Exit(1)
	}
	mutex := sync.Mutex{}
	fibParams, powParams := getParams()
	nReq = len(fibParams) + len(powParams)
	go sendRequest(fibParams, "Fibonacci", &mutex)
	go sendRequest(powParams, "Pow", &mutex)
	select {}
}

func sendRequest(params []int, service string, mutex *sync.Mutex) {
	loadBalancerAddress := os.Getenv("LOAD_BALANCER")
	loadBalancer, err := rpc.Dial("tcp", loadBalancerAddress)
	for _, param := range params {
		mutex.Lock()
		args := Args{}
		var ret Return
		args.Service = service
		args.Input = param
		if err != nil {
			fmt.Printf("An error occured : %s \n", err)
		}
		start := time.Now()
		done := loadBalancer.Go("LoadBalancer.ServeRequest", args, &ret, nil)
		go waitResult(done, args, &ret, mutex, start)
		mutex.Unlock()
	}
}

func waitResult(done *rpc.Call, args Args, p *Return, mutex *sync.Mutex, start time.Time) {
	done = <-done.Done
	mutex.Lock()
	end := time.Now()
	totalResponseTime = end.Sub(start)
	if done.Error != nil {
		fmt.Printf("An error occured %s \n", done.Error)
		mutex.Unlock()
		os.Exit(1)
	} else {
		fmt.Printf("The result of %s(%d) is %d \n", args.Service, args.Input, *p)
		nReq--
		if nReq == 0 {
			fmt.Printf("The total response time is %v \n", totalResponseTime)
		}
		totalResponseTime = time.Duration(0)
		mutex.Unlock()
	}
}

func getParams() ([]int, []int) {
	fibParams := make([]int, 0)
	powParams := make([]int, 0)
	for i, arg := range os.Args {
		if strings.Compare(arg, "fib") == 0 {
			n, _ := strconv.Atoi(os.Args[i+1])
			fibParams = append(fibParams, n)
		}
		if strings.Compare(arg, "pow") == 0 {
			n, _ := strconv.Atoi(os.Args[i+1])
			powParams = append(powParams, n)
		}
	}
	return fibParams, powParams
}
