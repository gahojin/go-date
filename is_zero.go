package date

// IsZeroer omitzero指定時に、IsZeroがtrueを返すフィールドをomitするためのインタフェース
type IsZeroer interface {
	IsZero() bool
}
