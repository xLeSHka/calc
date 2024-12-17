# Основной путь к исходникам
MAIN_DIR := ./cmd/main
# Имя бинарника

BINARY_NAME := main

GO := go

PROTO_DIR := ./proto
OUT_DIR := ./pkg/api

PROTOC := protoc
PROTO_FILES := ${PROTO_DIR}/calc/calc.proto

start: generate
generate: ${PROTO_FILES}
	$(PROTOC)  \
    --go_out $(OUT_DIR) --go_opt paths=source_relative \
    --go-grpc_out $(OUT_DIR) --go-grpc_opt paths=source_relative \
	--grpc-gateway_out ${OUT_DIR}  --grpc-gateway_opt paths=source_relative \
	--swagger_out ${OUT_DIR}  \
	--proto_path=${PROTO_DIR} ${PROTO_FILES}
build: 
	@echo "Building the binary..."
	${GO} build -o	${BINARY_NAME} ${MAIN_DIR}
run:	 build
	@echo "Running the server..."
	./${BINARY_NAME}
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	