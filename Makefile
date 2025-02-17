.PHONY: proto clean deps

# Переменные
PROTO_DIR = proto
GEN_DIR = proto/gen

# Установка зависимостей
deps:
	go mod init water_delivery
	go mod tidy
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Генерация proto файлов
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto

# Очистка сгенерированных файлов
clean:
	rm -rf $(GEN_DIR)

# Запуск сервера
run:
	go run server/main.go

# Установка всех необходимых зависимостей
install: deps proto 