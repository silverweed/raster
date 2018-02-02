package main

import "math"

type Vec2 struct {
	X float32
	Y float32
}

type Vec3 struct {
	X float32
	Y float32
	Z float32
}

type Vec4 struct {
	X float32
	Y float32
	Z float32
	W float32
}

type Vertex struct {
	Position Vec3
}

type Mesh struct {
	Vertices []Vertex
}

type Object struct {
	Mesh
	Transform Transform
}

type Transform struct {
	Position Vec3
	Rotation Quat
	Scale    Vec3
}

// q = a + bi + cj + dk
type Quat struct {
	A float32
	B float32
	C float32
	D float32
}

///// Functions
func MakeObject(m Mesh) Object {
	return Object{
		m,
		Transform{
			Position: Vec3{0, 0, 0},
			Rotation: QuatIdentity(),
			Scale:    Vec3{1, 1, 1},
		},
	}
}

func QuatIdentity() Quat {
	return Quat{1, 0, 0, 0}
}

func Vec4to3(v Vec4) Vec3 {
	return Vec3{v.X, v.Y, v.Z}
}

func Vec3to2(v Vec3) Vec2 {
	return Vec2{v.X, v.Y}
}

func (v Vec3) Normalized() Vec3 {
	norm := float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
	return Vec3{v.X / norm, v.Y / norm, v.Z / norm}
}

func (v Vec3) Add(v2 Vec3) Vec3 {
	return Vec3{v.X + v2.X, v.Y + v2.Y, v.Z + v2.Z}
}

func (v Vec3) Neg() Vec3 {
	return Vec3{-v.X, -v.Y, -v.Z}
}

func Cross(v1, v2 Vec3) Vec3 {
	return Vec3{v1.Y*v2.Z - v2.Y*v1.Z,
		v1.Z*v2.X - v2.Z*v1.X,
		v1.X*v2.Y - v2.X*v1.Y}
}

func Dot(v1, v2 Vec3) float32 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

func Mat44Identity() Mat44 {
	return Mat44{
		[4]float32{1, 0, 0, 0},
		[4]float32{0, 1, 0, 0},
		[4]float32{0, 0, 1, 0},
		[4]float32{0, 0, 0, 1},
	}
}

func Mat44Zero() Mat44 {
	return Mat44{
		[4]float32{0, 0, 0, 0},
		[4]float32{0, 0, 0, 0},
		[4]float32{0, 0, 0, 0},
		[4]float32{0, 0, 0, 0},
	}
}

func (m Mat44) Det() float32 {
	return m[0][3]*m[1][2]*m[2][1]*m[3][0] - m[0][2]*m[1][3]*m[2][1]*m[3][0] -
		m[0][3]*m[1][1]*m[2][2]*m[3][0] + m[0][1]*m[1][3]*m[2][2]*m[3][0] +
		m[0][2]*m[1][1]*m[2][3]*m[3][0] - m[0][1]*m[1][2]*m[2][3]*m[3][0] -
		m[0][3]*m[1][2]*m[2][0]*m[3][1] + m[0][2]*m[1][3]*m[2][0]*m[3][1] +
		m[0][3]*m[1][0]*m[2][2]*m[3][1] - m[0][0]*m[1][3]*m[2][2]*m[3][1] -
		m[0][2]*m[1][0]*m[2][3]*m[3][1] + m[0][0]*m[1][2]*m[2][3]*m[3][1] +
		m[0][3]*m[1][1]*m[2][0]*m[3][2] - m[0][1]*m[1][3]*m[2][0]*m[3][2] -
		m[0][3]*m[1][0]*m[2][1]*m[3][2] + m[0][0]*m[1][3]*m[2][1]*m[3][2] +
		m[0][1]*m[1][0]*m[2][3]*m[3][2] - m[0][0]*m[1][1]*m[2][3]*m[3][2] -
		m[0][2]*m[1][1]*m[2][0]*m[3][3] + m[0][1]*m[1][2]*m[2][0]*m[3][3] +
		m[0][2]*m[1][0]*m[2][1]*m[3][3] - m[0][0]*m[1][2]*m[2][1]*m[3][3] -
		m[0][1]*m[1][0]*m[2][2]*m[3][3] + m[0][0]*m[1][1]*m[2][2]*m[3][3]
}

func (m Mat44) Inverse() Mat44 {
	d := 1 / m.Det()
	return Mat44{
		[4]float32{(m[1][1]*m[2][2]*m[3][3] + m[1][2]*m[2][3]*m[3][1] + m[1][3]*m[2][1]*m[3][2] -
			m[1][1]*m[2][3]*m[3][2] - m[1][2]*m[2][1]*m[3][3] - m[1][3]*m[2][2]*m[3][1]) * d,
			(m[0][1]*m[2][3]*m[3][2] + m[0][2]*m[2][1]*m[3][3] + m[0][3]*m[2][2]*m[3][1] -
				m[0][1]*m[2][2]*m[3][3] - m[0][2]*m[2][3]*m[3][1] - m[0][3]*m[2][1]*m[3][2]) * d,
			(m[0][1]*m[1][2]*m[3][3] + m[0][2]*m[1][3]*m[3][1] + m[0][3]*m[1][1]*m[3][2] -
				m[0][1]*m[1][3]*m[3][2] - m[0][2]*m[1][1]*m[3][3] - m[0][3]*m[1][2]*m[3][1]) * d,
			(m[0][1]*m[1][3]*m[2][2] + m[0][2]*m[1][1]*m[2][3] + m[0][3]*m[1][2]*m[2][1] -
				m[0][1]*m[1][2]*m[2][3] - m[0][2]*m[1][3]*m[2][1] - m[0][3]*m[1][1]*m[2][2]) * d},
		[4]float32{(m[1][0]*m[2][3]*m[3][2] + m[1][2]*m[2][0]*m[3][3] + m[1][3]*m[2][2]*m[3][0] -
			m[1][0]*m[2][2]*m[3][3] - m[1][2]*m[2][3]*m[3][0] - m[1][3]*m[2][0]*m[3][2]) * d,
			(m[0][0]*m[2][2]*m[3][3] + m[0][2]*m[2][3]*m[3][0] + m[0][3]*m[2][0]*m[3][2] -
				m[0][0]*m[2][3]*m[3][2] - m[0][2]*m[2][0]*m[3][3] - m[0][3]*m[2][2]*m[3][0]) * d,
			(m[0][0]*m[1][3]*m[3][2] + m[0][2]*m[1][0]*m[3][3] + m[0][3]*m[1][2]*m[3][0] -
				m[0][0]*m[1][2]*m[3][3] - m[0][2]*m[1][3]*m[3][0] - m[0][3]*m[1][0]*m[3][2]) * d,
			(m[0][0]*m[1][2]*m[2][3] + m[0][2]*m[1][3]*m[2][0] + m[0][3]*m[1][0]*m[2][2] -
				m[0][0]*m[1][3]*m[2][2] - m[0][2]*m[1][0]*m[2][3] - m[0][3]*m[1][2]*m[2][0]) * d},
		[4]float32{(m[1][0]*m[2][1]*m[3][3] + m[1][1]*m[2][3]*m[3][0] + m[1][3]*m[2][0]*m[3][1] -
			m[1][0]*m[2][3]*m[3][1] - m[1][1]*m[2][0]*m[3][3] - m[1][3]*m[2][1]*m[3][0]) * d,
			(m[0][0]*m[2][3]*m[3][1] + m[0][1]*m[2][0]*m[3][3] + m[0][3]*m[2][1]*m[3][0] -
				m[0][0]*m[2][1]*m[3][3] - m[0][1]*m[2][3]*m[3][0] - m[0][3]*m[2][0]*m[3][1]) * d,
			(m[0][0]*m[1][1]*m[3][3] + m[0][1]*m[1][3]*m[3][0] + m[0][3]*m[1][0]*m[3][1] -
				m[0][0]*m[1][3]*m[3][1] - m[0][1]*m[1][0]*m[3][3] - m[0][3]*m[1][1]*m[3][0]) * d,
			(m[0][0]*m[1][3]*m[2][1] + m[0][1]*m[1][0]*m[2][3] + m[0][3]*m[1][1]*m[2][0] -
				m[0][0]*m[1][1]*m[2][3] - m[0][1]*m[1][3]*m[2][0] - m[0][3]*m[1][0]*m[2][1]) * d},
		[4]float32{(m[1][0]*m[2][2]*m[3][1] + m[1][1]*m[2][0]*m[3][2] + m[1][2]*m[2][1]*m[3][0] -
			m[1][0]*m[2][1]*m[3][2] - m[1][1]*m[2][2]*m[3][0] - m[1][2]*m[2][0]*m[3][1]) * d,
			(m[0][0]*m[2][1]*m[3][2] + m[0][1]*m[2][2]*m[3][0] + m[0][2]*m[2][0]*m[3][1] -
				m[0][0]*m[2][2]*m[3][1] - m[0][1]*m[2][0]*m[3][2] - m[0][2]*m[2][1]*m[3][0]) * d,
			(m[0][0]*m[1][2]*m[3][1] + m[0][1]*m[1][0]*m[3][2] + m[0][2]*m[1][1]*m[3][0] -
				m[0][0]*m[1][1]*m[3][2] - m[0][1]*m[1][2]*m[3][0] - m[0][2]*m[1][0]*m[3][1]) * d,
			(m[0][0]*m[1][1]*m[2][2] + m[0][1]*m[1][2]*m[2][0] + m[0][2]*m[1][0]*m[2][1] -
				m[0][0]*m[1][2]*m[2][1] - m[0][1]*m[1][0]*m[2][2] - m[0][2]*m[1][1]*m[2][0]) * d},
	}
}

// inputs in degrees
func FromEuler(yaw, pitch, roll float32) Quat {
	cy := float32(math.Cos(float64(yaw * 0.5)))
	sy := float32(math.Sin(float64(yaw * 0.5)))
	cr := float32(math.Cos(float64(roll * 0.5)))
	sr := float32(math.Sin(float64(roll * 0.5)))
	cp := float32(math.Cos(float64(pitch * 0.5)))
	sp := float32(math.Sin(float64(pitch * 0.5)))

	return Quat{
		cy*cr*cp + sy*sr*sp,
		cy*sr*cp - sy*cr*sp,
		cy*cr*sp + sy*sr*cp,
		sy*cr*cp - cy*sr*sp,
	}
}
