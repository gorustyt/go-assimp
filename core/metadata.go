package core

import (
	"github.com/gorustyt/go-assimp/common"
	"github.com/gorustyt/go-assimp/common/pb_msg"
	"google.golang.org/protobuf/proto"
)

type AiMetadataType int

const (
	AI_BOOL AiMetadataType = iota
	AI_INT32
	AI_UINT64
	AI_FLOAT
	AI_DOUBLE
	AI_AISTRING
	AI_AIVECTOR3D
	AI_AIMETADATA
	AI_INT64
	AI_UINT32
	AI_META_MAX
)

/**
 * Container for holding metadata.
 *
 * Metadata is a key-value store using string keys and values.
 */

type AiMetadata struct {
	/** Arrays of keys, may not be NULL. Entries in this array may not be NULL as well. */
	Keys []string

	/** Arrays of values, may not be NULL. Entries in this array may be NULL if the
	 * corresponding property key has no assigned value. */
	Values []*AiMetadataEntry
}

func (ai *AiMetadata) Clone() *AiMetadata {
	if ai == nil {
		return nil
	}
	r := &AiMetadata{}
	r.Keys = ai.Keys
	for _, v := range ai.Values {
		r.Values = append(r.Values, v.Clone())
	}
	return r
}
func (ai *AiMetadata) ToPbMsg() *pb_msg.AiMetadata {
	if ai == nil {
		return nil
	}
	r := pb_msg.AiMetadata{}
	r.Keys = ai.Keys
	for _, v := range ai.Values {
		r.Values = append(r.Values, v.ToPbMsg())
	}
	return &r
}
func (ai *AiMetadata) FromPbMsg(data *pb_msg.AiMetadata) *AiMetadata {
	if data == nil {
		return nil
	}
	ai.Keys = data.Keys
	for _, v := range data.Values {
		ai.Values = append(ai.Values, ((&AiMetadataEntry{}).FromPbMsg(v)))
	}
	return ai
}

/**
 * Metadata entry
 *
 * The type field uniquely identifies the underlying type of the data field
 */

type AiMetadataEntry struct {
	Type AiMetadataType
	Data any
}

func (ai *AiMetadataEntry) Clone() *AiMetadataEntry {
	if ai == nil {
		return nil
	}
	r := &AiMetadataEntry{}
	r.Type = ai.Type

	switch ai.Type {
	case AI_AIVECTOR3D:
		r.Data = ai.Data.(*common.AiVector3D).Clone()
	case AI_AIMETADATA:
		r.Data = ai.Data.(*AiMetadata).Clone()
	default:
		r.Data = ai.Data
	}
	return r
}
func (ai *AiMetadataEntry) ToPbMsg() *pb_msg.AiMetadataEntry {
	if ai == nil {
		return nil
	}
	r := pb_msg.AiMetadataEntry{}
	r.Type = int32(ai.Type)
	var v proto.Message
	switch ai.Type {
	case AI_BOOL:
		b := ai.Data.(bool)
		t := 0
		if b {
			t = 1
		}
		v = &pb_msg.AiMaterialPropertyInt64{Data: []int64{int64(t)}}
	case AI_INT32:
		v = &pb_msg.AiMaterialPropertyInt64{Data: []int64{int64(ai.Data.(int32))}}
	case AI_UINT64:
		v = &pb_msg.AiMaterialPropertyInt64{Data: []int64{int64(ai.Data.(uint64))}}
	case AI_FLOAT:
		v = &pb_msg.AiMaterialPropertyFloat64{Data: []float64{float64(ai.Data.(float32))}}
	case AI_DOUBLE:
		v = &pb_msg.AiMaterialPropertyFloat64{Data: []float64{ai.Data.(float64)}}
	case AI_AISTRING:
		v = &pb_msg.AiMaterialPropertyString{Data: []string{ai.Data.(string)}}
	case AI_AIVECTOR3D:
		v = ai.Data.(*common.AiVector3D).ToPbMsg()
	case AI_AIMETADATA:
		v = ai.Data.(*AiMetadata).ToPbMsg()
	case AI_INT64:
		v = &pb_msg.AiMaterialPropertyInt64{Data: []int64{ai.Data.(int64)}}
	case AI_UINT32:
		v = &pb_msg.AiMaterialPropertyInt64{Data: []int64{int64(ai.Data.(uint32))}}
	}
	r.Data, _ = proto.Marshal(v)
	return &r
}

func (ai *AiMetadataEntry) FromPbMsg(data *pb_msg.AiMetadataEntry) *AiMetadataEntry {
	if data == nil {
		return nil
	}
	ai.Type = AiMetadataType(data.Type)
	switch ai.Type {
	case AI_BOOL:
		tmp := &pb_msg.AiMaterialPropertyInt64{}
		err := proto.Unmarshal(data.Data, tmp)
		if err != nil {
			panic(err)
		}
		if tmp.Data[0] != 0 {
			ai.Data = true
		}
	case AI_INT32:
		tmp := &pb_msg.AiMaterialPropertyInt64{}
		err := proto.Unmarshal(data.Data, tmp)
		if err != nil {
			panic(err)
		}
		ai.Data = int32(tmp.Data[0])
	case AI_UINT64:
		tmp := &pb_msg.AiMaterialPropertyInt64{}
		err := proto.Unmarshal(data.Data, tmp)
		if err != nil {
			panic(err)
		}
		ai.Data = uint64(tmp.Data[0])
	case AI_FLOAT:
		tmp := &pb_msg.AiMaterialPropertyFloat64{}
		err := proto.Unmarshal(data.Data, tmp)
		if err != nil {
			panic(err)
		}
		ai.Data = float32(tmp.Data[0])
	case AI_DOUBLE:
		tmp := &pb_msg.AiMaterialPropertyFloat64{}
		err := proto.Unmarshal(data.Data, tmp)
		if err != nil {
			panic(err)
		}
		ai.Data = tmp.Data[0]
	case AI_AISTRING:
		tmp := &pb_msg.AiMaterialPropertyString{}
		err := proto.Unmarshal(data.Data, tmp)
		if err != nil {
			panic(err)
		}
		ai.Data = tmp.Data[0]
	case AI_AIVECTOR3D:
		tmp := &pb_msg.AiVector3D{}
		err := proto.Unmarshal(data.Data, tmp)
		if err != nil {
			panic(err)
		}
		ai.Data = (&common.AiVector3D{}).FromPbMsg(tmp)
	case AI_AIMETADATA:
		tmp := &pb_msg.AiMetadata{}
		err := proto.Unmarshal(data.Data, tmp)
		if err != nil {
			panic(err)
		}
		ai.Data = (&AiMetadata{}).FromPbMsg(tmp)

	case AI_INT64:
		tmp := &pb_msg.AiMaterialPropertyInt64{}
		err := proto.Unmarshal(data.Data, tmp)
		if err != nil {
			panic(err)
		}

		ai.Data = tmp.Data[0]
	case AI_UINT32:
		tmp := &pb_msg.AiMaterialPropertyInt64{}
		err := proto.Unmarshal(data.Data, tmp)
		if err != nil {
			panic(err)
		}
		ai.Data = uint32(tmp.Data[0])
	}
	return ai
}
