package cfn

import (
	"fmt"
	"time"
	"context"
)

func loader(ctx context.Context, ch chan bool) {
	for {
		select {
			case <- ch:
				fmt.Println("")
				return
			default:
				time.Sleep(500 * time.Millisecond)
				logger.ColorPrintf(ctx, ".")
		}
	}
}