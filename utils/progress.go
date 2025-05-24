package utils

import (
	"fmt"
	"strings"
	"time"
)

type ProgressBar struct {
	total     int
	current   int
	width     int
	startTime time.Time
	lastDraw  time.Time
}

func NewProgressBar(total int) *ProgressBar {
	return &ProgressBar{
		total:     total,
		current:   0,
		width:     50,
		startTime: time.Now(),
		lastDraw:  time.Now(),
	}
}

func (pb *ProgressBar) Update(current int) {
	pb.current = current

	// 限制更新频率，避免频繁刷新
	if time.Since(pb.lastDraw) < 100*time.Millisecond && current != pb.total {
		return
	}

	pb.Draw()
	pb.lastDraw = time.Now()
}

func (pb *ProgressBar) Increment() {
	pb.Update(pb.current + 1)
}

func (pb *ProgressBar) Draw() {
	percentage := float64(pb.current) / float64(pb.total) * 100
	filled := int(float64(pb.width) * float64(pb.current) / float64(pb.total))

	bar := strings.Repeat("█", filled) + strings.Repeat("░", pb.width-filled)

	elapsed := time.Since(pb.startTime)
	var eta time.Duration
	if pb.current > 0 {
		eta = time.Duration(float64(elapsed) / float64(pb.current) * float64(pb.total-pb.current))
	}

	fmt.Printf("\r[%s] %d/%d (%.1f%%) - 用时: %v - 预计剩余: %v",
		bar, pb.current, pb.total, percentage,
		elapsed.Round(time.Second), eta.Round(time.Second))

	if pb.current == pb.total {
		fmt.Println() // 完成后换行
	}
}

func (pb *ProgressBar) Finish() {
	pb.Update(pb.total)
}
