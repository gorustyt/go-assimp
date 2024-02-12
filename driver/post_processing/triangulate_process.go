package post_processing

import (
	"context"
	"github.com/gorustyt/go-assimp/common"
	"github.com/gorustyt/go-assimp/common/logger"
	"github.com/gorustyt/go-assimp/core"
	"github.com/gorustyt/go-assimp/driver/base/iassimp"
	"math"
)

var (
	_ iassimp.PostProcessing = (*TriangulateProcess)(nil)
)

type TriangulateProcess struct {
}

func NewTriangulateProcess() *TriangulateProcess {
	return &TriangulateProcess{}
}

/**
 * @brief Encode the current triangle, and make sure it is recognized as a triangle.
 *
 * This method will rotate indices in tri if needed in order to avoid tri to be considered
 * part of the previous ngon. This method is to be used whenever you want to emit a real triangle,
 * and make sure it is seen as a triangle.
 *
 * @param tri Triangle to encode.
 */
type NGONEncoder struct {
	LastNGONFirstIndex int
}

func NewNGONEncoder() *NGONEncoder {
	return &NGONEncoder{LastNGONFirstIndex: -1}
}
func (n *NGONEncoder) ngonEncodeTriangle(tri *core.AiFace) {
	if len(tri.Indices) != 3 {
		panic("len(tri.Indices) != 3")
	}

	// Rotate indices in new triangle to avoid ngon encoding false ngons
	// Otherwise, the new triangle would be considered part of the previous NGON.
	if n.isConsideredSameAsLastNgon(tri) {
		tri.Indices[0], tri.Indices[2] = tri.Indices[2], tri.Indices[0]
		tri.Indices[1], tri.Indices[2] = tri.Indices[2], tri.Indices[1]
	}

	n.LastNGONFirstIndex = int(tri.Indices[0])
}

/**
 * @brief Encode a quad (2 triangles) in ngon encoding, and make sure they are seen as a single ngon.
 *
 * @param tri1 First quad triangle
 * @param tri2 Second quad triangle
 *
 * @pre Triangles must be properly fanned from the most appropriate vertex.
 */
func (n *NGONEncoder) ngonEncodeQuad(tri1 *core.AiFace, tri2 *core.AiFace) {
	if len(tri1.Indices) != 3 {
		panic("len(tri1.Indices) != 3")
	}
	if len(tri2.Indices) != 3 {
		panic("len(tri2.Indices) != 3")
	}
	if tri1.Indices[0] != tri2.Indices[0] {
		panic("tri1.Indices[0] != tri2.Indices[0]")
	}

	// If the selected fanning vertex is the same as the previously
	// emitted ngon, we use the opposite vertex which also happens to work
	// for tri-fanning a concave quad.
	// ref: https://github.com/assimp/assimp/pull/3695#issuecomment-805999760
	if n.isConsideredSameAsLastNgon(tri1) {
		// Right-rotate indices for tri1 (index 2 becomes the new fanning vertex)
		tri1.Indices[0], tri1.Indices[2] = tri1.Indices[2], tri1.Indices[0]
		tri1.Indices[1], tri1.Indices[2] = tri1.Indices[2], tri1.Indices[1]

		// Left-rotate indices for tri2 (index 2 becomes the new fanning vertex)
		tri2.Indices[1], tri2.Indices[2] = tri2.Indices[2], tri2.Indices[1]
		tri2.Indices[0], tri2.Indices[2] = tri2.Indices[2], tri2.Indices[0]

		if tri1.Indices[0] != tri2.Indices[0] {
			panic("tri1.Indices[0] != tri2.Indices[0]")
		}
	}

	n.LastNGONFirstIndex = int(tri1.Indices[0])
}

/**
 * @brief Check whether this triangle would be considered part of the lastly emitted ngon or not.
 *
 * @param tri Current triangle.
 * @return true If used as is, this triangle will be part of last ngon.
 * @return false If used as is, this triangle is not considered part of the last ngon.
 */
func (n *NGONEncoder) isConsideredSameAsLastNgon(tri *core.AiFace) bool {
	if len(tri.Indices) != 3 {
		panic("len(tri.Indices) != 3")
	}
	return int(tri.Indices[0]) == n.LastNGONFirstIndex
}
func (t *TriangulateProcess) IsActive(pFlags int) bool {
	return (iassimp.AiPostProcessSteps(pFlags) & iassimp.AiProcess_Triangulate) != 0
}

func (t *TriangulateProcess) Execute(pScene *core.AiScene) {
	logger.DebugF("TriangulateProcess begin")

	bHas := false
	for a := 0; a < len(pScene.Meshes); a++ {
		if pScene.Meshes[a] != nil {
			if t.TriangulateMesh(pScene.Meshes[a]) {
				bHas = true
			}
		}
	}
	if bHas {
		logger.DebugF("TriangulateProcess finished. All polygons have been triangulated.")
	} else {
		logger.DebugF("TriangulateProcess finished. There was nothing to be done.")
	}
}

func (t *TriangulateProcess) SetupProperties(ctx context.Context) {}
func (t *TriangulateProcess) TriangulateMesh(pMesh *core.AiMesh) bool {
	// Now we have aiMesh::mPrimitiveTypes, so this is only here for test cases
	if pMesh.PrimitiveTypes == 0 {
		bNeed := false

		for a := 0; a < len(pMesh.Faces); a++ {
			face := pMesh.Faces[a]

			if len(face.Indices) != 3 {
				bNeed = true
			}
		}
		if !bNeed {
			return false
		}

	} else if (pMesh.PrimitiveTypes & core.AiPrimitiveType_POLYGON) == 0 {
		return false
	}

	// Find out how many output faces we'll get

	numOut := 0
	max_out := 0
	get_normals := true
	for a := 0; a < len(pMesh.Faces); a++ {
		face := pMesh.Faces[a]
		if len(face.Indices) <= 4 {
			get_normals = false
		}
		if len(face.Indices) <= 3 {
			numOut++

		} else {
			numOut += len(face.Indices) - 2
			max_out = max(max_out, len(face.Indices))
		}
	}

	// Just another check whether aiMesh::mPrimitiveTypes is correct
	if numOut == len(pMesh.Faces) {
		panic("numOut == pMesh.mNumFaces")
	}

	var nor_out []*common.AiVector3D

	// if we don't have normals yet, but expect them to be a cheap side
	// product of triangulation anyway, allocate storage for them.
	if pMesh.Normals == nil && get_normals {
		// XXX need a mechanism to inform the GenVertexNormals process to treat these normals as preprocessed per-face normals
		//  nor_out = pMesh.mNormals = new aiVector3D[pMesh.mNumVertices];
	}

	// the output mesh will contain triangles, but no polys anymore
	pMesh.PrimitiveTypes |= core.AiPrimitiveType_TRIANGLE
	pMesh.PrimitiveTypes &= ^core.AiPrimitiveType_POLYGON

	// The mesh becomes NGON encoded now, during the triangulation process.
	pMesh.PrimitiveTypes |= core.AiPrimitiveType_NGONEncodingFlag

	out := make([]*core.AiFace, numOut)
	for i := range out {
		out[i] = core.NewAiFace()
	}
	curOut := 0
	/* temporary storage for vertices */
	temp_verts3d := make([]*common.AiVector3D, max_out+2)
	temp_verts := make([]*common.AiVector2D, max_out+2)
	ngonEncoder := NewNGONEncoder()

	// Apply vertex colors to represent the face winding?
	//if (pMesh.Colors[0]==nil){
	//	pMesh.Colors[0] = make([]*common.AiColor4D,len(pMesh.Vertices))
	//} else{
	//	new(pMesh.mColors[0])
	//}
	//
	//
	// clr := pMesh.Colors[0];
	verts := pMesh.Vertices

	// use std::unique_ptr to avoid slow std::vector<bool> specialiations
	done := make([]bool, max_out)
	for a := 0; a < len(pMesh.Faces); a++ {
		face := pMesh.Faces[a]

		idx := face.Indices
		num := len(face.Indices)
		ear := 0
		tmp, prev := num-1, num-1
		next := 0
		maxValue := float32(num)

		// Apply vertex colors to represent the face winding?
		//for  i := 0; i < len(face.Indices); i++ {
		//	 c := clr[idx[i]];
		//	c.R =float32 (i+1) / maxValue;
		//	c.B = 1. - c.R;
		//}

		last_face := curOut
		// if it's a simple point,line or triangle: just copy it
		if len(face.Indices) <= 3 {
			nface := out[curOut]
			curOut++
			nface.Indices = face.Indices
			face.Indices = nil

			// points and lines don't require ngon encoding (and are not supported either!)
			if len(nface.Indices) == 3 {
				ngonEncoder.ngonEncodeTriangle(nface)
			}

			continue
		} else if len(face.Indices) == 4 { // optimized code for quadrilaterals

			// quads can have at maximum one concave vertex. Determine
			// this vertex (if it exists) and start tri-fanning from
			// it.
			start_vertex := 0
			for i := 0; i < 4; i++ {
				v0 := verts[face.Indices[(i+3)%4]]
				v1 := verts[face.Indices[(i+2)%4]]
				v2 := verts[face.Indices[(i+1)%4]]

				v := verts[face.Indices[i]]

				left := v0.Sub(v)
				diag := v1.Sub(v)
				right := v2.Sub(v)

				left.Normalize()
				diag.Normalize()
				right.Normalize()
				angle := math.Acos(left.MulAiVector3D(diag)) + math.Acos(right.MulAiVector3D(diag))
				if angle > math.Pi {
					// this is the concave point
					start_vertex = i
					break
				}
			}

			temp := []uint32{face.Indices[0], face.Indices[1], face.Indices[2], face.Indices[3]}

			nface := out[curOut]
			curOut++
			nface.Indices = face.Indices[:3]
			nface.Indices[0] = temp[start_vertex]
			nface.Indices[1] = temp[(start_vertex+1)%4]
			nface.Indices[2] = temp[(start_vertex+2)%4]

			sface := out[curOut]
			curOut++
			sface.Indices = make([]uint32, 3)

			sface.Indices[0] = temp[start_vertex]
			sface.Indices[1] = temp[(start_vertex+2)%4]
			sface.Indices[2] = temp[(start_vertex+3)%4]

			// prevent double deletion of the indices field
			face.Indices = nil

			ngonEncoder.ngonEncodeQuad(nface, sface)

			continue
		} else {
			// A polygon with more than 3 vertices can be either concave or convex.
			// Usually everything we're getting is convex and we could easily
			// triangulate by tri-fanning. However, LightWave is probably the only
			// modeling suite to make extensive use of highly concave, monster polygons ...
			// so we need to apply the full 'ear cutting' algorithm to get it right.

			// REQUIREMENT: polygon is expected to be simple and *nearly* planar.
			// We project it onto a plane to get a 2d triangle.

			// Collect all vertices of of the polygon.
			for tmp := 0; tmp < int(maxValue); tmp++ {
				temp_verts3d[tmp] = verts[idx[tmp]]
			}
			front := common.SliceAiVector3DToFloatArr(temp_verts3d)
			// Get newell normal of the polygon. Store it for future use if it's a polygon-only mesh
			n := common.NewellNormal(int(maxValue), front[0:], front[1:], front[2:], 3, 3, 3)
			if nor_out != nil {
				for tmp := 0; tmp < int(maxValue); tmp++ {
					nor_out[idx[tmp]] = n
				}

			}

			// Select largest normal coordinate to ignore for projection
			ax := -n.X
			if n.X > 0 {
				ax = n.X
			}
			ay := -n.Y
			if n.Y > 0 {
				ay = n.Y
			}
			az := -n.Z
			if n.Z > 0 {
				az = n.Z
			}
			ac := 0
			bc := 1 /* no z coord. projection to xy */
			inv := n.Z
			if ax > ay {
				if ax > az { /* no x coord. projection to yz */
					ac = 1
					bc = 2
					inv = n.X
				}
			} else if ay > az { /* no y coord. projection to zy */
				ac = 2
				bc = 0
				inv = n.Y
			}

			// Swap projection axes to take the negated projection vector into account
			if inv < 0. {
				ac, bc = bc, ac
			}

			for tmp := 0; tmp < int(maxValue); tmp++ {
				temp_verts[tmp].X = verts[idx[tmp]].Index(ac)
				temp_verts[tmp].Y = verts[idx[tmp]].Index(bc)
				done[tmp] = false
			}
			//
			// FIXME: currently this is the slow O(kn) variant with a worst case
			// complexity of O(n^2) (I think). Can be done in O(n).
			for num > 3 {

				// Find the next ear of the polygon
				num_found := 0
				ear = next
				fn := func() {
					prev = ear
					ear = next
				}
				for ; ; fn() {

					// break after we looped two times without a positive match

					next = ear + 1
					for {
						if next >= int(maxValue) {
							next = 0
						}
						if !done[next] {
							break
						}
						next++
					}

					if next < ear {
						num_found++
						if num_found == 2 {
							break
						}
					}
					pnt1 := temp_verts[ear]
					pnt0 := temp_verts[prev]
					pnt2 := temp_verts[next]

					// Must be a convex point. Assuming ccw winding, it must be on the right of the line between p-1 and p+1.
					if common.OnLeftSideOfLine2D(pnt0, pnt2, pnt1) == 1 {
						continue
					}

					// Skip when three point is in a line
					left := pnt0.Sub(pnt1)
					right := pnt2.Sub(pnt1)

					left = left.Normalize()
					right = right.Normalize()
					mul := left.Mul(right)

					// if the angle is 0 or 180
					if math.Abs(float64(mul)-1.) < common.AiEpsilon || math.Abs(float64(mul)+1.) < common.AiEpsilon {
						// skip this ear
						logger.Warn("Skip a ear, due to its angle is near 0 or 180.")
						continue
					}

					// and no other point may be contained in this triangle
					for tmp := 0; tmp < int(maxValue); tmp++ {

						// We need to compare the actual values because it's possible that multiple indexes in
						// the polygon are referring to the same position. concave_polygon.obj is a sample
						//
						// FIXME: Use 'epsiloned' comparisons instead? Due to numeric inaccuracies in
						// PointInTriangle() I'm guessing that it's actually possible to construct
						// input data that would cause us to end up with no ears. The problem is,
						// which epsilon? If we chose a too large value, we'd get wrong results
						vtmp := temp_verts[tmp]
						if vtmp != pnt1 && vtmp != pnt2 && vtmp != pnt0 && common.PointInTriangle2D(pnt0, pnt1, pnt2, vtmp) {
							break
						}
					}
					if tmp != int(maxValue) {
						continue
					}

					// this vertex is an ear
					break
				}
				if num_found == 2 {

					// Due to the 'two ear theorem', every simple polygon with more than three points must
					// have 2 'ears'. Here's definitely something wrong ... but we don't give up yet.
					//

					// Instead we're continuing with the standard tri-fanning algorithm which we'd
					// use if we had only convex polygons. That's life.
					logger.ErrorF("Failed to triangulate polygon (no ear found). Probably not a simple polygon?")
					num = 0
					break
				}

				nface := out[curOut]
				curOut++
				if nface.Indices == nil {
					nface.Indices = make([]uint32, 3)
				}

				// setup indices for the new triangle ...
				nface.Indices[0] = uint32(prev)
				nface.Indices[1] = uint32(ear)
				nface.Indices[2] = uint32(next)

				// exclude the ear from most further processing
				done[ear] = true
				num--
			}
			if num > 0 {
				// We have three indices forming the last 'ear' remaining. Collect them.
				nface := out[curOut]
				curOut++
				if nface.Indices == nil {
					nface.Indices = make([]uint32, 3)
				}

				for tmp := 0; done[tmp]; tmp++ {
					nface.Indices[0] = uint32(tmp)

				}

				tmp++
				for ; done[tmp]; tmp++ {
					nface.Indices[1] = uint32(tmp)
				}

				tmp++
				for ; done[tmp]; tmp++ {
					nface.Indices[2] = uint32(tmp)
				}

			}
		}

		for f := last_face; f != curOut; {
			i := out[f].Indices

			i[0] = idx[i[0]]
			i[1] = idx[i[1]]
			i[2] = idx[i[2]]

			// IMPROVEMENT: Polygons are not supported yet by this ngon encoding + triangulation step.
			//              So we encode polygons as regular triangles. No way to reconstruct the original
			//              polygon in this case.
			ngonEncoder.ngonEncodeTriangle(out[f])
			f++
		}

		face.Indices = nil
	}
	// ... and store the new ones
	pMesh.Faces = out
	return true
}
