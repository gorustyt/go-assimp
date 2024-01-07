package core

import (
	"github.com/gorustyt/go-assimp/common"
	"github.com/gorustyt/go-assimp/common/pb_msg"
)

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
	Duration float64

	/** Ticks per second. 0 if not specified in the imported file */
	TicksPerSecond float64
	/** The node animation channels. Each channel affects a single node.
	 *  The array is mNumChannels in size. */
	Channels []*AiNodeAnim

	/** The mesh animation channels. Each channel affects a single mesh.
	 *  The array is mNumMeshChannels in size. */
	MeshChannels []*AiMeshAnim
	/** The morph mesh animation channels. Each channel affects a single mesh.
	 *  The array is mNumMorphMeshChannels in size. */
	MorphMeshChannels []*AiMeshMorphAnim
}

func (ai *AiAnimation) FromPbMsg(p *pb_msg.AiAnimation) *AiAnimation {
	if p == nil {
		return nil
	}
	ai.Name = p.Name
	ai.Duration = p.Duration
	ai.TicksPerSecond = p.TicksPerSecond
	for _, v := range p.Channels {
		ai.Channels = append(ai.Channels, (&AiNodeAnim{}).FromPbMsg(v))
	}
	for _, v := range p.MeshChannels {
		ai.MeshChannels = append(ai.MeshChannels, (&AiMeshAnim{}).FromPbMsg(v))
	}
	for _, v := range p.MorphMeshChannels {
		ai.MorphMeshChannels = append(ai.MorphMeshChannels, (&AiMeshMorphAnim{}).FromPbMsg(v))
	}
	return ai
}
func (ai *AiAnimation) Clone() *AiAnimation {
	if ai == nil {
		return nil
	}
	r := &AiAnimation{}
	r.Name = ai.Name
	r.Duration = ai.Duration
	r.TicksPerSecond = ai.TicksPerSecond
	for _, v := range ai.Channels {
		r.Channels = append(r.Channels, v.Clone())
	}
	for _, v := range ai.MeshChannels {
		r.MeshChannels = append(r.MeshChannels, v.Clone())
	}
	for _, v := range ai.MorphMeshChannels {
		r.MorphMeshChannels = append(r.MorphMeshChannels, v.Clone())
	}
	return r
}
func (ai *AiAnimation) ToPbMsg() *pb_msg.AiAnimation {
	if ai == nil {
		return nil
	}
	r := &pb_msg.AiAnimation{}
	r.Name = ai.Name
	r.Duration = ai.Duration
	r.TicksPerSecond = ai.TicksPerSecond
	for _, v := range ai.Channels {
		r.Channels = append(r.Channels, v.ToPbMsg())
	}
	for _, v := range ai.MeshChannels {
		r.MeshChannels = append(r.MeshChannels, v.ToPbMsg())
	}
	for _, v := range ai.MorphMeshChannels {
		r.MorphMeshChannels = append(r.MorphMeshChannels, v.ToPbMsg())
	}
	return r
}

type AiMeshMorphAnim struct {
	/** Name of the mesh to be animated. An empty string is not allowed,
	 *  animated meshes need to be named (not necessarily uniquely,
	 *  the name can basically serve as wildcard to select a group
	 *  of meshes with similar animation setup)*/
	Name string
	/** Key frames of the animation. May not be nullptr. */
	Keys []*AiMeshMorphKey
}

func (ai *AiMeshMorphAnim) FromPbMsg(p *pb_msg.AiMeshMorphAnim) *AiMeshMorphAnim {
	if p == nil {
		return nil
	}
	ai.Name = p.Name
	for _, v := range p.Keys {
		ai.Keys = append(ai.Keys, (&AiMeshMorphKey{}).FromPbMsg(v))
	}
	return ai
}
func (ai *AiMeshMorphAnim) Clone() *AiMeshMorphAnim {
	if ai == nil {
		return nil
	}
	r := &AiMeshMorphAnim{}
	r.Name = ai.Name
	for _, v := range ai.Keys {
		r.Keys = append(r.Keys, v.Clone())
	}
	return r
}

func (ai *AiMeshMorphAnim) ToPbMsg() *pb_msg.AiMeshMorphAnim {
	if ai == nil {
		return nil
	}
	r := &pb_msg.AiMeshMorphAnim{}
	r.Name = ai.Name
	for _, v := range ai.Keys {
		r.Keys = append(r.Keys, v.ToPbMsg())
	}
	return r
}

type AiMeshMorphKey struct {
	/** The time of this key */
	Time float64

	/** The values and weights at the time of this key
	 *   - mValues: index of attachment mesh to apply weight at the same position in mWeights
	 *   - mWeights: weight to apply to the blend shape index at the same position in mValues
	 */
	Values  []uint32
	Weights []float64
}

func (ai *AiMeshMorphKey) FromPbMsg(p *pb_msg.AiMeshMorphKey) *AiMeshMorphKey {
	if p == nil {
		return nil
	}
	ai.Time = p.Time
	ai.Values = p.Values
	ai.Weights = p.Weights
	return ai
}
func (ai *AiMeshMorphKey) Clone() *AiMeshMorphKey {
	if ai == nil {
		return nil
	}
	r := &AiMeshMorphKey{}
	r.Time = ai.Time
	r.Values = ai.Values
	r.Weights = ai.Weights
	return r
}
func (ai *AiMeshMorphKey) ToPbMsg() *pb_msg.AiMeshMorphKey {
	if ai == nil {
		return nil
	}
	r := &pb_msg.AiMeshMorphKey{}
	r.Time = ai.Time
	r.Values = ai.Values
	r.Weights = ai.Weights
	return r
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

	/** The position keys of this animation channel. Positions are
	 * specified as 3D vector. The array is mNumPositionKeys in size.
	 *
	 * If there are position keys, there will also be at least one
	 * scaling and one rotation key.*/
	PositionKeys []*common.AiVectorKey
	/** The rotation keys of this animation channel. Rotations are
	 *  given as quaternions,  which are 4D vectors. The array is
	 *  mNumRotationKeys in size.
	 *
	 * If there are rotation keys, there will also be at least one
	 * scaling and one position key. */
	RotationKeys []*common.AiQuatKey
	/** The scaling keys of this animation channel. Scalings are
	 *  specified as 3D vector. The array is mNumScalingKeys in size.
	 *
	 * If there are scaling keys, there will also be at least one
	 * position and one rotation key.*/
	ScalingKeys []*common.AiVectorKey

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

func (ai *AiNodeAnim) FromPbMsg(p *pb_msg.AiNodeAnim) *AiNodeAnim {
	if p == nil {
		return nil
	}
	ai.NodeName = p.NodeName

	for _, v := range p.PositionKeys {
		ai.PositionKeys = append(ai.PositionKeys, (&common.AiVectorKey{}).FromPbMsg(v))
	}
	for _, v := range p.RotationKeys {
		ai.RotationKeys = append(ai.RotationKeys, (&common.AiQuatKey{}).FromPbMsg(v))
	}

	for _, v := range p.ScalingKeys {
		ai.ScalingKeys = append(ai.ScalingKeys, (&common.AiVectorKey{}).FromPbMsg(v))
	}
	ai.PreState = AiAnimBehaviour(p.PreState)
	ai.PostState = AiAnimBehaviour(p.PostState)
	return ai
}
func (ai *AiNodeAnim) Clone() *AiNodeAnim {
	if ai == nil {
		return nil
	}
	r := &AiNodeAnim{}
	r.NodeName = ai.NodeName

	for _, v := range ai.PositionKeys {
		r.PositionKeys = append(r.PositionKeys, v.Clone())
	}
	for _, v := range ai.RotationKeys {
		r.RotationKeys = append(r.RotationKeys, v.Clone())
	}

	for _, v := range ai.ScalingKeys {
		r.ScalingKeys = append(r.ScalingKeys, v.Clone())
	}
	r.PreState = ai.PreState
	r.PostState = ai.PostState
	return r
}
func (ai *AiNodeAnim) ToPbMsg() *pb_msg.AiNodeAnim {
	if ai == nil {
		return nil
	}
	r := &pb_msg.AiNodeAnim{}
	r.NodeName = ai.NodeName

	for _, v := range ai.PositionKeys {
		r.PositionKeys = append(r.PositionKeys, v.ToPbMsg())
	}
	for _, v := range ai.RotationKeys {
		r.RotationKeys = append(r.RotationKeys, v.ToPbMsg())
	}

	for _, v := range ai.ScalingKeys {
		r.ScalingKeys = append(r.ScalingKeys, v.ToPbMsg())
	}
	r.PreState = int32(ai.PreState)
	r.PostState = int32(ai.PostState)
	return r
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

	/** Key frames of the animation. May not be nullptr. */
	Keys []*AiMeshKey
}

func (ai *AiMeshAnim) FromPbMsg(p *pb_msg.AiMeshAnim) *AiMeshAnim {
	if p == nil {
		return nil
	}
	ai.Name = p.Name
	for _, v := range p.Keys {
		ai.Keys = append(ai.Keys, (&AiMeshKey{}).FromPbMsg(v))
	}
	return ai
}

func (ai *AiMeshAnim) Clone() *AiMeshAnim {
	if ai == nil {
		return nil
	}
	r := &AiMeshAnim{
		Name: ai.Name,
	}
	for _, v := range ai.Keys {
		r.Keys = append(r.Keys, v.Clone())
	}
	return r
}
func (ai *AiMeshAnim) ToPbMsg() *pb_msg.AiMeshAnim {
	if ai == nil {
		return nil
	}
	r := &pb_msg.AiMeshAnim{
		Name: ai.Name,
	}
	for _, v := range ai.Keys {
		r.Keys = append(r.Keys, v.ToPbMsg())
	}
	return r
}

/** Binds a anim-mesh to a specific point in time. */

type AiMeshKey struct {
	/** The time of this key */
	Time float64

	/** Index into the aiMesh::mAnimMeshes array of the
	 *  mesh corresponding to the #aiMeshAnim hosting this
	 *  key frame. The referenced anim mesh is evaluated
	 *  according to the rules defined in the docs for #aiAnimMesh.*/
	Value uint32
}

func (ai *AiMeshKey) FromPbMsg(p *pb_msg.AiMeshKey) *AiMeshKey {
	if p == nil {
		return nil
	}
	ai.Time = p.Time
	ai.Value = p.Value
	return ai
}
func (ai *AiMeshKey) Clone() *AiMeshKey {
	if ai == nil {
		return nil
	}
	r := &AiMeshKey{
		Time:  ai.Time,
		Value: ai.Value,
	}
	return r
}

func (ai *AiMeshKey) ToPbMsg() *pb_msg.AiMeshKey {
	if ai == nil {
		return nil
	}
	r := &pb_msg.AiMeshKey{
		Time:  ai.Time,
		Value: ai.Value,
	}
	return r
}
