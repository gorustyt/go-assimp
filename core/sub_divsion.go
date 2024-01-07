package core

import (
	"github.com/gorustyt/go-assimp/common"
	"github.com/gorustyt/go-assimp/common/logger"
	"time"
)

type Algorithm int

const (
	CATMULL_CLARKE Algorithm = 0x1
)

type SubDivision interface {
	// ---------------------------------------------------------------
	/** Create a subdivider of a specific type
	 *
	 *  @param algo Algorithm to be used for subdivision
	 *  @return Subdivider instance. */
	SubdivideOneMesh(mesh *AiMesh,
		num int,
		discard_input bool) *AiMesh
	Subdivide(smesh []*AiMesh,
		out []*AiMesh,
		num int,
		discard_input bool)
}

// ---------------------------------------------------------------
/** Create a subdivider of a specific type
 *
 *  @param algo Algorithm to be used for subdivision
 *  @return Subdivider instance. */

func NewSubDivision(algo Algorithm) SubDivision {
	switch algo {
	case CATMULL_CLARKE:
		return newCatmullClarkSubdivider()
	}
	return nil
}

// ---------------------------------------------------------------------------
/** Intermediate description of an edge between two corners of a polygon*/
// ---------------------------------------------------------------------------
type Edge struct {
	edge_point, midpoint *Vertex
	ref                  int
}

func NewEdge() *Edge {
	return &Edge{
		edge_point: NewVertex(),
		midpoint:   NewVertex(),
	}
}

type CatmullClarkSubdivider struct {
	EdgeMap    map[int]Edge
	UIntVector []int
}

// ---------------------------------------------------------------
/** Subdivide a mesh using the selected algorithm
 *
 *  @param mesh First mesh to be subdivided. Must be in verbose
 *    format.
 *  @param out Receives the output mesh, allocated by me.
 *  @param num Number of subdivisions to perform.
 *  @param discard_input If true is passed, the input mesh is
 *    deleted after the subdivision is complete. This can
 *    improve performance because it allows the optimization
 *    to reuse the existing mesh for intermediate results.
 *  @pre out!=mesh*/

func (div *CatmullClarkSubdivider) SubdivideOneMesh(mesh *AiMesh, num int, discard_input bool) *AiMesh {
	out := make([]*AiMesh, 1)
	div.Subdivide([]*AiMesh{mesh}, out, num, discard_input)
	return out[0]
}

// ---------------------------------------------------------------
/** Subdivide multiple meshes using the selected algorithm. This
 *  avoids erroneous smoothing on objects consisting of multiple
 *  per-material meshes. Usually, most 3d modellers smooth on a
 *  per-object base, regardless the materials assigned to the
 *  meshes.
 *
 *  @param smesh Array of meshes to be subdivided. Must be in
 *    verbose format.
 *  @param nmesh Number of meshes in smesh.
 *  @param out Receives the output meshes. The array must be
 *    sufficiently large (at least @c nmesh elements) and may not
 *    overlap the input array. Output meshes map one-to-one to
 *    their corresponding input meshes. The meshes are allocated
 *    by the function.
 *  @param discard_input If true is passed, input meshes are
 *    deleted after the subdivision is complete. This can
 *    improve performance because it allows the optimization
 *    of reusing existing meshes for intermediate results.
 *  @param num Number of subdivisions to perform.
 *  @pre nmesh != 0, smesh and out may not overlap*/

func (div *CatmullClarkSubdivider) Subdivide(smesh []*AiMesh,
	out []*AiMesh,
	num int,
	discard_input bool) {
	if smesh == nil || out == nil || len(smesh) < len(out) {
		logger.Fatal("smesh == nil || out == nil || smesh < out || smesh+nmesh > out+nmesh")
	}
	if num == 0 {
		// No subdivision at all. Need to copy all the meshes .. argh.
		if discard_input {
			for s := 0; s < len(smesh); s++ {
				out[s] = smesh[s]
				smesh[s] = nil
			}
		} else {
			for s := 0; s < len(smesh); s++ {
				out[s] = smesh[s].Clone()
			}
		}
		return
	}
	var inmeshes []*AiMesh
	var outmeshes []*AiMesh
	var maptbl []int
	// Remove pure line and point meshes from the working set to reduce the
	// number of edge cases the subdivider is forced to deal with. Line and
	// point meshes are simply passed through.
	for s := 0; s < len(smesh); s++ {
		i := smesh[s]
		// FIX - mPrimitiveTypes might not yet be initialized
		if i.PrimitiveTypes != 0 && (i.PrimitiveTypes&(AiPrimitiveType_LINE|AiPrimitiveType_POINT)) == i.PrimitiveTypes {
			logger.Debug("Catmull-Clark Subdivider: Skipping pure line/point mesh")

			if discard_input {
				out[s] = i
				smesh[s] = nil
			} else {
				out[s] = i.Clone()
			}
			continue
		}

		outmeshes = append(outmeshes, nil)
		inmeshes = append(inmeshes, i)
		maptbl = append(maptbl, s)
	}

	// Do the actual subdivision on the preallocated storage. InternSubdivide
	// *always* assumes that enough storage is available, it does not bother
	// checking any ranges.
	common.AiAssert(len(inmeshes) == len(outmeshes), "len(inmeshes) == len(outmeshes)")
	common.AiAssert(len(inmeshes) == len(maptbl), "len(inmeshes) = len(maptbl)")
	if len(inmeshes) == 0 {
		logger.Warn("Catmull-Clark Subdivider: Pure point/line scene, I can't do anything")
		return
	}
	div.InternSubdivide(inmeshes, outmeshes, num)
	for i := 0; i < len(maptbl); i++ {
		common.AiAssert(nil != outmeshes[i], "nil == outmeshes[i]")
		out[maptbl[i]] = outmeshes[i]
	}

	if discard_input {
		for s := 0; s < len(smesh); s++ {
			smesh[s] = nil
		}
	}

}

// ------------------------------------------------------------------------------------------------
// Note - this is an implementation of the standard (recursive) Cm-Cl algorithm without further
// optimizations (except we're using some nice LUTs). A description of the algorithm can be found
// here: http://en.wikipedia.org/wiki/Catmull-Clark_subdivision_surface
//
// The code is mostly O(n), however parts are O(nlogn) which is therefore the algorithm's
// expected total runtime complexity. The implementation is able to work in-place on the same
// mesh arrays. Calling #InternSubdivide() directly is not encouraged. The code can operate
// in-place unless 'smesh' and 'out' are equal (no strange overlaps or reorderings).
// Previous data is replaced/deleted then.
// ------------------------------------------------------------------------------------------------

func (div *CatmullClarkSubdivider) InternSubdivide(smesh []*AiMesh,
	out []*AiMesh,
	num int) {
	moffsets := make([]*common.Pair[int, int], len(smesh))
	totfaces := 0
	totvert := 0
	if num == 0 {
		return
	}
	spatial := NewSpatialSort()
	for t := 0; t < len(smesh); t++ {
		mesh := smesh[t]
		spatial.Append(mesh.Vertices, false)
		moffsets[t] = common.NewPair(totfaces, totvert)
		totfaces += len(mesh.Faces)
		totvert += len(mesh.Vertices)
	}
	spatial.Finalize()
	maptbl, num_unique := spatial.GenerateMappingTable(ComputePositionEpsilon(smesh))
	FLATTEN_VERTEX_IDX := func(mesh_idx, vert_idx int) int { return moffsets[mesh_idx].Second + vert_idx }
	FLATTEN_FACE_IDX := func(mesh_idx, face_idx int) int { return moffsets[mesh_idx].First + face_idx }

	// ---------------------------------------------------------------------
	// 1. Compute the centroid point for all faces
	// ---------------------------------------------------------------------
	centroids := make([]*Vertex, totfaces)
	for i := range centroids {
		centroids[i] = NewVertex()
	}
	nfacesout := 0
	t := 0
	n := 0
	for ; t < len(smesh); t++ {
		mesh := smesh[t]
		for i := 0; i < len(mesh.Faces); i++ {
			face := mesh.Faces[i]
			c := centroids[n]

			for a := 0; a < len(face.Indices); a++ {
				c = c.Add(c, NewVertexFromAiMesh(mesh, int(face.Indices[a])))
			}

			centroids[n] = c.Div(c, float32(len(face.Indices)))
			nfacesout += len(face.Indices)
			n++
		}
	}

	// we want edges to go away before the recursive calls so begin a new scope
	edges := map[common.Pair[int, int]]*Edge{}
	makeEdgeHash := func(v1, v2 int) common.Pair[int, int] {
		if v1 > v2 {
			v1, v2 = v2, v1
		}
		return *common.NewPair(v1, v2)
	}
	// ---------------------------------------------------------------------
	// 2. Set each edge point to be the average of all neighbouring
	// face points and original points. Every edge exists twice
	// if there is a neighboring face.
	// ---------------------------------------------------------------------
	for t := 0; t < len(smesh); t++ {
		mesh := smesh[t]

		for i := 0; i < len(mesh.Faces); i++ {
			face := mesh.Faces[i]

			for p := 0; p < len(face.Indices); p++ {
				tmp := p + 1
				if p == len(face.Indices)-1 {
					tmp = 0
				}
				id := []int{
					int(face.Indices[p]),
					int(face.Indices[tmp]),
				}
				mp := []int{
					maptbl[FLATTEN_VERTEX_IDX(t, id[0])],
					maptbl[FLATTEN_VERTEX_IDX(t, id[1])],
				}

				key := makeEdgeHash(mp[0], mp[1])
				e, ok := edges[key]
				if !ok {
					e = NewEdge()
					edges[key] = e
				}
				e.ref++
				if e.ref <= 2 {
					if e.ref == 1 { // original points (end points) - add only once
						v1 := NewVertexFromAiMesh(mesh, id[0])
						e.midpoint = v1.Add(v1, NewVertexFromAiMesh(mesh, id[1]))
						e.edge_point = e.midpoint
						e.midpoint = e.midpoint.Mul(e.midpoint, 0.5)
					}
					e.edge_point = e.edge_point.Add(e.edge_point, centroids[FLATTEN_FACE_IDX(t, i)])
				}
			}
		}
	}

	// ---------------------------------------------------------------------
	// 3. Normalize edge points
	// ---------------------------------------------------------------------

	bad_cnt := 0
	for i, v := range edges {
		if v.ref < 2 {
			logger.WarnF("fond bad cnt first:%v second:%v", i.First, i.Second)
			common.AiAssert(v.ref != 0)
			bad_cnt++
		}
		v.edge_point = v.edge_point.Mul(v.edge_point, 1./(float32(v.ref)+2.))
	}

	if bad_cnt != 0 {
		// Report the number of bad edges. bad edges are referenced by less than two
		// faces in the mesh. They occur at outer model boundaries in non-closed
		// shapes.
		logger.DebugF("Catmull-Clark Subdivider: got %v bad edges touching only one face (totally %v edges).", bad_cnt,
			len(edges))
	}

	// ---------------------------------------------------------------------
	// 4. Compute a vertex-face adjacency table. We can't reuse the code
	// from VertexTriangleAdjacency because we need the table for multiple
	// meshes and out vertex indices need to be mapped to distinct values
	// first.
	// ---------------------------------------------------------------------
	faceadjac := make([]int, nfacesout)
	cntadjfac := make([]int, len(maptbl))
	ofsadjvec := make([]int, len(maptbl)+1)

	for t := 0; t < len(smesh); t++ {
		minp := smesh[t]
		for i := 0; i < len(minp.Faces); i++ {

			f := minp.Faces[i]
			for n := 0; n < len(f.Indices); n++ {
				cntadjfac[maptbl[FLATTEN_VERTEX_IDX(t, int(f.Indices[n]))]]++
			}
		}
	}
	cur := 0
	for i := 0; i < len(cntadjfac); i++ {
		ofsadjvec[i+1] = cur
		cur += cntadjfac[i]
	}
	for t := 0; t < len(smesh); t++ {
		minp := smesh[t]
		for i := 0; i < len(minp.Faces); i++ {

			f := minp.Faces[i]
			for n := 0; n < len(f.Indices); n++ {

				faceadjac[ofsadjvec[1+maptbl[FLATTEN_VERTEX_IDX(t, int(f.Indices[n]))]]] = FLATTEN_FACE_IDX(t, i)
				ofsadjvec[1+maptbl[FLATTEN_VERTEX_IDX(t, int(f.Indices[n]))]]++
			}
		}
	}

	// check the other way round for consistency

	for t := 0; t < len(ofsadjvec)-1; t++ {
		for m := 0; m < cntadjfac[t]; m++ {
			fidx := faceadjac[ofsadjvec[t]+m]
			common.AiAssert(fidx < totfaces, "fidx < totfaces")
			for n := 1; n < len(smesh); n++ {

				if moffsets[n].First > fidx {
					n--
					msh := smesh[n]
					f := msh.Faces[fidx-moffsets[n].First]

					haveit := false
					for i := 0; i < len(f.Indices); i++ {
						if maptbl[FLATTEN_VERTEX_IDX(n, int(f.Indices[i]))] == t {
							haveit = true
							break
						}
					}
					common.AiAssert(haveit)
					if !haveit {
						logger.Debug("Catmull-Clark Subdivider: Index not used")
					}
					break
				}
			}
		}
	}

	GET_ADJACENT_FACES_AND_CNT := func(vidx int) (fstartout []int, numout int) {
		fstartout = faceadjac[ofsadjvec[vidx]:]
		numout = cntadjfac[vidx]
		return
	}

	new_points := make([]*common.Pair[bool, *Vertex], num_unique)
	for i := range new_points {
		new_points[i] = common.NewPair(false, NewVertex())
	}
	// ---------------------------------------------------------------------
	// 5. Spawn a quad from each face point to the corresponding edge points
	// the original points being the fourth quad points.
	// ---------------------------------------------------------------------
	for t := 0; t < len(smesh); t++ {
		minp := smesh[t]
		out[t] = NewAiMesh()
		mout := out[t]

		var numFaces int
		for a := 0; a < len(minp.Faces); a++ {
			numFaces += len(minp.Faces[a].Indices)
		}

		// We need random access to the old face buffer, so reuse is not possible.
		mout.Faces = make([]*AiFace, numFaces)
		for i := range mout.Faces {
			mout.Faces[i] = NewAiFace()
		}
		numVertices := numFaces * 4
		mout.Vertices = make([]*common.AiVector3D, numVertices)

		// quads only, keep material index
		mout.PrimitiveTypes = AiPrimitiveType_POLYGON
		mout.MaterialIndex = minp.MaterialIndex

		if minp.HasNormals() {
			mout.Normals = make([]*common.AiVector3D, numVertices)
		}

		if minp.HasTangentsAndBitangents() {
			mout.Tangents = make([]*common.AiVector3D, numVertices)
			mout.Bitangents = make([]*common.AiVector3D, numVertices)
		}

		for i := 0; minp.HasTextureCoords(uint32(i)); i++ {
			mout.TextureCoords[i] = make([]*common.AiVector3D, numVertices)
			mout.NumUVComponents[i] = minp.NumUVComponents[i]
		}

		for i := 0; minp.HasVertexColors(i); i++ {
			mout.Colors[i] = make([]*common.AiColor4D, numVertices)
		}

		numVertices = len(mout.Faces) << 2
		i := 0
		v := 0
		n := 0
		now1 := time.Now()
		for ; i < len(minp.Faces); i++ {

			face := minp.Faces[i]
			for a := 0; a < len(face.Indices); a++ {

				// Get a clean new face.
				faceOut := mout.Faces[n]
				n++
				faceOut.Indices = make([]uint32, 4)

				// Spawn a new quadrilateral (ccw winding) for this original point between:
				// a) face centroid

				faceOut.Indices[0] = uint32(v)
				centroids[FLATTEN_FACE_IDX(t, i)].SortBack(mout, v)
				v++
				// b) adjacent edge on the left, seen from the centroid
				tmp := a + 1
				if a == len(face.Indices)-1 {
					tmp = 0
				}

				key := makeEdgeHash(maptbl[FLATTEN_VERTEX_IDX(t, int(face.Indices[a]))],
					maptbl[FLATTEN_VERTEX_IDX(t, int(face.Indices[tmp]))])
				e0, ok := edges[key] // fixme: replace with mod face.mNumIndices?
				if !ok {
					e0 = NewEdge()
					edges[key] = e0
				}
				tmp = a - 1
				if a == 0 {
					tmp = len(face.Indices) - 1
				}
				// c) adjacent edge on the right, seen from the centroid
				key = makeEdgeHash(maptbl[FLATTEN_VERTEX_IDX(t, int(face.Indices[a]))], maptbl[FLATTEN_VERTEX_IDX(t, int(face.Indices[tmp]))])
				e1, ok := edges[key]
				if !ok {
					e1 = NewEdge()
					edges[key] = e1
				}
				// fixme: replace with mod face.mNumIndices?

				faceOut.Indices[3] = uint32(v)
				e0.edge_point.SortBack(mout, v)
				v++
				faceOut.Indices[1] = uint32(v)
				e1.edge_point.SortBack(mout, v)
				v++
				// d= original point P with distinct index i
				// F := 0
				// R := 0
				// n := 0
				// for each face f containing i
				//    F := F+ centroid of f
				//    R := R+ midpoint of edge of f from i to i+1
				//    n := n+1
				//
				// (F+2R+(n-3)P)/n
				org := maptbl[FLATTEN_VERTEX_IDX(t, int(face.Indices[a]))]
				ov := new_points[org]

				if !ov.First {
					ov.First = true

					adj, cnt := GET_ADJACENT_FACES_AND_CNT(org)

					if cnt < 3 {
						ov.Second = NewVertexFromAiMesh(minp, int(face.Indices[a]))
					} else {

						F, R := NewVertex(), NewVertex()
						for o := 0; o < cnt; o++ {
							common.AiAssert(adj[o] < totfaces)
							F = F.Add(F, centroids[adj[o]])

							// adj[0] is a global face index - search the face in the mesh list
							var mp *AiMesh
							var nidx int

							if adj[o] < moffsets[0].First {
								nidx = 0
								mp = smesh[nidx]
							} else {
								for nidx = 1; nidx <= len(smesh); nidx++ {
									if nidx == len(smesh) || moffsets[nidx].First > adj[o] {
										nidx--
										mp = smesh[nidx]

										break
									}
								}
							}

							common.AiAssert(adj[o]-moffsets[nidx].First < len(mp.Faces))
							f := mp.Faces[adj[o]-moffsets[nidx].First]
							haveit := false
							// find our original point in the face
							for m := 0; m < len(f.Indices); m++ {
								if maptbl[FLATTEN_VERTEX_IDX(nidx, int(f.Indices[m]))] == org {

									// add *both* edges. this way, we can be sure that we add
									// *all* adjacent edges to R. In a closed shape, every
									// edge is added twice - so we simply leave out the
									// factor 2.f in the amove formula and get the right
									// result.
									tmp := m - 1
									if m == 0 {
										tmp = len(f.Indices) - 1
									}

									key = makeEdgeHash(org, maptbl[FLATTEN_VERTEX_IDX(nidx, int(f.Indices[tmp]))])
									c0, ok := edges[key]
									if !ok {
										c0 = NewEdge()
										edges[key] = c0
									}
									// fixme: replace with mod face.mNumIndices?
									tmp = m + 1
									if m == len(f.Indices)-1 {
										tmp = 0
									}
									key = makeEdgeHash(org, maptbl[FLATTEN_VERTEX_IDX(
										nidx, int(f.Indices[tmp]))])
									c1, ok := edges[key]
									if !ok {
										c1 = NewEdge()
										edges[key] = c1
									}
									// fixme: replace with mod face.mNumIndices?
									R = R.Add(R, c0.midpoint.Add(c0.midpoint, c1.midpoint))
									haveit = true
									break
								}
							}
							// this invariant *must* hold if the vertex-to-face adjacency table is valid
							common.AiAssert(haveit)
							if !haveit {
								logger.Warn("OBJ: no name for material library specified.")
							}
						}
						d := float32(cnt)
						divsq := float32(1.) / (d * d)
						vertex1 := NewVertexFromAiMesh(minp, int(face.Indices[a]))
						vertex1 = vertex1.Mul(vertex1, (d-3.)/d)
						R = R.Mul(R, divsq)
						vertex1 = vertex1.Add(vertex1, R)
						F = F.Mul(F, divsq)
						ov.Second = vertex1.Add(vertex1, F)
					}
				}
				faceOut.Indices[2] = uint32(v)
				ov.Second.SortBack(mout, v)
				v++
			}
		}
		cost := time.Since(now1).Milliseconds()
		logger.DebugF("1 cost :%v Milliseconds", cost)
	}
	// end of scope for edges, freeing its memory

	// ---------------------------------------------------------------------
	// 7. Apply the next subdivision step.
	// ---------------------------------------------------------------------
	if num != 1 {
		tmp1 := make([]*AiMesh, len(smesh))
		div.InternSubdivide(out, tmp1, num-1)
		for i := 0; i < len(smesh); i++ {
			out[i] = tmp1[i]
		}
	}

}
func newCatmullClarkSubdivider() *CatmullClarkSubdivider {
	return &CatmullClarkSubdivider{
		EdgeMap: map[int]Edge{},
	}
}
