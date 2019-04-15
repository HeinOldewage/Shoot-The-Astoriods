package core

import (
	"github.com/ByteArena/box2d"
	
	"github.com/faiface/pixel"
)

type UnitDef struct {
	BodyDef  *box2d.B2BodyDef
	Fixtures []*box2d.B2FixtureDef
	Sprite *pixel.Sprite
	Trust float64
	RotationalImpulse float64
}

type Unit struct {
	Def * UnitDef
	Body *box2d.B2Body	
	Trust float64
	RotationalImpulse float64
}

func (ud *UnitDef) Create(e * Engine) *Unit {
	u := e.AddUnit(ud)
	for _,f := range ud.Fixtures {
		u.AddFixture(f)
	}
	u.Trust = ud.Trust
	u.RotationalImpulse = ud.RotationalImpulse
	return u
}

func (u *Unit) AddFixture( f * box2d.B2FixtureDef) {
	u.Body.CreateFixtureFromDef(f)
}