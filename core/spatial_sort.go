package core

import (
	"assimp/common"
	"math"
	"sort"
	"unsafe"
)

// ------------------------------------------------------------------------------------------------
/** A little helper class to quickly find all vertices in the epsilon environment of a given
 * position. Construct an instance with an array of positions. The class stores the given positions
 * by their indices and sorts them by their distance to an arbitrary chosen plane.
 * You can then query the instance for all vertices close to a given position in an average O(log n)
 * time, with O(n) worst case complexity when all vertices lay on the plane. The plane is chosen
 * so that it avoids common planes in usual data sets. */
// ------------------------------------------------------------------------------------------------
var (
	PlaneInit = *common.NewAiVector3D3(0.8523, 0.34321, 0.5736)
)

type SpatialSort struct {
	/** Normal of the sorting plane, normalized.
	 */
	PlaneNormal common.AiVector3D

	/** The centroid of the positions, which is used as a point on the sorting plane
	 * when calculating distance. This value is calculated in Finalize.
	 */
	Centroid common.AiVector3D

	/** An entry in a spatially sorted position array. Consists of a vertex index,
	 * its position and its pre-calculated distance from the reference plane */
	// all positions, sorted by distance to the sorting plane
	Positions []*SpatialSortEntry

	/// false until the Finalize method is called.
	Finalized bool
}

func NewSpatialSort() *SpatialSort {
	s := &SpatialSort{}
	s.PlaneNormal = PlaneInit
	s.PlaneNormal = *s.PlaneNormal.Normalize()
	return s
}

func NewSpatialSortWithPos(pPositions []*common.AiVector3D) *SpatialSort {
	s := NewSpatialSort()
	s.Fill(pPositions)
	return s
}

// ------------------------------------------------------------------------------------------------
func (s *SpatialSort) CalculateDistance(pPosition *common.AiVector3D) float64 {
	return (pPosition.Sub(&s.Centroid)).MulAiVector3D(&s.PlaneNormal)
}

// ------------------------------------------------------------------------------------------------
func (s *SpatialSort) Finalize() {
	scale := 1.0 / len(s.Positions)
	for i := 0; i < len(s.Positions); i++ {
		s.Centroid = *s.Centroid.Add(s.Positions[i].Position.Mul(float32(scale)))
	}
	for i := 0; i < len(s.Positions); i++ {
		s.Positions[i].Distance = s.CalculateDistance(&s.Positions[i].Position)
	}
	sort.Slice(s.Positions, func(i, j int) bool {
		return s.Positions[i].Less(s.Positions[j])
	})
	s.Finalized = true
}

// ------------------------------------------------------------------------------------------------
func (s *SpatialSort) Append(pPositions []*common.AiVector3D, pFinalizes ...bool) {
	pFinalize := true
	if len(pFinalizes) > 0 {
		pFinalize = pFinalizes[0]
	}
	common.AiAssert(!s.Finalized, "You cannot add positions to the SpatialSort object after it has been finalized.")
	// store references to all given positions along with their distance to the reference plane
	initial := len(s.Positions)
	for a := 0; a < len(pPositions); a++ {
		s.Positions = append(s.Positions, NewSpatialSortEntryWithIndex(a+initial, *pPositions[a]))
	}

	if pFinalize {
		// now sort the array ascending by distance.
		s.Finalize()
	}
}

// ------------------------------------------------------------------------------------------------
// Returns an iterator for all positions close to the given position.
func (s *SpatialSort) FindPositions(pPosition *common.AiVector3D,
	pRadius float64) (poResults []int) {
	common.AiAssert(s.Finalized, "The SpatialSort object must be finalized before FindPositions can be called.")
	dist := s.CalculateDistance(pPosition)
	minDist := dist - pRadius
	maxDist := dist + pRadius
	// quick check for positions outside the range
	if len(s.Positions) == 0 {
		return
	}

	if maxDist < s.Positions[0].Distance {
		return
	}

	if minDist > s.Positions[len(s.Positions)].Distance {
		return
	}

	// do a binary search for the minimal distance to start the iteration there
	index := len(s.Positions) / 2
	binaryStepSize := len(s.Positions) / 4
	for binaryStepSize > 1 {
		if s.Positions[index].Distance < minDist {
			index += binaryStepSize
		} else {
			index -= binaryStepSize

			binaryStepSize /= 2
		}

		// depending on the direction of the last step we need to single step a bit back or forth
		// to find the actual beginning element of the range
		for index > 0 && s.Positions[index].Distance > minDist {
			index--
		}
	}

	for index < (len(s.Positions)-1) && s.Positions[index].Distance < minDist {
		index++
	}

	// Mow start iterating from there until the first position lays outside of the distance range.
	// Add all positions inside the distance range within the given radius to the result array
	it := index
	pSquared := pRadius * pRadius
	for s.Positions[index].Distance < maxDist {
		if (s.Positions[index].Position.Sub(pPosition)).SquareLength() < pSquared {
			poResults = append(poResults, s.Positions[index].Index)
		}

		it++
		if it == len(s.Positions) {
			break
		}

	}

	// that's it
	return poResults
}

// ------------------------------------------------------------------------------------------------
func (s *SpatialSort) Fill(pPositions []*common.AiVector3D, pFinalizes ...bool) {
	pFinalize := true
	if len(pFinalizes) > 0 {
		pFinalize = pFinalizes[0]
	}
	s.Positions = s.Positions[:0]
	s.Finalized = false
	s.Append(pPositions, pFinalize)
	s.Finalized = pFinalize
}

// Binary, signed-integer representation of a single-precision floating-point value.
// IEEE 754 says: "If two floating-point numbers in the same format are ordered then they are
//
//	ordered the same way when their bits are reinterpreted as sign-magnitude integers."
//
// This allows us to convert all floating-point numbers to signed integers of arbitrary size
//
//	and then use them to work with ULPs (Units in the Last Place, for high-precision
//	computations) or to compare them (integer comparisons are faster than floating-point
//	comparisons on many platforms).
type BinFloat uint64

// --------------------------------------------------------------------------------------------
// Converts the bit pattern of a floating-point number to its signed integer representation.
func ToBinary(pValue float64) BinFloat {

	// If this assertion fails, signed int is not big enough to store a float on your platform.
	//  Please correct the declaration of BinFloat a few lines above - but do it in a portable,
	//  #ifdef'd manner!
	common.AiAssert(unsafe.Sizeof(BinFloat(0)) >= unsafe.Sizeof(float32(0)), "sizeof(BinFloat) >= sizeof(ai_real)")
	binValue := BinFloat(math.Float64bits(pValue))
	// floating-point numbers are of sign-magnitude format, so find out what signed number
	//  representation we must convert negative values to.
	// See http://en.wikipedia.org/wiki/Signed_number_representations.
	mask := BinFloat(1) << (8*unsafe.Sizeof(BinFloat(0)) - 1)

	// Two's complement?
	DefaultValue := ((-42 == (^42 + 1)) && (binValue&mask) != 0)
	OneComplement := ((-42 == ^42) && (binValue&mask) != 0)

	if DefaultValue {
		return mask - binValue
	} else if OneComplement {
		return BinFloat(-0) - binValue
	}
	// One's complement?

	// Sign-magnitude? -0 = 1000... binary
	return binValue
}

// ------------------------------------------------------------------------------------------------
// Fills an array with indices of all positions identical to the given position. In opposite to
// FindPositions(), not an epsilon is used but a (very low) tolerance of four floating-point units.
func (s *SpatialSort) FindIdenticalPositions(pPosition *common.AiVector3D) (poResults []int) {
	common.AiAssert(s.Finalized, "The SpatialSort object must be finalized before FindIdenticalPositions can be called.")
	// Epsilons have a huge disadvantage: they are of constant precision, while floating-point
	//  values are of log2 precision. If you apply e=0.01 to 100, the epsilon is rather small, but
	//  if you apply it to 0.001, it is enormous.

	// The best way to overcome this is the unit in the last place (ULP). A precision of 2 ULPs
	//  tells us that a float does not differ more than 2 bits from the "real" value. ULPs are of
	//  logarithmic precision - around 1, they are 1*(2^24) and around 10000, they are 0.00125.

	// For standard C math, we can assume a precision of 0.5 ULPs according to IEEE 754. The
	//  incoming vertex positions might have already been transformed, probably using rather
	//  inaccurate SSE instructions, so we assume a tolerance of 4 ULPs to safely identify
	//  identical vertex positions.
	toleranceInULPs := 4
	// An interesting point is that the inaccuracy grows linear with the number of operations:
	//  multiplying to numbers, each inaccurate to four ULPs, results in an inaccuracy of four ULPs
	//  plus 0.5 ULPs for the multiplication.
	// To compute the distance to the plane, a dot product is needed - that is a multiplication and
	//  an addition on each number.
	distanceToleranceInULPs := toleranceInULPs + 1
	// The squared distance between two 3D vectors is computed the same way, but with an additional
	//  subtraction.
	distance3DToleranceInULPs := distanceToleranceInULPs + 1

	// Convert the plane distance to its signed integer representation so the ULPs tolerance can be
	//  applied. For some reason, VC won't optimize two calls of the bit pattern conversion.
	minDistBinary := ToBinary(s.CalculateDistance(pPosition)) - BinFloat(distanceToleranceInULPs)
	maxDistBinary := minDistBinary + 2*BinFloat(distanceToleranceInULPs)
	// do a binary search for the minimal distance to start the iteration there
	index := len(s.Positions) / 2
	binaryStepSize := len(s.Positions) / 4
	for binaryStepSize > 1 {
		// Ugly, but conditional jumps are faster with integers than with floats
		if minDistBinary > ToBinary(s.Positions[index].Distance) {
			index += binaryStepSize
		} else {
			index -= binaryStepSize
		}

		binaryStepSize /= 2
	}

	// depending on the direction of the last step we need to single step a bit back or forth
	// to find the actual beginning element of the range
	for index > 0 && minDistBinary < ToBinary(s.Positions[index].Distance) {
		index--
	}

	for index < (len(s.Positions)-1) && minDistBinary > ToBinary(s.Positions[index].Distance) {
		index++
	}

	// Now start iterating from there until the first position lays outside of the distance range.
	// Add all positions inside the distance range within the tolerance to the result array
	it := index
	for ToBinary(s.Positions[it].Distance) < maxDistBinary {
		if BinFloat(distance3DToleranceInULPs) >= ToBinary((s.Positions[it].Position.Sub(pPosition)).SquareLength()) {
			poResults = append(poResults, s.Positions[it].Index)
		}
		it++
		if it == len(s.Positions) {
			break
		}

	}
	return poResults
	// that's it
}

// ------------------------------------------------------------------------------------------------
func (s *SpatialSort) GenerateMappingTable(pRadius float64) ([]int, int) {
	common.AiAssert(s.Finalized, "The SpatialSort object must be finalized before GenerateMappingTable can be called.")

	fill := make([]int, len(s.Positions))
	var dist, maxDist float64

	t := 0
	pSquared := pRadius * pRadius
	for i := 0; i < len(s.Positions); {
		dist = (s.Positions[i].Position.Sub(&s.Centroid)).MulAiVector3D(&s.PlaneNormal)
		maxDist = dist + pRadius

		fill[s.Positions[i].Index] = t
		oldpos := s.Positions[i].Position
		i++
		for ; i < len(fill) && s.Positions[i].Distance < maxDist && (s.Positions[i].Position.Sub(&oldpos)).SquareLength() < pSquared; i++ {
			fill[s.Positions[i].Index] = t
		}
		t++
	}

	// debug invariant: mPositions[i].mIndex values must range from 0 to mPositions.size()-1
	for i := 0; i < len(fill); i++ {
		common.AiAssert(fill[i] < len(s.Positions), "fill[i] < len(s.Positions)")
	}

	return fill, t
}

type SpatialSortEntry struct {
	Index    int               ///< The vertex referred by this entry
	Position common.AiVector3D ///< Position
	/// Distance of this vertex to the sorting plane. This is set by Finalize.
	Distance float64
}

func NewSpatialSortEntry() *SpatialSortEntry {
	return &SpatialSortEntry{
		Index:    math.MaxInt,
		Distance: math.MaxFloat32,
	}
}

func NewSpatialSortEntryWithIndex(pIndex int, pPosition common.AiVector3D) *SpatialSortEntry {
	return &SpatialSortEntry{
		Index:    pIndex,
		Position: pPosition,
		Distance: math.MaxFloat32,
	}
}

func (s *SpatialSortEntry) Less(e *SpatialSortEntry) bool {
	return s.Distance < e.Distance
}
