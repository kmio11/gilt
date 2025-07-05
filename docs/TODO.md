# 未使用ファイル削除機能 実装計画

## 実装方針

### アプローチ
Listerインターフェースを使用した汎用的なファイル管理機能の実装

### 設計原則
- Red-Green-Refactorサイクルでテストファーストの開発
- 既存アーキテクチャとの自然な統合
- シンプルで汎用性の高い設計

## 実装タスク

### Cycle 1: Listerインターフェース基礎
- [ ] `TestGoldenFile_ListFiles`失敗テストを作成
- [ ] `Lister`インターフェースを`golden.go`に定義
- [ ] `GoldenFile`構造体に`ListFiles`メソッドの最小実装を追加
- [ ] テストが通ることを確認
- [ ] `filepath.Glob`を使用した実際のファイル検索実装
- [ ] エラーハンドリングの追加
- [ ] リファクタリング後のテスト実行

### Cycle 2: ファイル追跡機能
- [ ] `TestFileTracker`失敗テストを作成
- [ ] `FileTracker`構造体を定義
- [ ] 基本的な追跡機能の実装
- [ ] テストが通ることを確認
- [ ] `Read`/`Write`メソッドとの統合
- [ ] 自動追跡機能の実装
- [ ] スレッドセーフティの考慮
- [ ] リファクタリング後のテスト実行

### Cycle 3: クリーン機能
- [ ] `TestCleanUnusedGoldenFiles`失敗テストを作成
- [ ] `CleanUnusedGoldenFiles`関数の基本実装
- [ ] テストが通ることを確認
- [ ] `CleanOption`構造体の定義
- [ ] dry-run機能の実装
- [ ] バックアップ機能の実装
- [ ] 安全性チェックの追加
- [ ] ログ出力機能の実装
- [ ] リファクタリング後のテスト実行

### Cycle 4: 統合とユーザビリティ
- [ ] 統合テストの作成
- [ ] エンドツーエンドテストの作成
- [ ] 既存APIとの統合
- [ ] テストが通ることを確認
- [ ] APIの改善
- [ ] エラーハンドリングの強化
- [ ] ドキュメントの作成
- [ ] 最終テスト実行

## 実装詳細

### 新しいインターフェース

```go
// ファイル一覧取得インターフェース
type Lister interface {
    ListFiles(t *testing.T, namespace string) []string
}

// ファイル追跡
type FileTracker struct {
    usedFiles map[string]bool
    mutex     sync.RWMutex
}

// クリーンオプション
type CleanOption struct {
    DryRun     bool
    Backup     bool
    Interactive bool
    Patterns   []string
}
```

### 主要関数

```go
// 未使用ファイル削除のメイン関数
func CleanUnusedGoldenFiles(golden *Golden, t *testing.T, namespace string, options ...CleanOption) error

// ファイル追跡のヘルパー関数
func TrackFileUsage(filePath string)
func IsFileUsed(filePath string) bool
```

## テスト戦略

### テストファイル
- `lister_test.go` - Listerインターフェースのテスト
- `filetracker_test.go` - ファイル追跡機能のテスト
- `clean_test.go` - クリーン機能のテスト
- `integration_test.go` - 統合テスト

### テストケース
- 正常系: ファイル一覧取得、未使用ファイル削除
- 異常系: 存在しないディレクトリ、権限エラー
- エッジケース: 空のディレクトリ、同名ファイル
- 統合テスト: 既存機能との組み合わせ

## 品質基準

- [ ] 全テストが通過する
- [ ] カバレッジが85%以上
- [ ] 既存のテストが破綻しない
- [ ] `go fmt`、`go vet`でエラーが出ない
- [ ] 既存のAPIとの後方互換性を保つ