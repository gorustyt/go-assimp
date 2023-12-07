package core

type AiScene struct {
	/** Any combination of the AI_SCENE_FLAGS_XXX flags. By default
	 * this value is 0, no flags are set. Most applications will
	 * want to reject all scenes with the AI_SCENE_FLAGS_INCOMPLETE
	 * bit set.
	 */
	Flags int

	/** The root node of the hierarchy.
	 *
	 * There will always be at least the root node if the import
	 * was successful (and no special flags have been set).
	 * Presence of further nodes depends on the format and content
	 * of the imported file.
	 */
	RootNode []*AiNode

	/** The number of meshes in the scene. */
	NumMeshes int

	/** The array of meshes.
	 *
	 * Use the indices given in the aiNode structure to access
	 * this array. The array is mNumMeshes in size. If the
	 * AI_SCENE_FLAGS_INCOMPLETE flag is not set there will always
	 * be at least ONE material.
	 */
	Meshes []*AiMesh

	/** The number of materials in the scene. */
	NumMaterials int

	/** The array of materials.
	 *
	 * Use the index given in each aiMesh structure to access this
	 * array. The array is mNumMaterials in size. If the
	 * AI_SCENE_FLAGS_INCOMPLETE flag is not set there will always
	 * be at least ONE material.
	 */
	Materials []*AiMaterial

	/** The number of animations in the scene. */
	NumAnimations int

	/** The array of animations.
	 *
	 * All animations imported from the given file are listed here.
	 * The array is mNumAnimations in size.
	 */
	Animations []*AiAnimation

	/** The number of textures embedded into the file */
	NumTextures int

	/** The array of embedded textures.
	 *
	 * Not many file formats embed their textures into the file.
	 * An example is Quake's MDL format (which is also used by
	 * some GameStudio versions)
	 */
	Textures []*AiTexture

	/** The number of light sources in the scene. Light sources
	 * are fully optional, in most cases this attribute will be 0
	 */
	NumLights int

	/** The array of light sources.
	 *
	 * All light sources imported from the given file are
	 * listed here. The array is mNumLights in size.
	 */
	Lights []*AiLight

	/** The number of cameras in the scene. Cameras
	 * are fully optional, in most cases this attribute will be 0
	 */
	NumCameras int

	/** The array of cameras.
	 *
	 * All cameras imported from the given file are listed here.
	 * The array is mNumCameras in size. The first camera in the
	 * array (if existing) is the default camera view into
	 * the scene.
	 */
	Cameras []*AiCamera

	/**
	 *  @brief  The global metadata assigned to the scene itself.
	 *
	 *  This data contains global metadata which belongs to the scene like
	 *  unit-conversions, versions, vendors or other model-specific data. This
	 *  can be used to store format-specific metadata as well.
	 */
	MetaData []*AiMetadata

	/** The name of the scene itself.
	 */
	Name string

	/**
	 *
	 */
	NumSkeletons int

	/**
	 *
	 */
	Skeletons []*AiSkeleton
}

// -------------------------------------------------------------------------------
/**
 * A node in the imported hierarchy.
 *
 * Each node has name, a parent node (except for the root node),
 * a transformation relative to its parent and possibly several child nodes.
 * Simple file formats don't support hierarchical structures - for these formats
 * the imported scene does consist of only a single root node without children.
 */
// -------------------------------------------------------------------------------
type AiNode struct {
	/** The name of the node.
	 *
	 * The name might be empty (length of zero) but all nodes which
	 * need to be referenced by either bones or animations are named.
	 * Multiple nodes may have the same name, except for nodes which are referenced
	 * by bones (see #aiBone and #aiMesh::mBones). Their names *must* be unique.
	 *
	 * Cameras and lights reference a specific node by name - if there
	 * are multiple nodes with this name, they are assigned to each of them.
	 * <br>
	 * There are no limitations with regard to the characters contained in
	 * the name string as it is usually taken directly from the source file.
	 *
	 * Implementations should be able to handle tokens such as whitespace, tabs,
	 * line feeds, quotation marks, ampersands etc.
	 *
	 * Sometimes assimp introduces new nodes not present in the source file
	 * into the hierarchy (usually out of necessity because sometimes the
	 * source hierarchy format is simply not compatible). Their names are
	 * surrounded by @verbatim <> @endverbatim e.g.
	 *  @verbatim<DummyRootNode> @endverbatim.
	 */
	Name string

	/** The transformation relative to the node's parent. */
	Transformation [16]float64

	/** Parent node. nullptr if this node is the root node. */
	Parent *AiNode

	/** The number of child nodes of this node. */
	NumChildren int

	/** The child nodes of this node. nullptr if mNumChildren is 0. */
	Children []*AiNode

	/** The number of meshes of this node. */
	NumMeshes int

	/** The meshes of this node. Each entry is an index into the
	 * mesh list of the #aiScene.
	 */
	Meshes []int

	/** Metadata associated with this node or nullptr if there is no metadata.
	 *  Whether any metadata is generated depends on the source file format. See the
	 * @link importer_notes @endlink page for more information on every source file
	 * format. Importers that don't document any metadata don't write any.
	 */
	MetaData []*AiMetadata
}
