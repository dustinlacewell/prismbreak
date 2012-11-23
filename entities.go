package main

type Floor struct { Tile }
func NewFloor(x, y int) *Floor {
    floor := Floor {
        Tile: Tile {
            X:x, Y:y,
            Fg: &LightAmber,
            Bg: &DarkestAmber,
            Glyph: '.',
            Roughness: PATH_MIN,
            Visibility: 1.0,
        },
    }
    return &floor
}

type Wall struct { 
    Tile 
    Light 
}
func NewWall(x, y int) *Wall {
    color := RGB {
        uint8(RandRangeInt(0, 255)),
        uint8(RandRangeInt(0, 255)),
        uint8(RandRangeInt(0, 255)),
    }
    wall := Wall {
        Tile: Tile {
            X:x, Y:y,
            Fg: &color,
            Bg: &DarkestAmber,
            Glyph: 'O',
            Roughness: PATH_MAX,
            Visibility: 0.0,
        },
        Light: Light {
            Level: 1.0,
            Color: color,
            Rad: 9,
        },
    }
    return &wall
}
