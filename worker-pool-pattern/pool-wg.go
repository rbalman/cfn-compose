package pool

import (
	"fmt"
	"time"
	"sync"
)

func main(){
	fmt.Println("Starting the Master")
	jobCounts := []int{3,4,2}
	var wg sync.WaitGroup
	for idx, count := range jobCounts {
		// for w := 0 ; w < count ; w++ {
		fmt.Printf("Sending work to Job: %d, Worker: %d\n", idx, count)
		wg.Add(1)
		go func() {
			defer wg.Done()
			work(idx, count)
		}()
		// }
		// fmt.Println()
	}
	wg.Wait()
}

func work(job, workerId int){
	fmt.Printf("Job: %d, Worker: %d in progress\n", job, workerId)
	time.Sleep(time.Second*2)
}