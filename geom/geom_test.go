package geom

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

func TestVectorNormalized(t *testing.T) {
	t.Parallel()
	err := quick.Check(func(v Vector) bool {
		return Float64Equals(v.Unit().SquaredMagnitude(), 1.0)
	}, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestVectorNormalize(t *testing.T) {
	t.Parallel()
	err := quick.Check(func(v Vector) bool {
		return Float64Equals(v.Unit().SquaredMagnitude(), 1.0)
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

func TestRay_PlaneIntersectionHit(t *testing.T) {
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

func TestRay_PlaneIntersectionMiss(t *testing.T) {
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
