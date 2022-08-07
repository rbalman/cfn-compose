package cfn

import (
	"context"
	"fmt"
	"time"
)

func loader(ctx context.Context, ch chan bool) {
	for {
		select {
		case <-ch:
			fmt.Println("")
			return
		default:
			time.Sleep(500 * time.Millisecond)
			fmt.Print(".")
		}
	}
}
