package wallet

import (
	"fmt"
	"reflect"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// --- Interface and Implementation --- //

// ScalarSerialize is an interface that can be implemented by any type that
// can be serialized to a slice of Scalars
type ScalarSerialize interface {
	// FromScalars deserializes a value from a slice of Scalars
	FromScalars(scalars *ScalarIterator) error
	// ToScalars serializes a value to a slice of Scalars
	ToScalars() ([]Scalar, error)
	// NumScalars returns the number of Scalars that will be serialized
	NumScalars() int
}

// FromScalars converts a `ScalarIterator` to
func (s *Scalar) FromScalars(scalars *ScalarIterator) error {
	scalar, err := scalars.Next()
	if err != nil {
		return err
	}
	*s = scalar
	return nil
}

// ToScalars converts a `Scalar` to a slice fo `Scalar`s
func (s *Scalar) ToScalars() ([]Scalar, error) {
	return []Scalar{*s}, nil
}

// NumScalars returns the number of `Scalar`s in the `Scalar`
func (s *Scalar) NumScalars() int {
	return 1
}

// Uint64 is a type that can be serialized to a slice of `Scalar`s
type Uint64 uint64

// FromScalars converts a `ScalarIterator` to a `Uint64`
func (s *Uint64) FromScalars(scalars *ScalarIterator) error {
	scalar, err := scalars.Next()
	if err != nil {
		return err
	}

	elt := fr.Element(scalar)
	*s = Uint64(elt.Uint64())
	return nil
}

// ToScalars converts a `Uint64` to a slice of `Scalar`s
func (s *Uint64) ToScalars() ([]Scalar, error) {
	elt := fr.NewElement(uint64(*s))
	return []Scalar{Scalar(elt)}, nil
}

// NumScalars returns the number of `Scalar`s in the `Uint64`
func (s *Uint64) NumScalars() int {
	return 1
}

// --- Serialization --- //

// ToScalarsRecursive is a helper function to serialize a value to a
// slice of scalars using reflection
func ToScalarsRecursive(s interface{}) ([]Scalar, error) {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("input must be a pointer type, got %T", s)
	}

	if ss, ok := s.(ScalarSerialize); ok {
		return ss.ToScalars()
	}

	elem := v.Elem()
	switch elem.Kind() {
	case reflect.Struct:
		return toScalarsStruct(elem)
	case reflect.Array:
		return toScalarsArray(elem)
	case reflect.Pointer:
		return ToScalarsRecursive(elem.Interface())
	default:
		return nil, fmt.Errorf("unsupported type: %T", s)
	}
}

// toScalarsStruct is a helper function to serialize a struct to a slice of scalars using reflection
func toScalarsStruct(v reflect.Value) ([]Scalar, error) {
	scalars := []Scalar{}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanInterface() {
			continue
		}

		// Check for the scalar_serialize="skip" tag
		if v.Type().Field(i).Tag.Get("scalar_serialize") == "skip" {
			continue
		}

		// Convert the field to a Scalar
		fieldScalars, err := ToScalarsRecursive(field.Addr().Interface())
		if err != nil {
			return nil, fmt.Errorf("error serializing field %d: %w", i, err)
		}
		scalars = append(scalars, fieldScalars...)
	}
	return scalars, nil
}

// toScalarsArray is a helper function to serialize an array to a slice of scalars using reflection
func toScalarsArray(v reflect.Value) ([]Scalar, error) {
	scalars := []Scalar{}
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		if !elem.CanAddr() {
			return nil, fmt.Errorf("cannot take address of element %d", i)
		}

		// Convert the element to a Scalar, passing a pointer
		fieldScalars, err := ToScalarsRecursive(elem.Addr().Interface())
		if err != nil {
			return nil, fmt.Errorf("error serializing element %d: %w", i, err)
		}
		scalars = append(scalars, fieldScalars...)
	}
	return scalars, nil
}

// --- Deserialization --- //

// ScalarIterator is a helper type that iterates over a slice of scalars
type ScalarIterator struct {
	scalars []Scalar
	index   int
}

// NewScalarIterator creates a new ScalarIterator
func NewScalarIterator(scalars []Scalar) *ScalarIterator {
	return &ScalarIterator{scalars: scalars, index: 0}
}

// Next returns the next scalar in the iterator
func (s *ScalarIterator) Next() (Scalar, error) {
	if s.index >= len(s.scalars) {
		return Scalar{}, fmt.Errorf("no more scalars")
	}
	scalar := s.scalars[s.index]
	s.index++
	return scalar, nil
}

// NumRemaining returns the remaining scalars in the iterator
func (s *ScalarIterator) NumRemaining() int {
	return len(s.scalars) - s.index
}

// FromScalarsRecursive is a helper function to deserialize a struct from a
// slice of scalars using reflection
func FromScalarsRecursive(s interface{}, scalars *ScalarIterator) error {
	// If the type implements ScalarSerialize, use the specialized method
	if ss, ok := s.(ScalarSerialize); ok {
		return ss.FromScalars(scalars)
	}

	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("non-pointer argument to FromScalarsRecursive")
	}
	v = v.Elem()

	switch v.Kind() {
	case reflect.Struct:
		return fromScalarsStruct(v, scalars)
	case reflect.Array:
		return fromScalarsArray(v, scalars)
	case reflect.Ptr:
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}

		return FromScalarsRecursive(v.Interface(), scalars)
	default:
		return fmt.Errorf("unsupported type: %v", v.Type())
	}
}

// fromScalarsStruct is a helper function to deserialize a struct from a
// slice of scalars using reflection
func fromScalarsStruct(v reflect.Value, scalars *ScalarIterator) error {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanSet() {
			continue
		}

		// Check for the scalar_serialize="skip" tag
		if v.Type().Field(i).Tag.Get("scalar_serialize") == "skip" {
			continue
		}

		if err := FromScalarsRecursive(field.Addr().Interface(), scalars); err != nil {
			return err
		}
	}
	return nil
}

// fromScalarsArray is a helper function to deserialize an array from a
// slice of scalars using reflection
func fromScalarsArray(v reflect.Value, scalars *ScalarIterator) error {
	for i := 0; i < v.Len(); i++ {
		if err := FromScalarsRecursive(v.Index(i).Addr().Interface(), scalars); err != nil {
			return err
		}
	}
	return nil
}
