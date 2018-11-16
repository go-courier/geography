package maptile

import (
	"github.com/go-courier/geography"
	"math"
)

func lonLatToPixelXY(lon, lat float64, zoom uint32) (x, y float64) {
	maxTiles := float64(uint32(1 << zoom))
	x = (lon/360.0 + 0.5) * maxTiles

	// bound it because we have a top of the world problem
	siny := math.Sin(lat * geography.D2R)

	if siny < -0.9999 {
		y = 0
	} else if siny > 0.9999 {
		y = maxTiles - 1
	} else {
		lat = 0.5 + 0.5*math.Log((1.0+siny)/(1.0-siny))/(-2*math.Pi)
		y = lat * maxTiles
	}
	return
}

func TrailingZeros32(x uint32) int {
	if x == 0 {
		return 32
	}
	return int(deBruijn32tab[(x&-x)*deBruijn32>>(32-5)])
}

// http://supertech.csail.mit.edu/papers/debruijn.pdf
const deBruijn32 = 0x077CB531

var deBruijn32tab = [32]byte{
	0, 1, 28, 2, 29, 14, 24, 3, 30, 22, 20, 15, 25, 17, 4, 8,
	31, 27, 13, 23, 21, 19, 16, 7, 26, 12, 18, 6, 11, 5, 10, 9,
}
