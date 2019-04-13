package core

import (
	"github.com/ByteArena/box2d"
	
	"github.com/faiface/pixel"
)

type UnitDef struct {
	BodyDef  *box2d.B2BodyDef
	Fixtures []*box2d.B2FixtureDef
	Sprite *pixel.Sprite
}

type Unit struct {
	Def * UnitDef
	Body *box2d.B2Body
}

func (ud *UnitDef) Create(e * Engine) *Unit {
	u := e.AddUnit(ud)
	for _,f := range ud.Fixtures {
		u.AddFixture(f)
	}
	return u
}

func (u *Unit) AddFixture( f * box2d.B2FixtureDef) {
	u.Body.CreateFixtureFromDef(f)
}