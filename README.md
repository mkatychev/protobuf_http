Example HTTP Notebook server

## Setup

* `brew install clang-format protoc-gen-go protobuf` to setup on macOS for developmen
* `./proto-format.sh` will autoformat the proto
* to generate/regenerate `*.pb.go` files: `protoc notebook.proto --go_out=.`
