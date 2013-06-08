// Package geom provides geometric primitives for 2-dimensional Euclidean space.
package geom

// This file contains geometry primitives that work in K dimensions.

import (
	"math"
)

const (
	// K is the number of dimensions of the geometric primitives.
	K = 2

	// Threshold is the amount by which two floating points must differ
	// to be considered different by the equality rountines in this package.
	//
	// The current value is the square root of the IEEE 64-bit floating point
	// epsilon value.  This is the value recommended in Numerical
	// Recipes.
	Threshold = 1.4901161193847656e-08
)

// Float64Equals returns true if the two floating point numbers are
// close enough to be considered equal.
func Float64Equals(a, b float64) bool {
	return math.Abs(a-b) < Threshold
}

// A Point is a location in K-space.
type Point [K]float64

// Plus returns the sum a point and a vector.
func (p Point) Plus(v Vector) Point {
	p.Add(v)
	return p
}

// Add adds a vector to a point.
func (p *Point) Add(v Vector) {
	for i, vi := range v {
		p[i] += vi
	}
}

// Minus returns the difference between two points.
func (a Point) Minus(b Point) Vector {
	for i, bi := range b {
		a[i] -= bi
	}
	return Vector(a)
}

// Times returns the component-wise product of a point and a vector.
func (p Point) Times(v Vector) Point {
	for i, vi := range v {
		p[i] *= vi
	}
	return p
}

// SquaredDistance returns the squared distance between two points.
func (a Point) SquaredDistance(b Point) float64 {
	dist := 0.0
	for i, ai := range a {
		bi := b[i]
		d := ai - bi
		dist += d * d
	}
	return dist
}

// Distance returns the distance between two points.
func (a Point) Distance(b Point) float64 {
	return math.Sqrt(a.SquaredDistance(b))
}

// Equals returns true if the points are close enough to be considered equal.
func (a Point) Equals(b Point) bool {
	for i, ai := range a {
		if !Float64Equals(ai, b[i]) {
			return false
		}
	}
	return true
}

// A Vector is a direction and magnitude in K-space.
type Vector [K]float64

// Plus returns the sum of two vectors.
func (a Vector) Plus(b Vector) Vector {
	a.Add(b)
	return a
}

// Add adds a vector to the receiver vector.
func (a *Vector) Add(b Vector) {
	for i, bi := range b {
		a[i] += bi
	}
}

// Minus returns the difference between two vectors.
func (a Vector) Minus(b Vector) Vector {
	for i, bi := range b {
		a[i] -= bi
	}
	return a
}

// Subtract subtracts a vector from the receiver
func (a *Vector) Subtract(b Vector) {
	for i, bi := range b {
		a[i] -= bi
	}
}

// Times returns the component-wise product of two vectors.
func (a Vector) Times(b Vector) Vector {
	for i, bi := range b {
		a[i] *= bi
	}
	return a
}

// ScaledBy returns the product of a vector and a scalar.
func (v Vector) ScaledBy(k float64) Vector {
	for i := range v {
		v[i] *= k
	}
	return v
}

// Dot returns the dot product of two vectors.
func (a Vector) Dot(b Vector) float64 {
	dot := 0.0
	for i, ai := range a {
		dot += ai * b[i]
	}
	return dot
}

// SquaredMagnitude returns the squared magnitude of the vector.
func (v Vector) SquaredMagnitude() float64 {
	m := 0.0
	for _, vi := range v {
		m += vi * vi
	}
	return m
}

// Magnitude returns the magnitude of the vector.
func (v Vector) Magnitude() float64 {
	return math.Sqrt(v.SquaredMagnitude())
}

// Unit returns the normalized unit form of the vector.
func (v Vector) Unit() Vector {
	m := v.Magnitude()
	for i := range v {
		v[i] /= m
	}
	return v
}

// Inverse returns the vector point in the opposite direction.
func (v Vector) Inverse() Vector {
	for i, vi := range v {
		v[i] = -vi
	}
	return v
}

// Equals returns true if the vectors are close enough to be considered equal.
func (a Vector) Equals(b Vector) bool {
	for i, ai := range a {
		if !Float64Equals(ai, b[i]) {
			return false
		}
	}
	return true
}

// A Plane represented by a point and its normal vector.
type Plane struct {
	Origin Point
	// Normal is the unit vector perpendicular to the plane.
	Normal Vector
}

// A Ray is an origin point and a direction vector.
type Ray struct {
	Origin Point
	// Direction is the unit vector giving the direction of the ray.
	Direction Vector
}

// PlaneIntersection returns the distance along the ray at which it intersects a
// plane. The second return value is true if they do intersect, and it is false if
// they do not intersect.
func (r Ray) PlaneIntersection(p Plane) (float64, bool) {
	d := -p.Normal.Dot(Vector(p.Origin))
	numer := p.Normal.Dot(Vector(r.Origin)) + d
	denom := r.Direction.Dot(p.Normal)
	if Float64Equals(denom, 0) {
		return 0, false
	}
	return -numer / denom, true
}

// SphereIntersection returns the distance along the ray at which it intersects a.
// sphere. The second return value is true if they do intersect, and it is false if
// they do not intersect.
func (r Ray) SphereIntersection(s Sphere) (float64, bool) {
	Q := s.Center.Minus(r.Origin)
	c := Q.Magnitude()
	v := Q.Dot(r.Direction)
	d := s.Radius*s.Radius - (c*c - v*v)
	if d < 0 {
		return 0, false
	}
	return v - math.Sqrt(d), true
}

// A Segment is the portion of a line between and including two points.
type Segment [2]Point

// Center returns the point at the center of the face.
func (s Segment) Center() Point {
	d := s[1].Minus(s[0]).Unit()
	l := s.Length()
	return s[0].Plus(d.ScaledBy(l / 2))
}

// Length returns the length of the face.
func (s Segment) Length() float64 {
	return s[0].Distance(s[1])
}

// NearestPoint returns the point on the face nearest to p.
func (s Segment) NearestPoint(p Point) Point {
	V := s[1].Minus(s[0])
	d := V.Magnitude()
	V = V.Unit()
	t := V.Dot(p.Minus(s[0]))

	switch {
	case t < 0:
		return s[0]
	case t > d:
		return s[1]
	}
	return s[0].Plus(V.ScaledBy(t))
}

// A Sphere is the set of all points at a fixed distance from a center point.
type Sphere struct {
	Center Point
	Radius float64
}

// An Ellipsoid is like a sphere, but it has one radius for each axis.
type Ellipsoid struct {
	Center Point
	Radii  Vector
}
