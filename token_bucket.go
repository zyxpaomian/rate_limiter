package main
 
import (
    "fmt"
    "time"
    "math"
)
 
// 令牌桶限流器
type BucketLimit struct {
    rate       float64 //令牌桶放令牌的速率
    bucketSize float64 //令牌桶最多能装的令牌数量
    lastAccessTime   time.Time //上次访问时间
    curTokenNum  float64 //桶里当前的令牌数量
}

func NewBucketLimit(rate float64, bucketSize int64) *BucketLimit {
    return &BucketLimit{
        bucketSize: float64(bucketSize),
        rate:       rate,
        curTokenNum:   0,
    }
}

func (b *BucketLimit) AllowControl(tokenNeed float64) bool {
    now := time.Now()
    pastTime := now.Sub(b.lastAccessTime)
 
    // 在距离上次访问期间一共可发放了多少令牌
    newTokenNum := float64(pastTime) / float64(b.rate)

    // 剩余令牌数量不能超过桶的总空间
    b.curTokenNum = math.Min(b.bucketSize, b.curTokenNum + newTokenNum)

    b.lastAccessTime = now
    // tokenNeed 指处理一个请求需要的令牌数量
    if tokenNeed > b.curTokenNum {
        return false
    } else {
        b.curTokenNum = b.curTokenNum - tokenNeed
        return true
    }
}
 
func main() {
    // 创建一个放令牌速率为10qps,桶的总量为20的限流器
    limit := NewBucketLimit(10, 20)

    // 在桶里放入100滴水
    for i := 0; i < 100; i++ {
        // 处理1滴水需要消耗19块令牌
        if i == 50 {
            time.Sleep(3 * time.Second)
        }
        allow := limit.AllowControl(19)
        if allow {
            fmt.Printf("第%d滴水, 顺利流出\n", i)
            continue
        } else {
            fmt.Printf("第%d滴水, 溢出丢弃\n", i)
            time.Sleep(time.Millisecond * 100)
        }
    }
}