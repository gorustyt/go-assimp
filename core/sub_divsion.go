package core

import "assimp/common"

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
	Subdivide(mesh []*AiMesh,
		out []*AiMesh,
		num int,
		discard_input bool)
	Subdivide1(smesh *[]*AiMesh,
		nmesh int,
		out *[]*AiMesh,
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
	edge_point, midpoint Vertex
	ref                  int
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

func (div *CatmullClarkSubdivider) Subdivide(mesh []*AiMesh,
	out []*AiMesh,
	num int,
	discard_input bool) {
	common.AiAssert(&mesh != &out)

	div.Subdivide1(&mesh, 1, &out, num, discard_input)
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

func (div *CatmullClarkSubdivider) Subdivide1(smesh *[]*AiMesh,
	nmesh int,
	out *[]*AiMesh,
	num int,
	discard_input bool) {

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

func (div *CatmullClarkSubdivider) InternSubdivide() {

}
func newCatmullClarkSubdivider() *CatmullClarkSubdivider {
	return &CatmullClarkSubdivider{
		EdgeMap: map[int]Edge{},
	}
}
