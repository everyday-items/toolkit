package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hexagon-codes/toolkit/infra/db/mysql"
	"github.com/hexagon-codes/toolkit/infra/db/redis"
)

func main() {
	fmt.Println("=== GoPkg DB 示例 ===")

	// MySQL 示例
	mysqlExample()

	// Redis 示例
	redisExample()

	// 分布式锁示例
	lockExample()

	fmt.Println("\n✅ 示例完成!")
}

func mysqlExample() {
	fmt.Println("📦 MySQL 示例:")

	// 初始化 MySQL（实际使用时需要有效的 DSN）
	config := mysql.DefaultConfig("root:password@tcp(localhost:3306)/test?parseTime=true")

	// 注意：这里会连接失败，因为没有真实的 MySQL 服务
	// 实际使用时请提供有效的 DSN
	_, err := mysql.New(config)
	if err != nil {
		fmt.Printf("  ⚠️  MySQL 连接失败（预期行为）: %v\n", err)
		return
	}

	// 示例代码（连接成功后执行）
	fmt.Println("  - 创建用户表")
	fmt.Println("  - 插入用户数据")
	fmt.Println("  - 查询用户列表")
	fmt.Println("  - 事务操作")
	fmt.Println()
}

func redisExample() {
	fmt.Println("📦 Redis 示例:")

	// 初始化 Redis（实际使用时需要有效的 Redis 地址）
	config := redis.DefaultConfig("localhost:6379")

	// 注意：这里会连接失败，因为没有真实的 Redis 服务
	// 实际使用时请提供有效的 Redis 地址
	client, err := redis.New(config)
	if err != nil {
		fmt.Printf("  ⚠️  Redis 连接失败（预期行为）: %v\n", err)
		return
	}
	defer client.Close()

	ctx := context.Background()

	// Set
	fmt.Println("  - Set key: name = Alice")
	client.Set(ctx, "name", "Alice", time.Minute)

	// Get
	val, _ := client.Get(ctx, "name").Result()
	fmt.Printf("  - Get key: name = %s\n", val)

	// Incr
	client.Incr(ctx, "counter")
	fmt.Println("  - Incr counter")

	// Hash
	client.HSet(ctx, "user:1", "name", "Bob", "age", 25)
	fmt.Println("  - HSet user:1")

	// List
	client.LPush(ctx, "queue", "task1", "task2")
	fmt.Println("  - LPush queue")

	fmt.Println()
}

func lockExample() {
	fmt.Println("🔒 分布式锁示例:")

	config := redis.DefaultConfig("localhost:6379")
	client, err := redis.New(config)
	if err != nil {
		fmt.Printf("  ⚠️  Redis 连接失败（预期行为）: %v\n", err)
		return
	}
	defer client.Close()

	ctx := context.Background()

	// 使用 WithLock 自动管理锁（使用 UniversalClient）
	err = redis.WithLock(ctx, client.UniversalClient, "lock:resource", 30*time.Second, func() error {
		fmt.Println("  - 获取锁成功")
		fmt.Println("  - 执行业务逻辑...")
		time.Sleep(100 * time.Millisecond)
		fmt.Println("  - 业务逻辑完成")
		return nil
	})

	if err != nil {
		log.Printf("  ❌ 锁操作失败: %v", err)
		return
	}

	fmt.Println("  - 锁已自动释放")
	fmt.Println()
}
