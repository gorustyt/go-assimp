package test

import (
	"github.com/gorustyt/go-assimp"
	"github.com/gorustyt/go-assimp/common/pb_msg"
	"github.com/gorustyt/go-assimp/core"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestProto(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_nonbsd_data/BLEND/fleurOptonl.blend")
	AssertError(t, err)
	pv1 := p.ToPbMsg()
	data, err := proto.Marshal(pv1)
	AssertError(t, err)
	var pv2 pb_msg.AiScene
	err = proto.Unmarshal(data, &pv2)
	AssertError(t, err)
	Assert(t, DeepEqual(p, (&core.AiScene{}).FromPbMsg(&pv2)))
}

func TestWriteProto(t *testing.T) {
	err := assimp.ParseToProtoFile("../example/example_nonbsd_data/BLEND/fleurOptonl.blend",
		"../example/example_data/protoBin/fleurOptonl.protobin")
	AssertError(t, err)
}

func TestReadProto(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_nonbsd_data/BLEND/fleurOptonl.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/protoBin/fleurOptonl.protobin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}
