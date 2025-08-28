gen_server:
	go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.9.0 -package spec -generate types ../swagger/altron_swagger.yaml > ./core/generated/spec/types.gen.go
	go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.9.0 -package spec -generate spec ../swagger/altron_swagger.yaml > ./core/generated/spec/spec.gen.go
	go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.9.0 -package spec -generate gin ../swagger/altron_swagger.yaml > ./core/generated/spec/gin_server.gen.go
	
	go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.9.0 -package spec -generate types ../swagger/plugin_swagger.yaml > ./plugin/generated/spec/types.gen.go
	go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.9.0 -package spec -generate spec ../swagger/plugin_swagger.yaml > ./plugin/generated/spec/spec.gen.go
	go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.9.0 -package spec -generate gin ../swagger/plugin_swagger.yaml > ./plugin/generated/spec/gin_server.gen.go

	go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.9.0 -package spec -generate types ../swagger/session_swagger.yaml > ./session/generated/spec/types.gen.go
	go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.9.0 -package spec -generate spec ../swagger/session_swagger.yaml > ./session/generated/spec/spec.gen.go
	go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.9.0 -package spec -generate gin ../swagger/session_swagger.yaml > ./session/generated/spec/gin_server.gen.go

	go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.9.0 -package spec -generate types ../swagger/connection_swagger.yaml > ./connection/generated/spec/types.gen.go
	go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.9.0 -package spec -generate spec ../swagger/connection_swagger.yaml > ./connection/generated/spec/spec.gen.go
	go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.9.0 -package spec -generate gin ../swagger/connection_swagger.yaml > ./connection/generated/spec/gin_server.gen.go

	go mod tidy

gen_orm:
	- go run -mod=mod entgo.io/ent/cmd/ent generate ./core/repositories/schema --target ./core/repositories/ent