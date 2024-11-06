package dependencies

import (
	_ "embed"
	"encoding/json"
)

//go:embed prompt_defences.json
var defencesFile []byte

type DefenceType = int

const (
	InContext DefenceType = iota
	SystemModeSelfReminder
	SandwichDefence
	XmlEncapsulation
	RandomSequenceEnclosure
)

type Defence struct {
	Id          DefenceType `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Link        string      `json:"link"`
}

type DefenceList struct {
	Defences []Defence `json:"defences"`
}

func ProvideDefences() []Defence {
	defences := DefenceList{}

	err := json.Unmarshal(defencesFile, &defences)

	if err != nil {
		panic(err)
	}

	return defences.Defences
}
