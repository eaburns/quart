package phys

// Basically an implementation of: http://www.paulnettle.com/pub/FluidStudios/CollisionDetection/Fluid_Studios_Generic_Collision_Detection_for_Games_Using_Ellipsoids.pdf

import (
	"math"

	. "github.com/eaburns/quart/geom"
)

// MoveCircle moves a circle with a given velocity, handling collision with segments.
func MoveCircle(c Circle, v Vector, segs []Segment) Circle {
	for !v.Equals(Vector{}) {
		d, v2 := moveCircle1(c, v, segs)
		c.Center.Add(v.Unit().ScaledBy(d))
		v = v2
	}
	return c
}

// moveCircle1 traces a circle along a velocity vector until the first collision
// with a Segment.  The return value is the distance that the circle moved,
// and the new velocity vector.
func moveCircle1(c Circle, v Vector, segs []Segment) (float64, Vector) {
	hitPt := Point{}
	dist := math.Inf(1)

	for _, s := range segs {
		if d, pt, hit := CircleSegmentHit(c, v, s); hit && d < dist {
			dist = d
			hitPt = pt
		}
	}
	if math.IsInf(dist, 1) {
		return v.Magnitude(), Vector{}
	}

	vUnit := v.Unit()
	c.Center.Add(vUnit.ScaledBy(dist))
	slide := Plane{Origin: hitPt, Normal: hitPt.Minus(c.Center).Unit()}

	dest := hitPt.Plus(vUnit.ScaledBy(v.Magnitude() - dist))
	r := Ray{Origin: dest, Direction: slide.Normal}
	d, hit := r.PlaneIntersection(slide)
	if !hit {
		panic("Couldn't project to the sliding plane!")
	}
	dest.Add(slide.Normal.Unit().ScaledBy(d))
	return dist - Threshold, dest.Minus(hitPt)
}

// CircleSegmentHit returns information about the collision of a circle
// and a Segment.  The return values are the distance along the velocity
// vector of the collision, the point on the polygon that collided, and a
// boolean that is true if there was a collision and false if not.
func CircleSegmentHit(c Circle, v Vector, s Segment) (float64, Point, bool) {
	planeHit, hit := circlePlaneHit(c, v, Plane(s.Line()))
	if !hit {
		return 0, Point{}, false
	}
	polyHit := s.NearestPoint(planeHit)

	r := Ray{Origin: polyHit, Direction: v.Inverse().Unit()}
	d, hit := r.SphereIntersection(Sphere(c))
	if !hit || d < 0 || d > v.Magnitude() {
		return 0, Point{}, false
	}
	return d, polyHit, true
}

// circlePlaneHit returns the point at which a circle traveling with a
// given velocity will intersect with a plane.  The second return value is true if
// there is an intersection, and false if not.
func circlePlaneHit(c Circle, v Vector, p Plane) (Point, bool) {
	r := Ray{Origin: c.Center, Direction: p.Normal.Inverse()}
	d, hit := r.PlaneIntersection(p)
	if !hit || d < 0 {
		return Point{}, false
	}

	// The circle is embedded in the plane.
	if d <= c.Radius {
		return c.Center.Plus(p.Normal.Inverse().ScaledBy(d)), true
	}

	r.Origin = c.Center.Plus(p.Normal.Inverse().ScaledBy(c.Radius))
	r.Direction = v.Unit()
	d, hit = r.PlaneIntersection(p)
	if !hit || d < 0 {
		return Point{}, false
	}
	return r.Origin.Plus(r.Direction.ScaledBy(d)), true
}
