package utils

import (
	"fmt"
	"sync"
	"time"
)

type PerformanceStats struct {
	mu               sync.RWMutex
	TranslationCount int           `json:"translation_count"`
	CacheHits        int           `json:"cache_hits"`
	CacheMisses      int           `json:"cache_misses"`
	TotalTime        time.Duration `json:"total_time"`
	TranslationTime  time.Duration `json:"translation_time"`
	FileOperations   int           `json:"file_operations"`
	Errors           int           `json:"errors"`
}

var globalStats = &PerformanceStats{}

func (ps *PerformanceStats) AddTranslation(duration time.Duration) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.TranslationCount++
	ps.TranslationTime += duration
}

func (ps *PerformanceStats) AddCacheHit() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.CacheHits++
}

func (ps *PerformanceStats) AddCacheMiss() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.CacheMisses++
}

func (ps *PerformanceStats) AddFileOperation() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.FileOperations++
}

func (ps *PerformanceStats) AddError() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.Errors++
}

func (ps *PerformanceStats) GetStats() PerformanceStats {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return *ps
}

func (ps *PerformanceStats) Reset() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	*ps = PerformanceStats{}
}

func (ps *PerformanceStats) String() string {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	cacheHitRate := 0.0
	if ps.CacheHits+ps.CacheMisses > 0 {
		cacheHitRate = float64(ps.CacheHits) / float64(ps.CacheHits+ps.CacheMisses) * 100
	}

	return fmt.Sprintf(`📊 性能统计:
   🔄 翻译次数: %d
   ⚡ 缓存命中率: %.1f%% (%d/%d)
   ⏱️  平均翻译时间: %v
   📁 文件操作: %d
   ❌ 错误次数: %d`,
		ps.TranslationCount,
		cacheHitRate, ps.CacheHits, ps.CacheHits+ps.CacheMisses,
		ps.getAverageTranslationTime(),
		ps.FileOperations,
		ps.Errors)
}

func (ps *PerformanceStats) getAverageTranslationTime() time.Duration {
	if ps.TranslationCount == 0 {
		return 0
	}
	return ps.TranslationTime / time.Duration(ps.TranslationCount)
}

// 全局统计函数
func RecordTranslation(duration time.Duration) {
	globalStats.AddTranslation(duration)
}

func RecordCacheHit() {
	globalStats.AddCacheHit()
}

func RecordCacheMiss() {
	globalStats.AddCacheMiss()
}

func RecordFileOperation() {
	globalStats.AddFileOperation()
}

func RecordError() {
	globalStats.AddError()
}

func GetGlobalStats() PerformanceStats {
	return globalStats.GetStats()
}

func ResetGlobalStats() {
	globalStats.Reset()
}
