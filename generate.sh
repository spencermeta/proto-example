#!/bin/bash

protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.
protoc calc/calcpb/calc.proto --go_out=plugins=grpc:.