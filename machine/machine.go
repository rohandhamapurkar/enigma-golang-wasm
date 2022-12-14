package machine

import (
	"golang-wasm/config"
	"golang-wasm/constants"
	"golang-wasm/helpers"
	"strings"
)

type Machine struct {
	EnigmaConfig *config.EnigmaConfig
}

func NewMachine(rc *config.EnigmaConfig) *Machine {
	m := &Machine{}
	m.SetConfig(rc)
	return m
}

func (m *Machine) SetConfig(rc *config.EnigmaConfig) {
	if rc == nil {
		rc = config.DefaultEnigmaConfig
	}

	cpy := &config.EnigmaConfig{}
	cpy.Reflector = helpers.CopyStringMap(rc.Reflector)
	cpy.Rotors = make([]config.Rotor, len(rc.Rotors))

	for i := 0; i < len(cpy.Rotors); i++ {
		cpy.Rotors[i].CrossConnections = helpers.CopyStringMap(rc.Rotors[i].CrossConnections)
		cpy.Rotors[i].Sequence = make([]string, len(rc.Rotors[i].Sequence))
		copy(cpy.Rotors[i].Sequence, rc.Rotors[i].Sequence)
		cpy.Rotors[i].RotationsDone = rc.Rotors[i].RotationsDone
	}

	m.EnigmaConfig = cpy
}

func (m *Machine) ScrambleCharacter(c string) string {
	charIndex := helpers.StringIndexOf(strings.ToUpper(c), constants.Characters)
	if charIndex == -1 {
		return c
	}

	finalScrambledIndex := charIndex

	reflector := m.EnigmaConfig.Reflector
	rotors := m.EnigmaConfig.Rotors

	// first run through the rotors
	for rIndex := 0; rIndex < len(rotors); rIndex++ {
		rotorLetter := rotors[rIndex].Sequence[finalScrambledIndex]
		mappedLetter := rotors[rIndex].CrossConnections[rotorLetter]
		finalScrambledIndex = helpers.StringSliceIndexOf(mappedLetter, rotors[rIndex].Sequence)
	}

	// get the reflected letter
	reflectedLetter := reflector[rotors[len(rotors)-1].Sequence[finalScrambledIndex]]

	finalScrambledIndex = helpers.StringSliceIndexOf(reflectedLetter, rotors[len(rotors)-1].Sequence)

	// run back through the rotors in opposite direction
	for rIndex := len(rotors) - 1; rIndex >= 0; rIndex-- {
		rotorLetter := rotors[rIndex].Sequence[finalScrambledIndex]
		mappedLetter := rotors[rIndex].CrossConnections[rotorLetter]
		finalScrambledIndex = helpers.StringSliceIndexOf(mappedLetter, rotors[rIndex].Sequence)
	}

	// rotate the rotors
	for rIndex := 0; rIndex <= len(rotors)-1; rIndex++ {
		if rIndex == 0 {
			rotors[rIndex].Sequence = helpers.StringSliceRotateRight(rotors[rIndex].Sequence)
			rotors[rIndex].RotationsDone = rotors[rIndex].RotationsDone + 1
		} else {
			if rotors[rIndex-1].RotationsDone != 0 && rotors[rIndex-1].RotationsDone%uint64(len(rotors[rIndex].Sequence)) == 0 {
				rotors[rIndex].Sequence = helpers.StringSliceRotateRight(rotors[rIndex].Sequence)
				rotors[rIndex].RotationsDone = rotors[rIndex].RotationsDone + 1
			}
		}
	}

	return constants.CharactersArray[finalScrambledIndex]
}
