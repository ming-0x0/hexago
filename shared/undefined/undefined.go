package undefined

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// Verify interface implementations
var (
	_ json.Marshaler           = (*Undefined[any])(nil)
	_ json.Unmarshaler         = (*Undefined[any])(nil)
	_ encoding.TextUnmarshaler = (*Undefined[any])(nil)
	_ driver.Valuer            = (*Undefined[any])(nil)
	_ sql.Scanner              = (*Undefined[any])(nil)
)

// Undefined is a generic wrapper type that can represent a value that may be explicitly
// undefined or unset, which is particularly useful for:
//   - JSON marshaling/unmarshaling where fields can be omitted
//   - Database operations where NULL values need to be distinguished from zero values
//
// Supported types for T include:
//   - Basic types: int64, float64, bool, string
//   - time.Time for timestamp handling
//   - []byte for binary data
//   - nil for explicit NULL values in databases
//
// The zero value of Undefined[T] is considered undefined (valid = false).
type Undefined[T any] struct {
	value T
	valid bool
}

func New[T any](value T) Undefined[T] {
	return Undefined[T]{
		value: value,
		valid: true,
	}
}

func (u *Undefined[T]) Set(value T) {
	u.value = value
	u.valid = true
}

func (u *Undefined[T]) Unset() {
	var v T
	u.value = v
	u.valid = false
}

// Implement json.Unmarshaler
func (u *Undefined[T]) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &u.value); err != nil {
		return err
	}

	u.valid = true
	return nil
}

// Implement json.Marshaler
func (u Undefined[T]) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(u.value)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Implement encoding.TextUnmarshaler
func (u *Undefined[T]) UnmarshalText(text []byte) error {
	u.valid = len(text) > 0
	if textUnmarshaler, ok := any(&u.value).(encoding.TextUnmarshaler); ok {
		if err := textUnmarshaler.UnmarshalText(text); err != nil {
			return err
		}
		u.valid = true
		return nil
	}

	return errors.New("Undefined: cannot unmarshal text: underlying value doesn't implement encoding.TextUnmarshaler")
}

// Implement driver.Valuer
func (u Undefined[T]) Value() (driver.Value, error) {
	if !u.valid {
		return nil, nil
	}

	if valuer, ok := any(u.value).(driver.Valuer); ok {
		v, err := valuer.Value()
		return v, err
	}
	return u.value, nil
}

// Implement sql.Scanner
func (u *Undefined[T]) Scan(src any) error {
	u.valid = true

	switch val := src.(type) {
	case nil:
		var t T
		u.value = t
	case Undefined[T]:
		u.value = val.value
	case *Undefined[T]:
		if val == nil {
			var t T
			u.value = t
		} else {
			u.value = val.value
		}
	case T:
		u.value = val
	case *T:
		if val == nil {
			var t T
			u.value = t
		} else {
			u.value = *val
		}
	default:
		if scanner, ok := any(&u.value).(sql.Scanner); ok {
			return scanner.Scan(src)
		}
		var t T
		return fmt.Errorf("Undefined: Scan() incompatible types (src: %T, dst: %T)", src, t)
	}
	return nil
}

func (u Undefined[T]) IsUndefined() bool {
	return !u.valid
}

func (u Undefined[T]) Ptr() *T {
	if u.valid {
		return &u.value
	}

	return nil
}

func (u Undefined[T]) Equal(other Undefined[T]) bool {
	if u.valid != other.valid {
		return false
	}
	if !u.valid {
		return true
	}

	return reflect.DeepEqual(u.value, other.value)
}
