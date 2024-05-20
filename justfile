default:
    just --list

proto-gen:
    protoc --go_out=wadoh-be --go-grpc_out=wadoh-be --proto_path=wadoh-be/proto wadoh.proto
