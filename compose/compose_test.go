package compose

import (
	"fmt"
	"testing"

	"github.com/rbalman/cfn-compose/cfn"
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

func TestGetWorkersCount(t *testing.T) {
	testCases := []struct {
		flowsCount     int
		countFromFlag  int
		expectedResult int
	}{
		{0, 5, 0},
		{10, -5, 10},
		{10, 15, 10},
		{5, 3, 3},
		{8, 0, 8},
		{7, 7, 7},
	}

	for _, testCase := range testCases {
		actualResult := getWorkersCount(testCase.flowsCount, testCase.countFromFlag)
		if actualResult != testCase.expectedResult {
			t.Errorf("Test case failed: expected %d but got %d (flowsCount=%d, countFromFlag=%d)", testCase.expectedResult, actualResult, testCase.flowsCount, testCase.countFromFlag)
		}
	}
}
