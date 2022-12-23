package compose

import (
	"context"
	"errors"
	"fmt"
	"github.com/rbalman/cfn-compose/cfn"
	"github.com/rbalman/cfn-compose/config"
	"github.com/rbalman/cfn-compose/logger"
)

type CfnTask struct {
	Flow       config.Flow
	DryRun     bool
	DeployMode bool
<<<<<<< HEAD
	CM         cfn.CFNManager
=======
	CM cfn.CFNManager
>>>>>>> 3d62c37a9f9b502078c157129367e517cb2b4a00
}

func (ct CfnTask) Execute(ctx context.Context) Result {
	name := ct.Flow.Name
	flow := ct.Flow
	dryRun := ct.DryRun
	deployMode := ct.DeployMode
	ctx = context.WithValue(ctx, "flow", name)
	ctx = context.WithValue(ctx, "order", flow.Order)

	var stacks []cfn.Stack
	if deployMode {
		stacks = flow.Stacks
<<<<<<< HEAD
	} else {
		stacks = reverseStackOrder(flow.Stacks)
	}

	for _, stack := range stacks {
		ctx := context.WithValue(ctx, "stack", stack.StackName)
		var err error
		if dryRun {
			if deployMode {
				err = stack.ApplyDryRun(ctx, ct.CM)
			} else {
				err = stack.DestroyDryRun(ctx, ct.CM)
			}
		} else {
			if deployMode {
				err = stack.ApplyChanges(ctx, ct.CM)
			} else {
=======
	}else{
		stacks = reverseStackOrder(flow.Stacks)
	}

	for _, stack := range stacks{
		ctx := context.WithValue(ctx, "stack", stack.StackName)
		var err error
		if dryRun {
			if deployMode{
				err = stack.ApplyDryRun(ctx, ct.CM)
			}else{
				err = stack.DestoryDryRun(ctx, ct.CM)
			}
		} else {
			if deployMode{
				err = stack.ApplyChanges(ctx, ct.CM)
			}else{
>>>>>>> 3d62c37a9f9b502078c157129367e517cb2b4a00
				err = stack.Destroy(ctx, ct.CM)
			}
		}

		if err != nil {
			errStr := fmt.Sprintf("[FLOW: %s] [STACK: %s]. Error: %s\n", name, stack.StackName, err)
			logger.Log.Infoln(errStr)
			return Result{
<<<<<<< HEAD
				Error:    errors.New(errStr),
=======
				Error:   errors.New(errStr),
>>>>>>> 3d62c37a9f9b502078c157129367e517cb2b4a00
				FlowName: name,
			}
		}
	}
	return Result{FlowName: name}
}
