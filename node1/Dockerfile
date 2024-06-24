# Establecer la imagen base
FROM golang:1.17-alpine AS builder

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar los archivos necesarios al contenedor
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compilar la aplicación
RUN go build -o main .

# Etapa de producción
FROM alpine:latest

WORKDIR /root/

# Copiar el binario compilado desde la etapa de compilación
COPY --from=builder /app/main .

# Exponer el puerto en el que escucha el servidor Go
EXPOSE 8000

# Comando para ejecutar la aplicación
CMD ["./main"]
