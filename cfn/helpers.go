package cfn

import (
	"fmt"
	"time"
)

func loader(ch chan bool) {
	for {
		select {
			case <- ch:
				fmt.Println("")
				return
			default:
				time.Sleep(500 * time.Millisecond)
				fmt.Printf(".")
		}
	}
}