//
// Created by go on 2023/12/20.
//
//g++ -I ./google/protobuf/include person.pb.cc  main.cpp  -L ./common/lib  -lprotobuf

#include "pb.h"
int  main(){
    PbConv("hello");
    return 0;
}