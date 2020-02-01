package tool

import "time"

// IntervalTimer 时间间隔计算器
type IntervalTimer struct {
	startTime time.Time
	endTime   time.Time
}

// NewIntervalTimer 初始化一个间隔计时器
func NewIntervalTimer() IntervalTimer {
	return IntervalTimer{
		startTime: time.Now(),
	}
}

// Millisecond 返回距离上次调用时间间隔计算器的时间
func (t *IntervalTimer) Millisecond() int64 {
	t.endTime = time.Now()
	defer func() { t.startTime = time.Now() }()
	return t.endTime.Sub(t.startTime).Milliseconds()
}
