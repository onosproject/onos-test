package model

import (
	"fmt"
	"strconv"
)

// NewTrace returns a new model Trace from the given values
func NewTrace(values ...interface{}) (*Trace, error) {
	traceValues := make([]*Value, len(values))
	for i, value := range values {
		traceValue, err := NewTraceValue(value)
		if err != nil {
			return nil, err
		}
		traceValues[i] = traceValue
	}
	return &Trace{
		Values: traceValues,
	}, nil
}

// NewTraceValue returns a new model value from the given value
func NewTraceValue(value interface{}) (*Value, error) {
	if stringValue, ok := value.(string); ok {
		return &Value{
			Type:  Value_STRING,
			Bytes: []byte(stringValue),
		}, nil
	}
	if intValue, ok := value.(int); ok {
		return &Value{
			Type:  Value_INTEGER,
			Bytes: []byte(strconv.Itoa(intValue)),
		}, nil
	}
	if intValue, ok := value.(int32); ok {
		return &Value{
			Type:  Value_INTEGER,
			Bytes: []byte(strconv.Itoa(int(intValue))),
		}, nil
	}
	if intValue, ok := value.(int64); ok {
		return &Value{
			Type:  Value_INTEGER,
			Bytes: []byte(strconv.Itoa(int(intValue))),
		}, nil
	}
	if intValue, ok := value.(uint); ok {
		return &Value{
			Type:  Value_INTEGER,
			Bytes: []byte(strconv.Itoa(int(intValue))),
		}, nil
	}
	if intValue, ok := value.(uint32); ok {
		return &Value{
			Type:  Value_INTEGER,
			Bytes: []byte(strconv.Itoa(int(intValue))),
		}, nil
	}
	if intValue, ok := value.(uint64); ok {
		return &Value{
			Type:  Value_INTEGER,
			Bytes: []byte(strconv.Itoa(int(intValue))),
		}, nil
	}
	if floatValue, ok := value.(float32); ok {
		return &Value{
			Type:  Value_DECIMAL,
			Bytes: []byte(strconv.FormatFloat(float64(floatValue), 'E', -1, 32)),
		}, nil
	}
	if floatValue, ok := value.(float64); ok {
		return &Value{
			Type:  Value_DECIMAL,
			Bytes: []byte(strconv.FormatFloat(floatValue, 'E', -1, 64)),
		}, nil
	}
	if boolValue, ok := value.(bool); ok {
		return &Value{
			Type:  Value_BOOLEAN,
			Bytes: []byte(strconv.FormatBool(boolValue)),
		}, nil
	}
	if bytesValue, ok := value.([]byte); ok {
		return &Value{
			Type:  Value_UNKNOWN,
			Bytes: []byte(bytesValue),
		}, nil
	}
	return nil, fmt.Errorf("cannot determine value type for %v", value)
}

func (v *Value) AsString() string {
	return string(v.Bytes)
}

func (v *Value) AsQuotedString() string {
	if v.Type == Value_STRING {
		return fmt.Sprintf("\"%s\"", v.AsString())
	}
	return v.AsString()
}
