package main

import (
	"fmt"
	"math"
	"os"
)

type Matrix [][]interface{}

type Point struct {
	X int
	Y int
}

func main() {
	in := Matrix{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	// center point
	p := Point{0, 0}
	for i := 0; true; i++ {
		old := p
		s := int(math.Abs(float64(p.X)) + math.Abs(float64(p.Y)))
		switch i % 4 {
		case 0:
			// step forward
			s += 1
			// left
			p.X -= s
		case 1:
			// down
			p.Y += s
		case 2:
			// right
			p.X += s
		case 3:
			// up
			p.Y -= s
		}
		printPath(in, old, p)
	}
}

func printPath(m Matrix, f, t Point) {
	xd, yd := 1, 1
	if t.X-f.X < 0 {
		xd = -1
	}
	if t.Y-f.Y < 0 {
		yd = -1
	}
	for f.X != t.X || f.Y != t.Y {
		fmt.Printf("%d ", m.valueAt(f))

		if f.X != t.X {
			f.X += xd
		}
		if f.Y != t.Y {
			f.Y += yd
		}
		if !m.isValid(f) {
			fmt.Print("\n")
			os.Exit(0)
		}
	}
}

// Matrix
//
func (m Matrix) isValid(rel Point) bool {
	p := m.relToAbs(rel)
	return p.X >= 0 && p.X < len(m) && p.Y >= 0 && p.Y < len(m[0])
}

func (m Matrix) valueAt(rel Point) interface{} {
	p := m.relToAbs(rel)
	return m[p.X][p.Y]
}

func (m Matrix) relToAbs(rel Point) Point {
	cx := (len(m) - 1) / 2
	cy := (len(m[0]) - 1) / 2
	return Point{
		X: rel.Y + cx,
		Y: rel.X + cy,
	}
}
