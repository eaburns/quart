// © 2012 the Quart Authors under the MIT license. See AUTHORS for the list of authors.

package geom

// This file contains geometry that is specific to 2 dimensions.
// This assignment will fail for K != 2.
var ensure2d [2]float64 = Vector{}

// A Line is a 2-dimensional Plane.
type Line Plane

// Direction returns a vector along the direction of the line.
func (l Line) Direction() Vector {
	d := l.Normal
	d[0], d[1] = -d[1], d[0]
	return d
}

// LineIntersection returns the point at which two lines intersect.
// The second return value is true if they do intersect, and it is
// false if they do not intersect.
func (a Line) LineIntersection(b Line) (Point, bool) {
	r := Ray{Origin: a.Origin, Direction: a.Direction()}
	d, hit := r.PlaneIntersection(Plane(b))
	if !hit {
		return Point{}, false
	}
	return r.Origin.Plus(r.Direction.ScaledBy(d)), true
}

// Normal returns the normal vector of the segment.
func (s Segment) Normal() Vector {
	n := s[1].Minus(s[0]).Unit()
	n[0], n[1] = -n[1], n[0]
	return n
}

// Line returns the line containing the segment.
func (s Segment) Line() Line {
	return Line{Origin: s[0], Normal: s.Normal()}
}

// A Circle is a 2-dimensional sphere.
type Circle Sphere

// An Ellipse is a 2-dimensional ellipsoid.
type Ellipse Ellipsoid

// A Rectangle represents a rectangular region of space.
type Rectangle struct {
	Min  Point
	Size Vector
}

// Max returns the point on the rectangle with the maximum x and y values.
func (r *Rectangle) Max() Point {
	return r.Min.Plus(r.Size)
}

// Center returns the point in the center of the rectangle.
func (r *Rectangle) Center() Point {
	return r.Min.Plus(r.Size.ScaledBy(0.5))
}
