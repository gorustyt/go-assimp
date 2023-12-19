// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.1
// source: common/pb_msg/ai_texture.proto

package pb_msg

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// * Helper structure to describe an embedded texture
//
// Normally textures are contained in external files but some file formats embed
// them directly in the model file. There are two types of embedded textures:
// 1. Uncompressed textures. The color data is given in an uncompressed format.
// 2. Compressed textures stored in a file format like png or jpg. The raw file
// bytes are given so the application must utilize an image decoder (e.g. DevIL) to
// get access to the actual color data.
//
// Embedded textures are referenced from materials using strings like "*0", "*1", etc.
// as the texture paths (a single asterisk character followed by the
// zero-based index of the texture in the aiScene::mTextures array).
type AiTexture struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// * Width of the texture, in pixels
	//
	// If mHeight is zero the texture is compressed in a format
	// like JPEG. In this case mWidth specifies the size of the
	// memory area pcData is pointing to, in bytes.
	Width uint32 `protobuf:"varint,1,opt,name=Width,proto3" json:"Width,omitempty"`
	// * Height of the texture, in pixels
	//
	// If this value is zero, pcData points to an compressed texture
	// in any format (e.g. JPEG).
	Height uint32 `protobuf:"varint,2,opt,name=Height,proto3" json:"Height,omitempty"`
	//   - A hint from the loader to make it easier for applications
	//     to determine the type of embedded textures.
	//
	// If mHeight != 0 this member is show how data is packed. Hint will consist of
	// two parts: channel order and channel bitness (count of the bits for every
	// color channel). For simple parsing by the viewer it's better to not omit
	// absent color channel and just use 0 for bitness. For example:
	// 1. Image contain RGBA and 8 bit per channel, achFormatHint == "rgba8888";
	// 2. Image contain ARGB and 8 bit per channel, achFormatHint == "argb8888";
	// 3. Image contain RGB and 5 bit for R and B channels and 6 bit for G channel, achFormatHint == "rgba5650";
	// 4. One color image with B channel and 1 bit for it, achFormatHint == "rgba0010";
	// If mHeight == 0 then achFormatHint is set set to '\\0\\0\\0\\0' if the loader has no additional
	// information about the texture file format used OR the
	// file extension of the format without a trailing dot. If there
	// are multiple file extensions for a format, the shortest
	// extension is chosen (JPEG maps to 'jpg', not to 'jpeg').
	// E.g. 'dds\\0', 'pcx\\0', 'jpg\\0'.  All characters are lower-case.
	// The fourth character will always be '\\0'.
	AchFormatHint []byte `protobuf:"bytes,3,opt,name=AchFormatHint,proto3" json:"AchFormatHint,omitempty"` // 8 for string + 1 for terminator.
	// * Data of the texture.
	//
	// Points to an array of mWidth * mHeight aiTexel's.
	// The format of the texture data is always ARGB8888 to
	// make the implementation for user of the library as easy
	// as possible. If mHeight = 0 this is a pointer to a memory
	// buffer of size mWidth containing the compressed texture
	// data. Good luck, have fun!
	PcData []*AiTexel `protobuf:"bytes,4,rep,name=pcData,proto3" json:"pcData,omitempty"`
	// * Texture original filename
	//
	// Used to get the texture reference
	Filename string `protobuf:"bytes,5,opt,name=Filename,proto3" json:"Filename,omitempty"`
}

func (x *AiTexture) Reset() {
	*x = AiTexture{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_pb_msg_ai_texture_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AiTexture) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AiTexture) ProtoMessage() {}

func (x *AiTexture) ProtoReflect() protoreflect.Message {
	mi := &file_common_pb_msg_ai_texture_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AiTexture.ProtoReflect.Descriptor instead.
func (*AiTexture) Descriptor() ([]byte, []int) {
	return file_common_pb_msg_ai_texture_proto_rawDescGZIP(), []int{0}
}

func (x *AiTexture) GetWidth() uint32 {
	if x != nil {
		return x.Width
	}
	return 0
}

func (x *AiTexture) GetHeight() uint32 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *AiTexture) GetAchFormatHint() []byte {
	if x != nil {
		return x.AchFormatHint
	}
	return nil
}

func (x *AiTexture) GetPcData() []*AiTexel {
	if x != nil {
		return x.PcData
	}
	return nil
}

func (x *AiTexture) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

type AiTexel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	B uint32 `protobuf:"varint,1,opt,name=B,proto3" json:"B,omitempty"`
	G uint32 `protobuf:"varint,2,opt,name=G,proto3" json:"G,omitempty"`
	R uint32 `protobuf:"varint,3,opt,name=R,proto3" json:"R,omitempty"`
	A uint32 `protobuf:"varint,4,opt,name=A,proto3" json:"A,omitempty"`
}

func (x *AiTexel) Reset() {
	*x = AiTexel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_pb_msg_ai_texture_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AiTexel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AiTexel) ProtoMessage() {}

func (x *AiTexel) ProtoReflect() protoreflect.Message {
	mi := &file_common_pb_msg_ai_texture_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AiTexel.ProtoReflect.Descriptor instead.
func (*AiTexel) Descriptor() ([]byte, []int) {
	return file_common_pb_msg_ai_texture_proto_rawDescGZIP(), []int{1}
}

func (x *AiTexel) GetB() uint32 {
	if x != nil {
		return x.B
	}
	return 0
}

func (x *AiTexel) GetG() uint32 {
	if x != nil {
		return x.G
	}
	return 0
}

func (x *AiTexel) GetR() uint32 {
	if x != nil {
		return x.R
	}
	return 0
}

func (x *AiTexel) GetA() uint32 {
	if x != nil {
		return x.A
	}
	return 0
}

var File_common_pb_msg_ai_texture_proto protoreflect.FileDescriptor

var file_common_pb_msg_ai_texture_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x70, 0x62, 0x5f, 0x6d, 0x73, 0x67, 0x2f,
	0x61, 0x69, 0x5f, 0x74, 0x65, 0x78, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x06, 0x70, 0x62, 0x5f, 0x6d, 0x73, 0x67, 0x22, 0xa4, 0x01, 0x0a, 0x09, 0x41, 0x69, 0x54,
	0x65, 0x78, 0x74, 0x75, 0x72, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x57, 0x69, 0x64, 0x74, 0x68, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x57, 0x69, 0x64, 0x74, 0x68, 0x12, 0x16, 0x0a, 0x06,
	0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x48, 0x65,
	0x69, 0x67, 0x68, 0x74, 0x12, 0x24, 0x0a, 0x0d, 0x41, 0x63, 0x68, 0x46, 0x6f, 0x72, 0x6d, 0x61,
	0x74, 0x48, 0x69, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0d, 0x41, 0x63, 0x68,
	0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x48, 0x69, 0x6e, 0x74, 0x12, 0x27, 0x0a, 0x06, 0x70, 0x63,
	0x44, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x62, 0x5f,
	0x6d, 0x73, 0x67, 0x2e, 0x41, 0x69, 0x54, 0x65, 0x78, 0x65, 0x6c, 0x52, 0x06, 0x70, 0x63, 0x44,
	0x61, 0x74, 0x61, 0x12, 0x1a, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x22,
	0x41, 0x0a, 0x07, 0x41, 0x69, 0x54, 0x65, 0x78, 0x65, 0x6c, 0x12, 0x0c, 0x0a, 0x01, 0x42, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x01, 0x42, 0x12, 0x0c, 0x0a, 0x01, 0x47, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x01, 0x47, 0x12, 0x0c, 0x0a, 0x01, 0x52, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x01, 0x52, 0x12, 0x0c, 0x0a, 0x01, 0x41, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x01, 0x41, 0x42, 0x0f, 0x5a, 0x0d, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x70, 0x62, 0x5f,
	0x6d, 0x73, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_pb_msg_ai_texture_proto_rawDescOnce sync.Once
	file_common_pb_msg_ai_texture_proto_rawDescData = file_common_pb_msg_ai_texture_proto_rawDesc
)

func file_common_pb_msg_ai_texture_proto_rawDescGZIP() []byte {
	file_common_pb_msg_ai_texture_proto_rawDescOnce.Do(func() {
		file_common_pb_msg_ai_texture_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_pb_msg_ai_texture_proto_rawDescData)
	})
	return file_common_pb_msg_ai_texture_proto_rawDescData
}

var file_common_pb_msg_ai_texture_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_common_pb_msg_ai_texture_proto_goTypes = []interface{}{
	(*AiTexture)(nil), // 0: pb_msg.AiTexture
	(*AiTexel)(nil),   // 1: pb_msg.AiTexel
}
var file_common_pb_msg_ai_texture_proto_depIdxs = []int32{
	1, // 0: pb_msg.AiTexture.pcData:type_name -> pb_msg.AiTexel
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_common_pb_msg_ai_texture_proto_init() }
func file_common_pb_msg_ai_texture_proto_init() {
	if File_common_pb_msg_ai_texture_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_pb_msg_ai_texture_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AiTexture); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_pb_msg_ai_texture_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AiTexel); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_common_pb_msg_ai_texture_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_pb_msg_ai_texture_proto_goTypes,
		DependencyIndexes: file_common_pb_msg_ai_texture_proto_depIdxs,
		MessageInfos:      file_common_pb_msg_ai_texture_proto_msgTypes,
	}.Build()
	File_common_pb_msg_ai_texture_proto = out.File
	file_common_pb_msg_ai_texture_proto_rawDesc = nil
	file_common_pb_msg_ai_texture_proto_goTypes = nil
	file_common_pb_msg_ai_texture_proto_depIdxs = nil
}