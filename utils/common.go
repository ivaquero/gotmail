package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Account 账户数据结构
type Account struct {
	ID        string    `json:"id"`
	Address   string    `json:"address"`
	Password  string    `json:"password"`
	Token     TokenData `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
}

// TokenData JWT令牌数据结构
type TokenData struct {
	Token string `json:"token"`
}

// Domain API返回的域名数据结构
type Domain struct {
	Domain string `json:"domain"`
}

// DomainResponse API域名响应
type DomainResponse struct {
	HydraMember []Domain `json:"hydra:member"`
}

// Message 邮件消息数据结构
type Message struct {
	ID   string `json:"id"`
	From string `json:"from"`
	To   string `json:"to"`
}

// MessageResponse API消息响应
type MessageResponse struct {
	HydraMember []Message `json:"hydra:member"`
}

// EmailDetail 邮件详情数据结构
type EmailDetail struct {
	ID      string   `json:"id"`
	HTML    []string `json:"html"`
	Subject string   `json:"subject"`
}

// Database 数据库操作结构
type Database struct {
	dataPath string
	data     *Account
}

// NewDatabase 创建新的数据库实例
func NewDatabase(dataPath string) *Database {
	return &Database{
		dataPath: dataPath,
	}
}

// Read 读取账户数据
func (db *Database) Read() error {
	data, err := os.ReadFile(db.dataPath)
	if err != nil {
		if os.IsNotExist(err) {
			db.data = nil
			return nil
		}
		return fmt.Errorf("failed to read account file: %w", err)
	}

	var account Account
	if err := json.Unmarshal(data, &account); err != nil {
		return fmt.Errorf("failed to unmarshal account data: %w", err)
	}

	db.data = &account
	return nil
}

// Write 写入账户数据
func (db *Database) Write() error {
	if db.data == nil {
		return nil
	}

	data, err := json.MarshalIndent(db.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal account data: %w", err)
	}

	// 确保目录存在
	dir := filepath.Dir(db.dataPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(db.dataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write account file: %w", err)
	}

	return nil
}

// GetData 获取账户数据
func (db *Database) GetData() *Account {
	return db.data
}

// SetData 设置账户数据
func (db *Database) SetData(data *Account) {
	db.data = data
}

// DeleteData 删除账户数据文件
func (db *Database) DeleteData() error {
	if err := os.Remove(db.dataPath); err != nil {
		return fmt.Errorf("failed to delete account file: %w", err)
	}
	db.data = nil
	return nil
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// 如果加密随机数生成失败，回退到时间种子
			b[i] = charset[i%len(charset)]
		} else {
			b[i] = charset[num.Int64()]
		}
	}
	return string(b)
}

// Spinner 简单的加载动画
type Spinner struct {
	message string
	done    chan bool
}

// NewSpinner 创建新的加载动画
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message: message,
		done:    make(chan bool),
	}
}

// Start 开始加载动画
func (s *Spinner) Start() {
	go func() {
		chars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-s.done:
				fmt.Printf("\r%s\r", strings.Repeat(" ", len(s.message)+10))
				return
			default:
				fmt.Printf("\r%s %s", chars[i%len(chars)], s.message)
				time.Sleep(100 * time.Millisecond)
				i++
			}
		}
	}()
}

// Stop 停止加载动画
func (s *Spinner) Stop() {
	close(s.done)
}

// Color 简单的颜色输出函数
type Color struct{}

// Red 红色输出
func (c Color) Red(text string) string {
	return fmt.Sprintf("\033[31m%s\033[0m", text)
}

// Green 绿色输出
func (c Color) Green(text string) string {
	return fmt.Sprintf("\033[32m%s\033[0m", text)
}

// Blue 蓝色输出
func (c Color) Blue(text string) string {
	return fmt.Sprintf("\033[34m%s\033[0m", text)
}

// Underline 下划线输出
func (c Color) Underline(text string) string {
	return fmt.Sprintf("\033[4m%s\033[0m", text)
}
