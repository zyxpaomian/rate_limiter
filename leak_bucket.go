package main

import (
    "fmt"
    "math"
    "time"
)

// 漏桶限流器
type BucketLimit struct {
    rate       float64 //漏桶中水的漏出速率, 即每秒流多少水
    bucketSize float64 //漏桶最多能装的水大小
    lastAccessTime   time.Time //上次访问时间
    curWater   float64 //当前桶里面的水
}

func NewBucketLimit(rate float64, bucketSize int64) *BucketLimit {
    return &BucketLimit{
        bucketSize: float64(bucketSize),
        rate:       rate,
        curWater:   0,
    }
}

func (b *BucketLimit) AllowControl() bool {
    now := time.Now()
    pastTime := now.Sub(b.lastAccessTime)

    // 当前剩余水量，当前水量减去距离上次访问的流出水量，如果流完了，即剩余水量为0
    b.curWater = math.Max(0, b.curWater - float64(pastTime) * b.rate)

    b.lastAccessTime = now

    // 当前水量必须小于桶的总量，不然则流出了
    if b.curWater < b.bucketSize {
        b.curWater = b.curWater + 1
        return true
    }
    return false
}

func main() {
    // 创建一个流出速率为1qps,桶的总量为2的限流器
    limit := NewBucketLimit(1, 2)

    // 在桶里放入1000滴水
    for i := 0; i < 1000; i++ {
        allow := limit.AllowControl()
        if allow {
            fmt.Printf("第%d滴水, 顺利流出\n", i)
            continue
        } else {
            fmt.Printf("第%d滴水, 溢出丢弃\n", i)
            time.Sleep(time.Millisecond * 100)
        }
    }
}