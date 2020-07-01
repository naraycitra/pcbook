FROM golang:1.14.4-alpine3.12 AS build-env

# Allow Go to retrieve the dependencies for the building step
RUN apk add --no-cache git

# Secure against running as root
RUN adduser -D -u 10000 naray
RUN mkdir /pcbook/ && chown naray /pcbook/
USER naray

WORKDIR /pcbook/
ADD . /pcbook/


# Compile the binary, we don't want to run the cgo resolver
RUN CGO_ENABLED=0 go build -o server cmd/server/main.go

# Final stage
FROM alpine:3.12

# Secure againts running as root
RUN adduser -D -u 10000 naray
USER naray

WORKDIR /
COPY --from=build-env /pcbook/server /

#EXPOSE 8080
ENTRYPOINT [ "/server" ]
CMD [ "-port 8080" ]
