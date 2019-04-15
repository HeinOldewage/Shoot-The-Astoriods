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
	//"github.com/faiface/pixel/imdraw"
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

var (
	camPos       = pixel.ZV
	camSpeed     = 100.0
	camZoom      = 50.0
	camZoomSpeed = 1.2
)

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

		last := time.Now()
		for !win.Closed() {

			dt := time.Since(last).Seconds()
			last = time.Now()
			up := box2d.B2Vec2{math.Cos(ship.Body.GetAngle()), math.Sin(ship.Body.GetAngle())}

			down := box2d.B2Vec2MulScalar(-1, up)
			for engine.World.IsLocked() {
				time.Sleep(time.Millisecond)
			}
			if win.Pressed(pixelgl.KeyW) {
				trust := up
				trust.OperatorScalarMulInplace(ship.Trust)
				ship.Body.ApplyLinearImpulseToCenter(trust, true)
			}
			if win.Pressed(pixelgl.KeyS) {
				trust := down
				trust.OperatorScalarMulInplace(ship.Trust)
				ship.Body.ApplyLinearImpulseToCenter(trust, true)
			}

			if win.MouseScroll().Y != 0 {
				camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)
				fmt.Println(camZoom)
			}

			if win.Pressed(pixelgl.KeyLeft) {
				camPos.X -= camSpeed * dt
			}
			if win.Pressed(pixelgl.KeyRight) {
				camPos.X += camSpeed * dt
			}
			if win.Pressed(pixelgl.KeyDown) {
				camPos.Y -= camSpeed * dt
			}
			if win.Pressed(pixelgl.KeyUp) {
				camPos.Y += camSpeed * dt
			}
			camPos = toPixelVec(ship.Body.GetPosition())

			cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
			win.SetMatrix(cam)

			directionToMouse := cam.Unproject(win.MousePosition()).Sub(toPixelVec(ship.Body.GetPosition()))

			fromUpToMouse := directionToMouse.Sub(toPixelVec(up))

			directionToTurn := toPixelVec(up).Cross(fromUpToMouse)

			var impulse float64

			if directionToTurn > 0 {
				impulse = ship.RotationalImpulse
			}
			if directionToTurn < 0 {
				impulse = -ship.RotationalImpulse
			}

			fmt.Println(ship.Body.GetLinearVelocity().Length())

			ship.Body.ApplyAngularImpulse(impulse, true)

			win.Clear(colornames.Skyblue)

			for _, u := range units {

				mat := pixel.IM
				mat = mat.Scaled(pixel.ZV, 1.0/50.0)
				mat = mat.Rotated(pixel.ZV, u.Body.GetAngle()-math.Pi/2)
				mat = mat.Moved(toPixelVec(u.Body.GetPosition()))
				u.Def.Sprite.Draw(win, mat)
			}
			/*
				imd := imdraw.New(nil)

				imd.Color = colornames.Blueviolet
				imd.EndShape = imdraw.RoundEndShape
				imd.Push(toPixelVec(ship.Body.GetPosition()), toPixelVec(ship.Body.GetPosition()).Add(toPixelVec(up).Scaled(100)))
				imd.Line(10)
				imd.Draw(win)

				imdMouse := imdraw.New(nil)

				imdMouse.Color = colornames.Red
				imdMouse.EndShape = imdraw.RoundEndShape
				imdMouse.Push(cam.Unproject(win.MousePosition()), cam.Unproject(win.MousePosition()))
				imdMouse.Line(10)
				imdMouse.Draw(win)

				imdMouseLine := imdraw.New(nil)

				imdMouseLine.Color = colornames.Red
				imdMouseLine.EndShape = imdraw.RoundEndShape
				imdMouseLine.Push(toPixelVec(ship.Body.GetPosition()), toPixelVec(ship.Body.GetPosition()).Add(directionToMouse).Scaled(1))
				imdMouseLine.Line(10)
				imdMouseLine.Draw(win)

				imdLine := imdraw.New(nil)
				imdLine.Color = colornames.Green
				imdLine.EndShape = imdraw.RoundEndShape
				imdLine.Push(toPixelVec(up), toPixelVec(up).Add(fromUpToMouse).Scaled(1))
				imdLine.Line(10)
				imdLine.Draw(win)
			*/

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
		BodyDef:           box2d.NewB2BodyDef(),
		Fixtures:          make([]*box2d.B2FixtureDef, 0),
		Trust:             0.1,
		RotationalImpulse: 0.01,
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
	shipMainFixture.Friction = 0
	shipMainFixture.Filter.GroupIndex = 1

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
	asteriodFixture.Filter.GroupIndex = 2

	asteriodDef.Fixtures = append(asteriodDef.Fixtures, asteriodFixture)

	asteriodPic, err := loadPicture("assest/a10000.png")
	if err != nil {
		panic(err)
	}
	asteriodDef.Sprite = pixel.NewSprite(asteriodPic, asteriodPic.Bounds())

	ship := shipDef.Create(engine)
	asteriods := make([]*core.Unit, 0, 10)
	for k := 0; k < 10; k++ {
		asteriodDef.BodyDef.Position.Set(rand.Float64()*20, rand.Float64()*20)
		asteriodDef.BodyDef.LinearVelocity.Set(rand.Float64(), rand.Float64())

		asteriodDef.BodyDef.AngularVelocity = rand.Float64()*2 - 1
		asteriods = append(asteriods, asteriodDef.Create(engine))
	}

	engine.World.SetContactFilter(&noncolide{})

	return engine, ship, asteriods
}

type noncolide struct{}

func (*noncolide) ShouldCollide(fixtureA *box2d.B2Fixture, fixtureB *box2d.B2Fixture) bool {
	return fixtureA.GetFilterData().GroupIndex != fixtureB.GetFilterData().GroupIndex
}
