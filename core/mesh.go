package core

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
	 * Bitwise combination of the members of the #aiPrimitiveType enum.
	 * This specifies which types of primitives are present in the mesh.
	 * The "SortByPrimitiveType"-Step can be used to make sure the
	 * output meshes consist of one primitive type each.
	 */
	PrimitiveTypes int

	/**
	 * The number of vertices in this mesh.
	 * This is also the size of all of the per-vertex data arrays.
	 * The maximum value for this member is #AI_MAX_VERTICES.
	 */
	NumVertices int

	/**
	 * The number of primitives (triangles, polygons, lines) in this  mesh.
	 * This is also the size of the mFaces array.
	 * The maximum value for this member is #AI_MAX_FACES.
	 */
	NumFaces int

	/**
	 * @brief Vertex positions.
	 *
	 * This array is always present in a mesh. The array is
	 * mNumVertices in size.
	 */
	Vertices []*AiVector3D

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
	 * #define IS_QNAN(f) (f != f)
	 * @endcode
	 * still dangerous because even 1.f == 1.f could evaluate to false! (
	 * remember the subtleties of IEEE754 artithmetics). Use stuff like
	 * @c fpclassify instead.
	 * @note Normal vectors computed by Assimp are always unit-length.
	 * However, this needn't apply for normals that have been taken
	 * directly from the model file.
	 */
	mNormals []*AiVector3D

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
	Tangents []*AiVector3D

	/**
	 * @brief Vertex bitangents.
	 *
	 * The bitangent of a vertex points in the direction of the positive
	 * Y texture axis. The array contains normalized vectors, nullptr if not
	 * present. The array is mNumVertices in size.
	 * @note If the mesh contains tangents, it automatically also contains
	 * bitangents.
	 */
	mBitangents []*AiVector3D

	/**
	 * @brief Vertex color sets.
	 *
	 * A mesh may contain 0 to #AI_MAX_NUMBER_OF_COLOR_SETS vertex
	 * colors per vertex. nullptr if not present. Each array is
	 * mNumVertices in size if present.
	 */
	Colors []AiColor4D

	/**
	 * @brief Vertex texture coordinates, also known as UV channels.
	 *
	 * A mesh may contain 0 to AI_MAX_NUMBER_OF_TEXTURECOORDS per
	 * vertex. nullptr if not present. The array is mNumVertices in size.
	 */
	TextureCoords []AiVector3D

	/**
	 * @brief Specifies the number of components for a given UV channel.
	 *
	 * Up to three channels are supported (UVW, for accessing volume
	 * or cube maps). If the value is 2 for a given channel n, the
	 * component p.z of mTextureCoords[n][p] is set to 0.0f.
	 * If the value is 1 for a given channel, p.y is set to 0.0f, too.
	 * @note 4D coordinates are not supported
	 */
	NumUVComponents []int

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
	 * The number of bones this mesh contains. Can be 0, in which case the mBones array is nullptr.
	 */
	NumBones int

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
	MaterialIndex int

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
	 * The number of attachment meshes.
	 * Currently known to work with loaders:
	 * - Collada
	 * - gltf
	 */
	NumAnimMeshes int

	/**
	 * Attachment meshes for this mesh, for vertex-based animation.
	 * Attachment meshes carry replacement data for some of the
	 * mesh'es vertex components (usually positions, normals).
	 * Currently known to work with loaders:
	 * - Collada
	 * - gltf
	 */
	mAnimMeshes []*AiAnimMesh

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

func NewAiMesh() *AiMesh {
	return &AiMesh{
		NumUVComponents: make([]int, AI_MAX_NUMBER_OF_TEXTURECOORDS),
		Colors:          make([]AiColor4D, AI_MAX_NUMBER_OF_COLOR_SETS),
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
	Vertices []AiVector3D

	/** Replacement for aiMesh::mNormals.  */
	Normals []AiVector3D

	/** Replacement for aiMesh::mTangents. */
	mTangents []AiVector3D

	/** Replacement for aiMesh::mBitangents. */
	mBitangents []AiVector3D

	/** Replacement for aiMesh::mColors */
	Colors []AiColor4D

	/** Replacement for aiMesh::mTextureCoords */
	TextureCoords []AiVector3D

	/** The number of vertices in the aiAnimMesh, and thus the length of all
	 * the member arrays.
	 *
	 * This has always the same value as the mNumVertices property in the
	 * corresponding aiMesh. It is duplicated here merely to make the length
	 * of the member arrays accessible even if the aiMesh is not known, e.g.
	 * from language bindings.
	 */
	NumVertices int

	/**
	 * Weight of the AnimMesh.
	 */
	Weight float64
}

func NewAiAnimMesh() *AiAnimMesh {
	return &AiAnimMesh{
		Colors:        make([]AiColor4D, AI_MAX_NUMBER_OF_COLOR_SETS),
		TextureCoords: make([]AiVector3D, AI_MAX_NUMBER_OF_TEXTURECOORDS),
	}
}

type AiFace struct {
	//! Number of indices defining this face.
	//! The maximum value for this member is #AI_MAX_FACE_INDICES.
	NumIndices []int

	//! Pointer to the indices array. Size of the array is given in numIndices.
	Indices []int
}

type AiVertexWeight struct {
	//! Index of the vertex which is influenced by the bone.
	VertexId int

	//! The strength of the influence in the range (0...1).
	//! The influence from all bones at one vertex amounts to 1.
	Weight float64
}

type AiBone struct {
	/**
	 * The name of the bone.
	 */
	Name string

	/**
	 * The number of vertices affected by this bone.
	 * The maximum value for this member is #AI_MAX_BONE_WEIGHTS.
	 */
	NumWeights int
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
	OffsetMatrix AiMatrix4x4
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
 * root->node1->node3
 * Each node is represented as a skeleton instance.
 */

type AiSkeleton struct {
	/**
	 *  @brief The name of the skeleton instance.
	 */
	Name string

	/**
	 *  @brief  The number of bones in the skeleton.
	 */
	NumBones int

	/**
	 *  @brief The bone instance in the skeleton.
	 */
	Bones []*AiSkeletonBone
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
 * Tree: s1 -> s2 -> s3
 * Offset-Matrix s3 = locale-s3 * locale-s2 * locale-s1
 */

type AiSkeletonBone struct {
	/// The parent bone index, is -1 one if this bone represents the root bone.
	Parent int
	/// @brief The bone armature node - used for skeleton conversion
	/// you must enable aiProcess_PopulateArmatureData to populate this
	Armature []*AiNode

	/// @brief The bone node in the scene - used for skeleton conversion
	/// you must enable aiProcess_PopulateArmatureData to populate this
	Node []*AiNode
	/// @brief The number of weights
	NumnWeights int

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
	OffsetMatrix AiMatrix4x4

	/// Matrix that transforms the locale bone in bind pose.
	LocalMatrix AiMatrix4x4
}
