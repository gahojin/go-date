# AGENT.md

このファイルは、本リポジトリ（`go-date`）で作業するAIエージェントおよび開発者のためのガイドラインです。

## プロジェクト概要

- **モジュール名**: `github.com/gahojin/go-date`
- **主要パッケージ**: `date` (`package date`)
- **Goバージョン**: Go 1.24 以上
- **説明**: Go用の日付操作ライブラリ

## 開発環境とツール

本プロジェクトではバージョン管理・ツール管理に **mise**、Gitフック管理に **lefthook**、コミットメッセージの検証に **commitlint** を採用しています。

- **mise 構成ツール**:
  - `go`: Go言語環境
  - `golangci-lint`: 静的解析およびコードフォーマッター
  - `gitleaks`: シークレット検出ツール
  - `lefthook`: Gitフックマネージャー
  - `node` / `@commitlint/cli` / `@commitlint/config-conventional`: Conventional Commits の検証

### シェル環境（mise のアクティベート）について
非対話型シェル（エージェント実行環境やCIなど）では、事前に `mise` をアクティベートすることで、`mise` 管理下のツール（`go`, `golangci-lint`, `lefthook` 等）を直接コマンド名で呼び出せるようになります。

```bash
eval "$(mise activate bash)"
```

※ または `mise exec --` を各コマンドに付与して実行することも可能です。

## コマンド一覧

`mise` をアクティベート（`eval "$(mise activate bash)"`）した状態であれば、直接各ツールを実行できます。

### テスト実行
```bash
# 基本実行
go test ./...

# 詳細出力付き
go test -v ./...

# カバレッジ測定
go test -cover ./...

# または mise exec 経由
mise exec -- go test ./...
```

### リント & フォーマット
```bash
# 静的解析（golangci-lint）
golangci-lint run ./...

# フォーマッター適用
golangci-lint fmt ./...

# または標準の go fmt
go fmt ./...
```

### Gitフック実行 / セキュリティチェック
```bash
# pre-commit フックの手動実行
lefthook run pre-commit

# gitleaks によるシークレットスキャン
gitleaks git --pre-commit --staged --no-color --no-banner --verbose
```

## 設計方針・コーディング規約

1. **イミュータブルな設計**:
   - `Date` 構造体（`Year`, `Month`, `Day`）は値レシーバを基本とし、副作用のない純粋なメソッドを提供します。
2. **標準ライブラリとの整合性**:
   - `time.Time` や `time.Month` など標準パッケージの型・定数との相互運用性を重視します。
   - 比較処理には `cmp.Compare` などを活用します。
3. **テスト方針**:
   - アサーションには `github.com/stretchr/testify/assert`（致命的な事前条件には `require`）を使用します。`t.Errorf` による手動比較は避けてください。
   - テーブル駆動テスト（Table-Driven Tests）を標準とします。
   - 原則としてパッケージ外テスト（`package date_test`）を採用し、公開APIのインターフェースおよび使いやすさを検証します。
   - ゼロ値（`Date{}`）や境界値（年末年始、うるう年など）に対するテストケースを網羅します。
4. **ドキュメント・コメント**:
   - 公開される型・関数・メソッドには適切なGoDocコメントを記述します。

## コミット規約

コミットメッセージは **Conventional Commits** 規約に準拠する必要があります（`commitlint` によって検証されます）。

- **フォーマット**: `<type>(<scope>): <subject>` （`scope` は省略可能）
- **主な type**:
  - `feat`: 新機能の追加
  - `fix`: バグ修正
  - `docs`: ドキュメントの変更
  - `test`: テストの追加・修正
  - `refactor`: リファクタリング（機能の変更やバグ修正を含まないコード変更）
  - `chore`: ビルドツールや補助設定の更新
