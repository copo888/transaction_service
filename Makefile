t := transaction

rpc proto:
	goctl rpc protoc rpc/$(t).proto --go_out=rpc --go-grpc_out=rpc --zrpc_out=rpc --home rpc/template/1.2.4-cli


