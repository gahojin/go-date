package date_test

import (
	"testing"
	"time"

	"github.com/gahojin/go-date"
	"github.com/stretchr/testify/assert"
)

// TestNewTime は NewTime コンストラクタが正しく Time を生成し正規化するかをテストする。
func TestNewTime(t *testing.T) {
	tests := []struct {
		name     string
		hour     int
		min      int
		sec      int
		nsec     int
		wantHour int
		wantMin  int
		wantSec  int
		wantNsec int
	}{
		{
			name:     "standard time",
			hour:     15,
			min:      30,
			sec:      45,
			nsec:     123456789,
			wantHour: 15,
			wantMin:  30,
			wantSec:  45,
			wantNsec: 123456789,
		},
		{
			name:     "zero time",
			hour:     0,
			min:      0,
			sec:      0,
			nsec:     0,
			wantHour: 0,
			wantMin:  0,
			wantSec:  0,
			wantNsec: 0,
		},
		{
			name:     "hour is not normalized (25 hour stays 25)",
			hour:     25,
			min:      0,
			sec:      0,
			nsec:     0,
			wantHour: 25,
			wantMin:  0,
			wantSec:  0,
			wantNsec: 0,
		},
		{
			name:     "hour beyond 24 (26 hour stays 26)",
			hour:     26,
			min:      0,
			sec:      0,
			nsec:     0,
			wantHour: 26,
			wantMin:  0,
			wantSec:  0,
			wantNsec: 0,
		},
		{
			name:     "normalize 61 minutes",
			hour:     10,
			min:      61,
			sec:      0,
			nsec:     0,
			wantHour: 11,
			wantMin:  1,
			wantSec:  0,
			wantNsec: 0,
		},
		{
			name:     "normalize negative minute",
			hour:     10,
			min:      -1,
			sec:      0,
			nsec:     0,
			wantHour: 9,
			wantMin:  59,
			wantSec:  0,
			wantNsec: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := date.NewTime(tt.hour, tt.min, tt.sec, tt.nsec)
			assert.Equal(t, tt.wantHour, tm.Hour())
			assert.Equal(t, tt.wantMin, tm.Minute())
			assert.Equal(t, tt.wantSec, tm.Second())
			assert.Equal(t, tt.wantNsec, tm.Nanosecond())
		})
	}
}

// TestTimeFromTime は TimeFromTime 関数が time.Time から Time を正しく生成するかをテストする。
func TestTimeFromTime(t *testing.T) {
	jst := time.FixedZone("JST", 9*60*60)

	tests := []struct {
		name string
		t    time.Time
		want date.Time
	}{
		{
			name: "UTC time",
			t:    time.Date(2024, 5, 1, 15, 30, 45, 123456789, time.UTC),
			want: date.NewTime(15, 30, 45, 123456789),
		},
		{
			name: "keeps the time.Time's own location (JST, no conversion)",
			t:    time.Date(2024, 5, 1, 23, 30, 0, 0, jst),
			want: date.NewTime(23, 30, 0, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := date.TimeFromTime(tt.t)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestTimeFromTimeIn は TimeFromTimeIn 関数をテストする。
func TestTimeFromTimeIn(t *testing.T) {
	jst := time.FixedZone("JST", 9*60*60)

	originalLocal := time.Local
	time.Local = time.FixedZone("UTC+12", 12*60*60)
	defer func() { time.Local = originalLocal }()

	tests := []struct {
		name string
		tm   time.Time
		loc  *time.Location
		want date.Time
	}{
		{
			name: "converts UTC to JST",
			tm:   time.Date(2024, 5, 1, 15, 0, 0, 0, time.UTC),
			loc:  jst,
			want: date.NewTime(0, 0, 0, 0), // 15:00 UTC = 翌日00:00 JST
		},
		{
			name: "converts JST to UTC",
			tm:   time.Date(2024, 5, 1, 9, 0, 0, 0, jst),
			loc:  time.UTC,
			want: date.NewTime(0, 0, 0, 0), // 09:00 JST = 当日00:00 UTC
		},
		{
			name: "nil location defaults to time.Local",
			tm:   time.Date(2024, 5, 1, 15, 30, 45, 0, time.UTC),
			loc:  nil,
			want: date.NewTime(3, 30, 45, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := date.TimeFromTimeIn(tt.tm, tt.loc)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestTime_Clock は Clock メソッドが時・分・秒を正しく返すかをテストする。
func TestTime_Clock(t *testing.T) {
	tm := date.NewTime(15, 30, 45, 100)
	h, m, s := tm.Clock()
	assert.Equal(t, 15, h)
	assert.Equal(t, 30, m)
	assert.Equal(t, 45, s)
}

// TestTime_CompactString は CompactString メソッドが HHMMSS 形式の文字列を返すかをテストする。
func TestTime_CompactString(t *testing.T) {
	tm := date.NewTime(9, 5, 3, 0)
	assert.Equal(t, "090503", tm.CompactString())

	tm2 := date.NewTime(15, 30, 45, 0)
	assert.Equal(t, "153045", tm2.CompactString())

	tm3 := date.NewTime(26, 63, 45, 0)
	assert.Equal(t, "270345", tm3.CompactString())
}

// TestTime_Add は Add メソッドが Duration を正しく加算した Time を返すかをテストする。
func TestTime_Add(t *testing.T) {
	tm := date.NewTime(10, 0, 0, 0)

	got := tm.Add(1*time.Hour + 30*time.Minute + 15*time.Second)
	assert.Equal(t, date.NewTime(11, 30, 15, 0), got)

	// 日跨ぎ
	got2 := tm.Add(15 * time.Hour)
	assert.Equal(t, date.NewTime(25, 0, 0, 0), got2)

	// 負の加算
	got3 := tm.Add(-2 * time.Hour)
	assert.Equal(t, date.NewTime(8, 0, 0, 0), got3)

	// 結果が負になる
	got4 := tm.Add(-12 * time.Hour)
	assert.Equal(t, date.NewTime(-2, 0, 0, 0), got4)
}

// TestTime_Sub は Sub メソッドが2つの Time の差（Duration）を正しく計算できるかをテストする。
func TestTime_Sub(t *testing.T) {
	t1 := date.NewTime(15, 30, 0, 0)
	t2 := date.NewTime(14, 0, 0, 0)

	assert.Equal(t, 1*time.Hour+30*time.Minute, t1.Sub(t2))
	assert.Equal(t, -(1*time.Hour + 30*time.Minute), t2.Sub(t1))
}

// TestTime_HourMinute は Time.HourMinute メソッドをテストする。
func TestTime_HourMinute(t *testing.T) {
	tests := []struct {
		name string
		tm   date.Time
		want string
	}{
		{
			name: "standard time",
			tm:   date.NewTime(15, 30, 45, 0),
			want: "15:30",
		},
		{
			name: "single digit hour",
			tm:   date.NewTime(9, 5, 0, 0),
			want: "09:05",
		},
		{
			name: "zero time",
			tm:   date.Time{},
			want: "00:00",
		},
		{
			name: "hour beyond 24",
			tm:   date.NewTime(26, 0, 0, 0),
			want: "26:00",
		},
		{
			name: "seconds and nanoseconds are ignored",
			tm:   date.NewTime(15, 30, 45, 123456789),
			want: "15:30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.tm.HourMinute())
		})
	}
}

// TestTime_IsZero は IsZero メソッドがゼロ値を正しく判定できるかをテストする。
func TestTime_IsZero(t *testing.T) {
	var zero date.Time
	assert.True(t, zero.IsZero())

	tm := date.NewTime(0, 0, 0, 0)
	assert.True(t, tm.IsZero())

	tm2 := date.NewTime(0, 0, 1, 0)
	assert.False(t, tm2.IsZero())
}

// TestTime_Equal は Equal メソッドが2つの Time が等しいかを正しく判定できるかをテストする。
func TestTime_Equal(t *testing.T) {
	t1 := date.NewTime(15, 30, 45, 100)
	t2 := date.NewTime(15, 30, 45, 100)
	t3 := date.NewTime(15, 30, 45, 200)

	assert.True(t, t1.Equal(t2))
	assert.False(t, t1.Equal(t3))
}

// TestTime_Compare は Compare メソッドが2つの Time の大小を正しく判定できるかをテストする。
func TestTime_Compare(t *testing.T) {
	tests := []struct {
		name string
		a    date.Time
		b    date.Time
		want int
	}{
		{
			name: "equal",
			a:    date.NewTime(15, 30, 45, 100),
			b:    date.NewTime(15, 30, 45, 100),
			want: 0,
		},
		{
			name: "earlier hour",
			a:    date.NewTime(14, 30, 45, 100),
			b:    date.NewTime(15, 30, 45, 100),
			want: -1,
		},
		{
			name: "later hour",
			a:    date.NewTime(16, 30, 45, 100),
			b:    date.NewTime(15, 30, 45, 100),
			want: 1,
		},
		{
			name: "earlier minute",
			a:    date.NewTime(15, 29, 45, 100),
			b:    date.NewTime(15, 30, 45, 100),
			want: -1,
		},
		{
			name: "later minute",
			a:    date.NewTime(15, 31, 45, 100),
			b:    date.NewTime(15, 30, 45, 100),
			want: 1,
		},
		{
			name: "earlier second",
			a:    date.NewTime(15, 30, 44, 100),
			b:    date.NewTime(15, 30, 45, 100),
			want: -1,
		},
		{
			name: "later second",
			a:    date.NewTime(15, 30, 46, 100),
			b:    date.NewTime(15, 30, 45, 100),
			want: 1,
		},
		{
			name: "earlier nsec",
			a:    date.NewTime(15, 30, 45, 99),
			b:    date.NewTime(15, 30, 45, 100),
			want: -1,
		},
		{
			name: "later nsec",
			a:    date.NewTime(15, 30, 45, 101),
			b:    date.NewTime(15, 30, 45, 100),
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.a.Compare(tt.b))
		})
	}
}

// TestTime_Before_After は Before および After メソッドをテストする。
func TestTime_Before_After(t *testing.T) {
	t1 := date.NewTime(10, 0, 0, 0)
	t2 := date.NewTime(11, 0, 0, 0)

	assert.True(t, t1.Before(t2))
	assert.False(t, t2.Before(t1))
	assert.False(t, t1.Before(t1))

	assert.True(t, t2.After(t1))
	assert.False(t, t1.After(t2))
	assert.False(t, t1.After(t1))
}

// TestTime_Order は OrderTime 関数および Order メソッドをテストする。
func TestTime_Order(t *testing.T) {
	t1 := date.NewTime(10, 0, 0, 0)
	t2 := date.NewTime(12, 0, 0, 0)

	first, sec := date.OrderTime(t2, t1)
	assert.Equal(t, t1, first)
	assert.Equal(t, t2, sec)

	first2, sec2 := t2.Order(t1)
	assert.Equal(t, t1, first2)
	assert.Equal(t, t2, sec2)

	first3, sec3 := date.OrderTime(t1, t2)
	assert.Equal(t, t1, first3)
	assert.Equal(t, t2, sec3)
}

// TestTime_ToTime は ToTime メソッドをテストする。
func TestTime_ToTime(t *testing.T) {
	tm := date.NewTime(15, 30, 45, 0)
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)

	tNil := tm.ToTime(nil)
	assert.Equal(t, time.Local, tNil.Location())
	assert.Equal(t, 15, tNil.Hour())
	assert.Equal(t, 30, tNil.Minute())
	assert.Equal(t, 45, tNil.Second())

	tJST := tm.ToTime(jst)
	assert.Equal(t, jst, tJST.Location())
}

// TestTime_Format_String は Format および String メソッドをテストする。
func TestTime_Format_String(t *testing.T) {
	tm := date.NewTime(15, 30, 45, 0)
	assert.Equal(t, "15:30:45", tm.String())
	assert.Equal(t, "15:30", tm.Format("15:04"))
	assert.Equal(t, "03:30 PM", tm.Format("03:04 PM"))

	tm = date.NewTime(25, 0, 0, 0)
	assert.Equal(t, "01:00:00", tm.Format("15:04:05"))

	tmNano := date.NewTime(15, 30, 45, 123456789)
	assert.Equal(t, "15:30:45.123456789", tmNano.String())
}

// TestParseTime は ParseTime 関数をテストする。
func TestParseTime(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    date.Time
		wantErr bool
	}{
		{
			name:    "valid standard time",
			value:   "15:30:45",
			want:    date.NewTime(15, 30, 45, 0),
			wantErr: false,
		},
		{
			name:    "valid time with nanoseconds",
			value:   "15:30:45.123456789",
			want:    date.NewTime(15, 30, 45, 123456789),
			wantErr: false,
		},
		{
			name:    "valid time with hour beyond 24",
			value:   "26:00:00",
			want:    date.NewTime(26, 0, 0, 0),
			wantErr: false,
		},
		{
			name:    "minute out of range is normalized",
			value:   "15:60:00",
			want:    date.NewTime(16, 0, 0, 0),
			wantErr: false,
		},
		{
			name:    "invalid format",
			value:   "invalid",
			want:    date.Time{},
			wantErr: true,
		},
		{
			name:    "hour:minute without seconds is not supported",
			value:   "09:15",
			want:    date.Time{},
			wantErr: true,
		},
		{
			name:    "empty fractional part",
			value:   "15:30:45.",
			want:    date.Time{},
			wantErr: true,
		},
		{
			name:    "fractional part too long",
			value:   "15:30:45.1234567890",
			want:    date.Time{},
			wantErr: true,
		},
		{
			name:    "fractional part with trailing junk",
			value:   "15:30:45.123456789junk",
			want:    date.Time{},
			wantErr: true,
		},
		{
			name:    "single digit fractional part",
			value:   "15:30:45.1",
			want:    date.NewTime(15, 30, 45, 100000000),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := date.ParseTime(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, got.IsZero())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
