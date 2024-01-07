package protoBin

import (
	"github.com/gorustyt/go-assimp/common/config"
	"github.com/gorustyt/go-assimp/common/pb_msg"
	"github.com/gorustyt/go-assimp/common/reader"
	"github.com/gorustyt/go-assimp/core"
	"github.com/gorustyt/go-assimp/driver/base/iassimp"
	"google.golang.org/protobuf/proto"
)

var (
	Desc = core.AiImporterDesc{
		"Assimp Binary Importer",
		"Gargaj / Conspiracy",
		"",
		"",
		0,
		0,
		0,
		0,
		0,
		[]string{"protobin"},
		"ASSIMP.binaryProto-dump.",
	}
)

type AssProtoImporter struct {
	reader.StreamReader
}

func (ai *AssProtoImporter) CanRead(checkSig bool) bool {
	data, err := ai.Peek(len(Desc.Magic))
	if err != nil {
		return false
	}
	if string(data) != Desc.Magic {
		return false
	}
	return true
}

func NewAssProtoImporter(data []byte) (iassimp.Loader, error) {
	r, err := reader.NewFileStreamReader(data)
	if err != nil {
		return nil, err
	}
	return &AssProtoImporter{StreamReader: r}, nil
}

func (ai *AssProtoImporter) Read(pScene *core.AiScene) (err error) {
	err = ai.Discard(len(Desc.Magic))
	if err != nil {
		return err
	}
	length, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	data, err := ai.GetNBytes(int(length))
	if err != nil {
		return err
	}
	var p pb_msg.AiScene
	err = proto.Unmarshal(data, &p)
	if err != nil {
		return err
	}
	pScene.FromPbMsg(&p)
	return
}

func (ai *AssProtoImporter) InitConfig(cfg *config.Config) {

}
