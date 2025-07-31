package main

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

func TestMainCommand(t *testing.T) {
	// Goバイナリをビルド
	cmd := exec.Command("go", "build", "-o", "test-binary", ".")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("test-binary")

	tests := []struct {
		name          string
		args          []string
		expectError   bool
		expectedLines int
		errorMessage  string
	}{
		{
			name:          "Default behavior (1 UUID)",
			args:          []string{},
			expectError:   false,
			expectedLines: 1,
		},
		{
			name:          "Generate 3 UUIDs",
			args:          []string{"--number", "3"},
			expectError:   false,
			expectedLines: 3,
		},
		{
			name:          "Generate 5 UUIDs short flag",
			args:          []string{"-n", "5"},
			expectError:   false,
			expectedLines: 5,
		},
		{
			name:          "Generate UUIDs without hyphens",
			args:          []string{"--number", "2", "--no-hyphens"},
			expectError:   false,
			expectedLines: 2,
		},
		{
			name:         "Negative number should fail",
			args:         []string{"--number=-1"},
			expectError:  true,
			errorMessage: "Number must be a positive integer",
		},
		{
			name:         "Zero should fail",
			args:         []string{"--number", "0"},
			expectError:  true,
			errorMessage: "Number must be a positive integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./test-binary", tt.args...)
			output, err := cmd.CombinedOutput()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but command succeeded")
				}
				if !strings.Contains(string(output), tt.errorMessage) {
					t.Errorf("Expected error message containing %q, got: %s", tt.errorMessage, string(output))
				}
			} else {
				if err != nil {
					t.Errorf("Command failed: %v, output: %s", err, string(output))
				}

				// 出力行数をチェック
				lines := strings.Split(strings.TrimSpace(string(output)), "\n")
				if len(lines) != tt.expectedLines {
					t.Errorf("Expected %d lines of output, got %d", tt.expectedLines, len(lines))
				}

				// UUIDフォーマットをチェック（ハイフンありの場合）
				if !contains(tt.args, "--no-hyphens") && !contains(tt.args, "-H") {
					for i, line := range lines {
						if len(line) != 36 {
							t.Errorf("Line %d: expected UUID with hyphens (36 chars), got %d chars: %s", i+1, len(line), line)
						}
					}
				} else {
					// ハイフンなしの場合
					for i, line := range lines {
						if len(line) != 32 {
							t.Errorf("Line %d: expected UUID without hyphens (32 chars), got %d chars: %s", i+1, len(line), line)
						}
					}
				}
			}
		})
	}
}

func TestMainHelp(t *testing.T) {
	// Goバイナリをビルド
	cmd := exec.Command("go", "build", "-o", "test-binary", ".")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("test-binary")

	// ヘルプメッセージのテスト
	cmd = exec.Command("./test-binary", "--help")
	output, err := cmd.CombinedOutput()

	// --helpは正常終了する
	if err != nil {
		t.Errorf("Help command failed: %v", err)
	}

	helpOutput := string(output)
	expectedStrings := []string{
		"Generate UUIDv7 (draft)",
		"--number",
		"--no-hyphens",
		"Number of UUIDs to generate",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(helpOutput, expected) {
			t.Errorf("Help output should contain %q, got: %s", expected, helpOutput)
		}
	}
}

func TestMainLargeNumber(t *testing.T) {
	// Goバイナリをビルド
	cmd := exec.Command("go", "build", "-o", "test-binary", ".")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("test-binary")

	// 大きな数でのテスト
	largeNum := 100
	cmd = exec.Command("./test-binary", "--number", strconv.Itoa(largeNum))
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Command with large number failed: %v, output: %s", err, string(output))
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) != largeNum {
		t.Errorf("Expected %d lines of output, got %d", largeNum, len(lines))
	}
}

// ヘルパー関数：スライスに文字列が含まれているかチェック
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
