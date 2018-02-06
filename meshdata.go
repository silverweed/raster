package main

var TRIANGLE_MESH = Mesh{
	[]Vertex{
		Vertex{Position: Vec3{0.5, -0.5, 0}},  /* vec3(0, 0, 1), vec2(1, 1) }*/
		Vertex{Position: Vec3{0, 0.5, 0}},     /* vec3(0, 0, 1), vec2(0, 1) }*/
		Vertex{Position: Vec3{-0.5, -0.5, 0}}, /* vec3(0, 0, 1), vec2(0, 0) }*/
	},
}
var QUAD_MESH = Mesh{
	[]Vertex{
		Vertex{Position: Vec3{0.5, 0.5, 0}},   /* vec3(0, 0, 1), vec2(1, 1) }*/
		Vertex{Position: Vec3{-0.5, 0.5, 0}},  /* vec3(0, 0, 1), vec2(0, 1) }*/
		Vertex{Position: Vec3{-0.5, -0.5, 0}}, /* vec3(0, 0, 1), vec2(0, 0) }*/
		Vertex{Position: Vec3{0.5, 0.5, 0}},   /* vec3(0, 0, 1), vec2(1, 1) }*/
		Vertex{Position: Vec3{-0.5, -0.5, 0}}, /* vec3(0, 0, 1), vec2(0, 0) }*/
		Vertex{Position: Vec3{0.5, -0.5, 0}},  /* vec3(0, 0, 1), vec2(1, 0) }*/
	},
}
var CUBE_MESH = Mesh{
	[]Vertex{
		Vertex{Position: Vec3{-0.5, -0.5, -0.5}}, /*,  [0.0,  0.0, -1.0],  [0.0, 0.0] }*/
		Vertex{Position: Vec3{0.5, 0.5, -0.5}},   /*,  [0.0,  0.0, -1.0],  [1.0, 1.0] }*/
		Vertex{Position: Vec3{0.5, -0.5, -0.5}},  /*,  [0.0,  0.0, -1.0],  [1.0, 0.0] }*/
		Vertex{Position: Vec3{0.5, 0.5, -0.5}},   /*,  [0.0,  0.0, -1.0],  [1.0, 1.0] }*/
		Vertex{Position: Vec3{-0.5, -0.5, -0.5}}, /*,  [0.0,  0.0, -1.0],  [0.0, 0.0] }*/
		Vertex{Position: Vec3{-0.5, 0.5, -0.5}},  /*,  [0.0,  0.0, -1.0],  [0.0, 1.0] }*/

		Vertex{Position: Vec3{-0.5, -0.5, 0.5}}, /*,  [0.0,  0.0, 1.0],   [0.0, 0.0] }*/
		Vertex{Position: Vec3{0.5, -0.5, 0.5}},  /*,  [0.0,  0.0, 1.0],   [1.0, 0.0] }*/
		Vertex{Position: Vec3{0.5, 0.5, 0.5}},   /*,  [0.0,  0.0, 1.0],   [1.0, 1.0] }*/
		Vertex{Position: Vec3{0.5, 0.5, 0.5}},   /*,  [0.0,  0.0, 1.0],   [1.0, 1.0] }*/
		Vertex{Position: Vec3{-0.5, 0.5, 0.5}},  /*,  [0.0,  0.0, 1.0],   [0.0, 1.0] }*/
		Vertex{Position: Vec3{-0.5, -0.5, 0.5}}, /*,  [0.0,  0.0, 1.0],   [0.0, 0.0] }*/

		Vertex{Position: Vec3{-0.5, 0.5, 0.5}},   /*,  [-1.0,  0.0, 0.0],  [1.0, 0.0] }*/
		Vertex{Position: Vec3{-0.5, 0.5, -0.5}},  /*,  [-1.0,  0.0, 0.0],  [1.0, 1.0] }*/
		Vertex{Position: Vec3{-0.5, -0.5, -0.5}}, /*,  [-1.0,  0.0, 0.0],  [0.0, 1.0] }*/
		Vertex{Position: Vec3{-0.5, -0.5, -0.5}}, /*,  [-1.0,  0.0, 0.0],  [0.0, 1.0] }*/
		Vertex{Position: Vec3{-0.5, -0.5, 0.5}},  /*,  [-1.0,  0.0, 0.0],  [0.0, 0.0] }*/
		Vertex{Position: Vec3{-0.5, 0.5, 0.5}},   /*,  [-1.0,  0.0, 0.0],  [1.0, 0.0] }*/

		Vertex{Position: Vec3{0.5, 0.5, 0.5}},   /*,  [1.0,  0.0,  0.0],  [1.0, 0.0] }*/
		Vertex{Position: Vec3{0.5, -0.5, -0.5}}, /*,  [1.0,  0.0,  0.0],  [0.0, 1.0] }*/
		Vertex{Position: Vec3{0.5, 0.5, -0.5}},  /*,  [1.0,  0.0,  0.0],  [1.0, 1.0] }*/
		Vertex{Position: Vec3{0.5, -0.5, -0.5}}, /*,  [1.0,  0.0,  0.0],  [0.0, 1.0] }*/
		Vertex{Position: Vec3{0.5, 0.5, 0.5}},   /*,  [1.0,  0.0,  0.0],  [1.0, 0.0] }*/
		Vertex{Position: Vec3{0.5, -0.5, 0.5}},  /*,  [1.0,  0.0,  0.0],  [0.0, 0.0] }*/

		Vertex{Position: Vec3{-0.5, -0.5, -0.5}}, /*,  [0.0, -1.0,  0.0],  [0.0, 1.0] }*/
		Vertex{Position: Vec3{0.5, -0.5, -0.5}},  /*,  [0.0, -1.0,  0.0],  [1.0, 1.0] }*/
		Vertex{Position: Vec3{0.5, -0.5, 0.5}},   /*,  [0.0, -1.0,  0.0],  [1.0, 0.0] }*/
		Vertex{Position: Vec3{0.5, -0.5, 0.5}},   /*,  [0.0, -1.0,  0.0],  [1.0, 0.0] }*/
		Vertex{Position: Vec3{-0.5, -0.5, 0.5}},  /*,  [0.0, -1.0,  0.0],  [0.0, 0.0] }*/
		Vertex{Position: Vec3{-0.5, -0.5, -0.5}}, /*,  [0.0, -1.0,  0.0],  [0.0, 1.0] }*/

		Vertex{Position: Vec3{0.5, 0.5, -0.5}},  /*,  [0.0,  1.0,  0.0],  [1.0, 1.0] }*/
		Vertex{Position: Vec3{-0.5, 0.5, -0.5}}, /*,  [0.0,  1.0,  0.0],  [0.0, 1.0] }*/
		Vertex{Position: Vec3{0.5, 0.5, 0.5}},   /*,  [0.0,  1.0,  0.0],  [1.0, 0.0] }*/
		Vertex{Position: Vec3{-0.5, 0.5, 0.5}},  /*,  [0.0,  1.0,  0.0],  [0.0, 0.0] }*/
		Vertex{Position: Vec3{0.5, 0.5, 0.5}},   /*,  [0.0,  1.0,  0.0],  [1.0, 0.0] }*/
		Vertex{Position: Vec3{-0.5, 0.5, -0.5}}, /*,  [0.0,  1.0,  0.0],  [0.0, 1.0] }*/
	},
}
