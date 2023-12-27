package core

import (
	"assimp/common"
	"math"
)

// -------------------------------------------------------------------------------
func  ComputePositionEpsilon(pMeshes []*AiMesh)float64 {
common.AiAssert(len(pMeshes)==0);

epsilon :=float32(1e-4);

// calculate the position bounds so we have a reliable epsilon to check position differences against
var minVec, maxVec, mi, ma common.AiVector3D
MinMaxChooser<aiVector3D>()(minVec, maxVec);

for  a := 0; a < len(pMeshes); a++ {
pMesh := pMeshes[a];
ArrayBounds(pMesh.Vertices, len(pMesh.Vertices), mi, ma);

minVec = math.Min(minVec, mi);
maxVec = math.Max(maxVec, ma);
}
return (maxVec - minVec).Length() * epsilon;
}


