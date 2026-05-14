package jwt

import (
	"strings"
	"testing"
	"time"
)

// TestGenToken 测试生成的 token
func TestGenToken(t *testing.T) {
	// 1.准备数据
	type data struct {
		name     string
		userid   int64
		username string
		wantErr  bool
	}

	input := []data{
		{"valid user", 1, "singdile", false},
		{"zero userid", 0, "monica", true},
		{"empty username", 2, "", true},
		{"negative userid", -1, "jimmy", true},
	}

	// 2.执行函数
	for _, tt := range input {
		t.Run(tt.name, func(t *testing.T) {
			tokenstring, err := GenToken(tt.userid, tt.username)

			//验证错误
			if tt.wantErr {
				if err == nil {
					t.Error("expected error,got nil")
				}
				return //期望错误的情况，但是没有错误，所以退出
			}

			//验证成功
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tokenstring == "" {
				t.Error("token is empty")
			}

			//3.验证结果 验证 jwt 格式
			parts := strings.Split(tokenstring, ".")
			if len(parts) != 3 {
				t.Errorf("invalid JWT format: got %d parts", len(parts))
			}

		})
	}

}

// TestParseToken
func TestParseToken(t *testing.T) {
	// 1.准备数据
	userid := int64(32)
	username := "pite"
	tokenstring, err := GenToken(userid, username)

	if err != nil {
		t.Fatalf("userid = %d, username = %s, GenToken failed:%v", userid, username, err)
	}

	input := []struct {
		name        string //子测试的名字
		tokenstring string //输入参数
		wantErr     bool   //期望错误
	}{
		{"valid token", tokenstring, false},
		{"empty token", "", true},
		{"invalid token", "not.ajwt.token", true},
		{"invalid signature", tokenstring + "wrong", true},
	}
	// 2.执行函数
	for _, tt := range input {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ParseToken(tt.tokenstring)

			// 3.验证结果
			//验证错误情况
			if tt.wantErr {
				if err == nil {
					t.Error("expected error,got nil")
				}
				return
			}

			//验证成功情况
			if err != nil {
				t.Fatalf("unexpected err:%v", err)
			}

			//验证结果
			if claims.UserID != userid {
				t.Errorf("userid = %d,want %d", claims.UserID, userid)
			}

			if claims.Username != username {
				t.Errorf("username = %s,want %s", claims.Username, username)
			}

		})
	}

}

// TestTokenFlow 测试完整 token 流程
func TestTokenFlow(t *testing.T) {
	// 1. 准备数据
	userID := int64(12345)
	username := "flowtest"

	// 2. 生成 token
	token, err := GenToken(userID, username)
	if err != nil {
		t.Fatalf("GenToken failed: %v", err)
	}
	t.Logf("Generated token: %s", token[:50]+"...")

	// 3. 解析 token
	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}

	// 4. 验证数据一致性
	if claims.UserID != userID {
		t.Errorf("UserID mismatch: got %d, want %d", claims.UserID, userID)
	}
	if claims.Username != username {
		t.Errorf("Username mismatch: got %s, want %s", claims.Username, username)
	}

	// 5. 验证 Issuer
	if claims.Issuer != "bluebell" {
		t.Errorf("Issuer mismatch: got %s, want bluebell", claims.Issuer)
	}

	// 6. 验证 token 未过期
	if claims.ExpiresAt.Time.Before(time.Now()) {
		t.Error("token is already expired")
	}

	t.Log("Token flow test passed!")
}
