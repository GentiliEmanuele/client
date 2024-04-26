package main

import (
	"fmt"
	"net/rpc"
	"os"
	"strconv"
)

type Args struct {
	Input   int
	Service string
}

type Return int

func main() {
	args := Args{}
	var ret Return
	loadBalancerAddress := os.Getenv("LOAD_BALANCER")
	if len(os.Args) < 3 {
		fmt.Printf("No args passed in\n")
		os.Exit(1)
	}
	args.Service = os.Args[1]
	args.Input, _ = strconv.Atoi(os.Args[2])
	loadBalancer, err := rpc.Dial("tcp", loadBalancerAddress)
	if err != nil {
		fmt.Printf("An error occured : %s \n", err)
	}
	defer func(loadBalancer *rpc.Client) {
		_ = loadBalancer.Close()
	}(loadBalancer)
	done := loadBalancer.Go("LoadBalancer.ServeRequest", args, &ret, nil)
	done = <-done.Done
	if done.Error != nil {
		fmt.Printf("An error occured : %s \n", done.Error)
	}
	fmt.Printf("The result of %s(%d) is %d\n", args.Service, args.Input, ret)
}
