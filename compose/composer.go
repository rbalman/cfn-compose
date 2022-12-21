package compose

import (
	"github.com/balmanrawat/cfn-compose/cfn"
	"github.com/balmanrawat/cfn-compose/logger"
	"github.com/balmanrawat/cfn-compose/libs"
	"github.com/balmanrawat/cfn-compose/config"
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"time"
)

// var colors []string = []string{log.Blue, log.Yellow, log.Green, log.Magenta, log.Cyan}

type Work struct {
	Flow        config.Flow
	DryRun     bool
	DeployMode bool
}

type Result struct {
	FlowName string
	Error   error
}

type Composer struct {
	Config config.ComposeConfig
	LogLevel string
	CherryPickedFlow string
	DeployMode bool
	DryRun bool
	ConfigFile string
}

func (c *Composer)Apply() {
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

	c.Config = cc
	logger.StartWithLabel(c.LogLevel)

	// Exporting AWS_PROFILE and AWS_REGION got from config
	if val, ok := cc.Vars["AWS_PROFILE"]; ok {
		os.Setenv("AWS_PROFILE", val)
	}

	if val, ok := cc.Vars["AWS_REGION"]; ok {
		os.Setenv("AWS_REGION", val)
	}

	var flowsMap map[int][]config.Flow
	if c.CherryPickedFlow != "" {
		flowsMap = cherryPickFlow(c.CherryPickedFlow, c.Config.Flows)
		if len(flowsMap) == 0 {
			fmt.Printf("Err: Cannot find the selected flow: %s in the config\n", c.CherryPickedFlow)
			os.Exit(1)
		}
	}else{
		flowsMap = SortFlows(c.Config.Flows)
	}
	
	orders := keys(flowsMap)
	if c.DeployMode {
		sort.Ints(orders)
	}else{
		sort.Sort(sort.Reverse(sort.IntSlice(orders)))
	}

	workChan := make(chan Work)
	resultsChan := make(chan Result)
	//Generate the worker pool as pre the flow counts
	for i := 0; i < len(c.Config.Flows); i++ {
		go ExecuteFlow(ctx, workChan, resultsChan, i)
	}
	logger.Log.Debugf("TOTAL FLOW COUNT: %d\n", len(c.Config.Flows))
	
	//Dispatch Flows in order
	for _, order := range orders {
		flows, ok := flowsMap[order]
		if !ok {
			continue
		}

		for _, flow := range flows {
			workChan <- Work{Flow: flow, DryRun: c.DryRun, DeployMode: c.DeployMode}
		}

		logger.Log.Debugf("Dispatched Order: %d, FlowCount: %d.\n", order, len(flows))

		//wait for flows in each order to complete
		for i := 0; i < len(flows); i++ {
			r := <-resultsChan
			if r.Error != nil {
				cancelCtx()
				logger.Log.Debugln("Graceful wait for cancelled flows")
				time.Sleep(time.Second * 5)
				logger.Log.Errorf("compose failed with Error: %s", r.Error)
				return
			}
		}
		logger.Log.Infof("All Flows completed for Order: %d\n\n", order)
	}

	time.Sleep(time.Second * 2)
	logger.Log.Infoln("Successfully Completed!!")
}

func SortFlows(flows map[string]config.Flow) (map[int][]config.Flow) {
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

func PrintFlowsMap(flowsMap map[int][]config.Flow) () {
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

func cherryPickFlow(flowName string, flows map[string]config.Flow) (map[int][]config.Flow) {
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
	for i := len(stacks) - 1 ; i >= 0; i-- {
		rs = append(rs, stacks[i])
	}
	return rs
}

func ExecuteFlow(ctx context.Context, workChan chan Work, resultsChan chan Result, workerId int) {
	defer func() {
		logger.Log.Debugf("Worker: %d exiting...\n", workerId)
	}()

	for {
		select {
		case work := <-workChan:
			//sleeping from readability
			time.Sleep(time.Millisecond * 500)
			name := work.Flow.Name
			flow := work.Flow
			dryRun := work.DryRun
			deployMode := work.DeployMode
			ctx := context.WithValue(ctx, "flow", name)
			ctx = context.WithValue(ctx, "order", flow.Order)

			sess, err := libs.GetAWSSession()
			if err != nil {
				logger.Log.Errorf("Failed while creating AWS Session: %s\n", err.Error())
				os.Exit(1)
			}
			cm := cfn.CFNManager{Session: sess}

			var stacks []cfn.Stack
			if deployMode {
				stacks = flow.Stacks
			}else{
				stacks = reverseStackOrder(flow.Stacks)
			}

			for _, stack := range stacks{
				ctx := context.WithValue(ctx, "stack", stack.StackName)
				var err error
				if dryRun {
					if deployMode{
						err = stack.ApplyDryRun(ctx, cm)
					}else{
						err = stack.DestoryDryRun(ctx, cm)
					}
				} else {
					if deployMode{
						err = stack.ApplyChanges(ctx, cm)
					}else{
						err = stack.Destroy(ctx, cm)
					}
				}

				if err != nil {
					errStr := fmt.Sprintf("[FLOW: %s] [STACK: %s]. Error: %s\n", name, stack.StackName, err)
					logger.Log.Infoln(errStr)
					resultsChan <- Result{
						Error:   errors.New(errStr),
						FlowName: name,
					}
					break
				}
			}

			resultsChan <- Result{FlowName: name}

		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				logger.Log.DebugCtxf(ctx, "Cancel signal received Worker: %d, Info: %s\n", workerId, err)
			}
			return
		}
	}
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
	fmt.Printf("DeployMode: %t\n\n", c.DeployMode)
}
