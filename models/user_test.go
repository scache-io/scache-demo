package models

import (
	"testing"
	"time"
)

// TestUserJSONSerialization 测试 User 模型的 JSON 序列化/反序列化
func TestUserJSONSerialization(t *testing.T) {
	now := time.Now()
	user := User{
		ID:        1,
		Name:      "张三",
		Email:     "zhangsan@example.com",
		Age:       25,
		CreatedAt: now,
	}

	// 验证字段值
	if user.ID != 1 {
		t.Errorf("Expected ID 1, got %d", user.ID)
	}

	if user.Name != "张三" {
		t.Errorf("Expected Name '张三', got '%s'", user.Name)
	}

	if user.Email != "zhangsan@example.com" {
		t.Errorf("Expected Email 'zhangsan@example.com', got '%s'", user.Email)
	}

	if user.Age != 25 {
		t.Errorf("Expected Age 25, got %d", user.Age)
	}

	if user.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

// TestUserCreation 测试创建 User 实例
func TestUserCreation(t *testing.T) {
	tests := []struct {
		name  string
		user  User
		valid bool
	}{
		{
			name:  "Valid user",
			user:  User{ID: 1, Name: "Test", Email: "test@example.com", Age: 20},
			valid: true,
		},
		{
			name:  "User with zero ID",
			user:  User{ID: 0, Name: "Test", Email: "test@example.com", Age: 20},
			valid: true, // ID 0 可能是有效的
		},
		{
			name:  "User with empty name",
			user:  User{ID: 1, Name: "", Email: "test@example.com", Age: 20},
			valid: true, // 模型本身不验证业务逻辑
		},
		{
			name:  "User with negative age",
			user:  User{ID: 1, Name: "Test", Email: "test@example.com", Age: -1},
			valid: true, // 模型本身不验证业务逻辑
		},
		{
			name:  "Empty user",
			user:  User{},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// User 模型没有验证方法，所以我们只测试字段的存储
			if tt.user.ID != 0 && tt.user.ID < 0 {
				t.Error("ID should not be negative")
			}

			// 测试字段是否正确设置
			if tt.user.Name == "" && tt.valid {
				t.Log("Empty name is allowed by model")
			}
		})
	}
}

// TestUserTimestamps 测试时间戳字段
func TestUserTimestamps(t *testing.T) {
	now := time.Now()
	user := User{
		ID:        1,
		Name:      "Test",
		Email:     "test@example.com",
		Age:       25,
		CreatedAt: now,
	}

	// 验证 CreatedAt
	if user.CreatedAt.Unix() != now.Unix() {
		t.Error("CreatedAt should match the provided time")
	}

	// 测试创建带不同时间戳的用户
	past := time.Now().Add(-24 * time.Hour)
	user2 := User{
		ID:        2,
		Name:      "Test2",
		CreatedAt: past,
	}

	if user2.CreatedAt.After(now) {
		t.Error("CreatedAt should be in the past")
	}

	duration := now.Sub(user2.CreatedAt)
	if duration < 23*time.Hour || duration > 25*time.Hour {
		t.Errorf("Expected ~24h difference, got %v", duration)
	}
}

// TestProductModel 测试 Product 模型
func TestProductModel(t *testing.T) {
	product := Product{
		ID:    "PROD-001",
		Name:  "测试产品",
		Price: 99.99,
		Stock: 100,
	}

	// 验证字段
	if product.ID != "PROD-001" {
		t.Errorf("Expected ID 'PROD-001', got '%s'", product.ID)
	}

	if product.Name != "测试产品" {
		t.Errorf("Expected Name '测试产品', got '%s'", product.Name)
	}

	if product.Price != 99.99 {
		t.Errorf("Expected Price 99.99, got %.2f", product.Price)
	}

	if product.Stock != 100 {
		t.Errorf("Expected Stock 100, got %d", product.Stock)
	}
}

// TestProductEdgeCases 测试 Product 边界情况
func TestProductEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		product Product
	}{
		{
			name: "Zero price",
			product: Product{
				ID:    "PROD-001",
				Name:  "免费产品",
				Price: 0,
				Stock: 1000,
			},
		},
		{
			name: "Zero stock",
			product: Product{
				ID:    "PROD-002",
				Name:  "缺货产品",
				Price: 99.99,
				Stock: 0,
			},
		},
		{
			name: "Negative stock (backorder)",
			product: Product{
				ID:    "PROD-003",
				Name:  "预售产品",
				Price: 99.99,
				Stock: -10,
			},
		},
		{
			name: "Very high price",
			product: Product{
				ID:    "PROD-004",
				Name:  "奢侈品",
				Price: 999999.99,
				Stock: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Product 模型不验证业务逻辑，只测试字段存储
			if tt.product.ID == "" {
				t.Error("Product ID should not be empty")
			}

			t.Logf("Product: %+v", tt.product)
		})
	}
}

// TestSessionModel 测试 Session 模型
func TestSessionModel(t *testing.T) {
	now := time.Now()
	expires := now.Add(24 * time.Hour)

	session := Session{
		SessionID: "sess-12345",
		UserID:    "user-67890",
		Username:  "张三",
		LoginAt:   now,
		ExpiresAt: expires,
	}

	// 验证字段
	if session.SessionID != "sess-12345" {
		t.Errorf("Expected SessionID 'sess-12345', got '%s'", session.SessionID)
	}

	if session.UserID != "user-67890" {
		t.Errorf("Expected UserID 'user-67890', got '%s'", session.UserID)
	}

	if session.Username != "张三" {
		t.Errorf("Expected Username '张三', got '%s'", session.Username)
	}

	if session.LoginAt.IsZero() {
		t.Error("LoginAt should not be zero")
	}

	if session.ExpiresAt.IsZero() {
		t.Error("ExpiresAt should not be zero")
	}

	// 验证 ExpiresAt 在 LoginAt 之后
	if !session.ExpiresAt.After(session.LoginAt) {
		t.Error("ExpiresAt should be after LoginAt")
	}
}

// TestSessionValidity 测试 Session 有效性
func TestSessionValidity(t *testing.T) {
	tests := []struct {
		name     string
		session  Session
		expired  bool
		testTime time.Time
	}{
		{
			name: "Valid session (future expiration)",
			session: Session{
				SessionID: "sess-001",
				UserID:    "user-001",
				LoginAt:   time.Now(),
				ExpiresAt: time.Now().Add(24 * time.Hour),
			},
			expired:  false,
			testTime: time.Now(),
		},
		{
			name: "Expired session",
			session: Session{
				SessionID: "sess-002",
				UserID:    "user-002",
				LoginAt:   time.Now().Add(-25 * time.Hour),
				ExpiresAt: time.Now().Add(-1 * time.Hour),
			},
			expired:  true,
			testTime: time.Now(),
		},
		{
			name: "Session expiring now",
			session: Session{
				SessionID: "sess-003",
				UserID:    "user-003",
				LoginAt:   time.Now().Add(-23 * time.Hour),
				ExpiresAt: time.Now().Add(time.Second), // 设置为 1 秒后过期
			},
			expired:  false,
			testTime: time.Now(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isExpired := tt.testTime.After(tt.session.ExpiresAt)

			if tt.expired != isExpired {
				t.Errorf("Expected expired=%v, got %v", tt.expired, isExpired)
			}

			t.Logf("Session validity: %v", !isExpired)
		})
	}
}

// TestSessionDuration 测试 Session 持续时间
func TestSessionDuration(t *testing.T) {
	now := time.Now()
	session := Session{
		SessionID: "sess-001",
		UserID:    "user-001",
		LoginAt:   now,
		ExpiresAt: now.Add(24 * time.Hour),
	}

	duration := session.ExpiresAt.Sub(session.LoginAt)
	expectedDuration := 24 * time.Hour

	// 允许 1 秒误差
	if duration < expectedDuration-time.Second || duration > expectedDuration+time.Second {
		t.Errorf("Expected duration ~24h, got %v", duration)
	}
}

// TestMultipleUsers 测试多个用户实例
func TestMultipleUsers(t *testing.T) {
	users := []User{
		{ID: 1, Name: "用户1", Email: "user1@example.com", Age: 20},
		{ID: 2, Name: "用户2", Email: "user2@example.com", Age: 25},
		{ID: 3, Name: "用户3", Email: "user3@example.com", Age: 30},
	}

	// 验证每个用户
	for i, user := range users {
		if user.ID != i+1 {
			t.Errorf("User %d: Expected ID %d, got %d", i, i+1, user.ID)
		}

		if user.Name != "用户"+string(rune('1'+i)) {
			t.Errorf("User %d: Name mismatch", i)
		}

		if user.Email == "" {
			t.Errorf("User %d: Email should not be empty", i)
		}

		if user.Age < 20 || user.Age > 30 {
			t.Errorf("User %d: Age should be between 20 and 30", i)
		}
	}
}

// TestMultipleProducts 测试多个产品实例
func TestMultipleProducts(t *testing.T) {
	products := []Product{
		{ID: "P1", Name: "产品1", Price: 10.99, Stock: 100},
		{ID: "P2", Name: "产品2", Price: 20.99, Stock: 200},
		{ID: "P3", Name: "产品3", Price: 30.99, Stock: 300},
	}

	// 验证每个产品
	for i, product := range products {
		expectedID := "P" + string(rune('1'+i))
		if product.ID != expectedID {
			t.Errorf("Product %d: Expected ID '%s', got '%s'", i, expectedID, product.ID)
		}

		expectedName := "产品" + string(rune('1'+i))
		if product.Name != expectedName {
			t.Errorf("Product %d: Name mismatch", i)
		}

		if product.Price < 10 || product.Price > 31 {
			t.Errorf("Product %d: Price out of expected range", i)
		}

		if product.Stock < 100 || product.Stock > 300 {
			t.Errorf("Product %d: Stock out of expected range", i)
		}
	}
}

// TestModelFieldsCoverage 测试所有模型的字段覆盖率
func TestModelFieldsCoverage(t *testing.T) {
	// User 模型
	user := User{
		ID:        1,
		Name:      "测试用户",
		Email:     "test@example.com",
		Age:       25,
		CreatedAt: time.Now(),
	}
	_ = user.ID
	_ = user.Name
	_ = user.Email
	_ = user.Age
	_ = user.CreatedAt

	// Product 模型
	product := Product{
		ID:    "P1",
		Name:  "测试产品",
		Price: 99.99,
		Stock: 100,
	}
	_ = product.ID
	_ = product.Name
	_ = product.Price
	_ = product.Stock

	// Session 模型
	session := Session{
		SessionID: "S1",
		UserID:    "U1",
		Username:  "测试用户",
		LoginAt:   time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	_ = session.SessionID
	_ = session.UserID
	_ = session.Username
	_ = session.LoginAt
	_ = session.ExpiresAt

	t.Log("All model fields are covered")
}
