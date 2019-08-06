package mvt

import (
	"errors"

	"github.com/davecgh/go-spew/spew"
)

var (
	InvalidMVTPoint      = errors.New("invalid point")
	InvalidMVTLineString = errors.New("invalid line string")
	InvalidMVTPolygon    = errors.New("invalid polygon")
)

func NewMVTGeometryReader(data []uint32) MVTGeometryReader {
	return &geometryReader{
		data: data,
	}
}

type geometryReader struct {
	data  []uint32
	prevX int32
	prevY int32
}

func (g *geometryReader) ReadPoints() ([]Coord, error) {
	c, points := g.Read()
	if c != CommandMoveTo {
		return nil, InvalidMVTPoint
	}
	return points, nil
}

func (g *geometryReader) ReadLineStrings() ([][]Coord, error) {
	spew.Dump(g.data)

	lines := make([][]Coord, 0)

	for len(g.data) > 0 {
		line, err := g.ReadLineString()
		if err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}

	return lines, nil
}

func (g *geometryReader) ReadLineString() ([]Coord, error) {
	c, points := g.Read()
	if c != CommandMoveTo {
		return nil, InvalidMVTLineString
	}
	c2, points2 := g.Read()
	if c2 != CommandLineTo {
		return nil, InvalidMVTLineString
	}
	return append(points, points2...), nil
}

func (g *geometryReader) ReadPolygon() ([][]Coord, error) {
	rings := make([][]Coord, 0)

	for len(g.data) > 0 {
		ring, err := g.ReadRing()
		if err != nil {
			return nil, err
		}
		rings = append(rings, ring)
	}

	return rings, nil
}

func (g *geometryReader) ReadRing() ([]Coord, error) {
	c, points := g.Read()
	if c != CommandMoveTo {
		return nil, InvalidMVTPolygon
	}
	c2, points2 := g.Read()
	if c2 != CommandLineTo {
		return nil, InvalidMVTPolygon
	}
	c3, _ := g.Read()
	if c3 != CommandClosePath {
		return nil, InvalidMVTPolygon
	}
	return append(points, append(points2, points[0])...), nil
}

func (g *geometryReader) readCoord() (int32, int32) {
	dx := fromV(g.data[0])
	dy := fromV(g.data[1])

	x, y := dx+g.prevX, dy+g.prevY

	g.prevX = x
	g.prevY = y

	g.data = g.data[2:]
	return x, y
}

func (g *geometryReader) readCmd() uint32 {
	cmdV := g.data[0]

	if len(g.data) > 0 {
		g.data = g.data[1:]
	}

	return cmdV
}

func (g *geometryReader) Read() (Command, []Coord) {
	i := g.readCmd()
	c := i % 8
	n := (i ^ c) >> 3

	points := make([]Coord, n)

	for i := uint32(0); i < n; i++ {
		if len(g.data) >= 2 {
			x, y := g.readCoord()
			points[i] = point{float64(x), float64(y)}
		}
	}

	return Command(c), points
}

type point [2]float64

func (p point) X() float64 {
	return p[0]
}

func (p point) Y() float64 {
	return p[1]
}
