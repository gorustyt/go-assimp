package BLEND

type ID struct {
	name string
	flag int16
	*ElemBase
}

// -------------------------------------------------------------------------------
type ListBase struct {
	*ElemBase
	first IElemBase
	last  IElemBase
}

// -------------------------------------------------------------------------------
type PackedFile struct {
	*ElemBase
	size int32
	seek int32
	data *FileOffset
}

// -------------------------------------------------------------------------------
type GroupObject struct {
	*ElemBase
	prev, next *GroupObject
	ob         Object
}

// -------------------------------------------------------------------------------
type Group struct {
	*ElemBase
	id      ID
	layer   int32
	gobject *GroupObject
}

// -------------------------------------------------------------------------------
type CollectionObject struct {
	*ElemBase
	//CollectionObject* prev;
	next *CollectionObject
	ob   *Object
}

// -------------------------------------------------------------------------------
type CollectionChild struct {
	*ElemBase
	next, prev *CollectionChild
	collection *Collection
}

// -------------------------------------------------------------------------------
type Collection struct {
	*ElemBase
	id       ID
	gobject  ListBase // CollectionObject
	children ListBase // CollectionChild
}

// -------------------------------------------------------------------------------
type MVert struct {
	*ElemBase
	co      [3]float32
	no      [3]float32 // read as and divided through / 32767.f
	flag    uint8
	mat_nr  int32
	bweight int32
}

// -------------------------------------------------------------------------------
type MEdge struct {
	*ElemBase
	v1, v2          int32
	crease, bweight uint8
	flag            int16
}

// -------------------------------------------------------------------------------
type MLoop struct {
	*ElemBase
	v, e int32
}

// -------------------------------------------------------------------------------
type MLoopUV struct {
	*ElemBase
	uv   [2]float32
	flag int32
}

// -------------------------------------------------------------------------------
type World struct {
	*ElemBase
	id ID
}

// -------------------------------------------------------------------------------
// Note that red and blue are not swapped, as with MCol
type MLoopCol struct {
	*ElemBase
	r, g, b, a uint8
}

// -------------------------------------------------------------------------------
type MPoly struct {
	*ElemBase
	loopstart int32
	totloop   int32
	mat_nr    int16
	flag      uint8
}

// -------------------------------------------------------------------------------
type MTexPoly struct {
	*ElemBase
	tpage           *Image
	flag, transp    uint8
	mode, tile, pad int16
}

// -------------------------------------------------------------------------------
type MCol struct {
	*ElemBase
	r, g, b, a uint8
}

// -------------------------------------------------------------------------------
type MFace struct {
	*ElemBase
	v1, v2, v3, v4 int32
	mat_nr         int32
	flag           uint8
}

// -------------------------------------------------------------------------------
type TFace struct {
	*ElemBase
	uv     [4][2]float64
	col    [4]int32
	flag   int8
	mode   int16
	tile   int16
	unwrap int16
}

// -------------------------------------------------------------------------------
type MTFace struct {
	*ElemBase
	uv     [4][2]float32
	flag   uint8
	mode   int16
	tile   int16
	unwrap int16
	// std::shared_ptr<Image> tpage;
}

// -------------------------------------------------------------------------------
type MDeformWeight struct {
	*ElemBase
	def_nr int32
	weight float32
}

// -------------------------------------------------------------------------------
type MDeformVert struct {
	*ElemBase
	dw        []*MDeformWeight
	totweight int32
}

// -------------------------------------------------------------------------------

const (
	MA_RAYMIRROR    = 0x40000
	MA_TRANSPARENCY = 0x10000
	MA_RAYTRANSP    = 0x20000
	MA_ZTRANSP      = 0x00040
)

type Material struct {
	*ElemBase
	id ID

	r, g, b             float32
	specr, specg, specb float32
	har                 int16
	ambr, ambg, ambb    float32
	mirr, mirg, mirb    float32
	emit                float32
	ray_mirror          float32
	alpha               float32
	ref                 float32
	translucency        float32
	mode                int32
	roughness           float32
	darkness            float32
	refrac              float32

	amb              float32
	ang              float32
	spectra          float32
	spec             float32
	zoffs            float32
	add              float32
	fresnel_mir      float32
	fresnel_mir_i    float32
	fresnel_tra      float32
	fresnel_tra_i    float32
	filter           float32
	tx_limit         float32
	tx_falloff       float32
	gloss_mir        float32
	gloss_tra        float32
	adapt_thresh_mir float32
	adapt_thresh_tra float32
	aniso_gloss_mir  float32
	dist_mir         float32
	hasize           float32
	flaresize        float32
	subsize          float32
	flareboost       float32
	strand_sta       float32
	strand_end       float32
	strand_ease      float32
	strand_surfnor   float32
	strand_min       float32
	strand_widthfade float32
	sbias            float32
	lbias            float32
	shad_alpha       float32
	param            float32
	rms              float32
	rampfac_col      float32
	rampfac_spec     float32
	friction         float32
	fh               float32
	reflect          float32
	fhdist           float32
	xyfrict          float32
	sss_radius       float32
	sss_col          float32
	sss_error        float32
	sss_scale        float32
	sss_ior          float32
	sss_colfac       float32
	sss_texfac       float32
	sss_front        float32
	sss_back         float32

	material_type   int16
	flag            int16
	ray_depth       int16
	ray_depth_tra   int16
	samp_gloss_mir  int16
	samp_gloss_tra  int16
	fadeto_mir      int16
	shade_flag      int16
	flarec          int16
	starc           int16
	linec           int16
	ringc           int16
	pr_lamp         int16
	pr_texture      int16
	ml_flag         int16
	texco           int16
	mapto           int16
	ramp_show       int16
	pad3            int16
	dynamode        int16
	pad2            int16
	sss_flag        int16
	sss_preset      int16
	shadowonly_flag int16
	index           int16
	vcol_alpha      int16
	pad4            int16
	seed1           uint8
	seed2           uint8

	group *Group

	diff_shader int16
	spec_shader int16

	mtex [18]*MTex
}

/*
CustomDataLayer 104

	int type 0 4
	int offset 4 4
	int flag 8 4
	int active 12 4
	int active_rnd 16 4
	int active_clone 20 4
	int active_mask 24 4
	int uid 28 4
	char name 32 64
	void *data 96 8
*/
type CustomDataLayer struct {
	*ElemBase
	Type         CustomDataType
	offset       int32
	flag         int32
	active       int32
	active_rnd   int32
	active_clone int32
	active_mask  int32
	uid          int32
	name         string
	data         IElemBase
}

/*
CustomData 208

	CustomDataLayer *layers 0 8
	int typemap 8 168
	int pad_i1 176 4
	int totlayer 180 4
	int maxlayer 184 4
	int totsize 188 4
	BLI_mempool *pool 192 8
	CustomDataExternal *external 200 8
*/
type CustomData struct {
	*ElemBase
	layers   []*CustomDataLayer
	typemap  [42]int32 // CD_NUMTYPES
	totlayer int32
	maxlayer int32
	totsize  int32
	/*
	   pool *BLI_mempool
	   external *CustomDataExternal
	*/
}

type Mesh struct {
	*ElemBase
	id ID

	totface int32
	totedge int32
	totvert int32
	totloop int32
	totpoly int32

	subdiv      int16
	subdivr     int16
	subsurftype int16
	smoothresh  int16

	mface    []*MFace
	mtface   []*MTFace
	tface    []*TFace
	mvert    []*MVert
	medge    []*MEdge
	mloop    []*MLoop
	mloopuv  []*MLoopUV
	mloopcol []*MLoopCol
	mpoly    []*MPoly
	mtpoly   []*MTexPoly
	dvert    []*MDeformVert
	mcol     []*MCol

	mat []*Material

	vdata *CustomData
	edata *CustomData
	fdata *CustomData
	pdata *CustomData
	ldata *CustomData
}

func NewMesh() *Mesh {
	return &Mesh{}
}

// -------------------------------------------------------------------------------
type Library struct {
	*ElemBase
	id ID

	name     string
	filename string
	parent   *Library
}

type CameraType int

const (
	CameraType_PERSP CameraType = 0
	CameraType_ORTHO            = 1
)

// -------------------------------------------------------------------------------
type Camera struct {
	*ElemBase

	id ID

	Type, flag       CameraType
	lens             float32
	sensor_x         float32
	clipsta, clipend float32
}

// -------------------------------------------------------------------------------
type LampFalloffType int

const (
	LampFalloffType_Constant  LampFalloffType = 0x0
	LampFalloffType_InvLinear LampFalloffType = 0x1
	LampFalloffType_InvSquare LampFalloffType = 0x2
	//,FalloffType_Curve    = 0x3
	//,FalloffType_Sliders  = 0x4
)

type LampType int

const (
	LampType_Local LampType = 0x0
	LampType_Sun   LampType = 0x1
	LampType_Spot  LampType = 0x2
	LampType_Hemi  LampType = 0x3
	LampType_Area  LampType = 0x4
	//,Type_YFPhoton    = 0x5
)

type Lamp struct {
	*ElemBase
	id ID
	//AnimData *adt;

	Type  LampType
	flags int16

	//int mode;

	colormodel, totex int16
	r, g, b, k        float32
	//float shdwr, shdwg, shdwb;

	energy, dist, spotsize, spotblend float32
	//float haint;

	constant_coefficient  float32
	linear_coefficient    float32
	quadratic_coefficient float32

	att1, att2 float32
	//struct CurveMapping *curfalloff;
	falloff_type LampFalloffType

	//float clipsta, clipend, shadspotsize;
	//float bias, soft, compressthresh;
	//short bufsize, samp, buffers, filtertype;
	//char bufflag, buftype;

	//short ray_samp, ray_sampy, ray_sampz;
	//short ray_samp_type;
	area_shape                        int16
	area_size, area_sizey, area_sizez float32
	//float adapt_thresh;
	//short ray_samp_method;

	//short texact, shadhalostep;

	//short sun_effect_type;
	//short skyblendtype;
	//float horizon_brightness;
	//float spread;
	sun_brightness float32
	//float sun_size;
	//float backscattered_light;
	//float sun_intensity;
	//float atm_turbidity;
	//float atm_inscattering_factor;
	//float atm_extinction_factor;
	//float atm_distance_factor;
	//float skyblendfac;
	//float sky_exposure;
	//short sky_colorspace;

	// int YF_numphotons, YF_numsearch;
	// short YF_phdepth, YF_useqmc, YF_bufsize, YF_pad;
	// float YF_causticblur, YF_ltradius;

	// float YF_glowint, YF_glowofs;
	// short YF_glowtype, YF_pad2;

	//struct Ipo *ipo;
	//struct MTex *mtex[18];
	// short pr_texture;

	//struct PreviewImage *preview;
}

// -------------------------------------------------------------------------------
type ModifierData_ModifierType int

const (
	eModifierType_None ModifierData_ModifierType = iota
	eModifierType_Subsurf
	eModifierType_Lattice
	eModifierType_Curve
	eModifierType_Build
	eModifierType_Mirror
	eModifierType_Decimate
	eModifierType_Wave
	eModifierType_Armature
	eModifierType_Hook
	eModifierType_Softbody
	eModifierType_Boolean
	eModifierType_Array
	eModifierType_EdgeSplit
	eModifierType_Displace
	eModifierType_UVProject
	eModifierType_Smooth
	eModifierType_Cast
	eModifierType_MeshDeform
	eModifierType_ParticleSystem
	eModifierType_ParticleInstance
	eModifierType_Explode
	eModifierType_Cloth
	eModifierType_Collision
	eModifierType_Bevel
	eModifierType_Shrinkwrap
	eModifierType_Fluidsim
	eModifierType_Mask
	eModifierType_SimpleDeform
	eModifierType_Multires
	eModifierType_Surface
	eModifierType_Smoke
	eModifierType_ShapeKey
)

type ModifierData struct {
	*ElemBase
	next IElemBase
	prev IElemBase

	Type, mode int32
	name       string
}

// ------------------------------------------------------------------------------------------------
type SharedModifierData struct {
	*ElemBase
	modifier *ModifierData
}

// -------------------------------------------------------------------------------
type SubsurfModifierDataType int
type SubsurfModifierDataFlags int

const (
	SubsurfModifierDataType_CatmullClarke SubsurfModifierDataType = 0x0
	SubsurfModifierDataTYPE_Simple        SubsurfModifierDataType = 0x1
	// some omitted

	FLAGS_SubsurfUV SubsurfModifierDataFlags = 1 << 3
)

type SubsurfModifierData struct {
	*ElemBase
	*SharedModifierData
	subdivType   int16
	levels       int16
	renderLevels int16
	flags        int16
}

// -------------------------------------------------------------------------------
type MirrorModifierDataFlags int

const (
	Flags_CLIPPING = 1 << 0
	Flags_MIRROR_U = 1 << 1
	Flags_MIRROR_V = 1 << 2
	Flags_AXIS_X   = 1 << 3
	Flags_AXIS_Y   = 1 << 4
	Flags_AXIS_Z   = 1 << 5
	Flags_VGROUP   = 1 << 6
)

type MirrorModifierData struct {
	*ElemBase
	*SharedModifierData
	axis, flag int16
	tolerance  float32
	mirror_ob  *Object
}

// -------------------------------------------------------------------------------
type ObjectType int

const (
	Type_EMPTY  ObjectType = 0
	Type_MESH   ObjectType = 1
	Type_CURVE  ObjectType = 2
	Type_SURF   ObjectType = 3
	Type_FONT   ObjectType = 4
	Type_MBALL  ObjectType = 5
	Type_LAMP   ObjectType = 10
	Type_CAMERA ObjectType = 11

	Type_WAVE    ObjectType = 21
	Type_LATTICE ObjectType = 22
)

type Object struct {
	*ElemBase
	id        ID
	Type      ObjectType
	obmat     [4][4]float32
	parentinv [4][4]float32
	parsubstr string

	parent *Object
	track  *Object

	proxy, proxy_from, proxy_group *Object
	dup_group                      *Group
	data                           IElemBase

	modifiers *ListBase
}

// -------------------------------------------------------------------------------
type Base struct {
	*ElemBase
	prev   *Base
	next   *Base
	object *Object
}

// -------------------------------------------------------------------------------
type Scene struct {
	*ElemBase
	id ID

	camera            *Object
	world             *World
	basact            *Base
	master_collection *Collection
	base              ListBase
}

// -------------------------------------------------------------------------------
type Image struct {
	*ElemBase
	id ID

	name string

	//struct anim *anim;

	ok, flag                int16
	source, Type, pad, pad1 int16
	lastframe               int32

	tpageflag, totbind int16
	xrep, yrep         int16
	twsta, twend       int16
	//unsigned int bindcode;
	//unsigned int *repbind;

	packedfile *PackedFile
	//struct PreviewImage * preview;

	lastupdate float32
	lastused   int32
	animspeed  int16

	gen_x, gen_y, gen_type int16
}

type TexType int

func (t TexType) GetTextureTypeDisplayString() string {
	switch t {
	case Type_CLOUDS:
		return "Clouds"
	case Type_WOOD:
		return "Wood"
	case Type_MARBLE:
		return "Marble"
	case Type_MAGIC:
		return "Magic"
	case Type_BLEND:
		return "Blend"
	case Type_STUCCI:
		return "Stucci"
	case Type_NOISE:
		return "Noise"
	case Type_PLUGIN:
		return "Plugin"
	case Type_MUSGRAVE:
		return "Musgrave"
	case Type_VORONOI:
		return "Voronoi"
	case Type_DISTNOISE:
		return "DistortedNoise"
	case Type_ENVMAP:
		return "EnvMap"
	case Type_IMAGE:
		return "Image"
	default:
		break
	}
	return "<Unknown>"
}

const (
	Type_CLOUDS       TexType = 1
	Type_WOOD         TexType = 2
	Type_MARBLE       TexType = 3
	Type_MAGIC        TexType = 4
	Type_BLEND        TexType = 5
	Type_STUCCI       TexType = 6
	Type_NOISE        TexType = 7
	Type_IMAGE        TexType = 8
	Type_PLUGIN       TexType = 9
	Type_ENVMAP       TexType = 10
	Type_MUSGRAVE     TexType = 11
	Type_VORONOI      TexType = 12
	Type_DISTNOISE    TexType = 13
	Type_POINTDENSITY TexType = 14
	Type_VOXELDATA    TexType = 15
)

type TexImageFlags int

const (
	ImageFlags_INTERPOL      TexImageFlags = 1
	ImageFlags_USEALPHA      TexImageFlags = 2
	ImageFlags_MIPMAP        TexImageFlags = 4
	ImageFlags_IMAROT        TexImageFlags = 16
	ImageFlags_CALCALPHA     TexImageFlags = 32
	ImageFlags_NORMALMAP     TexImageFlags = 2048
	ImageFlags_GAUSS_MIP     TexImageFlags = 4096
	ImageFlags_FILTER_MIN    TexImageFlags = 8192
	ImageFlags_DERIVATIVEMAP TexImageFlags = 16384
)

type Tex struct {
	*ElemBase
	id ID
	// AnimData *adt;

	//float noisesize, turbul;
	//float bright, contrast, rfac, gfac, bfac;
	//float filtersize;

	//float mg_H, mg_lacunarity, mg_octaves, mg_offset, mg_gain;
	//float dist_amount, ns_outscale;

	//float vn_w1;
	//float vn_w2;
	//float vn_w3;
	//float vn_w4;
	//float vn_mexp;
	//short vn_distm, vn_coltype;

	//short noisedepth, noisetype;
	//short noisebasis, noisebasis2;

	//short flag;
	imaflag TexImageFlags
	Type    TexType
	//short stype;

	//float cropxmin, cropymin, cropxmax, cropymax;
	//int texfilter;
	//int afmax;
	//short xrepeat, yrepeat;
	//short extend;

	//short fie_ima;
	//int len;
	//int frames, offset, sfra;

	//float checkerdist, nabla;
	//float norfac;

	//ImageUser iuser;

	//bNodeTree *nodetree;
	//Ipo *ipo;
	ima *Image
	//PluginTex *plugin;
	//ColorBand *coba;
	//EnvMap *env;
	//PreviewImage * preview;
	//PointDensity *pd;
	//VoxelData *vd;

	//char use_nodes;

}

func NewTex() *Tex {
	return &Tex{
		imaflag: ImageFlags_INTERPOL,
		Type:    Type_CLOUDS,
	}
}

type MTexFlag int

const (
	Flag_RGBTOINT  MTexFlag = 0x1
	Flag_STENCIL   MTexFlag = 0x2
	Flag_NEGATIVE  MTexFlag = 0x4
	Flag_ALPHAMIX  MTexFlag = 0x8
	Flag_VIEWSPACE MTexFlag = 0x10
)

type MTexProjection int

const (
	Proj_N MTexProjection = 0
	Proj_X MTexProjection = 1
	Proj_Y MTexProjection = 2
	Proj_Z MTexProjection = 3
)

type MTexBlendType int

const (
	BlendType_BLEND       MTexBlendType = 0
	BlendType_MUL         MTexBlendType = 1
	BlendType_ADD         MTexBlendType = 2
	BlendType_SUB         MTexBlendType = 3
	BlendType_DIV         MTexBlendType = 4
	BlendType_DARK        MTexBlendType = 5
	BlendType_DIFF        MTexBlendType = 6
	BlendType_LIGHT       MTexBlendType = 7
	BlendType_SCREEN      MTexBlendType = 8
	BlendType_OVERLAY     MTexBlendType = 9
	BlendType_BLEND_HUE   MTexBlendType = 10
	BlendType_BLEND_SAT   MTexBlendType = 11
	BlendType_BLEND_VAL   MTexBlendType = 12
	BlendType_BLEND_COLOR MTexBlendType = 13
)

type MTexMapType int

const (
	MapType_COL      MTexMapType = 1
	MapType_NORM     MTexMapType = 2
	MapType_COLSPEC  MTexMapType = 4
	MapType_COLMIR   MTexMapType = 8
	MapType_REF      MTexMapType = 16
	MapType_SPEC     MTexMapType = 32
	MapType_EMIT     MTexMapType = 64
	MapType_ALPHA    MTexMapType = 128
	MapType_HAR      MTexMapType = 256
	MapType_RAYMIRR  MTexMapType = 512
	MapType_TRANSLU  MTexMapType = 1024
	MapType_AMB      MTexMapType = 2048
	MapType_DISPLACE MTexMapType = 4096
	MapType_WARP     MTexMapType = 8192
)

type MTex struct {
	*ElemBase
	// short texco, maptoneg;
	mapto MTexMapType

	blendtype MTexBlendType
	object    *Object
	tex       *Tex
	uvname    string

	projx, projy, projz MTexProjection
	mapping             uint8
	ofs, size           [3]float32
	rot                 float32

	texflag                       int32
	colormodel, pmapto, pmaptoneg int16
	//short normapspace, which_output;
	//char brush_map_mode;
	r, g, b, k float32
	//float def_var, rt;

	//float colfac, varfac;

	norfac float32
	//float dispfac, warpfac;
	colspecfac, mirrfac, alphafac      float32
	difffac, specfac, emitfac, hardfac float32
	//float raymirrfac, translfac, ambfac;
	//float colemitfac, colreflfac, coltransfac;
	//float densfac, scatterfac, reflfac;

	//float timefac, lengthfac, clumpfac;
	//float kinkfac, roughfac, padensfac;
	//float lifefac, sizefac, ivelfac, pvelfac;
	//float shadowfac;
	//float zenupfac, zendownfac, blendfac;
}
