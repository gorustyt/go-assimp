package BLEND

import "assimp/common/logger"

type BlenderBMeshConverter struct {
	BMesh                            *Mesh
	triMesh                          *Mesh
	ASSIMP_BLEND_WITH_GLU_TESSELLATE bool
}

func NewBlenderBMeshConverter(mesh *Mesh) *BlenderBMeshConverter {
	b := &BlenderBMeshConverter{BMesh: mesh}
	return b
}

// ------------------------------------------------------------------------------------------------
func (b *BlenderBMeshConverter) AddTFace(uv1, uv2, uv3, uv4 []float64) {
	var mtface MTFace
	copy(mtface.uv[0][:], uv1[:2])
	copy(mtface.uv[1][:], uv2[:2])
	copy(mtface.uv[2][:], uv3[:2])

	if len(uv4) > 0 {
		copy(mtface.uv[3][:], uv4[:2])
	}

	b.triMesh.mtface = append(b.triMesh.mtface, mtface)
}

// ------------------------------------------------------------------------------------------------
func (b *BlenderBMeshConverter) AddFace(v1, v2, v3, v4 int32) {
	var face MFace
	face.v1 = v1
	face.v2 = v2
	face.v3 = v3
	face.v4 = v4
	face.flag = 0
	// TODO - Work out how materials work
	face.mat_nr = 0
	b.triMesh.mface = append(b.triMesh.mface, face)
	b.triMesh.totface = int32(len(b.triMesh.mface))
}

// ------------------------------------------------------------------------------------------------
func (b *BlenderBMeshConverter) ConvertPolyToFaces(poly *MPoly) {
	polyLoop := b.BMesh.mloop[poly.loopstart:]

	if poly.totloop == 3 || poly.totloop == 4 {
		tmp := int32(0)
		if poly.totloop == 4 {
			tmp = polyLoop[3].v
		}
		b.AddFace(polyLoop[0].v, polyLoop[1].v, polyLoop[2].v, tmp)

		// UVs are optional, so only convert when present.
		if len(b.BMesh.mloopuv) > 0 {
			if (poly.loopstart + poly.totloop) > int32(len(b.BMesh.mloopuv)) {
				logger.FatalF("BMesh uv loop array has incorrect size")
			}
			loopUV := b.BMesh.mloopuv[poly.loopstart:]
			var tmp1 [2]float64
			if poly.totloop == 4 {
				tmp1 = loopUV[3].uv
			}
			b.AddTFace(loopUV[0].uv[:], loopUV[1].uv[:], loopUV[2].uv[:], tmp1[:])
		}
	} else if poly.totloop > 4 {
		if b.ASSIMP_BLEND_WITH_GLU_TESSELLATE {
			//var tessGL BlenderTessellatorGL
			//tessGL.Tessellate(polyLoop, poly.totloop, b.triMesh.mvert)
		} else {
			var tessP2T BlenderTessellatorP2T
			tessP2T.Tessellate(polyLoop, poly.totloop, b.triMesh.mvert)
		}
	}
}

// ------------------------------------------------------------------------------------------------
func (b *BlenderBMeshConverter) AssertValidSizes() {
	if b.BMesh.totpoly != int32(len(b.BMesh.mpoly)) {
		logger.FatalF("BMesh poly array has incorrect size")
	}
	if b.BMesh.totloop != int32(len(b.BMesh.mloop)) {
		logger.FatalF("BMesh loop array has incorrect size")
	}
}

// ------------------------------------------------------------------------------------------------
func (b *BlenderBMeshConverter) PrepareTriMesh() {
	if b.triMesh != nil {
		b.triMesh = nil
	}

	b.triMesh = b.BMesh
	b.triMesh.totface = 0
	b.triMesh.mface = b.triMesh.mface[:0]
}

// ------------------------------------------------------------------------------------------------
func (b *BlenderBMeshConverter) AssertValidMesh() {
	if !b.ContainsBMesh() {
		logger.FatalF("BlenderBMeshConverter requires a BMesh with \"polygons\" - please call BlenderBMeshConverter::ContainsBMesh to check this first")
	}
}

// ------------------------------------------------------------------------------------------------
func (b *BlenderBMeshConverter) ContainsBMesh() bool {
	// TODO - Should probably do some additional verification here
	return b.BMesh.totpoly > 0 && b.BMesh.totloop > 0 && b.BMesh.totvert > 0
}

// ------------------------------------------------------------------------------------------------
func (b *BlenderBMeshConverter) TriangulateBMesh() *Mesh {
	b.AssertValidMesh()
	b.AssertValidSizes()
	b.PrepareTriMesh()

	for i := int32(0); i < b.BMesh.totpoly; i++ {
		poly := b.BMesh.mpoly[i]
		b.ConvertPolyToFaces(poly)
	}

	return b.triMesh
}
