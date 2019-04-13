package core


import (
	"github.com/ByteArena/box2d"
)

type Engine struct{
	World box2d.B2World
 }

 func (e *Engine) Init() {
	 e.World = box2d.MakeB2World(box2d.MakeB2Vec2(0,0))
 }

 func (e *Engine) AddUnit(ud * UnitDef) *Unit {
	u := &Unit{Def : ud}
	u.Body = e.World.CreateBody(ud.BodyDef)
	return u
 }

 func (e *Engine) Step(time float64) {
	e.World.Step(time, 8, 3);
 }
