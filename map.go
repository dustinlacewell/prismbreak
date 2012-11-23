package main

import (
	"os"
	"log"
	"encoding/gob"
	"image"
)

type Entity interface {
	Renderable
	IVisible
	Walkable
}

type Map struct {
	W, H          int
	Tiles map[Point]Entity
	VisibilityMap VisibilityMap
}

type TileProcessor func(*Map, Entity, int, int)

func NewMap(width, height int) *Map {
	return &Map{
		width, height, 
		make(map[Point]Entity, 0),
		NewVisibilityMap(0, 0),
	}
}
func Load(filename string) *Map {
	fp, err := os.Open(filename);
	if err != nil {
		log.Fatal(err);
	}
	defer fp.Close()

	var tilemap *Map
	dec := gob.NewDecoder(fp);
	err = dec.Decode(&tilemap);
	if err != nil {
		log.Fatal(err);
	}

	return tilemap
}

func (self *Map) Save(filename string) {
	fp, err := os.Create(filename);
	if err != nil {
		log.Fatal(err);
	}
	defer fp.Close()

	enc := gob.NewEncoder(fp);
	err = enc.Encode(self);
	if err != nil {
		log.Fatal(err);
	}
}

func (self *Map) Each(proc TileProcessor) {
	for x:=0; x<=self.W; x++ {
		for y:=0; y<=self.H; y++ {
			tile, ok := self.Tiles[Point{x, y}];
			if ok {
				proc(self, tile, x, y);
			} else {
				proc(self, nil, x, y);
			}
		}
	}
}


// In returns whether the coordinate is inside the map bounds.
func (self *Map) In(x, y int) bool {
	return x >= 0 && x < self.W && y >= 0 && y < self.H
}

func (self *Map) SetTile(e Entity, x, y int) {
	self.Tiles[Point{x, y}]	= e
}

func (self *Map) GetTile(x, y int) (Entity, bool) {
	e, ok := self.Tiles[Point{x, y}]
	return e, ok
}

func (self *Map) GetVisibility(x, y int) float64 {
	tile, ok := self.Tiles[Point{x, y}]
	switch ok {
		case true:
			return tile.GetVisibility()
		case false:
			return 0.0
	}; panic(nil)
}

func (self *Map) GetVisibilityMap() VisibilityMap {
	return self.VisibilityMap
}

func (self *Map) Roughness(x, y int) int {
	tile, ok := self.Tiles[Point{x, y}]
	switch ok {
		case true:
			return tile.GetRoughness()
		case false:
			return PATH_MAX
	}; panic(nil)
}

// Update runs the give fov alogrithm on the map.
func (self *Map) Fov(x, y, radius int, includeWalls bool) VisibilityMap {
	return FOVCircular(self, x, y, radius, includeWalls)
}

func (self *Map) GetLights(min, max Point) LightMap {
	lmap := make(LightMap, 0)
    RectRange(min, max, func(p Point){
        tile, ok := self.GetTile(p.X, p.Y)
        if ok {
        	light, ok := tile.(ILight)
        	if ok {
        		lmap[p] = light
        	}
        }
    })
    return lmap
}

// Performs astar and returns the list of cells on the path.
func (self *Map) Path(x0, y0, x1, y1 int) []image.Point {
	nodes := Astar(self, x0, y0, x1, y1, true)
	points := make([]image.Point, len(nodes))
	for i := 0; i < len(nodes); i++ {
		points[i] = image.Pt(nodes[i].X, nodes[i].Y)
	}
	return points
}

func (self *Map) Look(x, y int) float64 {
	return self.GetVisibility(x, y)
}

// Clear resets the map to completely unblocked but unviewable.
func (self *Map) Clear() {
	self.VisibilityMap = NewVisibilityMap(0, 0)
}

// Width returns the width in cells of the map.
func (self *Map) Width() int {
	return self.W
}

// Height returns the height in cells of the map.
func (self *Map) Height() int {
	return self.H
}
