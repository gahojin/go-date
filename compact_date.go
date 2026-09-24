package date

import (
	"encoding"
	"time"
)

// CompactDate 日付情報
type CompactDate struct {
	Date
}

// 型チェック
var (
	_ encoding.TextMarshaler   = (*CompactDate)(nil)
	_ encoding.TextUnmarshaler = (*CompactDate)(nil)
	_ IsZeroer                 = (*CompactDate)(nil)
)

func NewCompact(year int, month time.Month, day int) CompactDate {
	return CompactDate{New(year, month, day)}
}

// CompactFromTime は time.Time から CompactDate を生成して返します。
func CompactFromTime(t time.Time) CompactDate {
	return CompactDate{FromTime(t)}
}

// CompactFromTimeIn は time.Time を指定されたタイムゾーン（nil指定時はtime.Local）に変換した上で CompactDate を生成して返します。
func CompactFromTimeIn(t time.Time, loc *time.Location) CompactDate {
	return CompactDate{FromTimeIn(t, loc)}
}

// String は日付を "YYYYMMDD" 形式（例: "20240501"）の文字列として返します。
func (d CompactDate) String() string {
	return d.CompactString()
}

// ParseCompact は"YYYYMMDD" 形式に従って文字列をCompactDateにパースします。
func ParseCompact(value string) (CompactDate, error) {
	d, err := Parse("20060102", value)
	if err != nil {
		return CompactDate{}, err
	}
	return d.ToCompact(), nil
}
