package core

import (
	"fmt"
	"github.com/jasonkayzk/distributed-id-generator/mysql"
	"sync"
	"time"
)

// 号段：[left,right)
type Segment struct {
	offset int // 消费偏移
	left   int // 左区间
	right  int // 右区间
}

// 关联到AppTag的号码池
type AppIdAllocator struct {
	mutex        sync.Mutex
	appTag       string      // 业务标识
	segments     []*Segment  // 双Buffer, 最少0个, 最多2个号段在内存
	isAllocating bool        // 是否正在分配中(远程获取)
	waiting      []chan byte // 因号码池空而挂起等待的客户端
}

// 全局分配器, 管理所有App的号码分配
type IdAllocator struct {
	mutex  sync.Mutex
	appMap map[string]*AppIdAllocator
}

var GlobalIdAllocator *IdAllocator

func InitAlloc() (err error) {
	GlobalIdAllocator = &IdAllocator{
		appMap: map[string]*AppIdAllocator{},
	}
	return
}

func (i *IdAllocator) LeftCount(appTag string) int {
	i.mutex.Lock()
	appIdAllocator, _ := i.appMap[appTag]
	i.mutex.Unlock()

	if appIdAllocator != nil {
		return appIdAllocator.leftCountWithMutex()
	}
	return 0
}

func (a *AppIdAllocator) leftCountWithMutex() (count int) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.leftCount()
}

func (a *AppIdAllocator) leftCount() int {
	var count int
	for i := 0; i < len(a.segments); i++ {
		count += a.segments[i].right - a.segments[i].left - a.segments[i].offset
	}
	return count
}

func (i *IdAllocator) NextId(appTag string) (int, error) {
	var (
		appIdAlloc *AppIdAllocator
		exist      bool
	)

	i.mutex.Lock()
	if appIdAlloc, exist = i.appMap[appTag]; !exist {
		appIdAlloc = &AppIdAllocator{
			appTag:       appTag,
			segments:     make([]*Segment, 0),
			isAllocating: false,
			waiting:      make([]chan byte, 0),
		}
		i.appMap[appTag] = appIdAlloc
	}
	i.mutex.Unlock()

	return appIdAlloc.nextId()
}

func (a *AppIdAllocator) nextId() (int, error) {
	var (
		nextId    int
		waitChan  chan byte
		waitTimer *time.Timer
		hasId     = false
	)

	a.mutex.Lock()
	defer a.mutex.Unlock()

	// 1:有剩余号码, 立即分配
	if a.leftCount() != 0 {
		nextId = a.popNextId()
		hasId = true
	}

	// 2:段<=1个, 启动补偿线程
	if len(a.segments) <= 1 && !a.isAllocating {
		a.isAllocating = true
		go a.fillSegments()
	}

	if hasId {
		return nextId, nil
	}

	// 3:没有剩余号码, 此时补偿线程一定正在运行, 等待其至多一段时间
	waitChan = make(chan byte, 1)
	a.waiting = append(a.waiting, waitChan) // 排队等待唤醒

	// 释放锁, 等待补偿线程唤醒
	a.mutex.Unlock()

	waitTimer = time.NewTimer(2 * time.Second) // 最多等待2秒
	select {
	case <-waitChan:
	case <-waitTimer.C:
	}

	// 4:再次上锁尝试获取号码
	a.mutex.Lock()
	if a.leftCount() != 0 {
		nextId = a.popNextId()
		return nextId, nil
	} else {
		return 0, fmt.Errorf("no available id")
	}
}

func (a *AppIdAllocator) popNextId() int {
	var nextId = a.segments[0].left + a.segments[0].offset
	a.segments[0].offset++
	if nextId+1 >= a.segments[0].right {
		a.segments = append(a.segments[:0], a.segments[1:]...) // 弹出第一个seg, 后续seg向前移动
	}
	return nextId
}

// 分配号码段, 直到足够2个segment, 否则始终不会退出
func (a *AppIdAllocator) fillSegments() {
	var failTimes = 0
	for {
		a.mutex.Lock()
		if len(a.segments) <= 1 { // 只剩余<=1段, 那么继续获取新号段
			a.mutex.Unlock()
			// 请求mysql获取号段
			if seg, err := a.newSegment(); err != nil {
				failTimes++
				if failTimes > 3 { // 连续失败超过3次则停止分配
					a.mutex.Lock()
					a.wakeup() // 唤醒等待者, 让它们立马失败
					goto LEAVE
				}
			} else {
				failTimes = 0 // 分配成功则失败次数重置为0
				// 新号段补充进去
				a.mutex.Lock()
				a.segments = append(a.segments, seg)
				a.wakeup()               // 尝试唤醒等待资源的调用
				if len(a.segments) > 1 { // 已生成2个号段, 停止继续分配
					goto LEAVE
				} else {
					a.mutex.Unlock()
				}
			}
		} else {
			// never reach
			break
		}
	}

LEAVE:
	a.isAllocating = false
	a.mutex.Unlock()
}

func (a *AppIdAllocator) newSegment() (*Segment, error) {
	maxId, step, err := mysql.NextId(a.appTag)
	if err != nil {
		return nil, err
	}

	return &Segment{
		left:  maxId - step,
		right: maxId,
	}, nil
}

func (a *AppIdAllocator) wakeup() {
	var waitChan chan byte

	for _, waitChan = range a.waiting {
		close(waitChan)
	}
	a.waiting = a.waiting[:0]
}
