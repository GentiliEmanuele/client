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
	defer func(loadBalancer *rpc.Client) {
		_ = loadBalancer.Close()
	}(loadBalancer)
	for _, param := range params {
		args.Service = service
		args.Input = param
		if err != nil {
			fmt.Printf("An error occured : %s \n", err)
		}
		done := loadBalancer.Go("LoadBalancer.ServeRequest", args, &ret, nil)
		done = <-done.Done
		if done.Error != nil {
			fmt.Printf("An error occured : %s \n", done.Error)
		}
		fmt.Printf("The result of %s(%d) is %d\n", args.Service, args.Input, ret)
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
