package main

import (
    "image"
    "math"
    "math/rand"
)

func lerp(a, b, amt float64) float64 {
    return a + amt * (b - a)
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func clamp(val, lower, upper) int {
    return min(upper, max(lower, val))
}

func fmax(a, b float64) float64 {
    if a > b {
        return a
    }
    return b
}

func fmin(a, b float64) float64 {
    if a < b {
        return a
    }
    return b
}

func fclamp(val, lower, upper) float64 {
    return fmin(upper, fmax(lower, val))
}

func round(x float64, prec int) float64 {
    var rounder float64
    pow := math.Pow(10, float64(prec))
    intermed := x * pow

    if intermed < 0.0 {
        intermed -= 0.5
    } else {
        intermed += 0.5
    }
    rounder = float64(int64(intermed))
    return rounder / float64(pow)
}

type Point image.Point

type Object interface{}

func In(min, max, point Point) bool {
    if point.X >= min.X && point.Y >= min.Y {
        if point.X <= max.X && point.Y <= point.Y {
            return true
        }
    }
    return false
}

func RectRange(min, max Point, fn func(p Point)) {
    for x:=min.X; x<=max.X; x++ {
        for y:=min.Y; y<=max.Y; y++ {
            fn(Point{x, y})
        }
    }
}

func RandRangeFloat(min, max float64) float64 {
    return (rand.Float64() * (max - min)) + min
}

func RandRangeInt(min, max int) int {
    return int(RandRangeFloat(float64(min), float64(max)))
}

func RangeScaleFloat(min, max, val float64) float64 {
    return 1-(val/(min+max))
}
