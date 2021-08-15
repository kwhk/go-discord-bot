# Alpine is chosen for its small footprint
# compared to Ubuntu
FROM arm32v7/golang:alpine
ENV GOOS=linux \
    GOARCH=arm

WORKDIR /app

# Download Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy all project files to image
COPY . .

# Compile our Go application
RUN go build -o ./build/go_discord_bot ./cmd/go_discord_bot

CMD ["./build/go_discord_bot"]