############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add git
# Create appuser.
RUN adduser -D -g '' appuser
# WORKDIR $GOPATH/src/mypackage/myapp/
WORKDIR /myapp
ENV GO111MODULE=on
COPY go.mod go.sum ./

RUN go mod download

COPY . /myapp

RUN ls $GOPATH/src

# Build the binary.
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build -ldflags="-w -s" -o /go/bin/myapp week2/main.go

############################
# STEP 2 build a small image
############################
# FROM scratch
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
# Import the user and group files from the builder.
# COPY --from=builder /etc/passwd /etc/passwd

WORKDIR /root/
# Copy our static executable.
COPY --from=builder /go/bin/myapp .
COPY --from=builder /myapp/.env .
# COPY --from=builder /myapp/Makefile ./
# COPY --from=builder /myapp/week2/script.js ./

# Use an unprivileged user.
# USER appuser

# Run the hello binary. 
EXPOSE 8080

# CMD exec /bin/sh -c "trap : TERM INT; (while true; do sleep 1000; done) & wait"
ENTRYPOINT ["./myapp"]
# CMD -grpc-port=$GRPC_PORT
# CMD ["sh", "-c", "/root/myapp -grpc-port=$GRPC_PORT"]
