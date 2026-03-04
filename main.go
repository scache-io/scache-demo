package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/scache-io/scache-demo/models"
	"github.com/scache-io/scache/cache"
	"github.com/scache-io/scache/config"
)

func main() {
	fmt.Println("=== SCache 演示项目 ===")

	// 1. 基础缓存操作
	basicCacheDemo()

	// 2. TTL 过期演示
	ttlDemo()

	// 3. 多数据类型演示
	dataTypeDemo()

	// 4. 并发安全演示
	concurrencyDemo()

	// 5. 性能对比演示
	performanceDemo()

	fmt.Println("\n=== 演示完成 ===")
}

// basicCacheDemo 基础缓存操作演示
func basicCacheDemo() {
	fmt.Println("1. 基础缓存操作")
	fmt.Println("-----------------")

	// 创建局部缓存
	userCache := cache.NewLocalCache(config.DefaultEngineConfig())

	// 存储用户数据
	user := models.User{
		ID:        1,
		Name:      "张三",
		Email:     "zhangsan@example.com",
		Age:       25,
		CreatedAt: time.Now(),
	}

	err := userCache.Store("user:1", user, time.Hour)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ 存储用户数据: user:1")

	// 读取用户数据
	var loadedUser models.User
	err = userCache.Load("user:1", &loadedUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ 读取用户数据: %+v\n", loadedUser)

	// 检查键是否存在
	if userCache.Exists("user:1") {
		fmt.Println("✓ 用户缓存存在")
	}

	// 获取缓存大小
	fmt.Printf("✓ 缓存大小: %d\n", userCache.Size())

	fmt.Println()
}

// ttlDemo TTL 过期演示
func ttlDemo() {
	fmt.Println("2. TTL 过期演示")
	fmt.Println("---------------")

	// 创建局部缓存
	testCache := cache.NewLocalCache(config.DefaultEngineConfig())

	// 设置短期缓存（2秒过期）
	testCache.Store("temp:code", "123456", 2*time.Second)
	fmt.Println("✓ 设置验证码缓存（2秒过期）")

	// 立即读取
	if code, exists := testCache.GetString("temp:code"); exists {
		fmt.Printf("✓ 验证码（立即读取）: %s\n", code)
	}

	// 等待3秒后读取
	fmt.Println("⏳ 等待3秒...")
	time.Sleep(3 * time.Second)

	if _, exists := testCache.GetString("temp:code"); !exists {
		fmt.Println("✓ 验证码已过期")
	}

	// 设置并修改 TTL
	testCache.Store("user:profile", "张三", time.Hour)
	fmt.Println("✓ 设置用户信息缓存（1小时过期）")

	// 修改为5分钟过期
	testCache.Expire("user:profile", 5*time.Minute)
	fmt.Println("✓ 修改过期时间为5分钟")

	if ttl, exists := testCache.TTL("user:profile"); exists {
		fmt.Printf("✓ 剩余过期时间: %v\n", ttl)
	}

	fmt.Println()
}

// dataTypeDemo 多数据类型演示
func dataTypeDemo() {
	fmt.Println("3. 多数据类型演示")
	fmt.Println("----------------")

	// 创建局部缓存
	typeCache := cache.NewLocalCache(config.DefaultEngineConfig())

	// 字符串类型
	typeCache.SetString("app:name", "SCache Demo", time.Hour)
	fmt.Println("✓ 字符串: app:name = SCache Demo")

	if name, exists := typeCache.GetString("app:name"); exists {
		fmt.Printf("  读取: %s\n", name)
	}

	// 列表类型
	tags := []interface{}{"Go", "缓存", "高性能", "泛型"}
	typeCache.SetList("tags:go", tags, time.Hour)
	fmt.Println("✓ 列表: tags:go = [Go, 缓存, 高性能, 泛型]")

	if loadedTags, exists := typeCache.GetList("tags:go"); exists {
		fmt.Printf("  读取: %v\n", loadedTags)
	}

	// 哈希类型
	profile := map[string]interface{}{
		"name":   "李四",
		"email":  "lisi@example.com",
		"city":   "北京",
		"gender": "男",
	}
	typeCache.SetHash("profile:1001", profile, time.Hour)
	fmt.Println("✓ 哈希: profile:1001 = {name: 李四, email: lisi@example.com, city: 北京}")

	if loadedProfile, exists := typeCache.GetHash("profile:1001"); exists {
		fmt.Printf("  读取: %v\n", loadedProfile)
	}

	fmt.Println()
}

// concurrencyDemo 并发安全演示
func concurrencyDemo() {
	fmt.Println("4. 并发安全演示")
	fmt.Println("---------------")

	// 创建带统计的缓存
	cacheConfig := config.DefaultEngineConfig()
	cacheConfig.MaxSize = 100
	concurrentCache := cache.NewLocalCache(cacheConfig)

	var wg sync.WaitGroup
	routines := 10
	writesPerRoutine := 100

	fmt.Printf("✓ 启动 %d 个协程，每个写入 %d 次\n", routines, writesPerRoutine)

	// 并发写入
	for i := 0; i < routines; i++ {
		wg.Add(1)
		go func(routineID int) {
			defer wg.Done()
			for j := 0; j < writesPerRoutine; j++ {
				key := fmt.Sprintf("concurrent:%d:%d", routineID, j)
				user := models.User{
					ID:   routineID*1000 + j,
					Name: fmt.Sprintf("用户-%d-%d", routineID, j),
					Age:  20 + j%50,
				}
				concurrentCache.Store(key, user, time.Hour)
			}
		}(i)
	}

	// 并发读取
	for i := 0; i < routines; i++ {
		wg.Add(1)
		go func(routineID int) {
			defer wg.Done()
			for j := 0; j < writesPerRoutine; j++ {
				key := fmt.Sprintf("concurrent:%d:%d", routineID, j)
				var user models.User
				if concurrentCache.Load(key, &user) == nil {
					// 读取成功
				}
			}
		}(i)
	}

	wg.Wait()

	// 获取统计信息
	stats := concurrentCache.Stats()
	fmt.Printf("✓ 并发操作完成\n")
	fmt.Printf("  缓存大小: %v\n", stats.(map[string]interface{})["keys"])
	fmt.Printf("  写入次数: %v\n", stats.(map[string]interface{})["sets"])
	fmt.Printf("  命中次数: %v\n", stats.(map[string]interface{})["hits"])
	fmt.Printf("  未命中次数: %v\n", stats.(map[string]interface{})["misses"])

	fmt.Println()
}

// performanceDemo 性能对比演示
func performanceDemo() {
	fmt.Println("5. 性能对比演示")
	fmt.Println("---------------")

	// 创建缓存
	perfCache := cache.NewLocalCache(config.DefaultEngineConfig())

	// 准备测试数据
	testData := make([]models.User, 1000)
	for i := 0; i < 1000; i++ {
		testData[i] = models.User{
			ID:   i + 1,
			Name: fmt.Sprintf("用户-%d", i+1),
			Age:  20 + i%60,
		}
	}

	// 测试写入性能
	fmt.Println("测试写入性能...")
	writeStart := time.Now()
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("perf:user:%d", i+1)
		perfCache.Store(key, testData[i], time.Hour)
	}
	writeDuration := time.Since(writeStart)
	fmt.Printf("✓ 写入 1000 条数据，耗时: %v (%.2f 条/秒)\n",
		writeDuration, float64(1000)/writeDuration.Seconds())

	// 测试读取性能
	fmt.Println("测试读取性能...")
	readStart := time.Now()
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("perf:user:%d", i+1)
		var user models.User
		if perfCache.Load(key, &user) == nil {
			// 读取成功
		}
	}
	readDuration := time.Since(readStart)
	fmt.Printf("✓ 读取 1000 条数据，耗时: %v (%.2f 条/秒)\n",
		readDuration, float64(1000)/readDuration.Seconds())

	// 测试统计信息
	stats := perfCache.Stats()
	fmt.Printf("✓ 最终统计:\n")
	fmt.Printf("  命中率: %.2f%%\n", stats.(map[string]interface{})["hit_rate"].(float64)*100)
	fmt.Printf("  缓存大小: %v\n", stats.(map[string]interface{})["keys"])
	fmt.Printf("  内存使用: %v bytes\n", stats.(map[string]interface{})["memory"])

	fmt.Println()
}
