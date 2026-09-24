package date

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNewCompactDate は NewCompactDate コンストラクタが正しく CompactDate を生成し正規化するかをテストする。
func TestNewCompactDate(t *testing.T) {
	tests := []struct {
		name      string
		year      int
		month     time.Month
		day       int
		wantYear  int
		wantMonth time.Month
		wantDay   int
	}{
		{
			name:      "standard date",
			year:      2024,
			month:     time.May,
			day:       1,
			wantYear:  2024,
			wantMonth: time.May,
			wantDay:   1,
		},
		{
			name:      "leap year feb 29",
			year:      2024,
			month:     time.February,
			day:       29,
			wantYear:  2024,
			wantMonth: time.February,
			wantDay:   29,
		},
		{
			name:      "normalize Jan 32 to Feb 1",
			year:      2024,
			month:     time.January,
			day:       32,
			wantYear:  2024,
			wantMonth: time.February,
			wantDay:   1,
		},
		{
			name:      "normalize day 0 to previous month end",
			year:      2024,
			month:     time.March,
			day:       0,
			wantYear:  2024,
			wantMonth: time.February,
			wantDay:   29,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewCompact(tt.year, tt.month, tt.day)
			assert.Equal(t, tt.wantYear, d.Year())
			assert.Equal(t, tt.wantMonth, d.Month())
			assert.Equal(t, tt.wantDay, d.Day())
		})
	}
}

// TestCompactFromTimeIn は CompactFromTimeIn 関数をテストする。
func TestCompactFromTimeIn(t *testing.T) {
	jst := time.FixedZone("JST", 9*60*60)

	tests := []struct {
		name string
		tm   time.Time
		loc  *time.Location
		want CompactDate
	}{
		{
			name: "converts UTC to JST, crossing the date boundary",
			tm:   time.Date(2024, 5, 1, 23, 30, 0, 0, time.UTC),
			loc:  jst,
			want: NewCompact(2024, 5, 2),
		},
		{
			name: "converts JST to UTC, crossing the date boundary backwards",
			tm:   time.Date(2024, 5, 1, 1, 0, 0, 0, jst),
			loc:  time.UTC,
			want: NewCompact(2024, 4, 30),
		},
		{
			name: "nil location defaults to time.Local",
			tm:   time.Date(2024, 5, 1, 23, 30, 0, 0, jst),
			loc:  nil,
			want: NewCompact(2024, 5, 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CompactFromTimeIn(tt.tm, tt.loc)
			assert.True(t, tt.want.Equal(got.Date))
		})
	}
}

type testCompactDateData struct {
	Value    CompactDate
	Nullable *CompactDate `json:",omitempty"`
	Empty    CompactDate  `json:",omitzero"`
}

func TestCompactDate_Marshal(t *testing.T) {
	sut := NewCompact(2025, time.October, 24)
	data := testCompactDateData{Value: sut, Nullable: nil, Empty: CompactDate{}}

	actual, err := json.Marshal(data)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"Value":"20251024"}`), actual)
}

func TestCompactDate_Unmarshal(t *testing.T) {
	sut := NewCompact(2025, 10, 24)
	var actual testCompactDateData

	err := json.Unmarshal([]byte(`{"Value":"20251024"}`), &actual)
	assert.NoError(t, err)
	assert.Equal(t, testCompactDateData{Value: sut, Nullable: nil, Empty: CompactDate{}}, actual)
}

func TestCompactDate_Date(t *testing.T) {
	sut := NewCompact(2025, 10, 24)

	assert.Equal(t, New(2025, 10, 24), sut.Date)
}

func TestCompactDate_Equal(t *testing.T) {
	sut := NewCompact(2025, 10, 24)

	assert.True(t, sut.Equal(CompactFromTime(time.Date(2025, 10, 24, 0, 0, 0, 0, time.Local)).Date))
	assert.True(t, sut.Equal(CompactFromTime(time.Date(2025, 10, 24, 10, 20, 30, 4, time.Local)).Date))
	assert.False(t, sut.Equal(CompactFromTime(time.Date(2025, 10, 25, 0, 0, 0, 0, time.Local)).Date))
}

func TestCompactDate_String(t *testing.T) {
	sut := NewCompact(2025, 10, 24)

	assert.Equal(t, "20251024", sut.String())
}

// TestParseCompact は ParseCompact 関数が文字列をCompactDateにパースできるかをテストする。
func TestParseCompact(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    CompactDate
		wantErr bool
	}{
		{
			name:    "valid standard date format",
			value:   "20240501",
			want:    NewCompact(2024, time.May, 1),
			wantErr: false,
		},
		{
			name:    "invalid format",
			value:   "2024/05/01",
			want:    CompactDate{},
			wantErr: true,
		},
		{
			name:    "invalid date value",
			value:   "invalid-date",
			want:    CompactDate{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCompact(tt.value)
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
