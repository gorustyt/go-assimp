package BLEND

import (
	"assimp/common/logger"
	"assimp/common/reader"
	"assimp/core"
	"assimp/driver/base/iassimp"
	"errors"
	"io"
	"sort"
)

var Desc = core.AiImporterDesc{
	"Blender 3D Importer (http://www.blender3d.org)",
	"",
	"",
	"No animation support yet",
	core.AiImporterFlags_SupportBinaryFlavour,
	0,
	0,
	2,
	50,
	[]string{".blend"},
	"BLENDER",
}

type BlenderImporter struct {
	reader.StreamReader
}

func (b *BlenderImporter) checkMagic() ([]byte, bool, error) {
	magic, err := b.Peek(7)
	if err != nil {
		return magic, false, err
	}
	if string(magic[:]) == Desc.Magic {
		return magic, true, nil
	}
	return magic, false, nil
}

func (b *BlenderImporter) ParseMagic() error {
	magic, ok, err := b.checkMagic()
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	// Check for presence of the gzip header. If yes, assume it is a
	// compressed blend file and try uncompressing it, else fail. This is to
	// avoid uncompressing random files which our loader might end up with.
	if magic[0] != 0x1f || magic[1] != 0x8b {
		return errors.New("BLENDER magic bytes are missing, couldn't find GZIP header either")
	}

	logger.Info("Found no BLENDER magic word but a GZIP header, might be a compressed file")
	if magic[2] != 8 {
		return errors.New("Unsupported GZIP compression method")
	}
	err = b.ResetGzipReader()
	if err != nil {
		return err
	}
	magic, ok, err = b.checkMagic()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("magic not found ")
	}
	return nil
}

func (b *BlenderImporter) CanRead(checkSig bool) bool {
	return b.ParseMagic() == nil
}

func (b *BlenderImporter) Read(pScene *core.AiScene) (err error) {
	err = b.ParseMagic()
	if err != nil {
		return err
	}
	err = b.Discard(7)
	if err != nil {
		return err
	}
	file := NewFileDatabase(b.StreamReader)
	buffer, err := b.GetNBytes(1)
	if err != nil {
		return err
	}
	file.i64bit = buffer[0] == '-'
	buffer, err = b.GetNBytes(1)
	if err != nil {
		return err
	}
	file.little = buffer[0] == 'v'
	buffer, err = b.GetNBytes(3)
	if err != nil {
		return err
	}
	b.ChangeBytesOrder(file.little)
	logger.InfoF("Blender version is:%v.%v (64bit:%v , little endian:%v)",
		string(buffer[0]), uint8(buffer[0])+1,
		file.i64bit,
		file.little)

	err = b.ParseBlendFile(file)
	if err != nil {
		return err
	}
	var scene Scene
	err = b.ExtractScene(&scene, file)
	if err != nil {
		return err
	}
	return b.ConvertBlendFile(pScene, &scene, file)
}

func (b *BlenderImporter) ExtractScene(out *Scene, file *FileDatabase) error {
	var block *FileBlockHead
	it, ok := file.dna.indices["Scene"]
	if !ok {
		return errors.New("there is no `Scene` structure record")
	}

	ss := file.dna.structures[it]

	// we need a scene somewhere to start with.
	for _, bl := range file.entries {

		// Fix: using the DNA index is more reliable to locate scenes
		//if (bl.id == "SC") {

		if int(bl.dna_index) == it {
			block = bl
			break
		}
	}

	if block == nil {
		return errors.New("there is not a single `Scene` record to load")
	}
	file.SetCurPos(block.start)
	err := ss.Convert(out, file)
	if err != nil {
		return err
	}

	logger.InfoF(
		"(Stats) Fields read:%v, pointers resolved:%v, cache hits: %v, cached objects: ", file.stats().fields_read,
		file.stats().pointers_resolved,
		file.stats().cache_hits,
		file.stats().cached_objects)
	return nil
}

func (b *BlenderImporter) ConvertBlendFile(out *core.AiScene, in *Scene, file *FileDatabase) error {
	return nil
}
func (b *BlenderImporter) ParseBlendFile(file *FileDatabase) error {
	dna_reader := NewDNAParser(file, b.StreamReader)
	var dna *DNA
	parser := NewSectionParser(b.StreamReader, file.i64bit)
	for {
		err := parser.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		head := parser.current
		if head.id == "ENDB" {
			break // only valid end of the file
		} else if head.id == "DNA1" {
			err = dna_reader.Parse()
			if err != nil {
				return err
			}
			dna = dna_reader.GetDNA()
			continue
		}
		if err != nil {
			return nil
		}
		file.entries = append(file.entries, head)
	}
	if dna == nil {
		return errors.New("SDNA not found")
	}
	sort.Slice(file.entries, func(i, j int) bool {
		return file.entries[i].address.val < file.entries[j].address.val
	})
	return nil
}
func NewBlenderImporter(data []byte) (iassimp.Loader, error) {
	r, err := reader.NewFileStreamReader(data)
	if err != nil {
		return nil, err
	}
	return &BlenderImporter{StreamReader: r}, nil
}
