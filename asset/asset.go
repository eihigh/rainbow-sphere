package asset

import (
	"bytes"
	"embed"
	"image/color"
	"image/png"
	"math"

	"github.com/ebiten/emoji"
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed image
var FS embed.FS

const (
	VW, VH = 800, 600 // なんか変だけどここでグローバル参照可能にしてしまう
)

var (
	Debug  bool
	Screen *ebiten.Image // なんか変だけどここでグローバル参照可能にしてしまう

	RainbowColors = [7]color.RGBA{
		{0xff, 0x00, 0x00, 0xff},
		{0xff, 0x33, 0x00, 0xff},
		{0xff, 0xcc, 0x00, 0xff},
		{0x00, 0x99, 0x00, 0xff},
		{0x00, 0x66, 0xff, 0xff},
		{0x00, 0x00, 0xcc, 0xff},
		{0x99, 0x00, 0xcc, 0xff},
	}

	PlayerImage  *ebiten.Image
	ShootImage   *ebiten.Image
	SphereImages [7]*ebiten.Image
	HeartImage   *ebiten.Image

	mplusFaces = map[float64]font.Face{}
	mplusFont  *opentype.Font
)

func Load() error {
	// もっと複雑なゲームならurlを渡して羃等性のあるロードをしてもらう
	// もっと厳密なゲームならmodelが必要なアセットリストを出して静的に検証できるようにする

	var err error
	PlayerImage, err = newImage("image/[dot]touhou_gamejam_dot_pack_2021/Split/Chara_123.png")
	if err != nil {
		PlayerImage = ebiten.NewImage(48, 48)
		PlayerImage.Fill(color.White)
	}

	HeartImage = emoji.Image("❤")
	ShootImage = newDiamondImage(color.RGBA{255, 255, 255, 255}, RainbowColors[1])

	for i, clr := range RainbowColors {
		SphereImages[i] = newSphereImage(clr)
	}

	mplusFont, err = opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		return err
	}

	return nil
}

func Mplus(size float64) font.Face {
	if f, ok := mplusFaces[size]; ok {
		return f
	}
	f, err := opentype.NewFace(mplusFont, &opentype.FaceOptions{
		Size:    size,
		DPI:     72, // pt = px のために固定
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}
	mplusFaces[size] = f
	return f
}

func newImage(name string) (*ebiten.Image, error) {
	b, err := FS.ReadFile(name)
	if err != nil {
		return nil, err
	}
	i, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(i), nil
}

func newSphereImage(clr color.RGBA) *ebiten.Image {
	r := float64(clr.R) / 0xff
	g := float64(clr.G) / 0xff
	b := float64(clr.B) / 0xff

	dc := gg.NewContext(128, 128)
	// 陽（半円）
	dc.DrawArc(64, 64, 64, 0, math.Pi)
	dc.SetRGB(1, 1, 1)
	dc.Fill()
	// 陰（半円）
	dc.DrawArc(64, 64, 64, math.Pi, math.Pi*2)
	dc.SetRGB(r, g, b)
	dc.Fill()

	// 張り出した陽
	dc.DrawCircle(96, 64, 32)
	dc.SetRGB(1, 1, 1)
	dc.Fill()
	// 張り出した陰
	dc.DrawCircle(32, 64, 32)
	dc.SetRGB(r, g, b)
	dc.Fill()

	// 陽中の陰
	dc.DrawCircle(96, 64, 11)
	dc.SetRGB(r, g, b)
	dc.Fill()
	// 陰中の陽
	dc.DrawCircle(32, 64, 11)
	dc.SetRGB(1, 1, 1)
	dc.Fill()

	return ebiten.NewImageFromImage(dc.Image())
}

func newDiamondImage(fillClr, strokeClr color.RGBA) *ebiten.Image {
	fr := float64(fillClr.R) / 0xff
	fg := float64(fillClr.G) / 0xff
	fb := float64(fillClr.B) / 0xff
	sr := float64(strokeClr.R) / 0xff
	sg := float64(strokeClr.G) / 0xff
	sb := float64(strokeClr.B) / 0xff

	dc := gg.NewContext(64, 64)
	dc.LineTo(0, 32)
	dc.LineTo(32, 8)
	dc.LineTo(64, 32)
	dc.LineTo(32, 56)
	dc.LineTo(0, 32)
	dc.SetRGB(fr, fg, fb)
	dc.FillPreserve()
	dc.SetRGB(sr, sg, sb)
	dc.SetLineWidth(3)
	dc.Stroke()

	return ebiten.NewImageFromImage(dc.Image())
}
