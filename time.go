package date

import (
	"cmp"
	"encoding"
	"fmt"
	"strconv"
	"strings"
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

// divMod は a を b で割った商と、常に非負になる余りを返します。
// 例: divMod(-1, 60) は (-1, 59) を返します。
func divMod[T int | time.Duration](a, b T) (q, r T) {
	q = a / b
	r = a % b
	if r < 0 {
		q--
		r += b
	}
	return q, r
}

func durationToTime(d time.Duration) Time {
	h, rem := divMod(d, time.Hour)
	m, rem := divMod(rem, time.Minute)
	s, ns := divMod(rem, time.Second)
	return Time{hour: int(h), min: int(m), sec: int(s), nsec: int(ns), present: true}
}

// NewTime は指定された時・分・秒・ナノ秒からTimeを生成して返します。
// 正規化は分や秒、ナノ秒のみとし、時はそのまま保持する (25時は25時として保持)
func NewTime(hour, min, sec, nsec int) Time {
	var carry int
	carry, nsec = divMod(nsec, int(time.Second))
	sec += carry

	carry, sec = divMod(sec, 60)
	min += carry

	carry, min = divMod(min, 60)
	hour += carry

	return Time{hour: hour, min: min, sec: sec, nsec: nsec, present: true}
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

// TimeFromTimeIn は time.Time を指定されたタイムゾーン（nil指定時はtime.Local）に変換した上で Time を生成して返します。
func TimeFromTimeIn(t time.Time, loc *time.Location) Time {
	if loc == nil {
		loc = time.Local
	}
	return TimeFromTime(t.In(loc))
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
	return fmt.Sprintf("%02d%02d%02d", t.hour, t.min, t.sec)
}

// Duration は0時0分0秒からの経過時間をDurationとして返します。
func (t Time) Duration() time.Duration {
	return time.Duration(t.hour)*time.Hour + time.Duration(t.min)*time.Minute + time.Duration(t.sec)*time.Second + time.Duration(t.nsec)*time.Nanosecond
}

// Add はTimeに対して指定されたDurationを加算した新しいTimeを返します。
func (t Time) Add(d time.Duration) Time {
	return durationToTime(t.Duration() + d)
}

// Sub は自身とotherの間の時間差（Duration）を返します。
func (t Time) Sub(other Time) time.Duration {
	return t.Duration() - other.Duration()
}

// HourMinute は時刻を "HH:MM" 形式の文字列として返します
func (t Time) HourMinute() string {
	return fmt.Sprintf("%02d:%02d", t.hour, t.min)
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
	s := fmt.Sprintf("%02d:%02d:%02d", t.hour, t.min, t.sec)
	if t.nsec != 0 {
		s += strings.TrimRight(fmt.Sprintf(".%09d", t.nsec), "0")
	}
	return s
}

// ParseTime は"HH:MM:SS"（またはナノ秒付き "H:MM:SS.nnnnnnnnn"）形式の文字列をTimeにパースします。
// 24以上の時（例: "26:00:00"）もそのまま扱います。
func ParseTime(value string) (Time, error) {
	sec, nsecPart, hasFrac := strings.Cut(value, ".")

	hourPart, rest, ok := strings.Cut(sec, ":")
	if !ok {
		return Time{}, fmt.Errorf("date: invalid time %q", value)
	}
	minPart, secPart, ok := strings.Cut(rest, ":")
	if !ok {
		return Time{}, fmt.Errorf("date: invalid time %q", value)
	}

	var err error
	var h, m, s int
	if h, err = strconv.Atoi(hourPart); err != nil {
		return Time{}, fmt.Errorf("date: invalid time %q: %w", value, err)
	}
	if m, err = strconv.Atoi(minPart); err != nil {
		return Time{}, fmt.Errorf("date: invalid time %q: %w", value, err)
	}
	if s, err = strconv.Atoi(secPart); err != nil {
		return Time{}, fmt.Errorf("date: invalid time %q: %w", value, err)
	}

	ns := 0
	if hasFrac {
		if len(nsecPart) == 0 || len(nsecPart) > 9 {
			return Time{}, fmt.Errorf("date: invalid time %q: fractional part must be 1-9 digits", value)
		}
		padded := nsecPart + strings.Repeat("0", 9-len(nsecPart))
		ns, err = strconv.Atoi(padded)
		if err != nil {
			return Time{}, fmt.Errorf("date: invalid time %q: %w", value, err)
		}
	}

	return NewTime(h, m, s, ns), nil
}
