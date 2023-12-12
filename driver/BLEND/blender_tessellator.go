package BLEND

import (
	"assimp/common"
	"assimp/common/logger"
	"assimp/common/poly2tri"
	"math"
)

const (
	BLEND_TESS_MAGIC = 0x83ed9ac3
)

type PointP2T struct {
	point3D *common.AiVector3D
	point2D *poly2tri.Point
	magic   int
	index   int
}

type PlaneP2T struct {
	centre *common.AiVector3D
	normal *common.AiVector3D
}

type BlenderTessellatorP2T struct {
	converter *BlenderBMeshConverter
}

// ------------------------------------------------------------------------------------------------
// Adapted from: http://missingbytes.blogspot.co.uk/2012/06/fitting-plane-to-point-cloud.html
func (f *BlenderTessellatorP2T) FindLLSQPlane(points []*PointP2T) *PlaneP2T {
	var result PlaneP2T

	sum := &common.AiVector3D{}
	for i := 0; i < len(points); i++ {
		sum = sum.Add(points[i].point3D)

	}
	result.centre = sum.Mul(1.0 / float32(len(points)))

	sumXX := float32(0.0)
	sumXY := float32(0.0)
	sumXZ := float32(0.0)
	sumYY := float32(0.0)
	sumYZ := float32(0.0)
	sumZZ := float32(0.0)
	for i := 0; i < len(points); i++ {
		offset := points[i].point3D.Sub(result.centre)
		sumXX += offset.X * offset.X
		sumXY += offset.X * offset.Y
		sumXZ += offset.X * offset.Z
		sumYY += offset.Y * offset.Y
		sumYZ += offset.Y * offset.Z
		sumZZ += offset.Z * offset.Z
	}

	mtx := common.NewAiMatrix3x3WithValues(sumXX, sumXY, sumXZ, sumXY, sumYY, sumYZ, sumXZ, sumYZ, sumZZ)

	det := mtx.Determinant()
	if det == 0.0 {
		result.normal = common.NewAiVector3D1(0.0)
	} else {
		invMtx := mtx
		invMtx.Inverse()
		result.normal = f.GetEigenVectorFromLargestEigenValue(invMtx)
	}

	return &result
}

// ------------------------------------------------------------------------------------------------
// Adapted from: http://missingbytes.blogspot.co.uk/2012/06/fitting-plane-to-point-cloud.html
func (f *BlenderTessellatorP2T) GetEigenVectorFromLargestEigenValue(mtx *common.AiMatrix3x3) *common.AiVector3D {
	//scale := f.FindLargestMatrixElem( mtx );
	//mc := f.ScaleMatrix( mtx, 1.0 / scale );
	//mc = mc * mc * mc;

	v := common.NewAiVector3D1(1.0)
	lastV := v
	for i := 0; i < 100; i++ {
		//v = mc * v;
		v.Normalize()
		if (v.Sub(lastV)).SquareLength() < 1e-16 {
			break
		}
		lastV = v
	}
	return v
}

// ------------------------------------------------------------------------------------------------
func (f *BlenderTessellatorP2T) p2tMax(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// ------------------------------------------------------------------------------------------------
// Adapted from: http://missingbytes.blogspot.co.uk/2012/06/fitting-plane-to-point-cloud.html
func (f *BlenderTessellatorP2T) FindLargestMatrixElem(mtx *common.AiMatrix3x3) float64 {
	// result := 0.0;
	//
	//for  x := 0; x < 3; x ++{
	//for  y := 0; y < 3;y  ++{
	//result = f.p2tMax( math.Abs( mtx[ x ][ y ] ), result );
	//}
	//}

	//return result;
	return 0
}

// ------------------------------------------------------------------------------------------------
// Apparently Assimp doesn't have matrix scaling
func (f *BlenderTessellatorP2T) ScaleMatrix(mtx *common.AiMatrix3x3, scale float64) *common.AiMatrix3x3 {
	//var result common.AiMatrix3x3
	//
	//for x := 0; x < 3; x ++{
	//for  y := 0; y < 3; y ++{
	//result[ x ][ y ] = mtx.Index[ x ][ y ] * scale;
	//}
	//}

	//return &result;
	return nil
}

// ------------------------------------------------------------------------------------------------
func (f *BlenderTessellatorP2T) MakeFacesFromTriangles(triangles []*poly2tri.Triangle) {
	for i := 0; i < len(triangles); i++ {
		Triangle := *triangles[i]

		pointA := f.GetActualPointStructure(Triangle.GetPoint(0))
		pointB := f.GetActualPointStructure(Triangle.GetPoint(1))
		pointC := f.GetActualPointStructure(Triangle.GetPoint(2))

		f.converter.AddFace(int32(pointA.index), int32(pointB.index), int32(pointC.index), 0)
	}
}

// ------------------------------------------------------------------------------------------------
func (f *BlenderTessellatorP2T) GeneratePointTransformMatrix(plane *PlaneP2T) *common.AiMatrix4x4 {
	sideA := common.NewAiVector3D3(1.0, 0.0, 0.0)
	if math.Abs(float64(common.MulAiVector3D(plane.normal, sideA))) > 0.999 {
		sideA = common.NewAiVector3D3(0.0, 1.0, 0.0)
	}

	sideB := common.NegationOperationSymbol(plane.normal, sideA)
	sideB.Normalize()
	sideA = common.NegationOperationSymbol(sideB, plane.normal)

	var result common.AiMatrix4x4
	result.A1 = sideA.X
	result.A2 = sideA.Y
	result.A3 = sideA.Z
	result.B1 = sideB.X
	result.B2 = sideB.Y
	result.B3 = sideB.Z
	result.C1 = plane.normal.X
	result.C2 = plane.normal.Y
	result.C3 = plane.normal.Z
	result.A4 = plane.centre.X
	result.B4 = plane.centre.Y
	result.C4 = plane.centre.Z
	result.Inverse()

	return &result
}

// ------------------------------------------------------------------------------------------------
func (f *BlenderTessellatorP2T) TransformAndFlattenVectices(transform *common.AiMatrix4x4, vertices []*PointP2T) {
	for i := 0; i < len(vertices); i++ {
		point := vertices[i]
		point.point3D = common.Matrix4x4tMulAiVector3D(transform, point.point3D)
		point.point2D.SetPoint(float64(point.point3D.Y), float64(point.point3D.Z))
	}
}

// ------------------------------------------------------------------------------------------------
func (f *BlenderTessellatorP2T) ReferencePoints(points []*PointP2T, pointRefs []*poly2tri.Point) {
	//pointRefs.resize( len(points) );
	//for  i := 0; i < len(points) ; i++{
	//pointRefs[ i ] = points[ i ].point2D;
	//}
}

// ------------------------------------------------------------------------------------------------
func (f *BlenderTessellatorP2T) GetActualPointStructure(point *poly2tri.Point) *PointP2T {

	//PointP2T& pointStruct = *reinterpret_cast< PointP2T* >( reinterpret_cast< char* >( &point ) - pointOffset );
	//if ( pointStruct.magic != static_cast<int>( BLEND_TESS_MAGIC ) ) {
	//logger.FatalF( "Point returned by poly2tri was probably not one of ours. This indicates we need a new way to store vertex information" );
	//}
	//return pointStruct;
	return nil
}

// -

// ------------------------------------------------------------------------------------------------
func (f *BlenderTessellatorP2T) Tessellate(polyLoop []*MLoop, vertexCount int32, vertices []*MVert) {
	f.AssertVertexCount(vertexCount)

	// NOTE - We have to hope that points in a Blender polygon are roughly on the same plane.
	//        There may be some triangulation artifacts if they are wildly different.

	var points []*PointP2T
	f.Copy3DVertices(polyLoop, vertexCount, vertices, points)

	plane := f.FindLLSQPlane(points)

	transform := f.GeneratePointTransformMatrix(plane)

	f.TransformAndFlattenVectices(transform, points)

	var pointRefs []*poly2tri.Point
	f.ReferencePoints(points, pointRefs)

	cdt := poly2tri.NewCDT(pointRefs)

	cdt.Triangulate()
	triangles := cdt.GetTriangles()

	f.MakeFacesFromTriangles(triangles)
}

// ------------------------------------------------------------------------------------------------
func (f *BlenderTessellatorP2T) AssertVertexCount(vertexCount int32) {
	if vertexCount <= 4 {
		logger.FatalF("Expected more than 4 vertices for tessellation")
	}
}

// ------------------------------------------------------------------------------------------------
func (f *BlenderTessellatorP2T) Copy3DVertices(polyLoop []*MLoop, vertexCount int32, vertices []*MVert, points []*PointP2T) {
	//points.resize( vertexCount );
	for i := int32(0); i < vertexCount; i++ {
		loop := polyLoop[i]
		vert := vertices[loop.v]

		point := points[i]
		point.point3D.Set(float32(vert.co[0]), float32(vert.co[1]), float32(vert.co[2]))
		point.index = int(loop.v)
		point.magic = BLEND_TESS_MAGIC
	}
}
