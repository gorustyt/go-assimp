package core

import (
	"assimp/common"
	"assimp/common/pb_msg"
)

// -------------------------------------------------------------------------------
/**
 * Specifies that the scene data structure that was imported is not complete.
 * This flag bypasses some internal validations and allows the import
 * of animation skeletons, material libraries or camera animation paths
 * using Assimp. Most applications won't support such data.
 */
var AI_SCENE_FLAGS_INCOMPLETE = 0x1

type AiScene struct {
	/** Any combination of the AI_SCENE_FLAGS_XXX flags. By default
	 * this value is 0, no flags are set. Most applications will
	 * want to reject all scenes with the AI_SCENE_FLAGS_INCOMPLETE
	 * bit set.
	 */
	Flags uint32

	/** The root node of the hierarchy.
	 *
	 * There will always be at least the root node if the import
	 * was successful (and no special flags have been set).
	 * Presence of further nodes depends on the format and content
	 * of the imported file.
	 */
	RootNode *AiNode
	/** The array of meshes.
	 *
	 * Use the indices given in the aiNode structure to access
	 * this array. The array is mNumMeshes in size. If the
	 * AI_SCENE_FLAGS_INCOMPLETE flag is not set there will always
	 * be at least ONE material.
	 */
	Meshes []*AiMesh
	/** The array of materials.
	 *
	 * Use the index given in each aiMesh structure to access this
	 * array. The array is mNumMaterials in size. If the
	 * AI_SCENE_FLAGS_INCOMPLETE flag is not set there will always
	 * be at least ONE material.
	 */
	Materials []*AiMaterial
	/** The array of animations.
	 *
	 * All animations imported from the given file are listed here.
	 * The array is mNumAnimations in size.
	 */
	Animations []*AiAnimation
	/** The array of embedded textures.
	 *
	 * Not many file formats embed their textures into the file.
	 * An example is Quake's MDL format (which is also used by
	 * some GameStudio versions)
	 */
	Textures []*AiTexture
	/** The array of light sources.
	 *
	 * All light sources imported from the given file are
	 * listed here. The array is mNumLights in size.
	 */
	Lights []*AiLight
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
	Skeletons []*AiSkeleton
}

func (ai *AiScene) ToPbMsg() *pb_msg.AiScene {
	r := &pb_msg.AiScene{}
	r.Flags = ai.Flags
	r.RootNode = ai.RootNode.ToPbMsg()
	for _, v := range ai.Meshes {
		r.Meshes = append(r.Meshes, v.ToPbMsg())
	}
	for _, v := range ai.Materials {
		r.Materials = append(r.Materials, v.ToPbMsg())
	}
	for _, v := range ai.Animations {
		r.Animations = append(r.Animations, v.ToPbMsg())
	}
	for _, v := range ai.Textures {
		r.Textures = append(r.Textures, v.ToPbMsg())
	}
	for _, v := range ai.Lights {
		r.Lights = append(r.Lights, v.ToPbMsg())
	}
	for _, v := range ai.Cameras {
		r.Cameras = append(r.Cameras, v.ToPbMsg())
	}
	for _, v := range ai.MetaData {
		r.MetaData = append(r.MetaData, v.ToPbMsg())
	}
	r.Name = ai.Name
	for _, v := range ai.Skeletons {
		r.Skeletons = append(r.Skeletons, v.ToPbMsg())
	}
	return r
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
	Transformation *common.AiMatrix4x4
	/** Parent node. nullptr if this node is the root node. */
	Parent *AiNode
	/** The child nodes of this node. nullptr if mNumChildren is 0. */
	Children []*AiNode
	/** The meshes of this node. Each entry is an index into the
	 * mesh list of the #aiScene.
	 */
	Meshes []int32
	/** Metadata associated with this node or nullptr if there is no metadata.
	 *  Whether any metadata is generated depends on the source file format. See the
	 * @link importer_notes @endlink page for more information on every source file
	 * format. Importers that don't document any metadata don't write any.
	 */
	MetaData []*AiMetadata
}

func (node *AiNode) ToPbMsg() *pb_msg.AiNode {
	r := pb_msg.AiNode{}
	r.Name = node.Name
	/** The transformation relative to the node's parent. */
	r.Transformation = node.Transformation.ToPbMsg()
	/** Parent node. nullptr if this node is the root node. */
	r.Parent = node.Parent.ToPbMsg()
	/** The child nodes of this node. nullptr if mNumChildren is 0. */
	for _, v := range node.Children {
		r.Children = append(r.Children, v.ToPbMsg())
	}
	/** The meshes of this node. Each entry is an index into the
	 * mesh list of the #aiScene.
	 */
	r.Meshes = node.Meshes
	for _, v := range node.MetaData {
		r.MetaData = append(r.MetaData, v.ToPbMsg())
	}
	return &r
}
func NewAiNode(name string) *AiNode {
	return &AiNode{
		Name:           name,
		Transformation: &common.AiMatrix4x4{},
	}
}
