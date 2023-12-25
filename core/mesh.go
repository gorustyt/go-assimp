package core

import (
	"assimp/common"
	"assimp/common/pb_msg"
)

type AiMorphingMethod int

const (
	/** Morphing method to be determined */
	aiMorphingMethod_UNKNOWN AiMorphingMethod = 0x0

	/** Interpolation between morph targets */
	aiMorphingMethod_VERTEX_BLEND AiMorphingMethod = 0x1

	/** Normalized morphing between morph targets  */
	aiMorphingMethod_MORPH_NORMALIZED AiMorphingMethod = 0x2

	/** Relative morphing between morph targets  */
	aiMorphingMethod_MORPH_RELATIVE AiMorphingMethod = 0x3

	/** This value is not used. It is just here to force the
	 *  compiler to map this enum to a 32 Bit integer.
	 */
)

// ---------------------------------------------------------------------------
/** @brief A mesh represents a geometry or model with a single material.
 *
 * It usually consists of a number of vertices and a series of primitives/faces
 * referencing the vertices. In addition there might be a series of bones, each
 * of them addressing a number of vertices with a certain weight. Vertex data
 * is presented in channels with each channel containing a single per-vertex
 * information such as a set of texture coordinates or a normal vector.
 * If a data pointer is non-null, the corresponding data stream is present.
 * From C++-programs you can also use the comfort functions Has*() to
 * test for the presence of various data streams.
 *
 * A Mesh uses only a single material which is referenced by a material ID.
 * @note The mPositions member is usually not optional. However, vertex positions
 * *could* be missing if the #AI_SCENE_FLAGS_INCOMPLETE flag is set in
 * @code
 * aiScene::mFlags
 * @endcode
 */

type AiMesh struct {
	/**
	 * Bitwise combination of the members of the #AiPrimitiveType enum.
	 * This specifies which types of primitives are present in the mesh.
	 * The "SortByPrimitiveType"-Step can be used to make sure the
	 * output meshes consist of one primitive type each.
	 */
	PrimitiveTypes uint32
	/**
	 * @brief Vertex positions.
	 *
	 * This array is always present in a mesh. The array is
	 * mNumVertices in size.
	 */
	Vertices []*common.AiVector3D

	/**
	 * @brief Vertex normals.
	 *
	 * The array contains normalized vectors, nullptr if not present.
	 * The array is mNumVertices in size. Normals are undefined for
	 * point and line primitives. A mesh consisting of points and
	 * lines only may not have normal vectors. Meshes with mixed
	 * primitive types (i.e. lines and triangles) may have normals,
	 * but the normals for vertices that are only referenced by
	 * point or line primitives are undefined and set to QNaN (WARN:
	 * qNaN compares to inequal to *everything*, even to qNaN itself.
	 * Using code like this to check whether a field is qnan is:
	 * @code
	 * IS_QNAN(f) (f != f)
	 * @endcode
	 * still dangerous because even 1.f == 1.f could evaluate to false! (
	 * remember the subtleties of IEEE754 artithmetics). Use stuff like
	 * @c fpclassify instead.
	 * @note Normal vectors computed by Assimp are always unit-length.
	 * However, this needn't apply for normals that have been taken
	 * directly from the model file.
	 */
	Normals []*common.AiVector3D

	/**
	 * @brief Vertex tangents.
	 *
	 * The tangent of a vertex points in the direction of the positive
	 * X texture axis. The array contains normalized vectors, nullptr if
	 * not present. The array is mNumVertices in size. A mesh consisting
	 * of points and lines only may not have normal vectors. Meshes with
	 * mixed primitive types (i.e. lines and triangles) may have
	 * normals, but the normals for vertices that are only referenced by
	 * point or line primitives are undefined and set to qNaN.  See
	 * the #mNormals member for a detailed discussion of qNaNs.
	 * @note If the mesh contains tangents, it automatically also
	 * contains bitangents.
	 */
	Tangents []*common.AiVector3D

	/**
	 * @brief Vertex bitangents.
	 *
	 * The bitangent of a vertex points in the direction of the positive
	 * Y texture axis. The array contains normalized vectors, nullptr if not
	 * present. The array is mNumVertices in size.
	 * @note If the mesh contains tangents, it automatically also contains
	 * bitangents.
	 */
	Bitangents []*common.AiVector3D

	/**
	 * @brief Vertex color sets.
	 *
	 * A mesh may contain 0 to #AI_MAX_NUMBER_OF_COLOR_SETS vertex
	 * colors per vertex. nullptr if not present. Each array is
	 * mNumVertices in size if present.
	 */
	Colors [][]*common.AiColor4D

	/**
	 * @brief Vertex texture coordinates, also known as UV channels.
	 *
	 * A mesh may contain 0 to AI_MAX_NUMBER_OF_TEXTURECOORDS per
	 * vertex. nullptr if not present. The array is mNumVertices in size.
	 */
	TextureCoords [][]*common.AiVector3D

	/**
	 * @brief Specifies the number of components for a given UV channel.
	 *
	 * Up to three channels are supported (UVW, for accessing volume
	 * or cube maps). If the value is 2 for a given channel n, the
	 * component p.z of mTextureCoords[n][p] is set to 0.0f.
	 * If the value is 1 for a given channel, p.y is set to 0.0f, too.
	 * @note 4D coordinates are not supported
	 */
	NumUVComponents []uint32

	/**
	 * @brief The faces the mesh is constructed from.
	 *
	 * Each face refers to a number of vertices by their indices.
	 * This array is always present in a mesh, its size is given
	 *  in mNumFaces. If the #AI_SCENE_FLAGS_NON_VERBOSE_FORMAT
	 * is NOT set each face references an unique set of vertices.
	 */
	Faces []*AiFace
	/**
	 * @brief The bones of this mesh.
	 *
	 * A bone consists of a name by which it can be found in the
	 * frame hierarchy and a set of vertex weights.
	 */
	Bones []*AiBone

	/**
	 * @brief The material used by this mesh.
	 *
	 * A mesh uses only a single material. If an imported model uses
	 * multiple materials, the import splits up the mesh. Use this value
	 * as index into the scene's material list.
	 */
	MaterialIndex int32

	/**
	 *  Name of the mesh. Meshes can be named, but this is not a
	 *  requirement and leaving this field empty is totally fine.
	 *  There are mainly three uses for mesh names:
	 *   - some formats name nodes and meshes independently.
	 *   - importers tend to split meshes up to meet the
	 *      one-material-per-mesh requirement. Assigning
	 *      the same (dummy) name to each of the result meshes
	 *      aids the caller at recovering the original mesh
	 *      partitioning.
	 *   - Vertex animations refer to meshes by their names.
	 */
	Name string
	/**
	 * Attachment meshes for this mesh, for vertex-based animation.
	 * Attachment meshes carry replacement data for some of the
	 * mesh'es vertex components (usually positions, normals).
	 * Currently known to work with loaders:
	 * - Collada
	 * - gltf
	 */
	AnimMeshes []*AiAnimMesh

	/**
	 *  Method of morphing when anim-meshes are specified.
	 *  @see aiMorphingMethod to learn more about the provided morphing targets.
	 */
	Method AiMorphingMethod

	/**
	 *  The bounding box.
	 */
	AABB AiAABB

	/**
	 * Vertex UV stream names. Pointer to array of size AI_MAX_NUMBER_OF_TEXTURECOORDS
	 */
	TextureCoordsNames []string
}

func (ai *AiMesh) ToPbMsg() *pb_msg.AiMesh {
	r := pb_msg.AiMesh{}
	r.PrimitiveTypes = ai.PrimitiveTypes
	for _, v := range ai.Vertices {
		r.Vertices = append(r.Vertices, v.ToPbMsg())
	}

	for _, v := range ai.Normals {
		r.Normals = append(r.Normals, v.ToPbMsg())
	}

	for _, v := range ai.Tangents {
		r.Tangents = append(r.Tangents, v.ToPbMsg())
	}

	for _, v := range ai.Bitangents {
		r.Bitangents = append(r.Bitangents, v.ToPbMsg())
	}
	for _, v := range ai.Colors {
		var tmp pb_msg.AiMesh_ColorsArray
		for _, v1 := range v {
			tmp.Colors = append(tmp.Colors, v1.ToPbMsg())
		}
		r.Colors = append(r.Colors, &tmp)
	}

	for _, v := range ai.TextureCoords {
		var tmp pb_msg.AiMesh_TextureCoordsArray
		for _, v1 := range v {
			tmp.TextureCoords = append(tmp.TextureCoords, v1.ToPbMsg())
		}
		r.TextureCoords = append(r.TextureCoords, &tmp)
	}
	r.NumUVComponents = ai.NumUVComponents

	for _, v := range ai.Faces {
		r.Faces = append(r.Faces, v.ToPbMsg())
	}

	for _, v := range ai.Bones {
		r.Bones = append(r.Bones, v.ToPbMsg())
	}
	r.MaterialIndex = int32(ai.MaterialIndex)
	r.Name = ai.Name

	for _, v := range ai.AnimMeshes {
		r.AnimMeshes = append(r.AnimMeshes, v.ToPbMsg())
	}
	r.Method = int32(ai.Method)
	r.AABB = ai.AABB.ToPbMsg()
	r.TextureCoordsNames = ai.TextureCoordsNames
	return &r
}

// ! @brief Check whether the mesh contains positions. Provided no special
// !        scene flags are set, this will always be true
// ! @return true, if positions are stored, false if not.
func (ai *AiMesh) HasPositions() bool {
	return ai.Vertices != nil && len(ai.Vertices) > 0
}

// ! @brief Check whether the mesh contains faces. If no special scene flags
// !        are set this should always return true
// ! @return true, if faces are stored, false if not.
func (ai *AiMesh) HasFaces() bool {
	return ai.Faces != nil && len(ai.Faces) > 0
}

// ! @brief Check whether the mesh contains normal vectors
// ! @return true, if normals are stored, false if not.
func (ai *AiMesh) HasNormals() bool {
	return ai.Normals != nil && len(ai.Normals) > 0
}

// ! @brief Check whether the mesh contains tangent and bitangent vectors.
// !
// ! It is not possible that it contains tangents and no bitangents
// ! (or the other way round). The existence of one of them
// ! implies that the second is there, too.
// ! @return true, if tangents and bi-tangents are stored, false if not.
func (ai *AiMesh) HasTangentsAndBitangents() bool {
	return ai.Tangents != nil && ai.Bitangents != nil && len(ai.Vertices) > 0
}

// ! @brief Check whether the mesh contains a vertex color set
// ! @param index    Index of the vertex color set
// ! @return true, if vertex colors are stored, false if not.
func (ai *AiMesh) HasVertexColors(index int) bool {
	if index >= AI_MAX_NUMBER_OF_COLOR_SETS {
		return false
	}
	return ai.Colors[index] != nil && len(ai.Vertices) > 0
}

// ! @brief Check whether the mesh contains a texture coordinate set
// ! @param index    Index of the texture coordinates set
// ! @return true, if texture coordinates are stored, false if not.
func (ai *AiMesh) HasTextureCoords(index uint32) bool {
	if index >= AI_MAX_NUMBER_OF_TEXTURECOORDS {
		return false
	}
	return (ai.TextureCoords[index] != nil && len(ai.Vertices) > 0)
}

// ! @brief Get the number of UV channels the mesh contains.
// ! @return the number of stored uv-channels.
func (ai *AiMesh) GetNumUVChannels() int {
	var n uint32
	for n < AI_MAX_NUMBER_OF_TEXTURECOORDS && ai.TextureCoords[n] != nil {
		n++
	}

	return int(n)
}

// ! @brief Get the number of vertex color channels the mesh contains.
// ! @return The number of stored color channels.
func (ai *AiMesh) GetNumColorChannels() int {
	var n int
	for n < AI_MAX_NUMBER_OF_COLOR_SETS && ai.Colors[n] != nil {
		n++
	}
	return n
}

// ! @brief Check whether the mesh contains bones.
// ! @return true, if bones are stored.
func (ai *AiMesh) HasBones() bool {
	return ai.Bones != nil && len(ai.Bones) > 0
}

// ! @brief  Check whether the mesh contains a texture coordinate set name
// ! @param pIndex Index of the texture coordinates set
// ! @return true, if texture coordinates for the index exists.
func (ai *AiMesh) HasTextureCoordsName(pIndex uint32) bool {
	if ai.TextureCoordsNames == nil || pIndex >= AI_MAX_NUMBER_OF_TEXTURECOORDS {
		return false
	}
	return ai.TextureCoordsNames[pIndex] != ""
}

// ! @brief  Get a texture coordinate set name
// ! @param  pIndex Index of the texture coordinates set
// ! @return The texture coordinate name.
func (ai *AiMesh) GetTextureCoordsName(index uint32) string {
	if ai.TextureCoordsNames == nil || index >= AI_MAX_NUMBER_OF_TEXTURECOORDS {
		return ""
	}

	return ai.TextureCoordsNames[index]
}
func NewAiMesh() *AiMesh {
	return &AiMesh{
		NumUVComponents: make([]uint32, AI_MAX_NUMBER_OF_TEXTURECOORDS),
		Colors:          make([][]*common.AiColor4D, AI_MAX_NUMBER_OF_COLOR_SETS),
	}
}

// ---------------------------------------------------------------------------
/** @brief An AnimMesh is an attachment to an #aiMesh stores per-vertex
 *  animations for a particular frame.
 *
 *  You may think of an #aiAnimMesh as a `patch` for the host mesh, which
 *  replaces only certain vertex data streams at a particular time.
 *  Each mesh stores n attached attached meshes (#aiMesh::mAnimMeshes).
 *  The actual relationship between the time line and anim meshes is
 *  established by #aiMeshAnim, which references singular mesh attachments
 *  by their ID and binds them to a time offset.
 */

type AiAnimMesh struct {
	/**Anim Mesh name */
	Name string

	/** Replacement for aiMesh::mVertices. If this array is non-nullptr,
	 *  it *must* contain mNumVertices entries. The corresponding
	 *  array in the host mesh must be non-nullptr as well - animation
	 *  meshes may neither add or nor remove vertex components (if
	 *  a replacement array is nullptr and the corresponding source
	 *  array is not, the source data is taken instead)*/
	Vertices []*common.AiVector3D

	/** Replacement for aiMesh::mNormals.  */
	Normals []*common.AiVector3D

	/** Replacement for aiMesh::mTangents. */
	Tangents []*common.AiVector3D

	/** Replacement for aiMesh::mBitangents. */
	Bitangents []*common.AiVector3D

	/** Replacement for aiMesh::mColors */
	Colors [][]*common.AiColor4D

	/** Replacement for aiMesh::mTextureCoords */
	TextureCoords [][]*common.AiVector3D
	/**
	 * Weight of the AnimMesh.
	 */
	Weight float32
}

func (ai *AiAnimMesh) ToPbMsg() *pb_msg.AiAnimMesh {
	r := &pb_msg.AiAnimMesh{}
	/**Anim Mesh name */
	r.Name = ai.Name

	/** Replacement for aiMesh::mVertices. If this array is non-nullptr,
	 *  it *must* contain mNumVertices entries. The corresponding
	 *  array in the host mesh must be non-nullptr as well - animation
	 *  meshes may neither add or nor remove vertex components (if
	 *  a replacement array is nullptr and the corresponding source
	 *  array is not, the source data is taken instead)*/
	for _, v := range ai.Vertices {
		r.Vertices = append(r.Vertices, v.ToPbMsg())
	}
	for _, v := range ai.Normals {
		r.Normals = append(r.Normals, v.ToPbMsg())
	}
	for _, v := range ai.Tangents {
		r.Tangents = append(r.Tangents, v.ToPbMsg())
	}

	for _, v := range ai.Bitangents {
		r.Bitangents = append(r.Bitangents, v.ToPbMsg())
	}

	for _, v := range ai.Colors {
		var tmp pb_msg.AiAnimMesh_ColorsArray
		for _, v1 := range v {
			tmp.Colors = append(tmp.Colors, v1.ToPbMsg())
		}
		r.Colors = append(r.Colors, &tmp)
	}

	for _, v := range ai.TextureCoords {
		var tmp pb_msg.AiAnimMesh_TextureCoordsArray
		for _, v1 := range v {
			tmp.TextureCoords = append(tmp.TextureCoords, v1.ToPbMsg())
		}
		r.TextureCoords = append(r.TextureCoords, &tmp)
	}

	r.Weight = ai.Weight
	return r
}
func (ai *AiAnimMesh) HasNormals() bool {
	return ai.Normals != nil
}

/**
 *  @brief Check whether the anim-mesh overrides the vertex tangents
 *         and bitangents of its host mesh. As for aiMesh,
 *         tangents and bitangents always go together.
 *  @return true if tangents and bi-tangents are stored, false if not.
 */
func (ai *AiAnimMesh) HasTangentsAndBitangents() bool {
	return ai.Tangents != nil
}

/**
 *  @brief Check whether the anim mesh overrides a particular
 *         set of vertex colors on his host mesh.
 *  @param pIndex 0<index<AI_MAX_NUMBER_OF_COLOR_SETS
 *  @return true if vertex colors are stored, false if not.
 */

func (ai *AiAnimMesh) HasVertexColors(pIndex int) bool {
	if pIndex >= AI_MAX_NUMBER_OF_COLOR_SETS {
		return false
	}
	return ai.Colors[pIndex] != nil
}

/**
 *  @brief Check whether the anim mesh overrides a particular
 *        set of texture coordinates on his host mesh.
 *  @param pIndex 0<index<AI_MAX_NUMBER_OF_TEXTURECOORDS
 *  @return true if texture coordinates are stored, false if not.
 */
func (ai *AiAnimMesh) HasTextureCoords(pIndex uint32) bool {
	if pIndex >= AI_MAX_NUMBER_OF_TEXTURECOORDS {
		return false
	}
	return ai.TextureCoords[pIndex] != nil
}

/**
 *  @brief Check whether the anim-mesh overrides the vertex positions
 *         of its host mesh.
 *  @return true if positions are stored, false if not.
 */
func (ai *AiAnimMesh) HasPositions() bool {
	return ai.Vertices != nil
}

func NewAiAnimMesh() *AiAnimMesh {
	return &AiAnimMesh{
		Colors:        make([][]*common.AiColor4D, AI_MAX_NUMBER_OF_COLOR_SETS),
		TextureCoords: make([][]*common.AiVector3D, AI_MAX_NUMBER_OF_TEXTURECOORDS),
	}
}

type AiFace struct {
	//! Pointer to the indices array. Size of the array is given in numIndices.
	Indices []uint32
}

func (ai *AiFace) ToPbMsg() *pb_msg.AiFace {
	return &pb_msg.AiFace{
		Indices: ai.Indices,
	}
}
func NewAiFace() *AiFace {
	return &AiFace{}
}

type AiVertexWeight struct {
	//! Index of the vertex which is influenced by the bone.
	VertexId uint32

	//! The strength of the influence in the range (0...1).
	//! The influence from all bones at one vertex amounts to 1.
	Weight float32
}

func (ai *AiVertexWeight) ToPbMsg() *pb_msg.AiVertexWeight {
	return &pb_msg.AiVertexWeight{
		VertexId: ai.VertexId,
		Weight:   ai.Weight,
	}
}

type AiBone struct {
	/**
	 * The name of the bone.
	 */
	Name string
	/**
	 * The bone armature node - used for skeleton conversion
	 * you must enable aiProcess_PopulateArmatureData to populate this
	 */
	Armature []*AiNode

	/**
	 * The bone node in the scene - used for skeleton conversion
	 * you must enable aiProcess_PopulateArmatureData to populate this
	 */
	Node []*AiNode

	/**
	 * The influence weights of this bone, by vertex index.
	 */
	Weights []*AiVertexWeight

	/**
	 * Matrix that transforms from mesh space to bone space in bind pose.
	 *
	 * This matrix describes the position of the mesh
	 * in the local space of this bone when the skeleton was bound.
	 * Thus it can be used directly to determine a desired vertex position,
	 * given the world-space transform of the bone when animated,
	 * and the position of the vertex in mesh space.
	 *
	 * It is sometimes called an inverse-bind matrix,
	 * or inverse bind pose matrix.
	 */
	OffsetMatrix *common.AiMatrix4x4
}

func (ai *AiBone) ToPbMsg() *pb_msg.AiBone {
	r := &pb_msg.AiBone{}
	r.Name = ai.Name
	for _, v := range ai.Armature {
		r.Armature = append(r.Armature, v.ToPbMsg())
	}
	for _, v := range ai.Node {
		r.Node = append(r.Node, v.ToPbMsg())
	}
	/**
	 * The influence weights of this bone, by vertex index.
	 */
	for _, v := range ai.Weights {
		r.Weights = append(r.Weights, v.ToPbMsg())
	}
	r.OffsetMatrix = ai.OffsetMatrix.ToPbMsg()
	return r
}

/**
 * @brief A skeleton represents the bone hierarchy of an animation.
 *
 * Skeleton animations can be described as a tree of bones:
 *                  root
 *                    |
 *                  node1
 *                  /   \
 *               node3  node4
 * If you want to calculate the transformation of node three you need to compute the
 * transformation hierarchy for the transformation chain of node3:
 * root.node1.node3
 * Each node is represented as a skeleton instance.
 */

type AiSkeleton struct {
	/**
	 *  @brief The name of the skeleton instance.
	 */
	Name string
	/**
	 *  @brief The bone instance in the skeleton.
	 */
	Bones []*AiSkeletonBone
}

func (ai *AiSkeleton) ToPbMsg() *pb_msg.AiSkeleton {
	r := &pb_msg.AiSkeleton{
		Name: ai.Name,
	}
	for _, v := range ai.Bones {
		r.Bones = append(r.Bones, v.ToPbMsg())
	}
	return r
}

/**
 * @brief  A skeleton bone represents a single bone is a skeleton structure.
 *
 * Skeleton-Animations can be represented via a skeleton struct, which describes
 * a hierarchical tree assembled from skeleton bones. A bone is linked to a mesh.
 * The bone knows its parent bone. If there is no parent bone the parent id is
 * marked with -1.
 * The skeleton-bone stores a pointer to its used armature. If there is no
 * armature this value if set to nullptr.
 * A skeleton bone stores its offset-matrix, which is the absolute transformation
 * for the bone. The bone stores the locale transformation to its parent as well.
 * You can compute the offset matrix by multiplying the hierarchy like:
 * Tree: s1 . s2 . s3
 * Offset-Matrix s3 = locale-s3 * locale-s2 * locale-s1
 */

type AiSkeletonBone struct {
	/// The parent bone index, is -1 one if this bone represents the root bone.
	Parent int32
	/// @brief The bone armature node - used for skeleton conversion
	/// you must enable aiProcess_PopulateArmatureData to populate this
	Armature []*AiNode

	/// @brief The bone node in the scene - used for skeleton conversion
	/// you must enable aiProcess_PopulateArmatureData to populate this
	Node []*AiNode
	/// The mesh index, which will get influenced by the weight.
	MeshId []*AiMesh

	/// The influence weights of this bone, by vertex index.
	Weights []*AiVertexWeight

	/** Matrix that transforms from bone space to mesh space in bind pose.
	 *
	 * This matrix describes the position of the mesh
	 * in the local space of this bone when the skeleton was bound.
	 * Thus it can be used directly to determine a desired vertex position,
	 * given the world-space transform of the bone when animated,
	 * and the position of the vertex in mesh space.
	 *
	 * It is sometimes called an inverse-bind matrix,
	 * or inverse bind pose matrix.
	 */
	OffsetMatrix common.AiMatrix4x4

	/// Matrix that transforms the locale bone in bind pose.
	LocalMatrix common.AiMatrix4x4
}

func (ai *AiSkeletonBone) ToPbMsg() *pb_msg.AiSkeletonBone {
	r := &pb_msg.AiSkeletonBone{}
	r.Parent = ai.Parent
	for _, v := range ai.Armature {
		r.Armature = append(r.Armature, v.ToPbMsg())
	}
	for _, v := range ai.Node {
		r.Node = append(r.Node, v.ToPbMsg())
	}
	for _, v := range ai.MeshId {
		r.MeshId = append(r.MeshId, v.ToPbMsg())
	}
	for _, v := range ai.Weights {
		r.Weights = append(r.Weights, v.ToPbMsg())
	}
	r.OffsetMatrix = ai.OffsetMatrix.ToPbMsg()
	r.LocalMatrix = ai.LocalMatrix.ToPbMsg()
	return r
}

// ---------------------------------------------------------------------------
/** @brief Enumerates the types of geometric primitives supported by Assimp.
 *
 *  @see aiFace Face data structure
 *  @see aiProcess_SortByPType Per-primitive sorting of meshes
 *  @see aiProcess_Triangulate Automatic triangulation
 *  @see AI_CONFIG_PP_SBP_REMOVE Removal of specific primitive types.
 */
type AiPrimitiveType int

const (
	/**
	 * @brief A point primitive.
	 *
	 * This is just a single vertex in the virtual world,
	 * #aiFace contains just one index for such a primitive.
	 */
	AiPrimitiveType_POINT AiPrimitiveType = 0x1

	/**
	 * @brief A line primitive.
	 *
	 * This is a line defined through a start and an end position.
	 * #aiFace contains exactly two indices for such a primitive.
	 */
	AiPrimitiveType_LINE AiPrimitiveType = 0x2

	/**
	 * @brief A triangular primitive.
	 *
	 * A triangle consists of three indices.
	 */
	AiPrimitiveType_TRIANGLE AiPrimitiveType = 0x4

	/**
	 * @brief A higher-level polygon with more than 3 edges.
	 *
	 * A triangle is a polygon, but polygon in this context means
	 * "all polygons that are not triangles". The "Triangulate"-Step
	 * is provided for your convenience, it splits all polygons in
	 * triangles (which are much easier to handle).
	 */
	AiPrimitiveType_POLYGON AiPrimitiveType = 0x8

	/**
	 * @brief A flag to determine whether this triangles only mesh is NGON encoded.
	 *
	 * NGON encoding is a special encoding that tells whether 2 or more consecutive triangles
	 * should be considered as a triangle fan. This is identified by looking at the first vertex index.
	 * 2 consecutive triangles with the same 1st vertex index are part of the same
	 * NGON.
	 *
	 * At the moment, only quads (concave or convex) are supported, meaning that polygons are 'seen' as
	 * triangles, as usual after a triangulation pass.
	 *
	 * To get an NGON encoded mesh, please use the aiProcess_Triangulate post process.
	 *
	 * @see aiProcess_Triangulate
	 * @link https://github.com/KhronosGroup/glTF/pull/1620
	 */
	AiPrimitiveType_NGONEncodingFlag AiPrimitiveType = 0x10

	/**
	 * This value is not used. It is just here to force the
	 * compiler to map this enum to a 32 Bit integer.
	 */
) //! enum AiPrimitiveType
