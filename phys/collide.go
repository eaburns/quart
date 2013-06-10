package phys

// Basically an implementation of: http://www.paulnettle.com/pub/FluidStudios/CollisionDetection/Fluid_Studios_Generic_Collision_Detection_for_Games_Using_Ellipsoids.pdf

import (
	"math"

	. "github.com/eaburns/quart/geom"
)

const (
	// The factor of the height of a body that is considered to be the
	// bottom.  This is used to determine if the body is on the ground
	// or if it is in the air: if the bottom of the body collides, then the
	// body is "on the ground."
	//
	// In effect, this determines how steep of a hill a body can climb.
	// The steeper the hill, the higher up on the body that it will collide.
	// So, a higher value for bottomFactor means the body can climb
	// steeper hills.  For elliptical bodies, the steepness of climbable
	// hills is also determined by the height of the ellipse.
	bottomFactor = 0.05
)

// MoveEllipse moves an ellipse with a given velocity, handling collision with segments.
// The second return value is true if the ellipse collided with a segment beneath it,
// otherwise it is false.  This value can be used to decide if it is "on the ground."
func MoveEllipse(e Ellipse, v Vector, segs []Segment) (Ellipse, bool) {
	tr := Vector{}
	for i, r := range e.Radii {
		tr[i] = 1 / r
	}

	c := Circle{Center: e.Center.Times(tr), Radius: 1}
	v = v.Times(tr)
	trSegs := make([]Segment, len(segs))
	for i := range segs {
		trSegs[i][0] = segs[i][0].Times(tr)
		trSegs[i][1] = segs[i][1].Times(tr)
	}
	c2, onGround := MoveCircle(c, v, trSegs)
	return Ellipse{Center: c2.Center.Times(e.Radii), Radii: e.Radii}, onGround
}

// MoveCircle moves a circle with a given velocity, handling collision with segments.
// The second return value is true if the circle collided with a segment beneath it,
// otherwise it is false.  This value can be used to decide if it is "on the ground."
func MoveCircle(c Circle, v Vector, segs []Segment) (Circle, bool) {
	onGround := false
	for !v.NearZero() {
		mv := moveCircle1(c, v, segs)
		c.Center.Add(v.Unit().ScaledBy(mv.distance))
		low := c.Center[1] - c.Radius*(1-bottomFactor*2)
		hitGround := v[1] < 0 && mv.hit && mv.hitPoint[1] < low
		onGround = onGround || hitGround
		v = mv.newVelocity
	}
	return c, onGround
}

type move struct {
	distance    float64
	newVelocity Vector
	hit         bool
	hitPoint    Point
}

// moveCircle1 moves a circle along a vector until the first collision with a Segment.
func moveCircle1(c Circle, v Vector, segs []Segment) move {
	hitPt := Point{}
	dist := math.Inf(1)

	for _, s := range segs {
		if d, pt, hit := circleSegmentHit(c, v, s); hit && d < dist {
			dist = d
			hitPt = pt
		}
	}
	if math.IsInf(dist, 1) {
		return move{
			distance:    v.Magnitude(),
			newVelocity: Vector{},
		}
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

	return move{
		distance:    dist - Threshold,
		newVelocity: dest.Minus(hitPt),
		hit:         true,
		hitPoint:    hitPt,
	}
}

// circleSegmentHit returns information about the collision of a circle
// and a Segment.  The return values are the distance along the velocity
// vector of the collision, the point on the polygon that collided, and a
// boolean that is true if there was a collision and false if not.
func circleSegmentHit(c Circle, v Vector, s Segment) (float64, Point, bool) {
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
