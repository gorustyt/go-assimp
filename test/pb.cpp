//
// Created by Administrator on 2023/12/20.
//

#include "pb.h"
#include "common/pb_msg/ai_scene.pb.h"
#include <assimp/scene.h>
#include <fstream>
#include <iostream>
using namespace std;
#include "google/protobuf/io/zero_copy_stream_impl_lite.h"  // ArrayOutputStream
#include "google/protobuf/io/zero_copy_stream_impl.h"   // FileOutputStream    OstreamOutputStream
#include "google/protobuf/text_format.h"
#include "google/protobuf/util/delimited_message_util.h"

void pbWriteFile(pb_msg::AiScene *pbscene,char * fileName){
             // 写文件
            std::fstream output(fileName, ios::out | ios::binary);
            if (!output)
            {
                cout<< "output" << fileName << " error"<<endl;
            }
            pbscene->SerializeToOstream(&output);
            //PB对象转文件流
            output.close();
    return;
}

char* PbConv(char * fileName){
    pb_msg::AiScene pbscene;
    pbWriteFile(&pbscene,"./xx.bin");
    return nullptr;
}
