FROM golang:alpine

# Install git
RUN apk update && apk add --no-cache git

WORKDIR /app

# Define build arguments
ARG META_TITLE
ARG META_DESCRIPTION
ARG META_HOST
ARG META_BASEPATH
ARG META_VERSION
ARG GIN_HOST
ARG GIN_PORT
ARG GIN_MODE
ARG GIN_TIMEOUT
ARG GIN_SHUTDOWNTIMEOUT
ARG GIN_LOGREQUEST
ARG GIN_LOGRESPONSE
ARG GIN_CORS_MODE
ARG SQL_HOST
ARG SQL_USERNAME
ARG SQL_PASSWORD
ARG SQL_PORT
ARG SQL_DATABASE
ARG MIDTRANS_SERVERKEY

# Create the config directory
RUN mkdir -p /etc/cfg

# Create config.json file from build arguments
RUN printf '{\n  "Meta": {\n    "Title": "%s",\n    "Description": "%s",\n    "Host": "%s",\n    "Basepath": "%s",\n    "Version": "%s"\n  },\n  "Gin": {\n    "Host": "%s",\n    "Port": "%s",\n    "Mode": "%s",\n    "Timeout": "%s",\n    "ShutdownTimeout": "%s",\n    "LogRequest": "%s",\n    "LogResponse": "%s",\n    "CORS": {\n      "Mode": "%s"\n    }\n  },\n  "SQL": {\n    "Host": "%s",\n    "Username": "%s",\n    "Password": "%s",\n    "Port": "%s",\n    "Database": "%s"\n  },\n  "Midtrans": {\n    "ServerKey": "%s"\n  }\n}' "$META_TITLE" "$META_DESCRIPTION" "$META_HOST" "$META_BASEPATH" "$META_VERSION" "$GIN_HOST" "$GIN_PORT" "$GIN_MODE" "$GIN_TIMEOUT" "$GIN_SHUTDOWNTIMEOUT" "$GIN_LOGREQUEST" "$GIN_LOGRESPONSE" "$GIN_CORS_MODE" "$SQL_HOST" "$SQL_USERNAME" "$SQL_PASSWORD" "$SQL_PORT" "$SQL_DATABASE" "$MIDTRANS_SERVERKEY" > /etc/cfg/config.json

# Copy the Go application files
COPY . .

# install swagg
RUN go install github.com/swaggo/swag/cmd/swag@v1.6.7

#run swag init & move to path docs/swagger
RUN `go env GOPATH`/bin/swag init -g src/cmd/main.go -o docs/swagger --parseInternal

# Run go mod tidy to clean up the dependencies
RUN go mod tidy

# Build the Go application
RUN go build -o binary ./src/cmd

# Set the entry point to the compiled binary
ENTRYPOINT ["/app/binary"]
