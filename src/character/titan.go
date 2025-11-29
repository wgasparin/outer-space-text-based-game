package character

/*
Modifiers:

	Firepower +50%
	Integrity -20%
*/
func (c *Chararcter) NewTitan(name string, level int) *Chararcter {
	return &Chararcter{
		Name:          name,
		Level:         level,
		CharacterType: Titan,
		Integrity:     IntegrityBase,
		Firepower:     FirepowerBase,
		Shielding:     ShieldingBase * 1.5,
		Thrusters:     ThrustersBase - 20,
	}
}
