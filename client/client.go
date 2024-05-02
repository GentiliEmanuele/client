package main

import (
	"fmt"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

type Args struct {
	Input   int
	Service string
}

type Return int

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
		go waitResult(done, args, &ret)
	}
}

func waitResult(done *rpc.Call, args Args, p *Return) {
	done = <-done.Done
	if done.Error != nil {
		fmt.Printf("An error occured %s \n", done.Error)
		os.Exit(1)
	} else {
		fmt.Printf("The result of %s(%d) is %d\n", args.Service, args.Input, *p)
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
