# Основной путь к исходникам
MAIN_DIR := ./cmd/main
# Имя бинарника

BINARY_NAME := calc_service

GO := go

PROTO_DIR := ./proto
OUT_DIR := ./pkg/api

PROTOC := protoc
PROTO_FILES := ${PROTO_DIR}/calc/calc.proto

start: build run 
build: 
	@echo "Building the binary..."
	${GO} build -o${BINARY_NAME} ${MAIN_DIR}
run:	 build
	@echo "Running the server..."
	./${BINARY_NAME}
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	