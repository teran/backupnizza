//go:build proto

package proto

// go:generate protoc -I=. --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative task.proto

var (
	_ = (*RunTaskRequest)(nil)
	_ = (*RunTaskResponse)(nil)
)
