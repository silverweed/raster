package main

type Mat44 [4][4]float32

/* [ROT11] [ROT12] [ROT13] tx
 * [ROT21] [ROT22] [ROT23] ty
 * [ROT31] [ROT32] [ROT33] tz
 * 0       0       0       1
 */
func MakeModelMatrix(t Transform) Mat44 {
	q := t.Rotation
	return Mat44{
		[4]float32{1 - 2*q.C*q.C - 2*q.D*q.D, 2*q.B*q.C - 2*q.D*q.A, 2*q.B*q.D + 2*q.C*q.A, t.Position.X},
		[4]float32{2*q.B*q.C + 2*q.D*q.A, 1 - 2*q.B*q.B - 2*q.D*q.D, 2*q.C*q.D - 2*q.B*q.A, t.Position.Y},
		[4]float32{2*q.B*q.D - 2*q.C*q.A, 2*q.C*q.D + 2*q.B*q.A, 1 - 2*q.B*q.B - 2*q.C*q.C, t.Position.Z},
		[4]float32{0, 0, 0, 1},
	}
}

func MulMatVec4(mat Mat44, vec Vec4) Vec4 {
	return Vec4{
		mat[0][0]*vec.X + mat[0][1]*vec.Y + mat[0][2]*vec.Z + mat[0][3]*vec.W,
		mat[1][0]*vec.X + mat[1][1]*vec.Y + mat[1][2]*vec.Z + mat[1][3]*vec.W,
		mat[2][0]*vec.X + mat[2][1]*vec.Y + mat[2][2]*vec.Z + mat[2][3]*vec.W,
		mat[3][0]*vec.X + mat[3][1]*vec.Y + mat[3][2]*vec.Z + mat[3][3]*vec.W,
	}
}

func LocalToWorld(obj Object) Mesh {
	// Build the model matrix
	modelMat := MakeModelMatrix(obj.Transform)
	worldVertices := make([]Vertex, len(obj.Vertices))
	for i, v := range obj.Vertices {
		p := v.Position
		worldVertices[i] = Vertex{
			Position: Vec4to3(MulMatVec4(modelMat, Vec4{p.X, p.Y, p.Z, 1.0})),
			Color:    v.Color,
		}
	}
	return Mesh{worldVertices}
}

func WorldToView(obj Mesh, viewMat Mat44) Mesh {
	viewVertices := make([]Vertex, len(obj.Vertices))
	for i, v := range obj.Vertices {
		p := v.Position
		viewVertices[i] = Vertex{
			Position: Vec4to3(MulMatVec4(viewMat, Vec4{p.X, p.Y, p.Z, 1.0})),
			Color:    v.Color,
		}
	}
	return Mesh{viewVertices}
}

func ViewToClip(obj Mesh, projMat Mat44) Mesh {
	projVertices := make([]Vertex, len(obj.Vertices))
	for i, v := range obj.Vertices {
		p := v.Position
		cp := MulMatVec4(projMat, Vec4{p.X, p.Y, p.Z, 1.0})
		projVertices[i] = Vertex{
			Position: Vec3{
				cp.X / cp.W,
				cp.Y / cp.W,
				cp.Z / cp.W,
			},
			Color: v.Color,
		}
		// TODO perform frustum clipping
	}
	return Mesh{projVertices}
}

func Orthographic(left, right, bottom, top, near, far float32) Mat44 {
	if right-left == 0 {
		panic("Orthographic: right == left!")
	}
	if top-bottom == 0 {
		panic("Orthographic: top == bottom!")
	}
	if near-far == 0 {
		panic("Orthographic: near == far!")
	}
	ret := Mat44Zero()

	ret[0][0] = 2 / (right - left)
	ret[0][3] = -(right + left) / (right - left)
	ret[1][1] = 2 / (top - bottom)
	ret[1][3] = -(top + bottom) / (top - bottom)
	ret[2][2] = -2 / (far - near)
	ret[2][3] = -(far + near) / (far - near)
	ret[3][3] = 1

	return ret
}

func Perspective(left, right, bottom, top, near, far float32) Mat44 {
	if right-left == 0 {
		panic("Orthographic: right == left!")
	}
	if top-bottom == 0 {
		panic("Orthographic: top == bottom!")
	}
	if near-far == 0 {
		panic("Orthographic: near == far!")
	}
	ret := Mat44Zero()

	ret[0][0] = (2 * near) / (right - left)
	ret[0][2] = (right + left) / (right - left)
	ret[1][1] = (2 * near) / (top - bottom)
	ret[1][2] = (top + bottom) / (top - bottom)
	ret[2][2] = -(far + near) / (far - near)
	ret[2][3] = -(2 * far * near) / (far - near)
	ret[3][2] = -1

	return ret
}
