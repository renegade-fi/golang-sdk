package crypto

import (
	"errors"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// Poseidon2Sponge represents a sponge construction on top of the Poseidon2 permutation
// Modeled after the implementation in:
// https://github.com/renegade-fi/renegade/blob/main/renegade-crypto/src/hash/poseidon2.rs
// The original paper can be found at:
// https://eprint.iacr.org/2023/323
type Poseidon2Sponge struct {
	state     [WIDTH]fr.Element
	nextIndex int
	squeezing bool
}

// NewPoseidon2Sponge creates a new Poseidon2Sponge instance
func NewPoseidon2Sponge() *Poseidon2Sponge {
	return &Poseidon2Sponge{
		state:     [WIDTH]fr.Element{},
		nextIndex: 0,
		squeezing: false,
	}
}

// Hash hashes the given input and returns a single-squeeze
func (p *Poseidon2Sponge) Hash(seq []fr.Element) fr.Element {
	//nolint:errcheck
	p.AbsorbBatch(seq)
	return p.Squeeze()
}

// Absorb absorbs a single scalar into the sponge
func (p *Poseidon2Sponge) Absorb(x fr.Element) error {
	if p.squeezing {
		return errors.New("cannot absorb while squeezing")
	}

	if p.nextIndex == RATE {
		p.permute()
		p.nextIndex = 0
	}

	entry := p.nextIndex + CAPACITY
	(&p.state[entry]).Add(&p.state[entry], &x)
	p.nextIndex++
	return nil
}

// AbsorbBatch absorbs a batch of scalars into the sponge
func (p *Poseidon2Sponge) AbsorbBatch(x []fr.Element) error {
	for _, scalar := range x {
		if err := p.Absorb(scalar); err != nil {
			return err
		}
	}
	return nil
}

// Squeeze squeezes a single scalar from the sponge
func (p *Poseidon2Sponge) Squeeze() fr.Element {
	if !p.squeezing || p.nextIndex == RATE {
		p.permute()
		p.nextIndex = 0
		p.squeezing = true
	}

	entry := p.nextIndex + CAPACITY
	p.nextIndex++
	return p.state[entry]
}

// SqueezeBatch squeezes a batch of scalars from the sponge
func (p *Poseidon2Sponge) SqueezeBatch(n int) []fr.Element {
	result := make([]fr.Element, n)
	for i := 0; i < n; i++ {
		result[i] = p.Squeeze()
	}
	return result
}

// permute permutes the inner state
func (p *Poseidon2Sponge) permute() {
	p.externalMDS()

	half := R_F / 2
	for i := 0; i < half; i++ {
		p.externalRound(i)
	}

	for i := 0; i < R_P; i++ {
		p.internalRound(i)
	}

	for i := half; i < R_F; i++ {
		p.externalRound(i)
	}
}

// externalRound runs an external round on the state
func (p *Poseidon2Sponge) externalRound(roundNumber int) {
	p.externalAddRC(roundNumber)
	p.externalSbox()
	p.externalMDS()
}

// externalAddRC adds a round constant to the state in an external round
func (p *Poseidon2Sponge) externalAddRC(roundNumber int) {
	rc := FULL_ROUND_CONSTANTS[roundNumber]
	for i := range p.state {
		p.state[i].Add(&p.state[i], &rc[i])
	}
}

// externalSbox applies the S-box to the entire state in an external round
func (p *Poseidon2Sponge) externalSbox() {
	for i := range p.state {
		applySbox(&p.state[i])
	}
}

// externalMDS applies the external MDS matrix M_E to the state
func (p *Poseidon2Sponge) externalMDS() {
	var sum fr.Element
	for _, x := range p.state {
		sum.Add(&sum, &x)
	}

	for i := range p.state {
		p.state[i].Add(&p.state[i], &sum)
	}
}

// internalRound runs an internal round on the state
func (p *Poseidon2Sponge) internalRound(roundNumber int) {
	p.internalAddRC(roundNumber)
	p.internalSbox()
	p.internalMDS()
}

// internalAddRC adds a round constant to the first state element in an internal round
func (p *Poseidon2Sponge) internalAddRC(roundNumber int) {
	rc := PARTIAL_ROUND_CONSTANTS[roundNumber]
	p.state[0].Add(&p.state[0], &rc)
}

// internalSbox applies the S-box to the first state element in an internal round
func (p *Poseidon2Sponge) internalSbox() {
	applySbox(&p.state[0])
}

// internalMDS applies the internal MDS matrix M_I to the state
func (p *Poseidon2Sponge) internalMDS() {
	var sum fr.Element
	for _, x := range p.state {
		sum.Add(&sum, &x)
	}

	p.state[WIDTH-1].Double(&p.state[WIDTH-1])
	for i := range p.state {
		p.state[i].Add(&p.state[i], &sum)
	}
}

// applySbox applies the s-box to an element of the state
// We use the x^5 sbox
func applySbox(val *fr.Element) {
	var tmp fr.Element
	tmp.Square(val)
	tmp.Square(&tmp)
	val.Mul(val, &tmp)
}
