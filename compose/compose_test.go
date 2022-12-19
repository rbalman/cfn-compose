package compose

import (
	"testing"
	"github.com/balmanrawat/cfn-compose/cfn"
	"fmt"
)

func TestValidateComposeConfig(t *testing.T) {
	t.Log("When there are no items")
	{
		var stacks []cfn.Stack
		rs := reverse(stacks)
		if len(rs) != 0 {
			t.Fatal(fmt.Sprintf("Expected reverse stack length to be 0 but got %d", len(rs)))
		}
	}

	t.Log("When there are items")
	{
		stacks := []cfn.Stack{cfn.Stack{StackName: "first"},cfn.Stack{StackName: "second"}, cfn.Stack{StackName: "third"}}
		rs := reverse(stacks)
		j := 0
		for i := len(stacks) - 1; i >= 0 ; i-- {
			sName := stacks[i].StackName
			rsName := rs[j].StackName
			if  sName != rsName {
				t.Fatal(fmt.Sprintf("Expected stacks[%d] and rs[%d] to be equal but got %s and %s", i, j, sName, rsName ))
			}
			j++
		}
	}
}
