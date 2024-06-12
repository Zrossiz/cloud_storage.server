# Указываем базовый образ
FROM golang:1.21.3-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum файлы
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем остальные файлы проекта
COPY . .

# Собираем приложение
RUN go build -o /goapp cmd/app/main.go

# Указываем команду запуска контейнера
CMD ["/goapp"]
