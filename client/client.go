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

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("No args passed in\n")
		os.Exit(1)
	}
	fibParams, powParams := getParams()
	go sendRequest(fibParams, "Fibonacci")
	go sendRequest(powParams, "Pow")
	select {}
}

func sendRequest(params []int, service string) {
	args := Args{}
	var ret Return
	loadBalancerAddress := os.Getenv("LOAD_BALANCER")
	loadBalancer, err := rpc.Dial("tcp", loadBalancerAddress)
	for _, param := range params {
		args.Service = service
		args.Input = param
		if err != nil {
			fmt.Printf("An error occured : %s \n", err)
		}
		done := loadBalancer.Go("LoadBalancer.ServeRequest", args, &ret, nil)
		start := time.Now()
		go waitResult(done, args, &ret, &sync.Mutex{}, start)
	}
}

func waitResult(done *rpc.Call, args Args, p *Return, mutex *sync.Mutex, start time.Time) {
	mutex.Lock()
	done = <-done.Done
	end := time.Now()
	totalResponseTime = end.Sub(start)
	if done.Error != nil {
		fmt.Printf("An error occured %s \n", done.Error)
		mutex.Unlock()
		os.Exit(1)
	} else {
		fmt.Printf("The result of %s(%d) is %d. The actual response time is %v\n", args.Service, args.Input, *p, totalResponseTime/10)
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
