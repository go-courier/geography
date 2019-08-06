package mvt

import (
	"math"
	"strings"

	"github.com/go-courier/geography/encoding/mvt/vector_tile"
)

func TileGeomTypeFromGeoType(tpe string) vector_tile.Tile_GeomType {
	switch strings.ToUpper(tpe) {
	case "POINT", "MULTIPOINT":
		return vector_tile.Tile_POINT
	case "POLYGON", "MULTIPOLYGON":
		return vector_tile.Tile_POLYGON
	case "LINESTRING", "MULTILINESTRING":
		return vector_tile.Tile_LINESTRING
	}
	return vector_tile.Tile_UNKNOWN
}

func NewProjection(
	z uint32,
	x uint32,
	y uint32,
	extent uint32,
) *Projection {
	p := &Projection{
		Z:      z,
		X:      x,
		Y:      y,
		Extent: extent,
	}
	p.Init()
	return p
}

type Projection struct {
	Z        uint32
	X        uint32
	Y        uint32
	Extent   uint32
	minx     float64
	miny     float64
	maxTiles float64
}

var (
	D2R = math.Pi / 180.0
)

func (p *Projection) Init() {
	n := uint32(trailingZeros32(p.Extent))
	z := p.Z + n
	p.minx = float64(p.X << n)
	p.miny = float64(p.Y << n)
	p.maxTiles = float64(uint32(1 << z))
}

func (p *Projection) LonLatToTilePixelXY(lon, lat float64) (x, y float64) {
	x = (lon/360.0 + 0.5) * p.maxTiles

	// bound it because we have a top of the world problem
	siny := math.Sin(lat * D2R)

	if siny < -0.9999 {
		y = 0
	} else if siny > 0.9999 {
		y = p.maxTiles - 1
	} else {
		lat = 0.5 + 0.5*math.Log((1.0+siny)/(1.0-siny))/(-2*math.Pi)
		y = lat * p.maxTiles
	}

	return x - p.minx, y - p.miny
}

func (p *Projection) TilePixelXYToLonLat(x, y float64) (lon, lat float64) {
	y2 := 180 - (y+p.miny)*360/p.maxTiles

	lon = (x+p.minx)*360/p.maxTiles - 180
	lat = 360/math.Pi*math.Atan(math.Exp(y2*math.Pi/180)) - 90
	return
}

func trailingZeros32(x uint32) int {
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
