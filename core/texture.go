package core

import "assimp/common/pb_msg"

type AiTexture struct {
	/** Width of the texture, in pixels
	 *
	 * If mHeight is zero the texture is compressed in a format
	 * like JPEG. In this case mWidth specifies the size of the
	 * memory area pcData is pointing to, in bytes.
	 */
	Width uint32

	/** Height of the texture, in pixels
	 *
	 * If this value is zero, pcData points to an compressed texture
	 * in any format (e.g. JPEG).
	 */
	Height uint32

	/** A hint from the loader to make it easier for applications
	 *  to determine the type of embedded textures.
	 *
	 * If mHeight != 0 this member is show how data is packed. Hint will consist of
	 * two parts: channel order and channel bitness (count of the bits for every
	 * color channel). For simple parsing by the viewer it's better to not omit
	 * absent color channel and just use 0 for bitness. For example:
	 * 1. Image contain RGBA and 8 bit per channel, achFormatHint == "rgba8888";
	 * 2. Image contain ARGB and 8 bit per channel, achFormatHint == "argb8888";
	 * 3. Image contain RGB and 5 bit for R and B channels and 6 bit for G channel, achFormatHint == "rgba5650";
	 * 4. One color image with B channel and 1 bit for it, achFormatHint == "rgba0010";
	 * If mHeight == 0 then achFormatHint is set set to '\\0\\0\\0\\0' if the loader has no additional
	 * information about the texture file format used OR the
	 * file extension of the format without a trailing dot. If there
	 * are multiple file extensions for a format, the shortest
	 * extension is chosen (JPEG maps to 'jpg', not to 'jpeg').
	 * E.g. 'dds\\0', 'pcx\\0', 'jpg\\0'.  All characters are lower-case.
	 * The fourth character will always be '\\0'.
	 */
	AchFormatHint []uint8 // 8 for string + 1 for terminator.

	/** Data of the texture.
	 *
	 * Points to an array of mWidth * mHeight aiTexel's.
	 * The format of the texture data is always ARGB8888 to
	 * make the implementation for user of the library as easy
	 * as possible. If mHeight = 0 this is a pointer to a memory
	 * buffer of size mWidth containing the compressed texture
	 * data. Good luck, have fun!
	 */
	PcData []*AiTexel

	/** Texture original filename
	 *
	 * Used to get the texture reference
	 */
	Filename string
}

func (ai *AiTexture) FromPbMsg(p *pb_msg.AiTexture) *AiTexture {
	if p == nil {
		return nil
	}
	ai.Width = p.Width
	ai.Height = p.Height
	ai.AchFormatHint = p.AchFormatHint
	for _, v := range p.PcData {
		ai.PcData = append(ai.PcData, (&AiTexel{}).FromPbMsg(v))
	}
	ai.Filename = p.Filename
	return ai
}

func (ai *AiTexture) Clone() *AiTexture {
	if ai == nil {
		return nil
	}
	r := &AiTexture{}
	r.Width = ai.Width
	r.Height = ai.Height
	r.AchFormatHint = ai.AchFormatHint
	for _, v := range ai.PcData {
		r.PcData = append(r.PcData, v.Clone())
	}
	r.Filename = ai.Filename
	return r
}

func (ai *AiTexture) ToPbMsg() *pb_msg.AiTexture {
	if ai == nil {
		return nil
	}
	r := &pb_msg.AiTexture{}
	r.Width = ai.Width
	r.Height = ai.Height
	r.AchFormatHint = ai.AchFormatHint
	for _, v := range ai.PcData {
		r.PcData = append(r.PcData, v.ToPbMsg())
	}
	r.Filename = ai.Filename
	return r
}
func NewAiTexture() *AiTexture {
	return &AiTexture{
		AchFormatHint: make([]byte, HINTMAXTEXTURELEN),
	}
}

type AiTexel struct {
	B, G, R, A uint8
}

func (ai *AiTexel) FromPbMsg(p *pb_msg.AiTexel) *AiTexel {
	if p == nil {
		return nil
	}
	ai.R = uint8(p.R)
	ai.G = uint8(p.G)
	ai.B = uint8(p.B)
	ai.A = uint8(p.A)
	return ai
}
func (ai *AiTexel) Clone() *AiTexel {
	if ai == nil {
		return nil
	}
	return &AiTexel{
		B: ai.B,
		G: ai.G,
		R: ai.R,
		A: ai.A,
	}
}

func (ai *AiTexel) ToPbMsg() *pb_msg.AiTexel {
	if ai == nil {
		return nil
	}
	return &pb_msg.AiTexel{
		B: uint32(ai.B),
		G: uint32(ai.G),
		R: uint32(ai.R),
		A: uint32(ai.A),
	}
}

func BytesToAiTexel(bytes []byte) (res []*AiTexel) {
	for i := 0; i < len(bytes); {
		var v AiTexel
		res = append(res, &v)
		if i < len(bytes) {
			v.B = bytes[i]
		}
		i++

		if i < len(bytes) {
			v.G = bytes[i]
		}
		i++

		if i < len(bytes) {
			v.R = bytes[i]
		}
		i++

		if i < len(bytes) {
			v.A = bytes[i]
		}
		i++
	}
	return res
}
