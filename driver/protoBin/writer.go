package protoBin

import (
	"encoding/binary"
	"github.com/gorustyt/go-assimp/core"
	"google.golang.org/protobuf/proto"
	"os"
)

func WriteProto(path string, pScene *core.AiScene) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0664)
	if err != nil {
		return err
	}
	data, err := proto.Marshal(pScene.ToPbMsg())
	if err != nil {
		return err
	}
	err = writeMagic(f)
	if err != nil {
		return err
	}
	err = writeData(f, data)
	return err
}

func writeMagic(f *os.File) error {
	_, err := f.Write([]byte(Desc.Magic))
	return err
}

func writeData(f *os.File, data []byte) error {
	var buff [4]byte
	binary.LittleEndian.PutUint32(buff[:], uint32(len(data)))
	_, err := f.Write(buff[:])
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}
