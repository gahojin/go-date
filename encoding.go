package date

import (
	"encoding"
	"time"
)

var (
	_ encoding.TextMarshaler   = (*Date)(nil)
	_ encoding.TextUnmarshaler = (*Date)(nil)
	_ encoding.TextMarshaler   = (*Time)(nil)
	_ encoding.TextUnmarshaler = (*Time)(nil)
)

// MarshalText は日付を "YYYY-MM-DD" 形式のバイトスライスとしてマーシャルします。
// ゼロ値の場合は空のバイトスライスを返します。
func (d Date) MarshalText() ([]byte, error) {
	if d.IsZero() {
		return []byte{}, nil
	}
	return []byte(d.String()), nil
}

// UnmarshalText は "YYYY-MM-DD" 形式のテキストからDateをアンマーシャルします。
// テキストが空の場合はゼロ値を設定します。
func (d *Date) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*d = Date{}
		return nil
	}
	parsed, err := Parse(time.DateOnly, string(text))
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

// MarshalText は時刻を "HH:MM:SS" 形式のバイトスライスとしてマーシャルします。
// ゼロ値の場合は空のバイトスライスを返します。
func (t Time) MarshalText() ([]byte, error) {
	if t.IsZero() {
		return []byte{}, nil
	}
	return []byte(t.String()), nil
}

// UnmarshalText は "HH:MM:SS" 形式のテキストからTimeをアンマーシャルします。
// テキストが空の場合はゼロ値を設定します。
func (t *Time) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*t = Time{}
		return nil
	}
	parsed, err := ParseTime(time.TimeOnly, string(text))
	if err != nil {
		return err
	}
	*t = parsed
	return nil
}
