package core

import "github.com/gorustyt/go-assimp/common"

// -------------------------------------------------------------------------------
func ComputePositionEpsilon(pMeshes []*AiMesh) float64 {
	epsilon := 1e-4
	// calculate the position bounds so we have a reliable epsilon to check position differences against
	minVec, maxVec := common.ChooseMinMaxAiVector3D()
	for a := 0; a < len(pMeshes); a++ {
		pMesh := pMeshes[a]
		tmi, tma := common.ArrayBounds(pMesh.Vertices, func(a *common.AiVector3D, b *common.AiVector3D) *common.AiVector3D {
			return a.BoundMin(b)
		}, func(a *common.AiVector3D, b *common.AiVector3D) *common.AiVector3D {
			return a.BoundMax(b)
		})
		minVec = tmi.(*common.AiVector3D).BoundMin(minVec)
		maxVec = tma.(*common.AiVector3D).BoundMax(maxVec)
	}
	return maxVec.Sub(minVec).Length() * epsilon
}
