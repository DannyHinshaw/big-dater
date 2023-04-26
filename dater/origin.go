package dater

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Origin is the db where the zone was originally created.
type Origin int

const (
	OriginUnknown = iota
	OriginEST
	OriginPST
	OriginUTC
)

const (
	zedVal = "unknown"
	estVal = "db_est"
	pstVal = "db_pst"
	utcVal = "db_utc"
)

var originStringsToValues = map[string]Origin{
	zedVal: OriginUnknown,
	estVal: OriginEST,
	pstVal: OriginPST,
	utcVal: OriginUTC,
}

func OriginString(s string) (Origin, error) {
	if val, ok := originStringsToValues[s]; ok {
		return val, nil
	}

	return OriginUnknown, fmt.Errorf("%s does not belong to Origin values", s)
}

func (s Origin) String() string {
	switch s {
	case OriginUnknown:
		return zedVal
	case OriginEST:
		return estVal
	case OriginPST:
		return pstVal
	case OriginUTC:
		return utcVal
	default:
		return fmt.Sprintf("Unknown(%d)", s)
	}
}

func (s Origin) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Origin) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("expected a string for Origin, got %s", data)
	}

	var err error
	*s, err = OriginString(str)
	return err
}

func (s Origin) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *Origin) UnmarshalText(text []byte) error {
	var err error
	*s, err = OriginString(string(text))
	return err
}

func (s Origin) Value() (driver.Value, error) {
	return s.String(), nil
}

func (s *Origin) Scan(value any) error {
	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		bytes, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("value is not a byte slice")
		}

		str = string(bytes[:])
	}

	val, err := OriginString(str)
	if err != nil {
		return err
	}

	*s = val
	return nil
}
