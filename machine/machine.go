package machine

import (
	"golang-wasm/constants"
	"golang-wasm/helpers"
	"golang-wasm/rotor"
	"strings"
)

type Machine struct{}

func NewMachine() *Machine {
	return &Machine{}
}

func (m Machine) ScrambleCharacter(c string) string {
	charIndex := helpers.StringIndexOf(strings.ToUpper(c), constants.Characters)
	if charIndex == -1 {
		return c
	}

	finalScrambledIndex := charIndex

	reflector := rotor.DefaultConfig.Reflector
	rotors := rotor.DefaultConfig.Rotors

	// first run through the rotors
	for rIndex := 0; rIndex < len(rotors); rIndex++ {
		rotorLetter := rotors[rIndex].Sequence[finalScrambledIndex]
		mappedLetter := rotors[rIndex].CrossConnections[rotorLetter]
		finalScrambledIndex = helpers.SliceIndexOf(mappedLetter, rotors[rIndex].Sequence)
	}

	// get the reflected letter
	reflectedLetter := reflector[rotors[len(rotors)-1].Sequence[finalScrambledIndex]]

	finalScrambledIndex = helpers.SliceIndexOf(reflectedLetter, rotors[len(rotors)-1].Sequence)

	// run back through the rotors in opposite direction
	for rIndex := len(rotors) - 1; rIndex >= 0; rIndex-- {
		rotorLetter := rotors[rIndex].Sequence[finalScrambledIndex]
		mappedLetter := rotors[rIndex].CrossConnections[rotorLetter]
		finalScrambledIndex = helpers.SliceIndexOf(mappedLetter, rotors[rIndex].Sequence)
	}

	// rotate the rotors
	for rIndex := 0; rIndex <= len(rotors)-1; rIndex++ {
		if rIndex == 0 {
			rotors[rIndex].Sequence = helpers.SliceRotateRight(rotors[rIndex].Sequence)
			rotors[rIndex].RotationsDone = rotors[rIndex].RotationsDone + 1
		} else {
			if rotors[rIndex-1].RotationsDone != 0 && rotors[rIndex-1].RotationsDone % int64(len(rotors[rIndex].Sequence)) == 0 {
				rotors[rIndex].Sequence = helpers.SliceRotateRight(rotors[rIndex].Sequence)
				rotors[rIndex].RotationsDone = rotors[rIndex].RotationsDone + 1
			}
		}
	}

	return constants.CharactersArray[finalScrambledIndex]
}