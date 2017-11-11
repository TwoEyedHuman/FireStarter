package main

import (
	"os"
	"image"
	_"image/png"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"time"
	"math"
)

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
	pos pixel.Vec
	gear []item
	pack []item
	disp pixel.Vec
	dispTime float64
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

func updateDisp(plr *player, dec float64) {
	if plr.disp.X == 0 && plr.disp.Y > 0 && plr.dispTime > 0 {
		plr.disp.Y = plr.disp.Y - float64(64*dec)
		plr.dispTime -= float64(64*dec)
	} else if plr.disp.X == 0 && plr.disp.Y < 0 && plr.dispTime > 0 {
		plr.disp.Y = plr.disp.Y + float64(64*dec)
		plr.dispTime -= float64(64*dec)
	} else if plr.disp.Y == 0  && plr.disp.X > 0 && plr.dispTime > 0 {
		plr.disp.X = plr.disp.X - float64(64*dec)
		plr.dispTime -= float64(64*dec)
	} else if plr.disp.Y == 0 && plr.disp.X < 0 && plr.dispTime > 0{
		plr.disp.X = plr.disp.X + float64(64*dec)
		plr.dispTime -= float64(64*dec)
	} else {
		plr.disp = pixel.V(0.00, 0.00)
		plr.dispTime = 0
	}
}

func run() {
	TileCount := 8
	cfg := pixelgl.WindowConfig {
		Title: "Sunshine Walker",
		Bounds: pixel.R(0,0,512,512),
		VSync: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	spritesheet, err := loadPicture("MainGuySpriteSheet.png")
	if err != nil {
		panic(err)
	}

	var (
		tiles [8][8]*pixel.Sprite
		p1 player
	)

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(colornames.Aliceblue)
		for i := 0; i<TileCount; i++ {
			for j := 0; j<TileCount; j++ {
				tiles[i][j].Draw(win, pixel.IM.Scaled(pixel.ZV, 2).Moved(pixel.V(float64(32 + i*64), float64(32 + j*64))))
			}
		}

	p1.sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 1).Moved(p1.pos.Add(p1.disp)))

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
