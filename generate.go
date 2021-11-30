package main


//go:generate  stringer -type Code  -linecomment ./pkg/errcode
//go:generate sh -c "protoc -I=. -I=./third_party --gogo_out=paths=source_relative:. --validate_out=paths=source_relative,lang=go:. ./configs/*.proto"
//go:generate sh -c "protoc -I=. -I=./third_party --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. --grpc-gateway_out=paths=source_relative,logtostderr=true,generate_unbound_methods=true:. --validate_out=paths=source_relative,lang=go:. ./api/v1/*.proto"
//go:generate  wire ./cmd
