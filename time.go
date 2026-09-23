package date

import (
	"cmp"
	"encoding"
	"time"
)

// Time は時・分・秒・ナノ秒を表す構造体です。
type Time struct {
	hour    int
	min     int
	sec     int
	nsec    int
	present bool
}

// 型チェック
var (
	_ encoding.TextMarshaler   = (*Time)(nil)
	_ encoding.TextUnmarshaler = (*Time)(nil)
	_ IsZeroer                 = (*Time)(nil)
)

// NewTime は指定された時・分・秒・ナノ秒からTimeを生成して返します。
// time.Dateと同様に、正規化された時刻（例: 25時 -> 1時、-1分 -> 前の時間の59分）を生成します。
func NewTime(hour, min, sec, nsec int) Time {
	t := time.Date(2000, 1, 1, hour, min, sec, nsec, time.UTC)
	return Time{
		hour:    t.Hour(),
		min:     t.Minute(),
		sec:     t.Second(),
		nsec:    t.Nanosecond(),
		present: true,
	}
}

// TimeFromTime は time.Time から Time を生成して返します。
func TimeFromTime(t time.Time) Time {
	return Time{
		hour:    t.Hour(),
		min:     t.Minute(),
		sec:     t.Second(),
		nsec:    t.Nanosecond(),
		present: true,
	}
}

// Hour は時を返します。
func (t Time) Hour() int {
	return t.hour
}

// Minute は分を返します。
func (t Time) Minute() int {
	return t.min
}

// Second は秒を返します。
func (t Time) Second() int {
	return t.sec
}

// Nanosecond はナノ秒を返します。
func (t Time) Nanosecond() int {
	return t.nsec
}

// Clock は時・分・秒をそれぞれ返します。
func (t Time) Clock() (int, int, int) {
	return t.hour, t.min, t.sec
}

// CompactString は時刻を "HHMMSS" 形式（例: "150405"）の文字列として返します。
func (t Time) CompactString() string {
	return t.Format("150405")
}

// Add はTimeに対して指定されたDurationを加算した新しいTimeを返します。
func (t Time) Add(d time.Duration) Time {
	base := t.ToTime(time.UTC)
	return TimeFromTime(base.Add(d))
}

// Sub は自身とotherの間の時間差（Duration）を返します。
func (t Time) Sub(other Time) time.Duration {
	d1 := time.Duration(t.hour)*time.Hour + time.Duration(t.min)*time.Minute + time.Duration(t.sec)*time.Second + time.Duration(t.nsec)*time.Nanosecond
	d2 := time.Duration(other.hour)*time.Hour + time.Duration(other.min)*time.Minute + time.Duration(other.sec)*time.Second + time.Duration(other.nsec)*time.Nanosecond
	return d1 - d2
}

// IsZero は時刻がゼロ値（時・分・秒・ナノ秒がすべて0）であるかを判定します。
func (t Time) IsZero() bool {
	return t.hour == 0 && t.min == 0 && t.sec == 0 && t.nsec == 0
}

// Equal は2つのTimeが同じ時刻であるかを判定します。
func (t Time) Equal(other Time) bool {
	return t.hour == other.hour && t.min == other.min && t.sec == other.sec && t.nsec == other.nsec
}

// Before は自身の時刻が指定された時刻より前であるかを判定します。
func (t Time) Before(other Time) bool {
	return t.Compare(other) < 0
}

// After は自身の時刻が指定された時刻より後であるかを判定します。
func (t Time) After(other Time) bool {
	return other.Compare(t) < 0
}

// Compare は他のTimeと時刻を比較します。
// 自身がotherより前の場合は負の値、等しい場合は0、後の場合は正の値を返します。
func (t Time) Compare(other Time) int {
	c := cmp.Compare(t.hour, other.hour)
	if c == 0 {
		c = cmp.Compare(t.min, other.min)
	}
	if c == 0 {
		c = cmp.Compare(t.sec, other.sec)
	}
	if c == 0 {
		c = cmp.Compare(t.nsec, other.nsec)
	}
	return c
}

// OrderTime は2つの時刻の大小を比較し、昇順（早い時刻、遅い時刻）に並び替えて返します。
// aがbより後の場合は入れ替えて (b, a) を返します。
func OrderTime(a, b Time) (Time, Time) {
	if b.Before(a) {
		return b, a
	}
	return a, b
}

// Order は自身とotherの大小を比較し、昇順（早い時刻、遅い時刻）に並び替えて返します。
// 自身がotherより後の場合は入れ替えて (other, t) を返します。
func (t Time) Order(other Time) (Time, Time) {
	return OrderTime(t, other)
}

// ToTime はTimeを指定されたタイムゾーン（nil指定時はtime.Local、日付は0000-01-01）のtime.Timeに変換して返します。
func (t Time) ToTime(loc *time.Location) time.Time {
	if loc == nil {
		loc = time.Local
	}
	return time.Date(0, 1, 1, t.hour, t.min, t.sec, t.nsec, loc)
}

// Format は指定されたレイアウトに従って時刻を文字列にフォーマットします。
func (t Time) Format(layout string) string {
	return t.ToTime(nil).Format(layout)
}

// String は時刻を文字列として返します。ナノ秒がある場合は小数秒を含みます（例: "15:04:05.123456789"）。
func (t Time) String() string {
	if t.nsec != 0 {
		return t.Format("15:04:05.999999999")
	}
	return t.Format(time.TimeOnly)
}

// ParseTime は指定されたレイアウトに従って文字列をTimeにパースします。
func ParseTime(layout, value string) (Time, error) {
	parsed, err := time.Parse(layout, value)
	if err != nil {
		return Time{}, err
	}
	return Time{
		hour:    parsed.Hour(),
		min:     parsed.Minute(),
		sec:     parsed.Second(),
		nsec:    parsed.Nanosecond(),
		present: true,
	}, nil
}
