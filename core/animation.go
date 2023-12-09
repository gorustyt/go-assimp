package core

import "assimp/common"

type AiAnimBehaviour int

const (
	/** The value from the default node transformation is taken*/
	AiAnimBehaviour_DEFAULT AiAnimBehaviour = 0x0

	/** The nearest key value is used without interpolation */
	AiAnimBehaviour_CONSTANT AiAnimBehaviour = 0x1

	/** The value of the nearest two keys is linearly
	 *  extrapolated for the current time value.*/
	AiAnimBehaviour_LINEAR AiAnimBehaviour = 0x2

	/** The animation is repeated.
	 *
	 *  If the animation key go from n to m and the current
	 *  time is t, use the value at (t-n) % (|m-n|).*/
	AiAnimBehaviour_REPEAT AiAnimBehaviour = 0x3
)

type AiAnimation struct {
	/** The name of the animation. If the modeling package this data was
	 *  exported from does support only a single animation channel, this
	 *  name is usually empty (length is zero). */
	Name string

	/** Duration of the animation in ticks.  */
	Duration float32

	/** Ticks per second. 0 if not specified in the imported file */
	TicksPerSecond float32

	/** The number of bone animation channels. Each channel affects
	 *  a single node. */
	NumChannels int

	/** The node animation channels. Each channel affects a single node.
	 *  The array is mNumChannels in size. */
	mChannels []*AiNodeAnim

	/** The number of mesh animation channels. Each channel affects
	 *  a single mesh and defines vertex-based animation. */
	NumMeshChannels int

	/** The mesh animation channels. Each channel affects a single mesh.
	 *  The array is mNumMeshChannels in size. */
	MeshChannels []*AiMeshAnim

	/** The number of mesh animation channels. Each channel affects
	 *  a single mesh and defines morphing animation. */
	NumMorphMeshChannels int

	/** The morph mesh animation channels. Each channel affects a single mesh.
	 *  The array is mNumMorphMeshChannels in size. */
	MorphMeshChannels []*AiMeshMorphAnim
}

type AiMeshMorphAnim struct {
	/** Name of the mesh to be animated. An empty string is not allowed,
	 *  animated meshes need to be named (not necessarily uniquely,
	 *  the name can basically serve as wildcard to select a group
	 *  of meshes with similar animation setup)*/
	Name string

	/** Size of the #mKeys array. Must be 1, at least. */
	NumKeys int

	/** Key frames of the animation. May not be nullptr. */
	mKeys []*AiMeshMorphKey
}

type AiMeshMorphKey struct {
	/** The time of this key */
	Time float32

	/** The values and weights at the time of this key
	 *   - mValues: index of attachment mesh to apply weight at the same position in mWeights
	 *   - mWeights: weight to apply to the blend shape index at the same position in mValues
	 */
	Values  []int
	Weights []float32

	/** The number of values and weights */
	NumValuesAndWeights int
}

/** Describes the animation of a single node. The name specifies the
 *  bone/node which is affected by this animation channel. The keyframes
 *  are given in three separate series of values, one each for position,
 *  rotation and scaling. The transformation matrix computed from these
 *  values replaces the node's original transformation matrix at a
 *  specific time.
 *  This means all keys are absolute and not relative to the bone default pose.
 *  The order in which the transformations are applied is
 *  - as usual - scaling, rotation, translation.
 *
 *  @note All keys are returned in their correct, chronological order.
 *  Duplicate keys don't pass the validation step. Most likely there
 *  will be no negative time values, but they are not forbidden also ( so
 *  implementations need to cope with them! ) */

type AiNodeAnim struct {
	/** The name of the node affected by this animation. The node
	 *  must exist and it must be unique.*/
	NodeName string

	/** The number of position keys */
	NumPositionKeys int

	/** The position keys of this animation channel. Positions are
	 * specified as 3D vector. The array is mNumPositionKeys in size.
	 *
	 * If there are position keys, there will also be at least one
	 * scaling and one rotation key.*/
	PositionKeys []*AiVectorKey

	/** The number of rotation keys */
	NumRotationKeys int

	/** The rotation keys of this animation channel. Rotations are
	 *  given as quaternions,  which are 4D vectors. The array is
	 *  mNumRotationKeys in size.
	 *
	 * If there are rotation keys, there will also be at least one
	 * scaling and one position key. */
	RotationKeys []*AiQuatKey

	/** The number of scaling keys */
	NumScalingKeys int

	/** The scaling keys of this animation channel. Scalings are
	 *  specified as 3D vector. The array is mNumScalingKeys in size.
	 *
	 * If there are scaling keys, there will also be at least one
	 * position and one rotation key.*/
	ScalingKeys []*AiVectorKey

	/** Defines how the animation behaves before the first
	 *  key is encountered.
	 *
	 *  The default value is aiAnimBehaviour_DEFAULT (the original
	 *  transformation matrix of the affected node is used).*/
	PreState AiAnimBehaviour

	/** Defines how the animation behaves after the last
	 *  key was processed.
	 *
	 *  The default value is aiAnimBehaviour_DEFAULT (the original
	 *  transformation matrix of the affected node is taken).*/
	PostState AiAnimBehaviour
}

/** Describes vertex-based animations for a single mesh or a group of
 *  meshes. Meshes carry the animation data for each frame in their
 *  aiMesh::mAnimMeshes array. The purpose of aiMeshAnim is to
 *  define keyframes linking each mesh attachment to a particular
 *  point in time. */

type AiMeshAnim struct {
	/** Name of the mesh to be animated. An empty string is not allowed,
	 *  animated meshes need to be named (not necessarily uniquely,
	 *  the name can basically serve as wild-card to select a group
	 *  of meshes with similar animation setup)*/
	Name string

	/** Size of the #mKeys array. Must be 1, at least. */
	NumKeys int

	/** Key frames of the animation. May not be nullptr. */
	Keys []*AiMeshKey
}

/** Binds a anim-mesh to a specific point in time. */

type AiMeshKey struct {
	/** The time of this key */
	Time float32

	/** Index into the aiMesh::mAnimMeshes array of the
	 *  mesh corresponding to the #aiMeshAnim hosting this
	 *  key frame. The referenced anim mesh is evaluated
	 *  according to the rules defined in the docs for #aiAnimMesh.*/
	Value int
}

/** A time-value pair specifying a certain 3D vector for the given time. */

type AiVectorKey struct {
	/** The time of this key */
	Time float32

	/** The value of this key */
	Value common.AiVector3D
}

type AiQuatKey struct {
	/** The time of this key */
	Time float32

	/** The value of this key */
	Value common.AiQuaternion
}
