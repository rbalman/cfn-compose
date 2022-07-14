package pool

import (
	"fmt"
)

type Worker struct {

}

func main(){
	fmt.Println("Starting the Master")
	workch := make(chan string)
	results := make(chan string, 10)

	worker := 10
	for i := 0; i < worker ; i++ {
		go work(workch, results, i)
	}

	jobCounts := []int{3,20,2}
	for idx, count := range jobCounts {
		for i := 0; i < count; i++ {
			fmt.Printf("Order: %d, Job: %d\n", idx, i)
			workch <- fmt.Sprintf("%d/%d", idx, i)
		}

		//wait for each order
		for i := 0; i < count ; i++ {
			select {
				case result := <- results: {
					fmt.Println(result)
				}
			}
		}

		fmt.Printf("Order %d completed\n\n", idx)

	}
	
	fmt.Println("Sending close signal for workers")
	close(workch)
}

func work(workch, results chan string, workerId int){
	for info := range workch {
		fmt.Printf("Order/JobId: %s in progress..\n",info)
		results <- fmt.Sprintf("Completed: Order/JobId: %s", info)
	}
	fmt.Sprintf("Worker %d Exitted..", workerId)
}