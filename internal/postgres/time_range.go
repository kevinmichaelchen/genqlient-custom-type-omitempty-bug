package postgres

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	// matches gqlgen's expected time format
	// https://github.com/99designs/gqlgen/blob/v0.17.20/graphql/time.go#L22
	timeFormat = time.RFC3339Nano

	// Hasura formats timestamps as: 2022-11-07 14:30:00+00
	hasuraFormat = "2006-01-02 15:04:05-07"
)

// TimeRange is a time interval.
// e.g., "[2022-11-07 14:30, 2023-11-09 15:30)"
// A square bracket is inclusive (e.g., "greater than or equal to").
// A parenthesis is exclusive.
type TimeRange struct {
	Start          time.Time
	StartInclusive bool
	End            time.Time
	EndInclusive   bool
}

// NewTimeRange constructs a new range.
func NewTimeRange(start time.Time, duration time.Duration) TimeRange {
	return TimeRange{
		Start:          start,
		StartInclusive: false,
		End:            start.Add(duration),
		EndInclusive:   false,
	}
}

func (r *TimeRange) UnmarshalJSON(b []byte) error {
	return UnmarshalTimeRange(b, r)
}

// TimeRangeFromJSON - Returns a TimeRange unmarshalled from a field within an
// Event Trigger payload.
// Example: "(\"2023-03-15 00:00:00+00\",\"2027-10-23 00:00:00+00\")"
func TimeRangeFromJSON(s string) (*TimeRange, error) {
	var r TimeRange
	err := UnmarshalTimeRange([]byte(s), &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// MarshalTimeRange returns the JSON as a []byte
func MarshalTimeRange(v *TimeRange) ([]byte, error) {
	b, err := marshalTimeRangeNoQuotes(v)
	if err != nil {
		return nil, err
	}

	s := fmt.Sprintf(`"%s"`, string(b))

	return []byte(s), nil
}

func marshalTimeRangeNoQuotes(v *TimeRange) ([]byte, error) {
	if v == nil {
		return []byte(`null`), nil
	}

	var first, last string
	if v.StartInclusive {
		first = "["
	} else {
		first = "("
	}
	if v.EndInclusive {
		last = "]"
	} else {
		last = ")"
	}

	a := formatTime(v.Start)
	b := formatTime(v.End)

	s := fmt.Sprintf(`%s%s, %s%s`, first, a, b, last)

	return []byte(s), nil
}

func purgeBackSlashesQuotesAndWhitespaceTrim(s string) string {
	s = strings.ReplaceAll(s, `"`, ``)
	s = strings.ReplaceAll(s, `\`, ``)
	s = strings.TrimSpace(s)
	return s
}

func UnmarshalTimeRange(b []byte, v *TimeRange) error {
	s := string(b)
	if len(s) == 0 {
		// TODO how do we handle empty time ranges?
		return errors.New("empty time range")
	}

	s = purgeBackSlashesQuotesAndWhitespaceTrim(s)

	if len(s) == 0 {
		return errors.New("empty time range after trim")
	}

	if strings.HasPrefix(s, "[") {
		v.StartInclusive = true
		s = strings.TrimPrefix(s, "[")
	} else if strings.HasPrefix(s, "(") {
		s = strings.TrimPrefix(s, "(")
	} else {
		return fmt.Errorf("invalid prefix time interval character: %c", s[0])
	}

	if strings.HasSuffix(s, "]") {
		v.EndInclusive = true
		s = strings.TrimSuffix(s, "]")
	} else if strings.HasSuffix(s, ")") {
		s = strings.TrimSuffix(s, ")")
	} else {
		return fmt.Errorf("invalid suffix time interval character: %c", s[len(s)-1])
	}

	arr := strings.Split(s, `,`)
	if len(arr) == 0 {
		return errors.New("time range should have 2 elements in it")
	}

	startStr := parseTimeRangeElement(arr[0])
	start, err := time.Parse(hasuraFormat, startStr)
	if err != nil {
		return fmt.Errorf("failed to parse first time range element: %w", err)
	}

	endStr := parseTimeRangeElement(arr[1])
	end, err := time.Parse(hasuraFormat, endStr)
	if err != nil {
		return fmt.Errorf("failed to parse second time range element: %w", err)
	}

	v.Start = start
	v.End = end

	return nil
}

func parseTimeRangeElement(in string) string {
	s := strings.TrimSpace(in)
	s = purgeBackSlashesQuotesAndWhitespaceTrim(s)
	return s
}

// 2022-11-07 14:30
func formatTime(t time.Time) string {
	return t.Format(timeFormat)
}
