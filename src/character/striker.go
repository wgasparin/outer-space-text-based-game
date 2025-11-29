package character

/*
Modifiers:

	Firepower +50%
	Integrity -20%
*/
func (c *Chararcter) NewStriker(name string, level int) *Chararcter {
	return &Chararcter{
		Name:          name,
		Level:         level,
		CharacterType: Striker,
		Integrity:     IntegrityBase - 20,
		Firepower:     FirepowerBase * 1.5,
		Shielding:     ShieldingBase,
		Thrusters:     ThrustersBase,
	}
}
