cd pdfiles
protoc --go_out=plugins=grpc:../services/ *.proto
cd ..