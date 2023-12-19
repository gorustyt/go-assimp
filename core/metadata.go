package core

import "assimp/common/pb_msg"

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

func (ai *AiMetadata) ToPbMsg() *pb_msg.AiMetadata {
	r := pb_msg.AiMetadata{}
	r.Keys = ai.Keys
	for _, v := range ai.Values {
		r.Values = append(r.Values, v.ToPbMsg())
	}
	return &r
}

/**
 * Metadata entry
 *
 * The type field uniquely identifies the underlying type of the data field
 */

type AiMetadataEntry struct {
	Type AiMetadataType
	Data []byte
}

func (ai *AiMetadataEntry) ToPbMsg() *pb_msg.AiMetadataEntry {
	r := pb_msg.AiMetadataEntry{}
	r.Type = int32(ai.Type)
	r.Data = ai.Data
	return &r
}
