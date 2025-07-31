# go-uuidv7

UUIDv7（ドラフト版）を生成するGoコマンドラインツールです。

## 特徴

- UUIDv7（時間順序対応UUID）の生成
- 複数のUUID同時生成に対応
- ハイフンありまたはハイフンなしの形式選択可能
- 軽量で高速なコマンドラインインターフェース

## インストール

### Goを使ったインストール

```bash
go install github.com/kynmh69/go-uuidv7@latest
```

### ソースからビルド

```bash
git clone https://github.com/kynmh69/go-uuidv7.git
cd go-uuidv7
go build -o go-uuidv7
```

## 使用方法

### 基本的な使用方法

```bash
# 単一のUUIDを生成
go-uuidv7

# 5つのUUIDを生成
go-uuidv7 -n 5

# ハイフンなしで生成
go-uuidv7 -H

# 複数のUUIDをハイフンなしで生成
go-uuidv7 -n 3 -H
```

### オプション

| オプション | 短縮形 | デフォルト | 説明 |
|----------|-------|----------|------|
| `--number` | `-n` | 1 | 生成するUUIDの数 |
| `--no-hyphens` | `-H` | false | ハイフンを含めない |
| `--help` | `-h` | - | ヘルプを表示 |

## 使用例

```bash
# 1つのUUIDを生成
$ go-uuidv7
019e1a97-3b2c-7000-8000-123456789abc

# 3つのUUIDを生成
$ go-uuidv7 -n 3
019e1a97-3b2c-7000-8000-123456789abc
019e1a97-3b2c-7001-8000-123456789def
019e1a97-3b2c-7002-8000-123456789ghi

# ハイフンなしで生成
$ go-uuidv7 -H
019e1a973b2c7000800123456789abc
```

## UUIDv7について

UUIDv7は時間順序に基づくUUIDの新しい仕様（ドラフト版）です。従来のUUIDv4（ランダム）やUUIDv1（MACアドレス＋時間）と比較して、以下の利点があります：

- **時間順序性**: 生成時刻順にソートされる
- **プライバシー**: MACアドレスを含まない
- **データベース最適化**: インデックス効率が良い

## 技術仕様

- **Go バージョン**: 1.24.5+
- **依存関係**:
  - `github.com/google/uuid` v1.6.0
  - `github.com/alecthomas/kong` v1.12.1

## 開発

### テスト実行

```bash
go test ./...
```

### ビルド

```bash
go build -o go-uuidv7
```

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。詳細は [LICENSE](LICENSE) ファイルを参照してください。

## 貢献

プルリクエストやイシューの報告を歓迎します。貢献する前に、既存のイシューを確認してください。