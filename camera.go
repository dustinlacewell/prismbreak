package main

const (
    MASK_MAX_RADIUS = 20.0    
)

type Camera struct {
    hw, hh int
    x, y int
    ox, oy int

    tilemap *Map

    lighting LightingMap
    ambient_light LightInfo
    global_illumination float64

    masking *VisibilityMap
    mask_origin *Point
    mask_radius int
}

func NewCamera(tilemap *Map) *Camera {
    hw := app.width / 2;
    hh := app.height / 2;
    cam := &Camera {
        hw: hw, hh: hh,
        x: hw, y: hh,
        ox: 0, oy: 0,
        tilemap: tilemap,
        lighting: nil,
        ambient_light: LightInfo{0.0, White},
        global_illumination: 0.05,
        masking: nil,
        mask_origin: nil,
        mask_radius: MASK_MAX_RADIUS,
    }

    return cam;
}

func (self *Camera) boundPosition() {
    self.x = clamp(self.x, 0, self.tilemap.Width() - 1)
    self.ox = self.x - self.hw;

    self.y = clamp(self.y, 0, self.tilemap.Height() - 1)
    self.oy = self.y - self.hh;
}

func (self *Camera) SetAmbientLight(linfo LightInfo) {
    self.ambient_light = linfo
}

func (self *Camera) SetGlobalIllumination(level float64) {
    self.global_illumination = fmin(1.0, level)
}

func (self *Camera) SetMask(origin *Point, radius int) {
    self.mask_origin = origin
    self.mask_radius = min(
        MASK_MAX_RADIUS,
        radius,
    )
}

func (self *Camera) LookAt(x, y int) {
    self.x = x
    self.y = y
    self.ox = x - self.hw;
    self.oy = y - self.hh;

    self.boundPosition()
    self.SetMask(&Point{self.x, self.y}, 10)
    self.RenderLighting()
}

func (self *Camera) West() {
    self.Shift(-1, 0);
}

func (self *Camera) East() {
    self.Shift(1, 0);
}

func (self *Camera) North() {
    self.Shift(0, -1);
}

func (self *Camera) South() {
    self.Shift(0, 1);
}

func (self *Camera) Shift(dx, dy int) {
    self.LookAt(self.x + dx, self.y + dy);
}

func (self *Camera) Each(proc TileProcessor) {
    w := self.tilemap.Width();
    h := self.tilemap.Height();

    for x:=0; x<w; x++ {
        for y:=0; y<h; y++ {
            dx := self.ox + x;
            dy := self.oy + y;
            tile, ok := self.tilemap.GetTile(dx, dy);
            if ok {
                proc(self.tilemap, tile, x, y)
            } else {
                proc(self.tilemap, nil, x, y)
            }
        }
    }
}

func (self *Camera) RenderTile(x, y int, tile Renderable, linfo LightInfo) {
    fg := tile.GetFg()
    if fg != nil {
        lc := rgb_lerp(
            linfo.Color, 
            self.ambient_light.Color, 
            self.ambient_light.Intensity,
        )
        cfg := rgb_lerp(*fg, lc, linfo.Intensity * .5)
        fg = &cfg
    }

    bg := tile.GetBg()
    if bg != nil {
        lc := rgb_lerp(
            linfo.Color,
            self.ambient_light.Color,
            self.ambient_light.Intensity,
        )
        cbg := rgb_lerp(*bg, lc, linfo.Intensity * .5)
        bg = &cbg
    }

    light_level := fmax(linfo.Intensity, self.global_illumination)

    if fg == nil {
        Window.Set(
            x, y,
            (*tile.GetFg()).Scale(light_level),
            (*bg).Scale(light_level),
            string(tile.GetGlyph()),
        )
        return
    }

    if bg == nil {
        Window.Set(
            x, y,
            (*fg).Scale(light_level),
            (*tile.GetBg()).Scale(light_level),
            string(tile.GetGlyph()),
        )
        return
    }

    Window.Set(
        x, y,
        (*fg).Scale(light_level),
        (*bg).Scale(light_level),
        string(tile.GetGlyph()),
    )
    return

}

func (self *Camera) RenderLighting() {
    self.masking = nil

    lights := self.tilemap.GetLights(
        Point{self.ox, self.oy},
        Point{
            self.x + self.hw,
            self.y + self.hh,
        },
    )

    self.lighting = LightingMap{}
    for point, light := range lights {
        vmap := self.tilemap.Fov(
            point.X, point.Y, light.GetRadius(), true,
        )
        if (self.mask_origin != nil) &&
        (point == *self.mask_origin) &&
        (self.mask_radius ==  light.GetRadius()) {
            self.masking = &vmap
        }
        lmap := Illuminate(light, &vmap)
        self.lighting.Blend(lmap)
    }

    if (self.masking == nil) && 
    (self.mask_origin != nil) {
        mask := self.tilemap.Fov(
            self.mask_origin.X, 
            self.mask_origin.Y, 
            self.mask_radius, true,
        )
        self.masking = &mask
    }
}


func (self *Camera) CanSee(x, y int) bool {
    if self.masking == nil {
        return true
    } else {
        vinfo, ok := self.masking.Points[Point{x, y}]
        return ok && vinfo > 0.0
    }; panic(nil)
}

func (self *Camera) LightAt(x, y int) LightInfo {
    if self.lighting != nil {
        linfo, ok := self.lighting[Point{x, y}]
        if ok {
            return linfo
        }
    }
    return LightInfo{0.0, White}
}

func (self *Camera) Render() {
    // w := Window.GetWidth()
    // h := Window.GetHeight()

    Window.Clear(Black, Black, ' ')

    self.Each(func(tilemap *Map, ent Entity, x, y int){
            clear := true;
            switch ent != nil {
                case true:
                    dx := self.ox + x
                    dy := self.oy + y
                    if self.CanSee(dx, dy) {
                        linfo := self.LightAt(dx, dy)
                        self.RenderTile(x, y, ent, linfo)

                        if dx == self.x && dy == self.y {
                            Window.Set(
                                x, y,
                                Green.Scale(
                                    fmax(
                                        self.global_illumination,
                                        linfo.Intensity,
                                    )), nil, "웃",
                            );
                        }
                        clear = false;
                    }
            }
            if clear {
                Window.Set(
                    x, y,
                    Black, Black, " ",
                );
            }
    })
}