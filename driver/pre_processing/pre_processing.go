package pre_processing

import "assimp/core"

type ScenePreprocessor struct {
	scene *core.AiScene
}

func NewScenePreprocessor(scene *core.AiScene) *ScenePreprocessor {
	return &ScenePreprocessor{scene: scene}
}

func (s *ScenePreprocessor) SetScene(scene *core.AiScene) {
	s.scene = scene
}

func (s *ScenePreprocessor) ProcessScene() {

}

func (s *ScenePreprocessor) ProcessAnimation(anim *core.AiAnimation) {

}

func (s *ScenePreprocessor) ProcessMesh(mesh *core.AiMesh) {

}
