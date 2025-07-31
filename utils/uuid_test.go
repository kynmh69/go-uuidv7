package utils

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestGenerateUUIDv7(t *testing.T) {
	tests := []struct {
		name      string
		noHyphens bool
		wantRegex string
	}{
		{
			name:      "UUIDv7 with hyphens",
			noHyphens: false,
			wantRegex: `^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`,
		},
		{
			name:      "UUIDv7 without hyphens",
			noHyphens: true,
			wantRegex: `^[0-9a-f]{8}[0-9a-f]{4}7[0-9a-f]{3}[89ab][0-9a-f]{3}[0-9a-f]{12}$`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateUUIDv7(tt.noHyphens)

			// 正規表現でフォーマットをチェック
			matched, err := regexp.MatchString(tt.wantRegex, got)
			if err != nil {
				t.Fatalf("regex compilation failed: %v", err)
			}
			if !matched {
				t.Errorf("GenerateUUIDv7() = %v, does not match expected format %v", got, tt.wantRegex)
			}

			// ハイフンの有無をチェック
			if tt.noHyphens {
				if len(got) != 32 {
					t.Errorf("GenerateUUIDv7() with noHyphens=true should return 32 characters, got %d", len(got))
				}
				// ハイフンが含まれていないことを確認
				for _, char := range got {
					if char == '-' {
						t.Errorf("GenerateUUIDv7() with noHyphens=true should not contain hyphens, but got %v", got)
						break
					}
				}
			} else {
				if len(got) != 36 {
					t.Errorf("GenerateUUIDv7() with noHyphens=false should return 36 characters, got %d", len(got))
				}
				// 正しい位置にハイフンがあることを確認
				expectedHyphenPositions := []int{8, 13, 18, 23}
				for _, pos := range expectedHyphenPositions {
					if got[pos] != '-' {
						t.Errorf("GenerateUUIDv7() should have hyphen at position %d, but got %c", pos, got[pos])
					}
				}
			}

			// UUIDv7として有効かチェック（ハイフンありの場合のみ）
			if !tt.noHyphens {
				parsedUUID, err := uuid.Parse(got)
				if err != nil {
					t.Errorf("GenerateUUIDv7() returned invalid UUID: %v", err)
				}
				if parsedUUID.Version() != 7 {
					t.Errorf("GenerateUUIDv7() should return UUID version 7, got version %d", parsedUUID.Version())
				}
			}
		})
	}
}

func TestGenerateUUIDv7_Uniqueness(t *testing.T) {
	// 複数回生成して、すべて異なるUUIDが生成されることを確認
	const iterations = 100
	generated := make(map[string]bool)

	for i := 0; i < iterations; i++ {
		uuid1 := GenerateUUIDv7(false)
		uuid2 := GenerateUUIDv7(true)

		// 重複チェック
		if generated[uuid1] {
			t.Errorf("Duplicate UUID generated: %v", uuid1)
		}
		if generated[uuid2] {
			t.Errorf("Duplicate UUID generated: %v", uuid2)
		}

		generated[uuid1] = true
		generated[uuid2] = true
	}
}

func TestGenerateUUIDv7_ConsistentFormat(t *testing.T) {
	// 同じ引数で複数回呼び出した時に、フォーマットが一貫していることを確認
	for i := 0; i < 10; i++ {
		withHyphens := GenerateUUIDv7(false)
		withoutHyphens := GenerateUUIDv7(true)

		// ハイフンありの場合の長さチェック
		if len(withHyphens) != 36 {
			t.Errorf("Expected length 36 for UUID with hyphens, got %d", len(withHyphens))
		}

		// ハイフンなしの場合の長さチェック
		if len(withoutHyphens) != 32 {
			t.Errorf("Expected length 32 for UUID without hyphens, got %d", len(withoutHyphens))
		}
	}
}

// ベンチマークテスト
func BenchmarkGenerateUUIDv7_WithHyphens(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GenerateUUIDv7(false)
	}
}

func BenchmarkGenerateUUIDv7_WithoutHyphens(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GenerateUUIDv7(true)
	}
}

func BenchmarkGenerateMultipleUUIDs_Small(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GenerateMultipleUUIDs(10, false)
	}
}

func BenchmarkGenerateMultipleUUIDs_Medium(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GenerateMultipleUUIDs(100, false)
	}
}

func BenchmarkGenerateMultipleUUIDs_Large(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GenerateMultipleUUIDs(1000, false)
	}
}

func BenchmarkGenerateMultipleUUIDs_WithoutHyphens(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GenerateMultipleUUIDs(100, true)
	}
}

func TestGenerateMultipleUUIDs(t *testing.T) {
	tests := []struct {
		name      string
		count     int
		noHyphens bool
		wantLen   int
		wantRegex string
	}{
		{
			name:      "Generate 0 UUIDs",
			count:     0,
			noHyphens: false,
			wantLen:   0,
			wantRegex: "",
		},
		{
			name:      "Generate 1 UUID with hyphens",
			count:     1,
			noHyphens: false,
			wantLen:   1,
			wantRegex: `^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`,
		},
		{
			name:      "Generate 1 UUID without hyphens",
			count:     1,
			noHyphens: true,
			wantLen:   1,
			wantRegex: `^[0-9a-f]{8}[0-9a-f]{4}7[0-9a-f]{3}[89ab][0-9a-f]{3}[0-9a-f]{12}$`,
		},
		{
			name:      "Generate 5 UUIDs with hyphens",
			count:     5,
			noHyphens: false,
			wantLen:   5,
			wantRegex: `^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`,
		},
		{
			name:      "Generate 10 UUIDs without hyphens",
			count:     10,
			noHyphens: true,
			wantLen:   10,
			wantRegex: `^[0-9a-f]{8}[0-9a-f]{4}7[0-9a-f]{3}[89ab][0-9a-f]{3}[0-9a-f]{12}$`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateMultipleUUIDs(tt.count, tt.noHyphens)

			// スライスの長さをチェック
			if len(got) != tt.wantLen {
				t.Errorf("GenerateMultipleUUIDs() returned %d UUIDs, expected %d", len(got), tt.wantLen)
			}

			// 各UUIDのフォーマットをチェック（0個の場合はスキップ）
			if tt.count > 0 {
				for i, uuid := range got {
					matched, err := regexp.MatchString(tt.wantRegex, uuid)
					if err != nil {
						t.Fatalf("regex compilation failed: %v", err)
					}
					if !matched {
						t.Errorf("UUID at index %d: %v does not match expected format %v", i, uuid, tt.wantRegex)
					}

					// 長さもチェック
					expectedLen := 36
					if tt.noHyphens {
						expectedLen = 32
					}
					if len(uuid) != expectedLen {
						t.Errorf("UUID at index %d: expected length %d, got %d", i, expectedLen, len(uuid))
					}
				}
			}
		})
	}
}

func TestGenerateMultipleUUIDs_Uniqueness(t *testing.T) {
	// 大量のUUIDを生成して、すべて異なることを確認
	const count = 100
	uuids := GenerateMultipleUUIDs(count, false)

	if len(uuids) != count {
		t.Fatalf("Expected %d UUIDs, got %d", count, len(uuids))
	}

	// 重複チェック用のマップ
	seen := make(map[string]bool)
	for i, uuid := range uuids {
		if seen[uuid] {
			t.Errorf("Duplicate UUID found at index %d: %v", i, uuid)
		}
		seen[uuid] = true
	}
}

func TestGenerateMultipleUUIDs_NegativeCount(t *testing.T) {
	// 負の数や0を渡した場合の動作をテスト
	testCases := []struct {
		name  string
		count int
	}{
		{"negative count -1", -1},
		{"negative count -10", -10},
		{"zero count", 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uuids := GenerateMultipleUUIDs(tc.count, false)
			if len(uuids) != 0 {
				t.Errorf("Expected empty slice for count %d, got %d UUIDs", tc.count, len(uuids))
			}

			uuids = GenerateMultipleUUIDs(tc.count, true)
			if len(uuids) != 0 {
				t.Errorf("Expected empty slice for count %d (noHyphens=true), got %d UUIDs", tc.count, len(uuids))
			}
		})
	}
}

func TestGenerateMultipleUUIDs_LargeCount(t *testing.T) {
	// 大きな数でのテスト（パフォーマンスと正確性）
	const largeCount = 1000
	uuids := GenerateMultipleUUIDs(largeCount, false)

	if len(uuids) != largeCount {
		t.Fatalf("Expected %d UUIDs, got %d", largeCount, len(uuids))
	}

	// 最初と最後のUUIDのフォーマットをチェック
	uuidRegex := `^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`

	if matched, _ := regexp.MatchString(uuidRegex, uuids[0]); !matched {
		t.Errorf("First UUID does not match expected format: %v", uuids[0])
	}

	if matched, _ := regexp.MatchString(uuidRegex, uuids[largeCount-1]); !matched {
		t.Errorf("Last UUID does not match expected format: %v", uuids[largeCount-1])
	}
}

func TestGenerateMultipleUUIDs_ConsistentBehavior(t *testing.T) {
	// 同じパラメータで複数回呼び出した時の一貫性をテスト
	const count = 5

	for i := 0; i < 3; i++ {
		withHyphens := GenerateMultipleUUIDs(count, false)
		withoutHyphens := GenerateMultipleUUIDs(count, true)

		// 長さチェック
		if len(withHyphens) != count {
			t.Errorf("Iteration %d: expected %d UUIDs with hyphens, got %d", i, count, len(withHyphens))
		}
		if len(withoutHyphens) != count {
			t.Errorf("Iteration %d: expected %d UUIDs without hyphens, got %d", i, count, len(withoutHyphens))
		}

		// 各UUIDの長さチェック
		for j, uuid := range withHyphens {
			if len(uuid) != 36 {
				t.Errorf("Iteration %d, UUID %d with hyphens: expected length 36, got %d", i, j, len(uuid))
			}
		}
		for j, uuid := range withoutHyphens {
			if len(uuid) != 32 {
				t.Errorf("Iteration %d, UUID %d without hyphens: expected length 32, got %d", i, j, len(uuid))
			}
		}
	}
}

func TestPrintUUIDs(t *testing.T) {
	tests := []struct {
		name     string
		uuids    []string
		expected string
	}{
		{
			name:     "Empty slice",
			uuids:    []string{},
			expected: "",
		},
		{
			name:     "Single UUID",
			uuids:    []string{"550e8400-e29b-41d4-a716-446655440000"},
			expected: "550e8400-e29b-41d4-a716-446655440000\n",
		},
		{
			name: "Multiple UUIDs",
			uuids: []string{
				"550e8400-e29b-41d4-a716-446655440000",
				"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
				"6ba7b811-9dad-11d1-80b4-00c04fd430c8",
			},
			expected: "550e8400-e29b-41d4-a716-446655440000\n6ba7b810-9dad-11d1-80b4-00c04fd430c8\n6ba7b811-9dad-11d1-80b4-00c04fd430c8\n",
		},
		{
			name:     "UUIDs without hyphens",
			uuids:    []string{"550e8400e29b41d4a716446655440000", "6ba7b8109dad11d180b400c04fd430c8"},
			expected: "550e8400e29b41d4a716446655440000\n6ba7b8109dad11d180b400c04fd430c8\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 標準出力をキャプチャ
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// テスト対象の関数を実行
			PrintUUIDs(tt.uuids)

			// 出力をリストア
			if err := w.Close(); err != nil {
				t.Errorf("Failed to close writer: %v", err)
			}
			os.Stdout = old

			// 出力を読み取り
			var buf bytes.Buffer
			_, err := io.Copy(&buf, r)
			if err != nil {
				t.Fatalf("Failed to read output: %v", err)
			}
			got := buf.String()

			// 期待される出力と比較
			if got != tt.expected {
				t.Errorf("PrintUUIDs() output = %q, expected %q", got, tt.expected)
			}
		})
	}
}

func TestPrintUUIDs_WithRealUUIDs(t *testing.T) {
	// 実際にGenerateUUIDv7で生成したUUIDを使用してテスト
	uuids := make([]string, 3)
	for i := 0; i < 3; i++ {
		uuids[i] = GenerateUUIDv7(false)
	}

	// 標準出力をキャプチャ
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// テスト対象の関数を実行
	PrintUUIDs(uuids)

	// 出力をリストア
	if err := w.Close(); err != nil {
		t.Errorf("Failed to close writer: %v", err)
	}
	os.Stdout = old

	// 出力を読み取り
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}
	output := buf.String()

	// 各UUIDが出力に含まれていることを確認
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != len(uuids) {
		t.Errorf("Expected %d lines of output, got %d", len(uuids), len(lines))
	}

	for i, expectedUUID := range uuids {
		if i < len(lines) && lines[i] != expectedUUID {
			t.Errorf("Line %d: expected %q, got %q", i+1, expectedUUID, lines[i])
		}
	}
}

func BenchmarkPrintUUIDs(b *testing.B) {
	// ベンチマーク用のUUID生成
	uuids := make([]string, 100)
	for i := 0; i < 100; i++ {
		uuids[i] = GenerateUUIDv7(false)
	}

	// 標準出力を無効化（ベンチマーク結果に影響しないように）
	old := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = old }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PrintUUIDs(uuids)
	}
}
