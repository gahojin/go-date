package date_test

import (
	"encoding"
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gahojin/go-date"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_ encoding.TextMarshaler   = (*date.Date)(nil)
	_ encoding.TextUnmarshaler = (*date.Date)(nil)
	_ encoding.TextMarshaler   = date.Date{}

	_ encoding.TextMarshaler   = (*date.Time)(nil)
	_ encoding.TextUnmarshaler = (*date.Time)(nil)
	_ encoding.TextMarshaler   = date.Time{}
)

// TestDate_MarshalText は MarshalText メソッドがバイトスライスへのマーシャルを行えるかをテストする。
func TestDate_MarshalText(t *testing.T) {
	tests := []struct {
		name string
		d    date.Date
		want []byte
	}{
		{
			name: "normal date",
			d:    date.New(2024, time.May, 1),
			want: []byte("2024-05-01"),
		},
		{
			name: "zero date",
			d:    date.Date{},
			want: []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.MarshalText()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestDate_UnmarshalText は UnmarshalText メソッドがバイトスライスからDateへのアンマーシャルを行えるかをテストする。
func TestDate_UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    date.Date
		wantErr bool
	}{
		{
			name:    "valid date string",
			input:   []byte("2024-05-01"),
			want:    date.New(2024, time.May, 1),
			wantErr: false,
		},
		{
			name:    "empty slice",
			input:   []byte{},
			want:    date.Date{},
			wantErr: false,
		},
		{
			name:    "nil slice",
			input:   nil,
			want:    date.Date{},
			wantErr: false,
		},
		{
			name:    "invalid date format",
			input:   []byte("2024/05/01"),
			want:    date.Date{},
			wantErr: true,
		},
		{
			name:    "invalid string",
			input:   []byte("not-a-date"),
			want:    date.Date{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d date.Date
			err := d.UnmarshalText(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, d)
			}
		})
	}
}

// TestTime_MarshalText は Time の MarshalText メソッドをテストする。
func TestTime_MarshalText(t *testing.T) {
	tests := []struct {
		name string
		tm   date.Time
		want []byte
	}{
		{
			name: "normal time",
			tm:   date.NewTime(15, 30, 45, 0),
			want: []byte("15:30:45"),
		},
		{
			name: "zero time",
			tm:   date.Time{},
			want: []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tm.MarshalText()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestTime_UnmarshalText は Time の UnmarshalText メソッドをテストする。
func TestTime_UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    date.Time
		wantErr bool
	}{
		{
			name:    "valid time string",
			input:   []byte("15:30:45"),
			want:    date.NewTime(15, 30, 45, 0),
			wantErr: false,
		},
		{
			name:    "empty slice",
			input:   []byte{},
			want:    date.Time{},
			wantErr: false,
		},
		{
			name:    "nil slice",
			input:   nil,
			want:    date.Time{},
			wantErr: false,
		},
		{
			name:    "invalid time format",
			input:   []byte("15-30-45"),
			want:    date.Time{},
			wantErr: true,
		},
		{
			name:    "invalid string",
			input:   []byte("not-a-time"),
			want:    date.Time{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tm date.Time
			err := tm.UnmarshalText(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, tm)
			}
		})
	}
}

// TestJSON_RoundTrip は json.Marshal / json.Unmarshal によるシリアライズ・デシリアライズをテストする。
func TestJSON_RoundTrip(t *testing.T) {
	type item struct {
		Date date.Date `json:"date"`
		Time date.Time `json:"time"`
	}

	tests := []struct {
		name string
		val  item
		json string
	}{
		{
			name: "normal date and time",
			val: item{
				Date: date.New(2024, time.May, 1),
				Time: date.NewTime(15, 30, 45, 0),
			},
			json: `{"date":"2024-05-01","time":"15:30:45"}`,
		},
		{
			name: "zero date and time",
			val: item{
				Date: date.Date{},
				Time: date.Time{},
			},
			json: `{"date":"","time":""}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.val)
			require.NoError(t, err)
			assert.JSONEq(t, tt.json, string(data))

			var unmarshaled item
			err = json.Unmarshal(data, &unmarshaled)
			require.NoError(t, err)
			assert.Equal(t, tt.val, unmarshaled)
		})
	}
}

// TestJSON_Null は json.Marshal / json.Unmarshal における null の処理をテストする。
func TestJSON_Null(t *testing.T) {
	dVal := date.New(2024, time.May, 1)
	tVal := date.NewTime(15, 30, 45, 0)

	type itemPointer struct {
		Date *date.Date `json:"date"`
		Time *date.Time `json:"time"`
	}

	type itemPointerOmitempty struct {
		Date *date.Date `json:"date,omitempty"`
		Time *date.Time `json:"time,omitempty"`
	}

	type itemValue struct {
		Date date.Date `json:"date"`
		Time date.Time `json:"time"`
	}

	t.Run("nil pointers marshal to null", func(t *testing.T) {
		rec := itemPointer{
			Date: nil,
			Time: nil,
		}
		data, err := json.Marshal(rec)
		require.NoError(t, err)
		assert.JSONEq(t, `{"date":null,"time":null}`, string(data))

		var unmarshaled itemPointer
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)
		assert.Nil(t, unmarshaled.Date)
		assert.Nil(t, unmarshaled.Time)
	})

	t.Run("nil pointers with omitempty are omitted", func(t *testing.T) {
		rec := itemPointerOmitempty{
			Date: nil,
			Time: nil,
		}
		data, err := json.Marshal(rec)
		require.NoError(t, err)
		assert.JSONEq(t, `{}`, string(data))

		var unmarshaled itemPointerOmitempty
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)
		assert.Nil(t, unmarshaled.Date)
		assert.Nil(t, unmarshaled.Time)
	})

	t.Run("non-nil pointer round-trip", func(t *testing.T) {
		rec := itemPointer{
			Date: &dVal,
			Time: &tVal,
		}
		data, err := json.Marshal(rec)
		require.NoError(t, err)
		assert.JSONEq(t, `{"date":"2024-05-01","time":"15:30:45"}`, string(data))

		var unmarshaled itemPointer
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)
		require.NotNil(t, unmarshaled.Date)
		require.NotNil(t, unmarshaled.Time)
		assert.Equal(t, dVal, *unmarshaled.Date)
		assert.Equal(t, tVal, *unmarshaled.Time)
	})

	t.Run("unmarshal null into value struct", func(t *testing.T) {
		var unmarshaled itemValue
		err := json.Unmarshal([]byte(`{"date":null,"time":null}`), &unmarshaled)
		require.NoError(t, err)
		assert.True(t, unmarshaled.Date.IsZero())
		assert.True(t, unmarshaled.Time.IsZero())
	})

	t.Run("direct unmarshal null literal into Date and Time types", func(t *testing.T) {
		var d date.Date
		err := json.Unmarshal([]byte("null"), &d)
		require.NoError(t, err)
		assert.True(t, d.IsZero())

		var tm date.Time
		err = json.Unmarshal([]byte("null"), &tm)
		require.NoError(t, err)
		assert.True(t, tm.IsZero())

		var dp *date.Date
		err = json.Unmarshal([]byte("null"), &dp)
		require.NoError(t, err)
		assert.Nil(t, dp)

		var tmp *date.Time
		err = json.Unmarshal([]byte("null"), &tmp)
		require.NoError(t, err)
		assert.Nil(t, tmp)
	})

	t.Run("direct marshal nil pointer to Date and Time", func(t *testing.T) {
		var dp *date.Date
		data, err := json.Marshal(dp)
		require.NoError(t, err)
		assert.Equal(t, "null", string(data))

		var tmp *date.Time
		data, err = json.Marshal(tmp)
		require.NoError(t, err)
		assert.Equal(t, "null", string(data))
	})
}

// TestDate_DynamoDB は aws-sdk の attributevalue を用いたシリアライズ・デシリアライズをテストする。
func TestDate_DynamoDB(t *testing.T) {
	dVal := date.New(2024, time.May, 1)

	type recordPointer struct {
		ID   string     `dynamodbav:"id"`
		Date *date.Date `dynamodbav:"date"`
	}

	type recordPointerOmitempty struct {
		ID   string     `dynamodbav:"id"`
		Date *date.Date `dynamodbav:"date,omitempty"`
	}

	type recordValue struct {
		ID   string    `dynamodbav:"id"`
		Date date.Date `dynamodbav:"date"`
	}

	t.Run("nil Date pointer marshals to NULL AttributeValue", func(t *testing.T) {
		rec := recordPointer{
			ID:   "rec-1",
			Date: nil,
		}
		av, err := attributevalue.Marshal(rec)
		require.NoError(t, err)

		m, ok := av.(*types.AttributeValueMemberM)
		require.True(t, ok)
		assert.Equal(t, &types.AttributeValueMemberS{Value: "rec-1"}, m.Value["id"])
		assert.Equal(t, &types.AttributeValueMemberNULL{Value: true}, m.Value["date"])

		var unmarshaled recordPointer
		err = attributevalue.Unmarshal(av, &unmarshaled)
		require.NoError(t, err)
		assert.Equal(t, rec, unmarshaled)
		assert.Nil(t, unmarshaled.Date)
	})

	t.Run("nil Date pointer with omitempty is omitted", func(t *testing.T) {
		rec := recordPointerOmitempty{
			ID:   "rec-2",
			Date: nil,
		}
		av, err := attributevalue.Marshal(rec)
		require.NoError(t, err)

		m, ok := av.(*types.AttributeValueMemberM)
		require.True(t, ok)
		assert.Equal(t, &types.AttributeValueMemberS{Value: "rec-2"}, m.Value["id"])
		_, exists := m.Value["date"]
		assert.False(t, exists)

		var unmarshaled recordPointerOmitempty
		err = attributevalue.Unmarshal(av, &unmarshaled)
		require.NoError(t, err)
		assert.Equal(t, rec, unmarshaled)
		assert.Nil(t, unmarshaled.Date)
	})

	t.Run("non-nil Date pointer round-trip with TextMarshaler options", func(t *testing.T) {
		optEnc := func(eo *attributevalue.EncoderOptions) {
			eo.UseEncodingMarshalers = true
		}
		optDec := func(do *attributevalue.DecoderOptions) {
			do.UseEncodingUnmarshalers = true
		}

		rec := recordPointer{
			ID:   "rec-3",
			Date: &dVal,
		}
		av, err := attributevalue.MarshalWithOptions(rec, optEnc)
		require.NoError(t, err)

		var unmarshaled recordPointer
		err = attributevalue.UnmarshalWithOptions(av, &unmarshaled, optDec)
		require.NoError(t, err)
		require.NotNil(t, unmarshaled.Date)
		assert.Equal(t, dVal, *unmarshaled.Date)
	})

	t.Run("zero value Date struct marshals with inner omitempty fields", func(t *testing.T) {
		rec := recordValue{
			ID:   "rec-4",
			Date: date.Date{},
		}
		av, err := attributevalue.Marshal(rec)
		require.NoError(t, err)

		m, ok := av.(*types.AttributeValueMemberM)
		require.True(t, ok)
		assert.Equal(t, &types.AttributeValueMemberS{Value: "rec-4"}, m.Value["id"])

		// Unexported fields in Date will marshal into an empty map
		dateMap, ok := m.Value["date"].(*types.AttributeValueMemberM)
		require.True(t, ok)
		assert.Empty(t, dateMap.Value)

		var unmarshaled recordValue
		err = attributevalue.Unmarshal(av, &unmarshaled)
		require.NoError(t, err)
		assert.True(t, unmarshaled.Date.IsZero())
	})

	t.Run("non-zero value Date struct round-trip with TextMarshaler options", func(t *testing.T) {
		optEnc := func(eo *attributevalue.EncoderOptions) {
			eo.UseEncodingMarshalers = true
		}
		optDec := func(do *attributevalue.DecoderOptions) {
			do.UseEncodingUnmarshalers = true
		}

		rec := recordValue{
			ID:   "rec-5",
			Date: dVal,
		}
		av, err := attributevalue.MarshalWithOptions(rec, optEnc)
		require.NoError(t, err)

		var unmarshaled recordValue
		err = attributevalue.UnmarshalWithOptions(av, &unmarshaled, optDec)
		require.NoError(t, err)
		assert.Equal(t, rec, unmarshaled)
	})

	t.Run("direct unmarshal NULL AttributeValue into Date and *Date", func(t *testing.T) {
		var d date.Date
		err := attributevalue.Unmarshal(&types.AttributeValueMemberNULL{Value: true}, &d)
		require.NoError(t, err)
		assert.True(t, d.IsZero())

		var dp *date.Date
		err = attributevalue.Unmarshal(&types.AttributeValueMemberNULL{Value: true}, &dp)
		require.NoError(t, err)
		assert.Nil(t, dp)
	})

	t.Run("with UseEncodingMarshalers options (TextMarshaler round-trip)", func(t *testing.T) {
		optEnc := func(eo *attributevalue.EncoderOptions) {
			eo.UseEncodingMarshalers = true
		}
		optDec := func(do *attributevalue.DecoderOptions) {
			do.UseEncodingUnmarshalers = true
		}

		avZero, err := attributevalue.MarshalWithOptions(date.Date{}, optEnc)
		require.NoError(t, err)
		assert.Equal(t, &types.AttributeValueMemberS{Value: ""}, avZero)

		avVal, err := attributevalue.MarshalWithOptions(dVal, optEnc)
		require.NoError(t, err)
		assert.Equal(t, &types.AttributeValueMemberS{Value: "2024-05-01"}, avVal)

		var unmarshaled date.Date
		err = attributevalue.UnmarshalWithOptions(avVal, &unmarshaled, optDec)
		require.NoError(t, err)
		assert.Equal(t, dVal, unmarshaled)
	})
}
