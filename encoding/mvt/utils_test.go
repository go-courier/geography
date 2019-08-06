package mvt

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestProjection(t *testing.T) {
	p := NewProjection(5, 26, 13, 4096)

	x, y := p.LonLatToTilePixelXY(107.77, 31.72)
	spew.Dump(x, y)

	lon, lat := p.TilePixelXYToLonLat(x, y)
	spew.Dump(lon, lat)
}
