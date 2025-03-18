FROM golang:latest

WORKDIR /app
COPY . .

RUN apt-get update && apt-get install -y ffmpeg
RUN go mod download

CMD ["go", "run", "main.go"]
