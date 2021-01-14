package main

import (
	"fmt"
	"sync"
)


func setJob() <-chan int {
    var wg sync.WaitGroup
    jobChan := make(chan int, 50)
    wg.Add(1)
    go func(){
        for i := 0; i < 50; i++ {
            jobChan <- i
        }
        close(jobChan)
        wg.Done()
    }()
    wg.Wait()
    return jobChan
}

func main() {
	var wg sync.WaitGroup
    // 创建一个需要处理的数据来源, 假设是个channel
    jobChan := setJob()

    // 将结果存入一个channel
    resChan := make(chan int , 50)

    // 设置一个并发池
	buckets := make(chan bool, 10)
	for job := range jobChan {
		buckets <- true
		wg.Add(1)
		go func(job int) {
            res := 10 * job
            resChan <- res
			<-buckets
			wg.Done()
		}(job)
	}

    wg.Wait()
    close(resChan)
    tmpA := []int{}
	for r := range(resChan) {
		tmpA = append(tmpA, r)
	}
    fmt.Println(tmpA)
    
}