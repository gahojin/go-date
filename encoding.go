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

// MarshalText は時刻をバイトスライスとしてマーシャルします。
// 未設定値（ゼロ値）の場合は空のバイトスライスを返します。
// 有効な時刻の場合は "HH:MM:SS"（ナノ秒がある場合は小数秒を含む）形式で返します。
func (t Time) MarshalText() ([]byte, error) {
	if !t.present {
		return []byte{}, nil
	}
	return []byte(t.String()), nil
}

// UnmarshalText はテキストからTimeをアンマーシャルします。
// テキストが空の場合は未設定のゼロ値を設定します。
// "HH:MM:SS" および小数秒付きの "HH:MM:SS.sssssssss" 形式を受け付けます。
func (t *Time) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*t = Time{}
		return nil
	}
	parsed, err := ParseTime("15:04:05.999999999", string(text))
	if err != nil {
		return err
	}
	*t = parsed
	return nil
}
