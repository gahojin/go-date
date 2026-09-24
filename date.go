package date

import (
	"cmp"
	"encoding"
	"time"
)

// Date は年月日を表す構造体です。
type Date struct {
	year  int
	month time.Month
	day   int
}

// 型チェック
var (
	_ encoding.TextMarshaler   = (*Date)(nil)
	_ encoding.TextUnmarshaler = (*Date)(nil)
	_ IsZeroer                 = (*Date)(nil)
)

// New は指定された年・月・日からDateを生成して返します。
// time.Dateと同様に、正規化された日付（例: 2024年1月32日 -> 2024年2月1日）を生成します。
func New(year int, month time.Month, day int) Date {
	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return Date{
		year:  t.Year(),
		month: t.Month(),
		day:   t.Day(),
	}
}

// FromTime は time.Time から Date を生成して返します。
func FromTime(t time.Time) Date {
	return Date{
		year:  t.Year(),
		month: t.Month(),
		day:   t.Day(),
	}
}

// FromTimeIn は time.Time を指定されたタイムゾーン（nil指定時はtime.Local）に変換した上で Date を生成して返します。
func FromTimeIn(t time.Time, loc *time.Location) Date {
	if loc == nil {
		loc = time.Local
	}
	return FromTime(t.In(loc))
}

// Year は年を返します。
func (d Date) Year() int {
	return d.year
}

// Month は月を返します。
func (d Date) Month() time.Month {
	return d.month
}

// Day は日を返します。
func (d Date) Day() int {
	return d.day
}

// CompactString は日付を "YYYYMMDD" 形式（例: "20240501"）の文字列として返します。
func (d Date) CompactString() string {
	return d.Format("20060102")
}

// AddDate はDateに対して指定された年・月・日を加算した新しいDateを返します。
func (d Date) AddDate(years int, months int, days int) Date {
	t := d.ToTime(time.UTC).AddDate(years, months, days)
	return Date{
		year:  t.Year(),
		month: t.Month(),
		day:   t.Day(),
	}
}

// FirstDayOfMonth はその月の月初（1日）のDateを返します。
// ゼロ値の場合はゼロ値のDateを返します。
func (d Date) FirstDayOfMonth() Date {
	if d.IsZero() {
		return Date{}
	}
	return New(d.year, d.month, 1)
}

// LastDayOfMonth はその月の月末のDateを返します。
// ゼロ値の場合はゼロ値のDateを返します。
func (d Date) LastDayOfMonth() Date {
	if d.IsZero() {
		return Date{}
	}
	last := d.DaysIn()
	return Date{d.year, d.month, last}
}

// YearMonth は日付を "YYYYMM" 形式の文字列として返します。
func (d Date) YearMonth() string {
	return d.Format("200601")
}

// IsZero は日付がゼロ値（Year, Month, Dayがすべて0）であるかを判定します。
func (d Date) IsZero() bool {
	return d.year == 0 && d.month == 0 && d.day == 0
}

// Equal は2つのDateが同じ日付であるかを判定します。
func (d Date) Equal(other Date) bool {
	return d.year == other.year && d.month == other.month && d.day == other.day
}

// EqualMonth は2つのDateの年と月が同じであるかを判定します（日は比較しません）。
func (d Date) EqualMonth(other Date) bool {
	return d.year == other.year && d.month == other.month
}

// Before は自身の日付が指定された日付より前であるかを判定します。
func (d Date) Before(other Date) bool {
	c := d.Compare(other)
	return c < 0
}

// After は自身の日付が指定された日付より後であるかを判定します。
func (d Date) After(other Date) bool {
	c := other.Compare(d)
	return c < 0
}

// Compare は他のDateと日付を比較します。
// 自身がotherより前の場合は負の値、等しい場合は0、後の場合は正の値を返します。
func (d Date) Compare(other Date) int {
	c := cmp.Compare(d.year, other.year)
	if c == 0 {
		c = cmp.Compare(d.month, other.month)
	}
	if c == 0 {
		c = cmp.Compare(d.day, other.day)
	}
	return c
}

// Order は2つの日付の大小を比較し、昇順（早い日付、遅い日付）に並び替えて返します。
// aがbより後の場合は入れ替えて (b, a) を返します。
func Order(a, b Date) (Date, Date) {
	if b.Before(a) {
		return b, a
	}
	return a, b
}

// Order は自身とotherの大小を比較し、昇順（早い日付、遅い日付）に並び替えて返します。
// 自身がotherより後の場合は入れ替えて (other, d) を返します。
func (d Date) Order(other Date) (Date, Date) {
	return Order(d, other)
}

// SubMonths は自身とotherの間の月数（経過月数）を計算します。
// 日付（日）を考慮し、丸1ヶ月経過して初めて1ヶ月とカウントします（例: 1/15〜2/14は0、1/15〜2/15は1）。
// 2つの日付の順序に関わらず期間の月数（0以上の整数）を返します。
func (d Date) SubMonths(other Date) int {
	start, end := Order(d, other)

	m := (end.year-start.year)*12 + int(end.month-start.month)
	if end.day < start.day {
		daysInMonth := DaysIn(end.year, end.month)
		if end.day != daysInMonth || start.day < daysInMonth {
			m--
		}
	}
	return m
}

// SpanMonths は自身とotherが含まれるカレンダー上の月数（暦月数）を計算します。
// 日付（日）に関わらず開始月から終了月までの月数をカウントします（例: 1/15〜2/14は1月と2月で2）。
// 2つの日付の順序に関わらず期間の暦月数（正の整数）を返します。両方がゼロ値の場合は0を返します。
func (d Date) SpanMonths(other Date) int {
	if d.IsZero() && other.IsZero() {
		return 0
	}

	start, end := Order(d, other)

	return (end.year-start.year)*12 + int(end.month-start.month) + 1
}

// DaysIn は指定された年月の末日（日数）を返します。
func DaysIn(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// DaysIn は自身の日付の年月の末日（日数）を返します。
func (d Date) DaysIn() int {
	return DaysIn(d.year, d.month)
}

// Date は年・月・日をそれぞれ返します。
func (d Date) Date() (int, time.Month, int) {
	return d.year, d.month, d.day
}

// WithTime はDateに指定されたTimeおよびタイムゾーン（nil指定時はtime.Local）を組み合わせてtime.Timeを生成して返します。
func (d Date) WithTime(t Time, loc *time.Location) time.Time {
	if loc == nil {
		loc = time.Local
	}
	return time.Date(d.year, d.month, d.day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
}

// ToTime はDateを指定されたタイムゾーン（nil指定時はtime.Local、時刻は00:00:00）のtime.Timeに変換して返します。
func (d Date) ToTime(loc *time.Location) time.Time {
	if loc == nil {
		loc = time.Local
	}
	return time.Date(d.year, d.month, d.day, 0, 0, 0, 0, loc)
}

func (d Date) ToCompact() CompactDate {
	return CompactDate{d}
}

// Format は指定されたレイアウトに従って日付を文字列にフォーマットします。
func (d Date) Format(layout string) string {
	return d.ToTime(nil).Format(layout)
}

// String は日付を "YYYY-MM-DD" 形式（time.DateOnly）の文字列として返します。
func (d Date) String() string {
	return d.Format(time.DateOnly)
}

// Parse は指定されたレイアウトに従って文字列をDateにパースします。
func Parse(layout, value string) (Date, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return Date{}, err
	}
	return Date{
		year:  t.Year(),
		month: t.Month(),
		day:   t.Day(),
	}, nil
}
