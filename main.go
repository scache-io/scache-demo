package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/scache-io/scache"
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

	// 6. 键管理演示
	keyManagementDemo()

	// 7. 缓存配置演示
	cacheConfigDemo()

	// 8. 多缓存实例演示
	multipleCacheDemo()

	// 9. 实际场景演示：用户会话管理
	userSessionDemo()

	// 10. 实际场景演示：API 响应缓存
	apiCacheDemo()

	// 11. 实际场景演示：配置缓存
	configCacheDemo()

	// 12. 全局缓存 API 演示
	globalCacheAPIDemo()

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

// keyManagementDemo 键管理演示
func keyManagementDemo() {
	fmt.Println("6. 键管理演示")
	fmt.Println("-------------")

	// 创建缓存
	keyCache := cache.NewLocalCache(config.DefaultEngineConfig())

	// 添加多个键
	keyCache.SetString("key:1", "value1", time.Hour)
	keyCache.SetString("key:2", "value2", time.Hour)
	keyCache.SetString("key:3", "value3", time.Hour)
	fmt.Println("✓ 添加 3 个键")

	// 获取所有键
	keys := keyCache.Keys()
	fmt.Printf("✓ 所有键: %v (共 %d 个)\n", keys, len(keys))

	// 检查键是否存在
	fmt.Printf("✓ key:1 存在: %v\n", keyCache.Exists("key:1"))
	fmt.Printf("✓ key:4 存在: %v\n", keyCache.Exists("key:4"))

	// 删除键
	deleted := keyCache.Delete("key:2")
	fmt.Printf("✓ 删除 key:2: %v\n", deleted)

	// 获取删除后的所有键
	keys = keyCache.Keys()
	fmt.Printf("✓ 删除后的键: %v (共 %d 个)\n", keys, len(keys))

	// 清空缓存
	err := keyCache.Flush()
	if err == nil {
		fmt.Println("✓ 清空所有缓存")
	}

	// 获取清空后的缓存大小
	fmt.Printf("✓ 清空后缓存大小: %d\n", keyCache.Size())

	fmt.Println()
}

// cacheConfigDemo 缓存配置演示
func cacheConfigDemo() {
	fmt.Println("7. 缓存配置演示")
	fmt.Println("-------------")

	// 小型缓存配置
	smallConfig := &config.EngineConfig{
		MaxSize:         10,
		MemoryThreshold: 1024, // 1KB
	}
	smallCache := cache.NewLocalCache(smallConfig)
	fmt.Println("✓ 小型缓存 (MaxSize: 10, MemoryThreshold: 1KB)")

	// 写入超过限制
	for i := 0; i < 15; i++ {
		smallCache.SetString(fmt.Sprintf("small:%d", i), fmt.Sprintf("value%d", i), time.Hour)
	}
	fmt.Printf("✓ 写入 15 条数据后，实际缓存大小: %d (LRU 淘汰生效)\n", smallCache.Size())

	// 大型缓存配置
	largeConfig := &config.EngineConfig{
		MaxSize:         10000,
		MemoryThreshold: 10 * 1024 * 1024, // 10MB
	}
	_ = cache.NewLocalCache(largeConfig)
	fmt.Println("✓ 大型缓存 (MaxSize: 10000, MemoryThreshold: 10MB)")

	fmt.Println()
}

// multipleCacheDemo 多缓存实例演示
func multipleCacheDemo() {
	fmt.Println("8. 多缓存实例演示")
	fmt.Println("----------------")

	// 用户缓存
	userCache := cache.NewLocalCache(config.DefaultEngineConfig())
	userCache.Store("user:1", models.User{ID: 1, Name: "张三"}, time.Hour)
	fmt.Println("✓ 用户缓存: user:1 = 张三")

	// 产品缓存
	productCache := cache.NewLocalCache(config.DefaultEngineConfig())
	productCache.Store("product:1", models.Product{ID: "P001", Name: "iPhone", Price: 5999.99}, time.Hour)
	fmt.Println("✓ 产品缓存: product:1 = iPhone")

	// 会话缓存
	sessionCache := cache.NewLocalCache(config.DefaultEngineConfig())
	sessionCache.Store("session:abc123", models.Session{
		SessionID: "abc123",
		UserID:    "user_001",
		Username:  "张三",
		LoginAt:   time.Now(),
		ExpiresAt: time.Now().Add(time.Hour),
	}, time.Hour)
	fmt.Println("✓ 会话缓存: session:abc123")

	// 三个缓存实例相互独立
	fmt.Printf("✓ 用户缓存大小: %d\n", userCache.Size())
	fmt.Printf("✓ 产品缓存大小: %d\n", productCache.Size())
	fmt.Printf("✓ 会话缓存大小: %d\n", sessionCache.Size())

	fmt.Println()
}

// userSessionDemo 用户会话管理演示
func userSessionDemo() {
	fmt.Println("9. 用户会话管理演示")
	fmt.Println("------------------")

	sessionCache := cache.NewLocalCache(config.DefaultEngineConfig())

	// 模拟用户登录，创建会话
	sessionID := "sess_abc123"
	userSession := models.Session{
		SessionID: sessionID,
		UserID:    "user_001",
		Username:  "张三",
		LoginAt:   time.Now(),
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}
	sessionCache.Store(sessionID, userSession, 30*time.Minute)
	fmt.Println("✓ 用户登录，创建会话:", sessionID)

	// 模拟验证会话
	var loadedSession models.Session
	if sessionCache.Load(sessionID, &loadedSession) == nil {
		fmt.Printf("✓ 会话有效: 用户=%s, 剩余时间=%.0f分钟\n",
			loadedSession.Username, time.Until(loadedSession.ExpiresAt).Minutes())
	}

	// 模拟延长会话
	sessionCache.Expire(sessionID, 60*time.Minute)
	if ttl, exists := sessionCache.TTL(sessionID); exists {
		fmt.Printf("✓ 会话延期: 剩余时间=%.0f分钟\n", ttl.Minutes())
	}

	// 模拟用户登出
	deleted := sessionCache.Delete(sessionID)
	fmt.Printf("✓ 用户登出: 删除会话=%v\n", deleted)

	fmt.Println()
}

// apiCacheDemo API 响应缓存演示
func apiCacheDemo() {
	fmt.Println("10. API 响应缓存演示")
	fmt.Println("-------------------")

	apiCache := cache.NewLocalCache(config.DefaultEngineConfig())

	// 模拟 API 请求
	apiKey := "api:products:hot"

	// 第一次请求 - 缓存未命中
	var products []models.Product
	if apiCache.Load(apiKey, &products) != nil {
		// 模拟从数据库获取数据
		products = []models.Product{
			{ID: "P001", Name: "iPhone 15", Price: 5999.99, Stock: 100},
			{ID: "P002", Name: "MacBook Pro", Price: 12999.99, Stock: 50},
			{ID: "P003", Name: "AirPods Pro", Price: 1999.99, Stock: 200},
		}
		apiCache.Store(apiKey, products, 10*time.Minute)
		fmt.Println("✓ 第一次请求：从数据库获取数据，缓存 10 分钟")
	}

	fmt.Printf("✓ 热门商品: %d 个\n", len(products))
	for _, p := range products {
		fmt.Printf("  - %s: ¥%.2f\n", p.Name, p.Price)
	}

	// 第二次请求 - 缓存命中
	var cachedProducts []models.Product
	if apiCache.Load(apiKey, &cachedProducts) == nil {
		fmt.Println("✓ 第二次请求：从缓存获取数据（响应更快）")
	}

	// 获取剩余缓存时间
	if ttl, exists := apiCache.TTL(apiKey); exists {
		fmt.Printf("✓ 缓存剩余时间: %.2f 分钟\n", ttl.Minutes())
	}

	fmt.Println()
}

// configCacheDemo 配置缓存演示
func configCacheDemo() {
	fmt.Println("11. 配置缓存演示")
	fmt.Println("---------------")

	configCache := cache.NewLocalCache(config.DefaultEngineConfig())

	// 缓存应用配置
	configCache.SetString("app:name", "SCache Demo", 24*time.Hour)
	configCache.SetString("app:version", "1.0.0", 24*time.Hour)
	configCache.SetString("app:mode", "production", 24*time.Hour)
	configCache.SetString("app:debug", "false", 24*time.Hour)
	fmt.Println("✓ 缓存应用配置")

	// 读取配置
	if name, exists := configCache.GetString("app:name"); exists {
		fmt.Printf("✓ 应用名称: %s\n", name)
	}
	if version, exists := configCache.GetString("app:version"); exists {
		fmt.Printf("✓ 应用版本: %s\n", version)
	}
	if mode, exists := configCache.GetString("app:mode"); exists {
		fmt.Printf("✓ 运行模式: %s\n", mode)
	}

	// 批量获取所有配置
	allKeys := configCache.Keys()
	fmt.Printf("✓ 共缓存 %d 个配置项\n", len(allKeys))

	fmt.Println()
}

// globalCacheAPIDemo 全局缓存 API 演示
func globalCacheAPIDemo() {
	fmt.Println("12. 全局缓存 API 演示")
	fmt.Println("-------------------")

	// 使用全局 API（不需要手动创建缓存实例）
	scache.SetString("global:user:1", "全局用户1", time.Hour)
	fmt.Println("✓ 使用全局 API 存储数据")

	if user, exists := scache.GetString("global:user:1"); exists {
		fmt.Printf("✓ 使用全局 API 读取: %s\n", user)
	}

	// 使用全局 API 存储结构体
	globalUser := models.User{ID: 1, Name: "全局张三"}
	err := scache.Store("global:user:obj", globalUser, time.Hour)
	if err == nil {
		fmt.Println("✓ 使用全局 API 存储结构体")
	}

	var loadedGlobalUser models.User
	err = scache.Load("global:user:obj", &loadedGlobalUser)
	if err == nil {
		fmt.Printf("✓ 使用全局 API 读取结构体: %+v\n", loadedGlobalUser)
	}

	// 使用全局 API 删除
	deleted := scache.Delete("global:user:1")
	fmt.Printf("✓ 使用全局 API 删除: %v\n", deleted)

	// 获取全局统计信息
	stats := scache.Stats()
	statsMap := stats.(map[string]interface{})
	fmt.Printf("✓ 全局缓存统计:\n")
	fmt.Printf("  - 缓存大小: %v\n", statsMap["keys"])
	fmt.Printf("  - 写入次数: %v\n", statsMap["sets"])
	fmt.Printf("  - 命中率: %.2f%%\n", statsMap["hit_rate"].(float64)*100)

	fmt.Println()
}
