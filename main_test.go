package main

import (
	"sync"
	"testing"
	"time"

	"github.com/scache-io/scache-demo/models"
	"github.com/scache-io/scache/cache"
	"github.com/scache-io/scache/config"
)

// TestBasicCacheOperations 测试基础缓存操作
func TestBasicCacheOperations(t *testing.T) {
	c := cache.NewLocalCache(config.DefaultEngineConfig())

	// 测试 Store
	user := models.User{
		ID:        1,
		Name:      "张三",
		Email:     "zhangsan@example.com",
		Age:       25,
		CreatedAt: time.Now(),
	}

	err := c.Store("user:1", user, time.Hour)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	// 测试 Load
	var loadedUser models.User
	err = c.Load("user:1", &loadedUser)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loadedUser.ID != 1 || loadedUser.Name != "张三" {
		t.Errorf("Loaded user mismatch: got %+v", loadedUser)
	}

	// 测试 Exists
	if !c.Exists("user:1") {
		t.Error("Exists should return true for existing key")
	}

	// 测试 Size
	size := c.Size()
	if size != 1 {
		t.Errorf("Size should be 1, got %d", size)
	}

	// 测试不存在的键
	if c.Exists("user:nonexistent") {
		t.Error("Exists should return false for non-existing key")
	}
}

// TestTTLExpiration 测试 TTL 过期
func TestTTLExpiration(t *testing.T) {
	c := cache.NewLocalCache(config.DefaultEngineConfig())

	// 设置短期缓存
	c.Store("temp:code", "123456", 500*time.Millisecond)

	// 立即读取
	if code, exists := c.GetString("temp:code"); exists {
		// GetString 可能返回带引号的 JSON 字符串，移除引号进行比较
		if len(code) >= 2 && code[0] == '"' && code[len(code)-1] == '"' {
			code = code[1 : len(code)-1]
		}
		if code != "123456" {
			t.Errorf("Expected '123456', got '%s'", code)
		}
	} else {
		t.Error("Key should exist immediately after storage")
	}

	// 等待过期
	time.Sleep(600 * time.Millisecond)

	// 读取已过期的键
	if _, exists := c.GetString("temp:code"); exists {
		t.Error("Key should be expired after TTL")
	}
}

// TestExpireModification 测试修改 TTL
func TestExpireModification(t *testing.T) {
	c := cache.NewLocalCache(config.DefaultEngineConfig())

	c.Store("user:profile", "张三", time.Hour)

	// 修改为5分钟过期
	c.Expire("user:profile", 5*time.Minute)

	ttl, exists := c.TTL("user:profile")
	if !exists {
		t.Error("TTL should exist for existing key")
	}

	// TTL 应该接近 5 分钟（允许 1 秒误差）
	expectedTTL := 5 * time.Minute
	if ttl < expectedTTL-time.Second || ttl > expectedTTL+time.Second {
		t.Errorf("Expected TTL ~5m, got %v", ttl)
	}
}

// TestDataTypeStrings 测试字符串类型
func TestDataTypeStrings(t *testing.T) {
	c := cache.NewLocalCache(config.DefaultEngineConfig())

	// SetString
	c.SetString("app:name", "SCache Demo", time.Hour)

	// GetString
	if name, exists := c.GetString("app:name"); exists {
		if name != "SCache Demo" {
			t.Errorf("Expected 'SCache Demo', got '%s'", name)
		}
	} else {
		t.Error("String key should exist")
	}

	// 测试不存在的键
	if _, exists := c.GetString("nonexistent"); exists {
		t.Error("Non-existing string key should not exist")
	}
}

// TestDataTypeLists 测试列表类型
func TestDataTypeLists(t *testing.T) {
	c := cache.NewLocalCache(config.DefaultEngineConfig())

	tags := []interface{}{"Go", "缓存", "高性能", "泛型"}
	c.SetList("tags:go", tags, time.Hour)

	// GetList
	loadedTags, exists := c.GetList("tags:go")
	if !exists {
		t.Error("List key should exist")
	}

	if len(loadedTags) != 4 {
		t.Errorf("Expected 4 items, got %d", len(loadedTags))
	}

	// 验证内容
	expectedTags := []string{"Go", "缓存", "高性能", "泛型"}
	for i, tag := range loadedTags {
		if tag != expectedTags[i] {
			t.Errorf("Item %d: expected '%s', got '%v'", i, expectedTags[i], tag)
		}
	}
}

// TestDataTypeHashes 测试哈希类型
func TestDataTypeHashes(t *testing.T) {
	c := cache.NewLocalCache(config.DefaultEngineConfig())

	profile := map[string]interface{}{
		"name":   "李四",
		"email":  "lisi@example.com",
		"city":   "北京",
		"gender": "男",
	}
	c.SetHash("profile:1001", profile, time.Hour)

	// GetHash
	loadedProfile, exists := c.GetHash("profile:1001")
	if !exists {
		t.Error("Hash key should exist")
	}

	// 验证字段
	if loadedProfile["name"] != "李四" {
		t.Errorf("Expected name '李四', got '%v'", loadedProfile["name"])
	}

	if loadedProfile["city"] != "北京" {
		t.Errorf("Expected city '北京', got '%v'", loadedProfile["city"])
	}
}

// TestConcurrencySafety 测试并发安全
func TestConcurrencySafety(t *testing.T) {
	cacheConfig := config.DefaultEngineConfig()
	cacheConfig.MaxSize = 1000
	c := cache.NewLocalCache(cacheConfig)

	var wg sync.WaitGroup
	routines := 10
	writesPerRoutine := 100

	// 并发写入
	for i := 0; i < routines; i++ {
		wg.Add(1)
		go func(routineID int) {
			defer wg.Done()
			for j := 0; j < writesPerRoutine; j++ {
				key := "concurrent:" + string(rune(routineID)) + ":" + string(rune(j))
				user := models.User{
					ID:   routineID*1000 + j,
					Name: "用户",
					Age:  20 + j%50,
				}
				c.Store(key, user, time.Hour)
			}
		}(i)
	}

	// 并发读取
	for i := 0; i < routines; i++ {
		wg.Add(1)
		go func(routineID int) {
			defer wg.Done()
			for j := 0; j < writesPerRoutine; j++ {
				key := "concurrent:" + string(rune(routineID)) + ":" + string(rune(j))
				var user models.User
				_ = c.Load(key, &user)
			}
		}(i)
	}

	wg.Wait()

	// 验证没有 panic，测试通过
	stats := c.Stats()
	if stats == nil {
		t.Error("Stats should return non-nil value")
	}
}

// TestPerformanceWrite 测试写入性能
func TestPerformanceWrite(t *testing.T) {
	c := cache.NewLocalCache(config.DefaultEngineConfig())

	testData := make([]models.User, 100)
	for i := 0; i < 100; i++ {
		testData[i] = models.User{
			ID:   i + 1,
			Name: "用户",
			Age:  20 + i%60,
		}
	}

	start := time.Now()
	for i := 0; i < 100; i++ {
		key := "perf:user:" + string(rune(i))
		c.Store(key, testData[i], time.Hour)
	}
	duration := time.Since(start)

	t.Logf("写入 100 条数据，耗时: %v", duration)

	// 验证所有数据都写入了
	if c.Size() != 100 {
		t.Errorf("Expected size 100, got %d", c.Size())
	}
}

// TestPerformanceRead 测试读取性能
func TestPerformanceRead(t *testing.T) {
	c := cache.NewLocalCache(config.DefaultEngineConfig())

	// 先写入数据
	for i := 0; i < 100; i++ {
		user := models.User{
			ID:   i + 1,
			Name: "用户",
			Age:  20 + i%60,
		}
		key := "perf:user:" + string(rune(i))
		c.Store(key, user, time.Hour)
	}

	// 测试读取性能
	start := time.Now()
	for i := 0; i < 100; i++ {
		key := "perf:user:" + string(rune(i))
		var user models.User
		_ = c.Load(key, &user)
	}
	duration := time.Since(start)

	t.Logf("读取 100 条数据，耗时: %v", duration)
}

// TestStats 测试统计信息
func TestStats(t *testing.T) {
	c := cache.NewLocalCache(config.DefaultEngineConfig())

	// 写入一些数据
	for i := 0; i < 10; i++ {
		user := models.User{
			ID:   i + 1,
			Name: "用户",
			Age:  20,
		}
		c.Store("user:"+string(rune(i)), user, time.Hour)
	}

	// 读取一些数据
	var user models.User
	c.Load("user:0", &user)
	c.Load("user:1", &user)
	c.Load("nonexistent", &user) // 未命中

	stats := c.Stats()
	if stats == nil {
		t.Fatal("Stats should return non-nil value")
	}

	statsMap, ok := stats.(map[string]interface{})
	if !ok {
		t.Fatal("Stats should return a map")
	}

	// 检查必需的字段
	requiredFields := []string{"keys", "sets", "hits", "misses", "hit_rate", "memory"}
	for _, field := range requiredFields {
		if _, exists := statsMap[field]; !exists {
			t.Errorf("Stats should contain field '%s'", field)
		}
	}

	t.Logf("Stats: %+v", stats)
}

// TestMultipleDataTypes 混合测试多种数据类型
func TestMultipleDataTypes(t *testing.T) {
	c := cache.NewLocalCache(config.DefaultEngineConfig())

	// 字符串
	c.SetString("str", "test", time.Hour)
	if val, _ := c.GetString("str"); val != "test" {
		t.Error("String value mismatch")
	}

	// 列表
	c.SetList("list", []interface{}{1, 2, 3}, time.Hour)
	if list, _ := c.GetList("list"); len(list) != 3 {
		t.Error("List length mismatch")
	}

	// 哈希
	c.SetHash("hash", map[string]interface{}{"key": "value"}, time.Hour)
	if hash, _ := c.GetHash("hash"); hash["key"] != "value" {
		t.Error("Hash value mismatch")
	}

	// 对象
	user := models.User{ID: 1, Name: "Test"}
	c.Store("obj", user, time.Hour)
	var loaded models.User
	c.Load("obj", &loaded)
	if loaded.ID != 1 {
		t.Error("Object value mismatch")
	}
}

// TestCacheEviction 测试缓存淘汰（当达到最大大小时）
func TestCacheEviction(t *testing.T) {
	cacheConfig := config.DefaultEngineConfig()
	cacheConfig.MaxSize = 5
	c := cache.NewLocalCache(cacheConfig)

	// 写入超过最大大小的数据
	for i := 0; i < 10; i++ {
		user := models.User{
			ID:   i + 1,
			Name: "用户",
			Age:  20,
		}
		c.Store("user:"+string(rune(i)), user, time.Hour)
	}

	// 缓存大小应该接近 MaxSize
	size := c.Size()
	if size > 5 {
		t.Errorf("Cache size should not exceed MaxSize, got %d", size)
	}

	t.Logf("Cache size after eviction: %d", size)
}

// TestEmptyOperations 测试空操作
func TestEmptyOperations(t *testing.T) {
	c := cache.NewLocalCache(config.DefaultEngineConfig())

	// 读取不存在的键
	var user models.User
	err := c.Load("nonexistent", &user)
	if err == nil {
		t.Error("Load should return error for non-existing key")
	}

	// 获取不存在的 TTL
	_, exists := c.TTL("nonexistent")
	if exists {
		t.Error("TTL should not exist for non-existing key")
	}

	// Exists 不存在的键
	if c.Exists("nonexistent") {
		t.Error("Exists should return false for non-existing key")
	}
}
