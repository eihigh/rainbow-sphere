package model

// ゲームのコア部分

import (
	"math"
	"math/rand"

	"github.com/eihigh/rainbow-sphere/asset"
	"github.com/eihigh/rainbow-sphere/util"
	"github.com/eihigh/zu/geom"
	"github.com/eihigh/zu/tick"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	sphereRadius = 32
	shootRadius  = 8
	playerRadius = 16

	playerSpeed         = 5
	shootSpeed          = 7
	coolTime            = 4
	invincibleTick      = 36
	startInvincibleTick = 90
	tau                 = math.Pi * 2
)

type T = tick.Tick

type Stage struct {
	tick      T
	bounds    geom.Rectangle
	area      geom.Rectangle
	areaImage *ebiten.Image
	player    *player
	spheres   []*sphere
	shoots    []*shoot
}

type Config struct {
	HP       int
	SphereHP int
	MinSpeed float64
	AmpSpeed float64
	MinScale float64
	AmpScale float64
}

type Input struct {
	Left, Up, Right, Down bool
	Shoot                 bool
	Cursor                geom.Point
}

type player struct {
	pos       geom.Point
	left      bool // 左方向移動中フラグ
	hp        int
	mode      string // "" | damage | dead
	start     int
	lastShoot int
}

type sphere struct {
	index    int
	pos, vel geom.Point
	scale    float64
	hp       int
	mode     string // "" | damage | break
	start    int    // modeになったtick
}

type shoot struct {
	pos, vel geom.Point
	dead     bool
}

// ==================================================
//  Stage
// ==================================================

func NewStage(c *Config) *Stage {
	// がんばる
	s := &Stage{
		bounds: geom.Rect(0, 0, asset.VW, asset.VH),
		area:   geom.Rect(200, 100, 600, 500),
		player: &player{
			pos: geom.Pt(util.Col(6), util.Row(6)),
			// pos: geom.Point{},
			hp: c.HP,
		},
	}

	s.areaImage = ebiten.NewImage(int(s.area.Dx()), int(s.area.Dy()))
	s.areaImage.Fill(asset.RainbowColors[6])

	// 陰陽玉７個生成
	for i := 0; i < 7; i++ {

		// 位置
		pos := geom.Pt(util.Col(6), util.Row(6))
		r := util.Row(5)
		t := tau / 7 * float64(i)

		// 拡大率
		scale := c.MinScale + c.AmpScale*rand.Float64()

		// 外向きにランダムなベクトル
		vr := c.MinSpeed + c.AmpSpeed*rand.Float64()
		vr /= math.Sqrt(scale) // 大きな陰陽玉は気持ち遅くする
		vt := t + tau/2*(rand.Float64()-0.5)

		sp := &sphere{
			index: i,
			pos:   pos.Add(geom.PtFromRect(r, t)),
			vel:   geom.PtFromRect(vr, vt),
			scale: scale,
			hp:    c.SphereHP,
		}
		s.spheres = append(s.spheres, sp)
	}
	return s
}

func (s *Stage) UpdateMain(i *Input) string {
	defer s.tick.Advance(1)

	clear := true
	for _, sp := range s.spheres {
		if sp.hittable() {
			clear = false
			break
		}
	}
	failed := s.player.hp <= 0

	// 当たり判定 => Update（移動） => 描画、の順番
	// 見た目で当たっていることを表示してから当たり判定を行いたい

	if !clear && !failed {
		s.collision()
	}

	// プレイヤーUpdate（射撃、移動）
	s.shoot(i)
	if i.Left {
		s.player.pos.X -= playerSpeed
		s.player.left = true
	}
	if i.Up {
		s.player.pos.Y -= playerSpeed
	}
	if i.Right {
		s.player.pos.X += playerSpeed
		s.player.left = false
	}
	if i.Down {
		s.player.pos.Y += playerSpeed
	}
	if s.player.pos.X < s.area.Min.X {
		s.player.pos.X = s.area.Min.X
	}
	if s.player.pos.Y < s.area.Min.Y {
		s.player.pos.Y = s.area.Min.Y
	}
	if s.area.Max.X < s.player.pos.X {
		s.player.pos.X = s.area.Max.X
	}
	if s.area.Max.Y < s.player.pos.Y {
		s.player.pos.Y = s.area.Max.Y
	}

	// 陰陽玉Update
	// そんなに数ないのでGCしない
	for _, sp := range s.spheres {
		next := sp.pos.Add(sp.vel)
		if next.X < s.bounds.Min.X {
			sp.vel.X = -sp.vel.X
		}
		if s.bounds.Max.X < next.X {
			sp.vel.X = -sp.vel.X
		}
		if next.Y < s.bounds.Min.Y {
			sp.vel.Y = -sp.vel.Y
		}
		if s.bounds.Max.Y < next.Y {
			sp.vel.Y = -sp.vel.Y
		}
		sp.pos = sp.pos.Add(sp.vel)
	}

	// 射撃Update
	// 不要になったらGCする
	next := make([]*shoot, 0, len(s.shoots))
	for _, shoot := range s.shoots {
		if shoot.garbage() {
			continue
		}
		b := geom.Rect(s.bounds.Min.X-50, s.bounds.Min.Y-50, s.bounds.Max.X+50, s.bounds.Max.Y+50)
		if !shoot.pos.In(b) {
			continue
		}
		next = append(next, shoot)
		shoot.pos = shoot.pos.Add(shoot.vel)
	}
	s.shoots = next

	if clear {
		return "clear"
	}
	if failed {
		return "failed"
	}

	return ""
}

func (s *Stage) collision() {
	// 当たり判定
	for _, sphere := range s.spheres {
		if !sphere.hittable() {
			continue
		}
		for _, shoot := range s.shoots {
			if !shoot.hittable() {
				continue
			}

			r2 := math.Pow(sphereRadius*sphere.scale+shootRadius, 2)
			d2 := sphere.pos.Sub(shoot.pos).LengthSq()
			if d2 < r2 {
				// 陰陽玉が喰らう処理
				sphere.hit(s.tick)
				// 弾が消滅する処理
				shoot.dead = true
			}
		}
	}

	if s.player.hittable(s.tick) {
		for _, sphere := range s.spheres {
			if !sphere.hittable() {
				continue
			}
			r2 := math.Pow(sphereRadius*sphere.scale+playerRadius, 2)
			d2 := sphere.pos.Sub(s.player.pos).LengthSq()
			if d2 < r2 {
				// プレイヤーが喰らう処理
				s.player.hit(s.tick)
			}
		}
	}
}

func (s *Stage) shoot(i *Input) {
	if !i.Shoot {
		return
	}
	if s.tick.Elapsed()-s.player.lastShoot < coolTime {
		return
	}
	s.player.lastShoot = s.tick.Elapsed()

	aim := i.Cursor.Sub(s.player.pos).Angle()
	aim2 := aim - tau/35
	aim3 := aim + tau/35
	s.shoots = append(s.shoots, &shoot{
		pos: s.player.pos,
		vel: geom.PtFromRect(shootSpeed, aim),
	})
	s.shoots = append(s.shoots, &shoot{
		pos: s.player.pos,
		vel: geom.PtFromRect(shootSpeed, aim2),
	})
	s.shoots = append(s.shoots, &shoot{
		pos: s.player.pos,
		vel: geom.PtFromRect(shootSpeed, aim3),
	})
}

// ==================================================
//  player
// ==================================================

func (p *player) hittable(t T) bool {
	if p.mode == "dead" {
		return false
	}
	if p.mode == "damage" && t.SubTick(p.start).Elapsed() < invincibleTick {
		// ダメージ受けてから無敵時間以内は当たり判定消失
		return false
	}
	if p.mode == "" && t.SubTick(p.start).Elapsed() < startInvincibleTick {
		// 開幕無敵
		return false
	}
	return true
}

func (p *player) hit(t T) {
	p.hp--
	p.start = t.Elapsed()
	if p.hp <= 0 {
		p.mode = "dead"
	} else {
		p.mode = "damage"
	}
}

// ==================================================
//  sphere
// ==================================================

// hittable - 陰陽玉の当たり判定が有効かどうか返す
func (s *sphere) hittable() bool {
	return s.mode != "break"
}

func (s *sphere) hit(t T) {
	s.hp--
	s.start = t.Elapsed()
	if s.hp <= 0 {
		s.mode = "break"
	} else {
		s.mode = "damage"
	}
}

// ==================================================
//  shoot
// ==================================================

func (s *shoot) hittable() bool {
	return !s.dead
}

// garbage - 不要なオブジェクトになったかどうか返す
func (s *shoot) garbage() bool {
	return s.dead
}
