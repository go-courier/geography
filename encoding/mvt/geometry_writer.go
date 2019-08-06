package mvt

import (
	"fmt"
)

type Command uint32

const (
	CommandMoveTo    Command = 1
	CommandLineTo    Command = 2
	CommandClosePath Command = 7
)

func NewMVTGeometryWriter(cap int) MVTGeometryWriter {
	return &geometryWriter{
		data: make([]uint32, 0, cap),
	}
}

type geometryWriter struct {
	data         []uint32
	prevX, prevY int32
}

func (w *geometryWriter) Data() []uint32 {
	return w.data
}

func (w *geometryWriter) MoveTo(l int, getCoord func(i int) Coord) {
	w.data = append(w.data, (uint32(l)<<3)|uint32(CommandMoveTo))
	for i := 0; i < l; i++ {
		w.writeCoord(getCoord(i))
	}
}

func (w *geometryWriter) LineTo(l int, getCoord func(i int) Coord) {
	w.data = append(w.data, (uint32(l)<<3)|uint32(CommandLineTo))
	for i := 0; i < l; i++ {
		w.writeCoord(getCoord(i))
	}
}

func (w *geometryWriter) ClosePath() {
	w.data = append(w.data, (1<<3)|uint32(CommandClosePath))
}

func (w *geometryWriter) writeCoord(coord Coord) {
	x0 := int32(coord.X())
	y0 := int32(coord.Y())

	dx := x0 - w.prevX
	dy := y0 - w.prevY

	w.prevX = x0
	w.prevY = y0

	fmt.Println(dy, toV(dy), fromV(toV(dy)))

	w.data = append(
		w.data,
		toV(dx),
		toV(dy),
	)
}

func toV(u int32) uint32 {
	return uint32((u << 1) ^ (u >> 31))
}

func fromV(u uint32) int32 {
	v := u / 2
	sign := u % 2
	if sign == 1 {
		return -int32(v) - 1
	}
	return int32(v)
}
