# Основной путь к исходникам
MAIN_DIR := ./cmd/main
# Имя бинарника

BINARY_NAME := calc_service

GO := go


OUT_DIR := ./internal/pkg/api

PROTOC := C:\Users\Nikita\AppData\Local\Microsoft\WinGet\Packages\Google.Protobuf_Microsoft.Winget.Source_8wekyb3d8bbwe
PROTO_FILES := ./proto/calc
generate:
	protoc \
	--go_out $(OUT_DIR) --go_opt paths=source_relative \
    --go-grpc_out $(OUT_DIR) --go-grpc_opt paths=source_relative \
    --proto_path=${PROTO_FILES} \
    --proto_path=${PROTOC}/include \
    ${PROTO_FILES}/calc.proto
start: build run
build: 
	@echo "Building the binary..."
	${GO} build -o ${BINARY_NAME} ${MAIN_DIR}
run:	 build
	@echo "Running the server..."
	./${BINARY_NAME}
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	