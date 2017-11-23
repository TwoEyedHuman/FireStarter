package main

import (
	"os"
	"image"
	_"image/png"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"time"
)

const tileCount int = 32
const pixelPerGrid int = 16

type intVec struct { //used to represent positions in the grid
	X int
	Y int
}

type item struct { //structure that is an item (held or field)
	sprite *pixel.Sprite
	pos intVec
	health int
}

type tile struct { //game board squares
	sprite *pixel.Sprite
	isAccessible bool
}
/*
type move struct {
	direction string
	pixelLeft int
}
*/
type player struct {
	sprite *pixel.Sprite
	pos intVec
	gear []item
	pack []item
	disp pixel.Vec
	dispTime float64
	facing intVec
}

type spriteFaceDirection struct {
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

func buildSpriteMoves(sm * spriteFaceDirection, spritesheet pixel.Picture) {
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

func initializePlayer(plr * player, ms * spriteFaceDirection) {
	plr.pos.X = tileCount/2//float64(pixelPerGrid/2 + ((tileCount-1)*pixelPerGrid))
	plr.pos.Y = tileCount/2//float64(pixelPerGrid/2 + ((tileCount-1)*pixelPerGrid))
	plr.disp.X = float64(0)
	plr.disp.Y = float64(0)
	plr.dispTime = 0
	plr.sprite = ms.standDown
}

func initializeWolf(wlf * player) {
	wolfPic, err := loadPicture("hound.png")
	if err!= nil {
		panic(err)
	}
	wlf.pos = intVec{5, 5}
	wlf.disp.X = float64(0)
	wlf.disp.Y = float64(0)
	wlf.dispTime = float64(0)
	wlf.sprite = pixel.NewSprite(wolfPic, wolfPic.Bounds())
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
	//house2
	for i:=8; i<16; i++ {
		for j:=20; j<29; j++ {
			vs[j][i] = 0
		}
	}
	vs[13][18] = 0
	vs[14][18] = 0
	vs[12][18] = 0

	return vs
}

func moveUpdate(plr * player, newDir intVec, moveSheet * spriteFaceDirection, vs [tileCount][tileCount]int) {
	if plr.disp.Len() > 0.01 {
		return
	}
	newLoc := addIntVec(plr.pos, newDir)
	if newLoc.X >= tileCount || newLoc.Y >= tileCount {
		return
	} else if vs[newLoc.X][newLoc.Y] == 0 {
		return
	}
	plr.pos = addIntVec(plr.pos, newDir)
	plr.disp.X = float64(-1*pixelPerGrid*newDir.X)
	plr.disp.Y = float64(-1*pixelPerGrid*newDir.Y)
	plr.dispTime = float64(pixelPerGrid)
	if newDir.X == 0 && newDir.Y == 1 {
		plr.sprite = moveSheet.walkUp
	} else if newDir.X == 0 && newDir.Y == -1 {
		plr.sprite = moveSheet.walkDown
	} else if newDir.X == -1 && newDir.Y == 0 {
		plr.sprite = moveSheet.walkLeft
	} else if newDir.X == 1 && newDir.Y == 0 {
		plr.sprite = moveSheet.walkRight
	}
}

func posToVec(pos intVec) (v pixel.Vec) {
	v.X = float64(pixelPerGrid * pos.X)
	v.Y = float64(pixelPerGrid * pos.Y)
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

func imageToSprite(filePath string) (spr *pixel.Sprite) {
	spriteImage, err := loadPicture(filePath)
	if err != nil {
		panic(err)
	}
	spr = pixel.NewSprite(spriteImage, spriteImage.Bounds())
	return
}

func wolfChase(wolf * player, plr player, fi []item, vs [tileCount][tileCount]int) {
	newDir := preferMoveDirection(wolf.pos, plr.pos)
	newLoc := addIntVec(wolf.pos, newDir)


	for _, itm := range fi {
		if newLoc.X == itm.pos.X && newLoc.Y == itm.pos.Y {
			return
		}
	}
	if vs[newLoc.X][newLoc.Y] == 0 {
		return
	}

	wolf.pos = newLoc
	wolf.disp.X = float64(-1*pixelPerGrid*newDir.X)
	wolf.disp.Y = float64(-1*pixelPerGrid*newDir.Y)
	wolf.dispTime = float64(pixelPerGrid)
	wolf.facing = newDir
	return
}

func myFieldItems() (fi []item) {
	var tmpItem item
	tmpItem.sprite = imageToSprite("firePotion.png")
	tmpItem.health = 3
	tmpItem.pos.X = 6
	tmpItem.pos.Y = 6
	fi = append(fi, tmpItem)
	tmpItem.pos.X = 9
	tmpItem.pos.Y = 6
	fi = append(fi, tmpItem)
	tmpItem.pos.X = 12
	tmpItem.pos.Y = 10
	fi = append(fi, tmpItem)
	return
}

func addIntVec(v1 intVec, v2 intVec) (v3 intVec) {
	v3.X = v1.X + v2.X
	v3.Y = v1.Y + v2.Y
	return
}

func multIntVec(v1 intVec, sclr int) (v2 intVec) {
	v2.X = sclr * v1.X
	v2.Y = sclr * v1.Y
	return
}

func preferMoveDirection(from intVec, to intVec) (v3 intVec) {
	diffX := from.X - to.X
	diffY := from.Y - to.Y
	if diffX <= diffY {
		if from.Y > to.Y {
			v3 = intVec{0, -1}
		} else if from.Y < to.Y {
			v3 = intVec{0, 1}
		} else {
			if from.X > to.X {
				v3 = intVec{-1,0}
			} else if from.X < to.X {
				v3 = intVec{1,0}
			} else {
				v3 = intVec{0,0}
			}
		}
	} else if diffX > diffY {
		if from.X > to.X {
			v3 = intVec{-1,0}
		} else if from.X < to.X {
			v3 = intVec{1,0}
		} else {
			v3 = intVec{0,0}
		}
	}
	return
}

func itemPickup(plr * player, fi []item) (newFi []item) {
	itemHasBeenPickedUp := false
	for _, itm := range fi {
		if plr.pos == itm.pos {
			plr.pack = append(plr.pack, itm)
			itemHasBeenPickedUp = true
		} else {
			newFi = append(newFi, itm)
		}
	}
	if !itemHasBeenPickedUp && len(plr.pack) > 0 {
		plr.pack[0].pos.X = plr.pos.X
		plr.pack[0].pos.Y = plr.pos.Y
		newFi = append(newFi, plr.pack[0])
		plr.pack = plr.pack[1:]
	}
	return
}

func isWinLose(plrpos intVec, houndpos intVec, runTime float64) (res bool) {
	if plrpos.X == houndpos.X && plrpos.Y == houndpos.Y {
		res = true
	} else if runTime > 60 {
		res = true
	} else {
		res = false
	}
	return
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
	var playerMoves spriteFaceDirection
	buildSpriteMoves(&playerMoves, playerSpritesheet)

	var (
		plr player
		validSpaces [tileCount][tileCount]int
		fieldItems []item
		showMenu bool
		wolf player
		endCondition bool
		startTime time.Time
		attackTimer float64
	)

	//Initialize the player
	initializePlayer(&plr, &playerMoves)
	initializeWolf(&wolf)

	validSpaces = initializeValidSpaces()

	last := time.Now() //Initialize the time for determine time difference
	dt := time.Since(last).Seconds()

	fieldItems = myFieldItems()

	showMenu = false
	endCondition = false
	startTime = time.Now()
	attackTimer = float64(0)

	for !win.Closed() && !endCondition{
		dt = time.Since(last).Seconds()
		attackTimer += dt
		last = time.Now()

		//Check for user input and react
		if win.Pressed(pixelgl.KeyUp) {
			moveUpdate(&plr, intVec{0,1}, &playerMoves, validSpaces)
		} else if win.Pressed(pixelgl.KeyDown) {
			moveUpdate(&plr, intVec{0,-1}, &playerMoves, validSpaces)
		} else if win.Pressed(pixelgl.KeyLeft) {
			moveUpdate(&plr, intVec{-1,0}, &playerMoves, validSpaces)
		} else if win.Pressed(pixelgl.KeyRight) {
			moveUpdate(&plr, intVec{1,0}, &playerMoves, validSpaces)
		}
		if win.JustPressed(pixelgl.KeySpace) {
			fieldItems = itemPickup(&plr, fieldItems)
		}
		if win.Pressed(pixelgl.Key1) {
			showMenu = true
		} else {
			showMenu = false
		}

		if wolf.dispTime <  0.0001 {
			wolfChase(&wolf, plr, fieldItems, validSpaces)
			endCondition = isWinLose(plr.pos, wolf.pos, time.Since(startTime).Seconds())
		}
		if attackTimer > 0.5 {
			attackTimer = 0
			var newFi []item
			for i, itm := range fieldItems {
				newFi = append(newFi, itm)
				if addIntVec(wolf.pos, wolf.facing) == itm.pos {
					newFi[i].health -= 1
				}
				if newFi[i].health <= 0 {
					newFi[i] = newFi[len(newFi)-1]
					newFi = newFi[:len(newFi)-1]
				}
			}
			fieldItems = newFi
			newFi = nil
		}


		//Update player displacement
		updateDisp(&wolf, 3*dt)
		updateDisp(&plr, 8*dt)

		win.Clear(colornames.Aliceblue)

		//Draw everything to screen		
		mapBase.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		plr.sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 1).Moved(posToVec(plr.pos).Add(plr.disp)))
		wolf.sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 1).Moved(posToVec(wolf.pos).Add(wolf.disp)))

		for _, fItem := range fieldItems {
			fItem.sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 1).Moved(posToVec(fItem.pos)))
		}
		if showMenu {
			for i, pItem := range plr.pack {
				pItem.sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 1).Moved(pixel.V(float64(tileCount*pixelPerGrid - 32), float64(win.Bounds().Center().Y) + float64(i*pixelPerGrid*2))))
			}
		}
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
