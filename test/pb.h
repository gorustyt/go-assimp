//
// Created by Administrator on 2023/12/20.
//

#ifndef ASSIMP_PB_H
#define ASSIMP_PB_H
#include "common/pb_msg/ai_scene.pb.h"
class PB {
};

#endif // ASSIMP_PB_H
void pbWriteFile(pb_msg::AiScene *pbscene,char * fileName);
char* PbConv(char * fileName);