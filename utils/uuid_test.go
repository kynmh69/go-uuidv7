package utils

import (
	"regexp"
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
