package character

import (
	"outerspace/weapon"
	"strconv"
)

type CharacterType int

//Base Stats for all: Integrity: 100, Firepower: 10, Shielding: 10, Thrusters: 100.

const (
	Undefined CharacterType = iota
	Striker                 = 1
	Titan                   = 2
	Spectre                 = 3
	Engineer                = 4
	Navigator               = 5
)

const IntegrityBase = 100
const FirepowerBase = 10
const ShieldingBase = 10
const ThrustersBase = 100

type ChararcterManagement interface {
	// 0 Energy Cost. Always available.
	PrimaryWeapon() string
	// Tactical Module: High Energy Cost. High Impact.
	TacticalModule(energy int)
}

type Chararcter struct {
	Name          string
	Level         int
	CharacterType CharacterType
	Integrity     int
	Firepower     int
	Shielding     int
	Thrusters     int
	PrimaryWeapon weapon.Weapon
}

func (c *Chararcter) String() string {
	return "Player: " + c.Name + " (" + c.CharacterType.String() + ") - Level " + strconv.Itoa(c.Level)
}

func (c CharacterType) String() string {
	switch c {
	case Striker:
		return "Striker"
	case Titan:
		return "Titan"
	case Spectre:
		return "Spectre"
	case Engineer:
		return "Engineer"
	case Navigator:
		return "Navigator"
	default:
		return "Undefined"
	}
}
