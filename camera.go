package main

type ProjType int

const (
	PROJ_ORTHO ProjType = iota
	PROJ_PERSP
)

type Camera struct {
	Position   Vec3
	Up         Vec3
	Right      Vec3
	Yaw        float32
	Pitch      float32
	Fov        float32
	Near       float32
	Far        float32
	Width      float32
	Height     float32
	Projection ProjType
}

func (c Camera) ViewMatrix() Mat44 {
	return Mat44{
		[4]float32{1, 0, 0, c.Position.X},
		[4]float32{0, 1, 0, c.Position.Y},
		[4]float32{0, 0, 1, c.Position.Z},
		[4]float32{0, 0, 0, 1},
	}
	/*front := Vec3{
		float32(math.Cos(float64(c.Yaw)) * math.Cos(float64(c.Pitch))),
		float32(math.Sin(float64(c.Pitch))),
		float32(math.Sin(float64(c.Yaw)) * math.Cos(float64(c.Pitch))),
	}.Normalized()
	return LookAt(c.Position, c.Position.Add(front), c.Up)*/
}

func LookAt(eye, target, up Vec3) Mat44 {
	look_dir := target.Add(eye.Neg()).Normalized()
	up_dir := up.Normalized()

	right_dir := Cross(look_dir, up_dir).Normalized()
	perp_up_dir := Cross(right_dir, look_dir)

	ret := Mat44Identity()
	ret[0] = [4]float32{right_dir.X, right_dir.Y, right_dir.Z, 0}
	ret[1] = [4]float32{perp_up_dir.X, perp_up_dir.Y, perp_up_dir.Z, 0}
	nl := look_dir.Neg()
	ret[2] = [4]float32{nl.X, nl.Y, nl.Z, 0}

	ret[0][3] = -Dot3(eye, right_dir)
	ret[1][3] = -Dot3(eye, perp_up_dir)
	ret[2][3] = Dot3(eye, look_dir)

	return ret
}

func (c Camera) ProjMatrix() Mat44 {
	if c.Projection == PROJ_PERSP {
		return Perspective(-c.Width/2, c.Width/2,
			-c.Height/2, c.Height/2, c.Near, c.Far)
	} else {
		return Orthographic(-c.Width/2, c.Width/2,
			-c.Height/2, c.Height/2, c.Near, c.Far)
	}
}
