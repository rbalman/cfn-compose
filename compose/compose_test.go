package compose

import (
	"fmt"
	"github.com/rbalman/cfn-compose/cfn"
	"testing"
)

func TestValidateComposeConfig(t *testing.T) {
	t.Log("When there are no items")
	{
		var stacks []cfn.Stack
		rs := reverseStackOrder(stacks)
		if len(rs) != 0 {
			t.Fatal(fmt.Sprintf("Expected reverse stack length to be 0 but got %d", len(rs)))
		}
	}

	t.Log("When there are items")
	{
		stacks := []cfn.Stack{{StackName: "first"}, {StackName: "second"}, {StackName: "third"}}
		rs := reverseStackOrder(stacks)
		j := 0
		for i := len(stacks) - 1; i >= 0; i-- {
			sName := stacks[i].StackName
			rsName := rs[j].StackName
			if sName != rsName {
				t.Fatal(fmt.Sprintf("Expected stacks[%d] and rs[%d] to be equal but got %s and %s", i, j, sName, rsName))
			}
			j++
		}
	}
}
