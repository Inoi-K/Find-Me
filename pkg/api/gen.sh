protoc --go_out=pkg/api \
    --go-grpc_out=pkg/api \
    pkg/api/proto/"$1".proto
