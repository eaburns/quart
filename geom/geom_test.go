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
		if NearEqual(s, test.dist) {
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
		if NearEqual(d, test.dot) {
			continue
		}
		t.Errorf("Expected %v dot %v to be %f, got %f", test.a, test.b, test.dot, d)
	}
}

func TestVectorUnit(t *testing.T) {
	t.Parallel()
	err := quick.Check(func(v Vector) bool {
		return NearEqual(v.Unit().SquaredMagnitude(), 1.0)
	}, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestVectorInverse(t *testing.T) {
	t.Parallel()
	err := quick.Check(func(v Vector) bool {
		return v.Inverse().Plus(v).NearlyEquals(Vector{})
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
		if NearEqual(d, test.d) {
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
		if NearEqual(d, test.d) {
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
		if c.NearlyEquals(test.c) {
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
		if n.NearlyEquals(test.n) {
			continue
		}
		t.Errorf("Expected nearest point to %v on %v to %v to be %v, got %v",
			test.p, test.s0, test.s1, test.n, n)
	}
}

func BenchmarkPointPlus(b *testing.B) {
	p, v := Point{1, 1}, Vector{2, 2}
	for i := 0; i < b.N; i++ {
		p.Plus(v)
	}
}

func BenchmarkPointAdd(b *testing.B) {
	p, v := Point{1, 1}, Vector{2, 2}
	for i := 0; i < b.N; i++ {
		p.Add(v)
	}
}

func BenchmarkPointMinus(b *testing.B) {
	p0, p1 := Point{1, 1}, Point{2, 2}
	for i := 0; i < b.N; i++ {
		p0.Minus(p1)
	}
}

func BenchmarkPointTimes(b *testing.B) {
	p, v := Point{1, 1}, Vector{2, 2}
	for i := 0; i < b.N; i++ {
		p.Times(v)
	}
}

func BenchmarkPointSquaredDistance(b *testing.B) {
	p0, p1 := Point{1, 1}, Point{2, 2}
	for i := 0; i < b.N; i++ {
		p0.SquaredDistance(p1)
	}
}

func BenchmarkPointDistance(b *testing.B) {
	p0, p1 := Point{1, 1}, Point{2, 2}
	for i := 0; i < b.N; i++ {
		p0.Distance(p1)
	}
}

func BenchmarkPointNearlyEqualsDiff(b *testing.B) {
	p0, p1 := Point{1, 1}, Point{2, 2}
	for i := 0; i < b.N; i++ {
		p0.NearlyEquals(p1)
	}
}

func BenchmarkPointNearlyEqualsSame(b *testing.B) {
	p := Point{1, 1}
	for i := 0; i < b.N; i++ {
		p.NearlyEquals(p)
	}
}

func BenchmarkVectorPlus(b *testing.B) {
	v0, v1 := Vector{1, 1}, Vector{2, 2}
	for i := 0; i < b.N; i++ {
		v0.Plus(v1)
	}
}

func BenchmarkVectorAdd(b *testing.B) {
	v0, v1 := Vector{1, 1}, Vector{2, 2}
	for i := 0; i < b.N; i++ {
		v0.Add(v1)
	}
}

func BenchmarkVectorMinus(b *testing.B) {
	v0, v1 := Vector{1, 1}, Vector{2, 2}
	for i := 0; i < b.N; i++ {
		v0.Minus(v1)
	}
}

func BenchmarkVectorSubtract(b *testing.B) {
	v0, v1 := Vector{1, 1}, Vector{2, 2}
	for i := 0; i < b.N; i++ {
		v0.Subtract(v1)
	}
}

func BenchmarkVectorTimes(b *testing.B) {
	v0, v1 := Vector{1, 1}, Vector{2, 2}
	for i := 0; i < b.N; i++ {
		v0.Times(v1)
	}
}

func BenchmarkVectorScaledBy(b *testing.B) {
	v0 := Vector{1, 1}
	for i := 0; i < b.N; i++ {
		v0.ScaledBy(2)
	}
}

func BenchmarkVectorDot(b *testing.B) {
	v0, v1 := Vector{1, 1}, Vector{2, 2}
	for i := 0; i < b.N; i++ {
		v0.Dot(v1)
	}
}

func BenchmarkVectorSquaredMagnitude(b *testing.B) {
	v := Vector{1, 1}
	for i := 0; i < b.N; i++ {
		v.SquaredMagnitude()
	}
}

func BenchmarkVectorMagnitude(b *testing.B) {
	v := Vector{1, 1}
	for i := 0; i < b.N; i++ {
		v.Magnitude()
	}
}

func BenchmarkVectorUnit(b *testing.B) {
	v := Vector{1, 1}
	for i := 0; i < b.N; i++ {
		v.Unit()
	}
}

func BenchmarkVectorInverse(b *testing.B) {
	v := Vector{1, 1}
	for i := 0; i < b.N; i++ {
		v.Inverse()
	}
}

func BenchmarkVectorNearlyEqualssDiff(b *testing.B) {
	v0, v1 := Vector{1, 1}, Vector{2, 2}
	for i := 0; i < b.N; i++ {
		v0.NearlyEquals(v1)
	}
}

func BenchmarkVectorNearlyEqualsSame(b *testing.B) {
	v := Vector{1, 1}
	for i := 0; i < b.N; i++ {
		v.NearlyEquals(v)
	}
}

func BenchmarkRayPlaneIntersectionHit(b *testing.B) {
	r := Ray{Origin: Point{0, 0}, Direction: Vector{1, 0}}
	p := Plane{Origin: Point{1, 0}, Normal: Vector{-1, 0}}
	for i := 0; i < b.N; i++ {
		r.PlaneIntersection(p)
	}
}

func BenchmarkRayPlaneIntersectionMiss(b *testing.B) {
	r := Ray{Origin: Point{0, 0}, Direction: Vector{1, 0}}
	p := Plane{Origin: Point{0, 1}, Normal: Vector{0, -1}}
	for i := 0; i < b.N; i++ {
		r.PlaneIntersection(p)
	}
}

func BenchmarkRaySphereIntersectionHit(b *testing.B) {
	r := Ray{Origin: Point{0, 0}, Direction: Vector{1, 0}}
	s := Sphere{Center: Point{1, 0}, Radius: 1}
	for i := 0; i < b.N; i++ {
		r.SphereIntersection(s)
	}
}

func BenchmarkRaySphereIntersectionMiss(b *testing.B) {
	r := Ray{Origin: Point{0, 0}, Direction: Vector{1, 0}}
	s := Sphere{Center: Point{2, 2}, Radius: 1}
	for i := 0; i < b.N; i++ {
		r.SphereIntersection(s)
	}
}

func BenchmarkSegmentCenter(b *testing.B) {
	s := Segment{Point{0, 0}, Point{2, 2}}
	for i := 0; i < b.N; i++ {
		s.Center()
	}
}

func BenchmarkSegmentLength(b *testing.B) {
	s := Segment{Point{0, 0}, Point{2, 2}}
	for i := 0; i < b.N; i++ {
		s.Length()
	}
}

func BenchmarkSegmentNearestPoint(b *testing.B) {
	s := Segment{Point{0, 0}, Point{2, 2}}
	p := Point{1, 0}
	for i := 0; i < b.N; i++ {
		s.NearestPoint(p)
	}
}
