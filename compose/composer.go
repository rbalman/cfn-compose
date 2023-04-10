package compose

import (
	"context"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/rbalman/cfn-compose/cfn"
	"github.com/rbalman/cfn-compose/config"
	"github.com/rbalman/cfn-compose/libs"
	"github.com/rbalman/cfn-compose/logger"
)

// var colors []string = []string{log.Blue, log.Yellow, log.Green, log.Magenta, log.Cyan}

type Task interface {
	Execute(context.Context) Result
}

type Result struct {
	FlowName string
	Error    error
}

type Composer struct {
	LogLevel         string
	CherryPickedFlow string
	DeployMode       bool
	DryRun           bool
	ConfigFile       string
	WorkersCount     int
}

func (c *Composer) Apply() {
	ctx := context.Background()
	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	cc, err := config.GetComposeConfig(c.ConfigFile)
	if err != nil {
		fmt.Printf("Err: Failed to Parse Compose Config: %s\n", err)
		os.Exit(1)
	}

	err = cc.Validate()
	if err != nil {
		fmt.Printf("Err: Failed While Validating Compose Config: %s\n", err)
		os.Exit(1)
	}

	logger.StartWithLabel(c.LogLevel)

	var flowsMap map[int][]config.Flow
	if c.CherryPickedFlow != "" {
		flowsMap = cherryPickFlow(c.CherryPickedFlow, cc.Flows)
		if len(flowsMap) == 0 {
			fmt.Printf("Err: Cannot find the selected flow: %s in the config\n", c.CherryPickedFlow)
			os.Exit(1)
		}
	} else {
		flowsMap = SortFlows(cc.Flows)
	}

	orders := keys(flowsMap)
	if c.DeployMode {
		sort.Ints(orders)
	} else {
		sort.Sort(sort.Reverse(sort.IntSlice(orders)))
	}

	// Exporting AWS_PROFILE and AWS_REGION got from config
	if val, ok := cc.Vars["AWS_PROFILE"]; ok {
		os.Setenv("AWS_PROFILE", val)
	}

	if val, ok := cc.Vars["AWS_REGION"]; ok {
		os.Setenv("AWS_REGION", val)
	}

	sess, err := libs.GetAWSSession()
	if err != nil {
		logger.Log.Errorf("Failed while creating AWS Session: %s\n", err.Error())
		os.Exit(1)
	}
	cm := cfn.CFNManager{Session: sess}

	cfnTask := make(chan Task)
	resultsChan := make(chan Result)
	numWorkers := getWorkersCount(len(cc.Flows), c.WorkersCount)
	//Generate the worker pool as pre the flow counts
	for i := 0; i < numWorkers; i++ {
		go executeFlow(ctx, cfnTask, resultsChan, i)
	}
	logger.Log.Debugf("TOTAL FLOW COUNT: %d\n", len(cc.Flows))
	logger.Log.Debugf("TOTAL WORKERS SPUN: %d\n", numWorkers)
	//Dispatch Flows based on the Order
	for _, order := range orders {
		flows, ok := flowsMap[order]
		if !ok {
			continue
		}

		for _, flow := range flows {
			cfnTask <- CfnTask{Flow: flow, DryRun: c.DryRun, DeployMode: c.DeployMode, CM: cm}
		}

		logger.Log.Debugf("Dispatched Order: %d, FlowCount: %d.\n", order, len(flows))

		//Wait for dispatched flows
		for i := 0; i < len(flows); i++ {
			//TODO: Add some form of timer for timeout
			r := <-resultsChan
			if r.Error != nil {
				cancelCtx()
				logger.Log.Debugln("Graceful wait for cancelled flows")
				time.Sleep(time.Second * 5)
				logger.Log.Errorf("Compose failed with Error: %s", r.Error)
				return
			}
		}
		logger.Log.Infof("All Flows completed for Order: %d\n\n", order)
	}

	logger.Log.Infoln("Successfully Completed!!")
}

func (c *Composer) PrintConfig() {
	fmt.Println("##########################")
	fmt.Println("# Compose Configuration #")
	fmt.Println("##########################")
	fmt.Printf("ConfigFile: %s\n", c.ConfigFile)
	if c.CherryPickedFlow != "" {
		fmt.Printf("Selected Flow: %s\n", c.CherryPickedFlow)
	}
	fmt.Printf("DryRun: %t\n", c.DryRun)
	fmt.Printf("LogLevel: %s\n", c.LogLevel)
	fmt.Printf("WorkersCount: %d\n", c.WorkersCount)
	fmt.Printf("DeployMode: %t\n\n", c.DeployMode)
}

func SortFlows(flows map[string]config.Flow) map[int][]config.Flow {
	sortedFlows := make(map[int][]config.Flow)
	for name, flow := range flows {
		flow.Name = name

		flows, ok := sortedFlows[flow.Order]
		if ok {
			flows = append(flows, flow)
			sortedFlows[flow.Order] = flows
		} else {
			sortedFlows[flow.Order] = []config.Flow{flow}
		}
	}

	return sortedFlows
}

func VisualizeFlowsMap(flowsMap map[int][]config.Flow) {
	orders := keys(flowsMap)
	sort.Ints(orders)

	for _, order := range orders {
		flows := flowsMap[order]
		fmt.Printf("ORDER: %d\n", order)
		for _, flow := range flows {
			fmt.Printf("  FLOW: %s\n", flow.Name)
			for _, stack := range flow.Stacks {
				fmt.Printf("    Stack: %s\n", stack.StackName)
			}
		}
	}
}

func keys(flowMap map[int][]config.Flow) []int {
	var keys []int
	for key := range flowMap {
		keys = append(keys, key)
	}
	return keys
}

func cherryPickFlow(flowName string, flows map[string]config.Flow) map[int][]config.Flow {
	cherryPickedFlow := make(map[int][]config.Flow)
	for name, flow := range flows {
		if name == flowName {
			flow.Name = name
			cherryPickedFlow[flow.Order] = []config.Flow{flow}
		}
	}
	return cherryPickedFlow
}

func reverseStackOrder(stacks []cfn.Stack) []cfn.Stack {
	var rs []cfn.Stack
	if len(stacks) == 0 {
		return rs
	}
	for i := len(stacks) - 1; i >= 0; i-- {
		rs = append(rs, stacks[i])
	}
	return rs
}

func executeFlow(ctx context.Context, taskC chan Task, resultsChan chan Result, workerId int) {
	defer func() {
		logger.Log.Debugf("Worker: %d exiting...\n", workerId)
	}()

	for {
		select {
		case task := <-taskC:
			time.Sleep(time.Millisecond * 500)
			result := task.Execute(ctx)
			resultsChan <- result

		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				logger.Log.DebugCtxf(ctx, "Cancel signal received Worker: %d, Info: %s\n", workerId, err)
			}
			return
		}
	}
}

func getWorkersCount(flowsCount, countFromFlag int) int {
	if countFromFlag <= 0 || countFromFlag > flowsCount {
		return flowsCount
	}
	return countFromFlag
}
