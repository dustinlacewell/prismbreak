package main

type Renderable interface {
    GetX() int
    GetY() int
    SetX(int)
    SetY(int)

    GetFg() *RGB
    GetBg() *RGB
    GetGlyph() string
    Render(dx, dy int, linfo LightInfo, vinfo LightInfo)
}

type Tileable interface {
    Renderable
    IVisible
}

type Tile struct {
    X, Y int
    Fg, Bg *RGB
    Glyph rune
    Roughness int
    Visibility float64
}

func (self *Tile) GetX() int { return self.X }
func (self *Tile) GetY() int { return self.Y }

func (self *Tile) SetX(v int) { self.X = v }
func (self *Tile) SetY(v int) { self.Y = v }

func (self *Tile) GetFg() *RGB { return self.Fg }
func (self *Tile) GetBg() *RGB { return self.Bg }
func (self *Tile) GetGlyph() string { return string(self.Glyph) }
func (self *Tile) GetRoughness() int { return self.Roughness }
func (self *Tile) GetVisibility() float64 { return self.Visibility }
func (self *Tile) Render(x, y int, linfo LightInfo, vinfo LightInfo) {
    fg := self.Fg
    if fg != nil {
        cfg := rgb_lerp(*fg, linfo.Color, linfo.Intensity * .5)
        fg = &cfg
    }

    bg := self.Bg
    if bg != nil {
        cbg := rgb_lerp(*bg, linfo.Color, linfo.Intensity * .5)
        bg = &cbg
    }

    if fg == nil {
        Window.Set(
            x, y,
            (*self.Fg).Scale(fmax(linfo.Intensity, vinfo.Intensity)), 
            (*bg).Scale(fmax(linfo.Intensity, vinfo.Intensity)),        
            string(self.GetGlyph()),
        )
        return
    }

    if bg == nil {
        Window.Set(
            x, y,
            (*fg).Scale(fmax(linfo.Intensity, vinfo.Intensity)), 
            (*self.Bg).Scale(fmax(linfo.Intensity, vinfo.Intensity)), 
            string(self.GetGlyph()),
        )
        return
    }

    Window.Set(
        x, y,
        (*fg).Scale(fmax(linfo.Intensity, vinfo.Intensity)), 
        (*bg).Scale(fmax(linfo.Intensity, vinfo.Intensity)),
        string(self.GetGlyph()),
    )
    return

}
