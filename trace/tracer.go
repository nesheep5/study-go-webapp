package trace

// ログ出力インターフェース
type Tracer interface {
	Trace(...interface{})
}
