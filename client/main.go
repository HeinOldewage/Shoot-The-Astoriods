package main

import (
	"fmt"
	"image"
	"math"
	"math/rand"
	"os"
	"sta/core"
	"time"

	_ "image/png"

	"github.com/ByteArena/box2d"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func main() {
	engine, ship, asteriods := SetupUnits()

	go func() {
		goalTime := time.Millisecond * 15
		for {
			start := time.Now()
			engine.Step(1.0 / 60.0)
			duration := time.Since(start)
			time.Sleep(goalTime - duration)
		}
	}()

	fmt.Println("Shoot the astoriods")
	pixelgl.Run(run(engine, ship, append(asteriods, ship)))
	fmt.Println("Done shooting")
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func run(engine *core.Engine, ship *core.Unit, units []*core.Unit) func() {
	return func() {
		cfg := pixelgl.WindowConfig{
			Title:  "Shoot the astoriods",
			Bounds: pixel.R(0, 0, 1024, 768),
			VSync:  true,
		}
		win, err := pixelgl.NewWindow(cfg)
		if err != nil {
			panic(err)
		}

		for !win.Closed() {
			up := box2d.B2Vec2{math.Cos(ship.Body.GetAngle()), math.Sin(ship.Body.GetAngle())}
			down := box2d.B2Vec2MulScalar(-1, up)
			for engine.World.IsLocked() {
				time.Sleep(time.Millisecond)
			}
			if win.Pressed(pixelgl.KeyW) {
				ship.Body.ApplyLinearImpulseToCenter(up, true)
			}
			if win.Pressed(pixelgl.KeyA) {
				ship.Body.ApplyLinearImpulseToCenter(box2d.B2Vec2{X: -1, Y: 0}, true)
			}
			if win.Pressed(pixelgl.KeyS) {
				ship.Body.ApplyLinearImpulseToCenter(down, true)
			}
			if win.Pressed(pixelgl.KeyD) {
				ship.Body.ApplyLinearImpulseToCenter(box2d.B2Vec2{X: 1, Y: 0}, true)
			}

			angleToTurnTo := win.MousePosition().Sub(toPixelVec(ship.Body.GetPosition())).Angle()
			angleToTurn := angleToTurnTo - ship.Body.GetAngle()
			if angleToTurn > 0.1 {
				angleToTurn = 0.1
			}
			if angleToTurn < -0.1 {
				angleToTurn = -0.1
			}

			ship.Body.ApplyAngularImpulse(angleToTurn*0.1, true)
			win.Clear(colornames.Skyblue)
			for _, u := range units {

				mat := pixel.IM
				mat = mat.Rotated(pixel.ZV, u.Body.GetAngle()-math.Pi/2)
				mat = mat.Moved(toPixelVec(u.Body.GetPosition()))
				u.Def.Sprite.Draw(win, mat)
			}
			win.Update()
		}
	}
}

func toPixelVec(v box2d.B2Vec2) pixel.Vec {
	return pixel.Vec{
		X: v.X,
		Y: v.Y,
	}
}

func SetupUnits() (*core.Engine, *core.Unit, []*core.Unit) {
	engine := &core.Engine{}
	engine.Init()

	shipDef := core.UnitDef{
		BodyDef:  box2d.NewB2BodyDef(),
		Fixtures: make([]*box2d.B2FixtureDef, 0),
	}

	shipDef.BodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	shipDef.BodyDef.Position.SetZero()
	shipDef.BodyDef.AngularDamping = 2

	shipPoly := box2d.NewB2PolygonShape()
	shipPoly.Set([]box2d.B2Vec2{
		box2d.B2Vec2{
			0,
			1,
		},
		box2d.B2Vec2{
			1,
			0,
		}, box2d.B2Vec2{
			0,
			-1,
		},
	}, 3)

	shipMainFixture := &box2d.B2FixtureDef{}
	shipMainFixture.Shape = shipPoly
	shipMainFixture.Density = 1
	shipMainFixture.Friction = 1

	shipDef.Fixtures = append(shipDef.Fixtures, shipMainFixture)

	pic, err := loadPicture("assest/SF08.png")
	if err != nil {
		panic(err)
	}
	shipDef.Sprite = pixel.NewSprite(pic, pic.Bounds())

	asteriodDef := core.UnitDef{
		BodyDef:  box2d.NewB2BodyDef(),
		Fixtures: make([]*box2d.B2FixtureDef, 0),
	}

	asteriodDef.BodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	asteriodDef.BodyDef.Position.SetZero()

	asteriodPoly := &box2d.B2CircleShape{}
	asteriodPoly.SetRadius(1)

	asteriodFixture := &box2d.B2FixtureDef{}
	asteriodFixture.Shape = asteriodPoly
	asteriodFixture.Density = 1

	asteriodDef.Fixtures = append(asteriodDef.Fixtures, asteriodFixture)

	asteriodPic, err := loadPicture("assest/a10000.png")
	if err != nil {
		panic(err)
	}
	asteriodDef.Sprite = pixel.NewSprite(asteriodPic, asteriodPic.Bounds())

	ship := shipDef.Create(engine)
	asteriods := make([]*core.Unit, 0, 10)
	for k := 0; k < 10; k++ {
		asteriodDef.BodyDef.Position.Set(rand.Float64()*100, rand.Float64()*100)
		asteriodDef.BodyDef.LinearVelocity.Set(rand.Float64()*30, rand.Float64()*30)

		asteriodDef.BodyDef.AngularVelocity = rand.Float64()*2 - 1
		asteriods = append(asteriods, asteriodDef.Create(engine))
	}

	return engine, ship, asteriods
}
