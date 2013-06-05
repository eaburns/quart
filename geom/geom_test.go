package geom

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

func TestSquaredDistance(t *testing.T) {
	t.Parallel()
	tests := []struct {
		a, b Point
		dist float64
	}{
		{Point{0, 0}, Point{0, 0}, 0},
		{Point{0, 0}, Point{0, 1}, 1},
		{Point{0, 0}, Point{1, 0}, 1},
		{Point{0, 0}, Point{math.Cos(math.Pi / 4), math.Sin(math.Pi / 4)}, 1},
		{Point{0, 0}, Point{1, 1}, 2},
		{Point{0, 0}, Point{2, 2}, 8},
	}
	for _, test := range tests {
		s := test.a.SquaredDistance(test.b)
		if Float64Equals(s, test.dist) {
			continue
		}
		t.Errorf("Expected squared distance of %f between %v and %v, got %f",
			test.dist, test.a, test.b, s)
	}
}

func TestVectorDot(t *testing.T) {
	t.Parallel()
	tests := []struct {
		a, b Vector
		dot  float64
	}{
		{Vector{1, 0}, Vector{0, 1}, 0},
		{Vector{0, 1}, Vector{1, 0}, 0},
		{Vector{1, 0}, Vector{math.Cos(math.Pi / 4), math.Sin(math.Pi / 4)}, math.Cos(math.Pi / 4)},
	}
	for _, test := range tests {
		d := test.a.Dot(test.b)
		if Float64Equals(d, test.dot) {
			continue
		}
		t.Errorf("Expected %v dot %v to be %f, got %f", test.a, test.b, test.dot, d)
	}
}

func TestVectorUnit(t *testing.T) {
	t.Parallel()
	err := quick.Check(func(v Vector) bool {
		return Float64Equals(v.Unit().SquaredMagnitude(), 1.0)
	}, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestVectorInverse(t *testing.T) {
	t.Parallel()
	err := quick.Check(func(v Vector) bool {
		return v.Inverse().Plus(v).Equals(Vector{})
	}, nil)
	if err != nil {
		t.Error(err)
	}
}

func (v Vector) Generate(r *rand.Rand, _ int) reflect.Value {
	for i := 0; i < K; i++ {
		v[i] = r.Float64()
	}
	return reflect.ValueOf(v)
}

func TestRayPlaneIntersectionHit(t *testing.T) {
	t.Parallel()
	tests := []struct {
		r Ray
		p Plane
		d float64
	}{
		{Ray{Point{}, Vector{1, 0}}, Plane{Point{1, 0}, Vector{-1, 0}}, 1},
		{Ray{Point{}, Vector{1, 0}}, Plane{Point{1, 1}, Vector{-1, 0}}, 1},
		{Ray{Point{}, Vector{1, 0}}, Plane{Point{0.5, 0}, Vector{-1, 0}}, 0.5},
		{Ray{Point{}, Vector{1, 0}}, Plane{Point{0.5, 0}, Vector{1, 0}}, 0.5},
		{Ray{Point{}, Vector{0, 1}}, Plane{Point{0, 1}, Vector{0, -1}}, 1},
		{
			Ray{Point{}, Vector{math.Cos(math.Pi / 4), math.Cos(math.Pi / 4)}},
			Plane{Point{1, 1}, Vector{-math.Cos(math.Pi / 4), -math.Cos(math.Pi / 4)}},
			math.Sqrt2,
		},
	}

	for _, test := range tests {
		d, _ := test.r.PlaneIntersection(test.p)
		if Float64Equals(d, test.d) {
			continue
		}
		t.Errorf("Expected %v to hit %v at %g, got %g", test.r, test.p, test.d, d)
	}
}

func TestRayPlaneIntersectionMiss(t *testing.T) {
	t.Parallel()
	tests := []struct {
		r Ray
		p Plane
	}{
		{Ray{Point{}, Vector{-1, 0}}, Plane{Point{1, 0}, Vector{-1, 0}}},
		{Ray{Point{}, Vector{0, 1}}, Plane{Point{1, 0}, Vector{-1, 0}}},
	}

	for _, test := range tests {
		d, ok := test.r.PlaneIntersection(test.p)
		if !ok || d < 0 {
			continue
		}
		t.Errorf("Expected %v to miss %v.  Got a hit at %g", test.r, test.p, d)
	}
}

func TestRaySphereIntersectionHit(t *testing.T) {
	t.Parallel()
	tests := []struct {
		r Ray
		s Sphere
		d float64
	}{
		{Ray{Point{}, Vector{1, 0}}, Sphere{Point{1, 0}, 1}, 0},
		{Ray{Point{}, Vector{1, 0}}, Sphere{Point{2, 0}, 1}, 1},
		{Ray{Point{}, Vector{1, 0}}, Sphere{Point{2, 0}, 2}, 0},
		{Ray{Point{}, Vector{-1, 0}}, Sphere{Point{-1, 0}, 1}, 0},
		{Ray{Point{}, Vector{-1, 0}}, Sphere{Point{-2, 0}, 1}, 1},
		{Ray{Point{}, Vector{-1, 0}}, Sphere{Point{-2, 0}, 2}, 0},
		{Ray{Point{}, Vector{0, 1}}, Sphere{Point{0, 1}, 1}, 0},
		{Ray{Point{}, Vector{0, 1}}, Sphere{Point{0, 2}, 1}, 1},
		{Ray{Point{}, Vector{0, 1}}, Sphere{Point{0, 2}, 2}, 0},
		{Ray{Point{}, Vector{0, -1}}, Sphere{Point{0, -1}, 1}, 0},
		{Ray{Point{}, Vector{0, -1}}, Sphere{Point{0, -2}, 1}, 1},
		{Ray{Point{}, Vector{0, -1}}, Sphere{Point{0, -2}, 2}, 0},
		{
			Ray{Point{}, Vector{1, 1}.Unit()},
			Sphere{Point{2, 2}, math.Sqrt(2)},
			math.Sqrt(2),
		},
	}

	for _, test := range tests {
		d, _ := test.r.SphereIntersection(test.s)
		if Float64Equals(d, test.d) {
			continue
		}
		t.Errorf("Expected %v to hit %v at %g, got %g", test.r, test.s, test.d, d)
	}
}

func TestRaySphereIntersectionMiss(t *testing.T) {
	t.Parallel()
	tests := []struct {
		r Ray
		s Sphere
	}{
		{Ray{Point{}, Vector{-1, 0}}, Sphere{Point{1, 0}, 0.5}},
		{Ray{Point{}, Vector{0, 1}}, Sphere{Point{1, 0}, 0.5}},
	}

	for _, test := range tests {
		d, ok := test.r.SphereIntersection(test.s)
		if !ok || d < 0 {
			continue
		}
		t.Errorf("Expected %v to miss %v.  Got a hit at %g", test.r, test.s, d)
	}
}

func TestSegmentCenter(t *testing.T) {
	t.Parallel()
	tests := []struct {
		s0, s1, c Point
	}{
		{Point{0, 0}, Point{1, 0}, Point{0.5, 0}},
		{Point{0, 0}, Point{1, 1}, Point{0.5, 0.5}},
	}

	for _, test := range tests {
		c := Segment{test.s0, test.s1}.Center()
		if c.Equals(test.c) {
			continue
		}
		t.Errorf("Expected %v to be center of %v to %v, got %v", test.c, test.s0, test.s1, c)
	}
}

func TestSegmentNearestPoint(t *testing.T) {
	t.Parallel()
	tests := []struct {
		s0, s1, p, n Point
	}{
		{Point{-1, 0}, Point{1, 0}, Point{2, 0}, Point{1, 0}},
		{Point{-1, 0}, Point{1, 0}, Point{0, 1}, Point{0, 0}},
		{Point{-1, 0}, Point{1, 0}, Point{0, -1}, Point{0, 0}},
		{Point{-1, -1}, Point{1, 1}, Point{-1, 1}, Point{0, 0}},
		{Point{-1, -1}, Point{1, 1}, Point{1, -1}, Point{0, 0}},
	}

	for _, test := range tests {
		n := Segment{test.s0, test.s1}.NearestPoint(test.p)
		if n.Equals(test.n) {
			continue
		}
		t.Errorf("Expected nearest point to %v on %v to %v to be %v, got %v",
			test.p, test.s0, test.s1, test.n, n)
	}
}
