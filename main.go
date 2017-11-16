package main

import (
	"os"
	"image"
	_"image/png"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"time"
	"fmt"
)

const tileCount int = 32
const pixelPerGrid int = 16

type item struct {
	sprite *pixel.Sprite
	pos pixel.Vec
}

type tile struct {
	sprite *pixel.Sprite
	isAccessible bool
}

type move struct {
	direction string
	pixelLeft int
}

type player struct {
	sprite *pixel.Sprite
	pos struct {
		X int
		Y int
	}
	gear []item
	pack []item
	disp pixel.Vec
	dispTime float64
}

type spriteMove struct {
	standUp *pixel.Sprite
	standDown *pixel.Sprite
	standLeft *pixel.Sprite
	standRight *pixel.Sprite
	walkUp *pixel.Sprite
	walkDown *pixel.Sprite
	walkLeft *pixel.Sprite
	walkRight *pixel.Sprite
	walkUpAlt *pixel.Sprite
	walkDownAlt *pixel.Sprite
	walkLeftAlt *pixel.Sprite
	walkRightAlt *pixel.Sprite
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

func buildSpriteMoves(sm * spriteMove, spritesheet pixel.Picture) {
	sm.standUp = pixel.NewSprite(spritesheet, pixel.R(36+10, 36, 72+10, 72))
	sm.standDown = pixel.NewSprite(spritesheet, pixel.R(36+10, 108, 72+10, 144))
	sm.standLeft = pixel.NewSprite(spritesheet, pixel.R(36+10, 0, 72+10, 36))
	sm.standRight = pixel.NewSprite(spritesheet, pixel.R(36+10, 72, 72+10, 108))
	sm.walkUp = pixel.NewSprite(spritesheet, pixel.R(0, 36, 36, 72))
	sm.walkDown = pixel.NewSprite(spritesheet, pixel.R(0, 108, 36, 144))
	sm.walkLeft = pixel.NewSprite(spritesheet, pixel.R(0, 0, 36, 36))
	sm.walkRight = pixel.NewSprite(spritesheet, pixel.R(0, 72, 36, 108))
	sm.walkUpAlt = pixel.NewSprite(spritesheet, pixel.R(72+20, 36, 108+20, 72))
	sm.walkDownAlt = pixel.NewSprite(spritesheet, pixel.R(72+20, 108, 108+20, 144))
	sm.walkLeftAlt = pixel.NewSprite(spritesheet, pixel.R(72+20, 0, 108+20, 36))
	sm.walkRightAlt = pixel.NewSprite(spritesheet, pixel.R(72+20, 72, 108+20, 108))
}

func initializePlayer(plr * player, ms * spriteMove) {
	plr.pos.X = tileCount/2//float64(pixelPerGrid/2 + ((tileCount-1)*pixelPerGrid))
	plr.pos.Y = tileCount/2//float64(pixelPerGrid/2 + ((tileCount-1)*pixelPerGrid))
	plr.disp.X = float64(0)
	plr.disp.Y = float64(0)
	plr.dispTime = 0
	plr.sprite = ms.standDown
}

func initializeValidSpaces() [tileCount][tileCount]int {
	var vs [tileCount][tileCount]int
	//initialize all squares to valid edges
	for i:= 0; i<tileCount; i++ {
		for j:=0; j<tileCount; j++ {
			vs[i][j] = 1
		}
	}
	//outer edge boundary
	for i:=0; i<tileCount; i++ {
		for j:=0; j<2; j++ {
			vs[j][i] = 0
			vs[i][j] = 0
			vs[i][tileCount-j-1] = 0
			vs[tileCount-j-1][i] = 0
		}
	}
	//house1
	for i:=19; i<27; i++ {
		for j:=6; j<15; j++ {
			vs[j][i] = 0
		}
	}
	vs[13][18] = 0
	vs[14][18] = 0
	vs[12][18] = 0

	return vs
}

func moveUpdate(plr * player, direction string, moveSheet * spriteMove, vs [tileCount][tileCount]int) {
	if plr.disp.Len() > 0.01 {
		return
	}
	if direction == "U" {
		if plr.pos.Y+1 >= tileCount {
			fmt.Println("Too high!")
			return
		} else if vs[plr.pos.X][plr.pos.Y+1] == 0 {
			return
		}
		plr.disp.Y = -1 * float64(pixelPerGrid)
		plr.disp.X = 0
		plr.dispTime = float64(pixelPerGrid)
		plr.pos.Y += 1
		plr.sprite = moveSheet.walkUp
	} else if direction == "D" {
		if plr.pos.Y-1 < 0 {
			return
		} else if vs[plr.pos.X][plr.pos.Y-1] == 0 {
			return
		}
		plr.disp.Y = float64(pixelPerGrid)
		plr.disp.X = 0
		plr.dispTime = float64(pixelPerGrid)
		plr.pos.Y -= 1
		plr.sprite = moveSheet.walkDown
	} else if direction == "L" {
		if plr.pos.X-1 < 0 {
			return
		} else if vs[plr.pos.X-1][plr.pos.Y] == 0 {
			return
		}
		plr.disp.X = float64(pixelPerGrid)
		plr.disp.Y = 0
		plr.dispTime = float64(pixelPerGrid)
		plr.pos.X -= 1
		plr.sprite = moveSheet.walkLeft
	} else if direction == "R" {
		if plr.pos.X+1 >= tileCount {
			return
		} else if vs[plr.pos.X+1][plr.pos.Y] == 0 {
			return
		}
		plr.disp.X = -1 * float64(pixelPerGrid)
		plr.disp.Y = 0
		plr.dispTime = float64(pixelPerGrid)
		plr.pos.X += 1
		plr.sprite = moveSheet.walkRight
	}
}

func posToVec(plr player) (v pixel.Vec) {
	v.X = float64(pixelPerGrid * plr.pos.X)
	v.Y = float64(pixelPerGrid * plr.pos.Y)
	return
}

func updateDisp(plr *player, dec float64) {
	if plr.disp.X == 0 && plr.disp.Y > 0 && plr.dispTime > 0 {
		plr.disp.Y = plr.disp.Y - float64(pixelPerGrid)*dec
		plr.dispTime -= float64(pixelPerGrid)*dec
	} else if plr.disp.X == 0 && plr.disp.Y < 0 && plr.dispTime > 0 {
		plr.disp.Y = plr.disp.Y + float64(pixelPerGrid)*dec
		plr.dispTime -= float64(pixelPerGrid)*dec
	} else if plr.disp.Y == 0  && plr.disp.X > 0 && plr.dispTime > 0 {
		plr.disp.X = plr.disp.X - float64(pixelPerGrid)*dec
		plr.dispTime -= float64(pixelPerGrid)*dec
	} else if plr.disp.Y == 0 && plr.disp.X < 0 && plr.dispTime > 0{
		plr.disp.X = plr.disp.X + float64(pixelPerGrid)*dec
		plr.dispTime -= float64(pixelPerGrid)*dec
	} else {
		plr.disp = pixel.V(0.00, 0.00)
		plr.dispTime = 0
	}
}

func run() {
	cfg := pixelgl.WindowConfig {
		Title: "Fire Starter",
		Bounds: pixel.R(0,0,512,512),
		VSync: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	//Build the base of the map
	mapPic, err := loadPicture("mymap.png")
	if err != nil {
		panic(err)
	}
	mapBase := pixel.NewSprite(mapPic, mapPic.Bounds())

	//Build the structure of player movement sprites
	playerSpritesheet, err := loadPicture("MainGuySpriteSheet.png")
	if err!= nil {
		panic(err)
	}
	var playerMoves spriteMove
	buildSpriteMoves(&playerMoves, playerSpritesheet)

	//Initialize the player
	var plr player
	initializePlayer(&plr, &playerMoves)

	var validSpaces [tileCount][tileCount]int
	validSpaces = initializeValidSpaces()

	last := time.Now() //Initialize the time for determine time difference
	dt := time.Since(last).Seconds()

	for !win.Closed() {
		dt = time.Since(last).Seconds()
		last = time.Now()

		//Check for user input and react
		if win.Pressed(pixelgl.KeyUp) {
			moveUpdate(&plr, "U", &playerMoves, validSpaces)
		} else if win.Pressed(pixelgl.KeyDown) {
			moveUpdate(&plr, "D", &playerMoves, validSpaces)
		} else if win.Pressed(pixelgl.KeyLeft) {
			moveUpdate(&plr, "L", &playerMoves, validSpaces)
		} else if win.Pressed(pixelgl.KeyRight) {
			moveUpdate(&plr, "R", &playerMoves, validSpaces)
		}

		//Update player displacement
		updateDisp(&plr, 16*dt)

		win.Clear(colornames.Aliceblue)

		//Draw everything to screen		
		mapBase.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		plr.sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 1).Moved(posToVec(plr).Add(plr.disp)))
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
