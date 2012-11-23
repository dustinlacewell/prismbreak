package main


type IVisible interface {
	GetVisibility() float64
}

type IViewableMap interface {
	In(x, y int) bool
	Width() int
	Height() int
	GetVisibility(x, y int) float64
	GetVisibilityMap() VisibilityMap
}

type VisibilityMap struct {
	Points map[Point]float64
	Origin Point
}

func NewVisibilityMap(x, y int) (VisibilityMap) {
	return VisibilityMap {
		make(map[Point]float64, 0),
		Point{x, y},
	}
}

func (self VisibilityMap) Get(x, y int) (float64, bool) {
	t, ok := self.Points[Point{x, y}]
	return t, ok
}

func (self VisibilityMap) Set(x, y int, v float64) {
	self.Points[Point{x, y}] = v
}

func (self VisibilityMap) Update(other VisibilityMap) {
	for k, v := range other.Points{
		self.Points[k] = v
	}
}

type FOVAlgo func(IViewableMap, int, int, int, bool) VisibilityMap

func fovCircularCastRay(fov IViewableMap, xo, yo, xd, yd, r2 int, walls bool) VisibilityMap {
	data := NewVisibilityMap(xo, yo)
	curx := xo;	cury := yo
	in := false; blocked := false

	if fov.In(curx, cury) {
		in = true
		data.Points[Point{curx, cury}] = fov.GetVisibility(curx, cury)
	}

	for _, p := range Line(xo, yo, xd, yd) {
		curx = p.X; cury = p.Y

		if r2 > 0 {
			curRadius := (curx-xo)*(curx-xo) + (cury-yo)*(cury-yo)
			if curRadius > r2 {
				break
			}
		}

		if fov.In(curx, cury) {
			in = true
			visibility := fov.GetVisibility(curx, cury)
			if !blocked && visibility == 0.0 {
				if !(curx == xo && cury == yo) {
					blocked = true
				}
			} else if blocked {
				break
			}
			if walls || !blocked {
				data.Points[Point{curx, cury}] = visibility;
			}
		} else if in {
			break
		}
	}
	return data
}

func fovCircularPostProc(fov IViewableMap, vdata VisibilityMap, x0, y0, x1, y1, dx, dy int) {
	data := NewVisibilityMap(vdata.Origin.X, vdata.Origin.Y)
	for cx := x0; cx <= x1; cx++ {
		for cy := y0; cy <= y1; cy++ {
			x2 := cx + dx
			y2 := cy + dy
			seen := false
			visibility, ok := vdata.Points[Point{cx, cy}];
			if ok {
				if visibility > 0.0 {
					seen = true
				}
			}

			if fov.In(cx, cy) && seen && fov.GetVisibility(cx, cy) > 0.0 {
				if x2 >= x0 && x2 <= x1 {
					if fov.In(x2, cy) && fov.GetVisibility(x2, cy) == 0.0 {
						data.Points[Point{x2, cy}] = visibility;
					}
				}
				if y2 >= y0 && y2 <= y1 {
					if fov.In(cx, y2) && fov.GetVisibility(cx, y2) == 0.0 {
						data.Points[Point{cx, y2}] = visibility;
					}
				}
				if x2 >= x0 && x2 <= x1 && y2 >= y0 && y2 <= y1 {
					if fov.In(x2, y2) && fov.GetVisibility(x2, y2) == 0.0 {
						data.Points[Point{x2, y2}] = visibility;
					}
				}
			}
		}
	}
	vdata.Update(data)
}

// FOVCicular raycasts out from the vantage in a circle.
func FOVCircular(fov IViewableMap, x, y, r int, walls bool) VisibilityMap {
	data := NewVisibilityMap(x, y)
	xmin := 0; ymin := 0
	xmax := fov.Width()
	ymax := fov.Height()
	r2 := r * r
	if r > 0 {
		xmin = max(0, x-r)
		ymin = max(0, y-r)
		xmax = min(fov.Width(), x+r+1)
		ymax = min(fov.Height(), y+r+1)
	}
	xo := xmin; yo := ymin

	for xo < xmax {
		data.Update(fovCircularCastRay(fov, x, y, xo, yo, r2, walls))
		xo++
	}
	xo = xmax - 1
	yo = ymin + 1
	for yo < ymax {
		data.Update(fovCircularCastRay(fov, x, y, xo, yo, r2, walls))
		yo++
	}
	xo = xmax - 2
	yo = ymax - 1
	for xo >= 0 {
		data.Update(fovCircularCastRay(fov, x, y, xo, yo, r2, walls))
		xo--
	}
	xo = xmin
	yo = ymax - 2
	for yo > 0 {
		data.Update(fovCircularCastRay(fov, x, y, xo, yo, r2, walls))
		yo--
	}
	if walls {
		fovCircularPostProc(fov, data, xmin, ymin, x, y, -1, -1)
		fovCircularPostProc(fov, data, x, ymin, xmax-1, y, 1, -1)
		fovCircularPostProc(fov, data, xmin, y, x, ymax-1, -1, 1)
		fovCircularPostProc(fov, data, x, y, xmax-1, ymax-1, 1, 1)
	}
	return data
}
