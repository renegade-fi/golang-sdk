package wallet

import (
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/stretchr/testify/assert"
)

type TestNestedStruct struct {
	NestedScalar Scalar
	NestedUint64 Uint64
}

type TestStruct struct {
	ScalarField  Scalar
	Uint64Field  Uint64
	NestedStruct TestNestedStruct
	ArrayField   [2]Scalar
}

func randomScalar() Scalar {
	elt := fr.Element{}
	// nolint: errcheck
	elt.SetRandom()
	return Scalar(elt)
}

func TestToFromScalarsBasic(t *testing.T) {
	scalar := randomScalar()
	assert.Equal(t, 1, scalar.NumScalars())

	// Serialize to scalars
	scalars, err := scalar.ToScalars()
	if err != nil {
		t.Fatalf("ToScalars failed: %v", err)
	}

	assert.Equal(t, 1, len(scalars))
	assert.Equal(t, scalar, scalars[0])

	// Deserialize from scalars
	var reconstructed Scalar
	err = reconstructed.FromScalars(NewScalarIterator(scalars))
	assert.NoError(t, err)

	assert.Equal(t, scalar, reconstructed)
}

func TestToFromScalarsArray(t *testing.T) {
	// Create an array of random scalars
	original := [3]Scalar{randomScalar(), randomScalar(), randomScalar()}

	// Serialize to scalars
	scalars, err := ToScalarsRecursive(&original)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(scalars))

	// Deserialize from scalars
	var reconstructed [3]Scalar
	err = FromScalarsRecursive(&reconstructed, NewScalarIterator(scalars))
	assert.NoError(t, err)

	// Compare original and reconstructed
	assert.Equal(t, original, reconstructed)
}

func TestToFromScalarsStruct(t *testing.T) {
	original := TestNestedStruct{
		NestedScalar: randomScalar(),
		NestedUint64: Uint64(42),
	}

	// Serialize to scalars
	scalars, err := ToScalarsRecursive(&original)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(scalars))

	// Deserialize from scalars
	var reconstructed TestNestedStruct
	err = FromScalarsRecursive(&reconstructed, NewScalarIterator(scalars))
	assert.NoError(t, err)

	// Compare original and reconstructed
	assert.Equal(t, original, reconstructed)
}

func TestToFromScalarsNestedStruct(t *testing.T) {
	original := TestStruct{
		ScalarField:  randomScalar(),
		Uint64Field:  Uint64(1),
		NestedStruct: TestNestedStruct{NestedScalar: randomScalar(), NestedUint64: Uint64(2)},
		ArrayField:   [2]Scalar{randomScalar(), randomScalar()},
	}

	// Serialize to scalars
	scalars, err := ToScalarsRecursive(&original)
	assert.NoError(t, err)
	assert.Equal(t, 6, len(scalars))

	// Deserialize from scalars
	var reconstructed TestStruct
	err = FromScalarsRecursive(&reconstructed, NewScalarIterator(scalars))
	assert.NoError(t, err)

	// Compare original and reconstructed
	assert.Equal(t, original, reconstructed)
}
