t := transaction

rpc proto:
	goctl rpc proto -src rpc/$(t).proto -dir rpc --home rpc/template/1.2.4-cli -consul grpc

