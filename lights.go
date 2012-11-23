package main

import "fmt"
import "math"

type LightInfo struct {
    Intensity float64
    Color RGB
}

type ILight interface {
    GetIntensity() float64
    GetColor() RGB
    GetRadius() int
}

type Light struct {
    Level float64
    Color RGB 
    Rad int 
}

func (self *Light) GetIntensity() float64 {
    return self.Level
}

func (self *Light) GetColor() RGB {
    return self.Color
}

func (self *Light) GetRadius() int {
    return self.Rad
}

type ILightableMap interface {
    In(x, y int) bool
    Width() int
    Height() int
    GetLights(min, max Point) LightMap
}

type LightMap map[Point]ILight
type LightingMap map[Point]LightInfo

func (self LightingMap) Blend(other LightingMap) {
    for point, linfo := range other {
        mine, ok := self[point]
        if ok {
            if math.Abs(linfo.Intensity - mine.Intensity) <= 0.01 {
                 linfo.Color = rgb_lerp(mine.Color, linfo.Color, .5)
            } else if linfo.Intensity - mine.Intensity > .001 {
                linfo.Color = rgb_lerp(linfo.Color, mine.Color, mine.Intensity)
                linfo.Intensity = linfo.Intensity
            } else if mine.Intensity - linfo.Intensity > .001 {
                linfo.Color = rgb_lerp(mine.Color, linfo.Color, linfo.Intensity)
                linfo.Intensity = mine.Intensity
            } 
        }
        self[point] = linfo
    }
}

func intensity(px, py, cx, cy, r int) float64 {
    r2 := float64(r * r)
    squaredDist := float64((px-cx)*(px-cx) + (py-cy)*(py-cy))
    coef1 := 1.0 / (1.0 + squaredDist / float64(r))
    coef2 := coef1 - 1.0 / (1.0 + r2)
    result := coef2 / (1.0 - 1.0/(1.0 + r2))
    fmt.Println(result)
    return fmin(1.0 , fmax(result, 0.0))
}

func Illuminate(light ILight, vmap *VisibilityMap) LightingMap {
    lmap := make(LightingMap, 0)
    origin := vmap.Origin
    for k, _ := range vmap.Points {
        level := intensity(origin.X, origin.Y, k.X, k.Y, light.GetRadius() - 1)
        linfo := LightInfo{level, light.GetColor()}
        lmap[Point{k.X, k.Y}] = linfo;
    }
    return lmap
}
                // linfo, ok := data[Point{curx, cury}];
                // c := color;
                // i := intensity(xo, yo, curx, cury, 10);
                // if ok {
                //     c.Overlay(linfo.Color);
                //     i = i * linfo.Intensity;
                // }
