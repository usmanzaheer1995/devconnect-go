FROM golang:alpine as base
WORKDIR /app

FROM base as dev
COPY go.* ./
RUN go mod download

RUN go env

RUN go get github.com/go-delve/delve/cmd/dlv \
    && go get github.com/githubnemo/CompileDaemon

EXPOSE 5000 2345

FROM dev as test

# 13. Copy the remaining api code into /api in the image's filesystem
COPY . .

# 14. Disable CGO and run unit tests
RUN export CGO_ENABLED=0 && \
    go test -v ./...

# 15. Extend the test stage and create a new stage named build-stage
FROM test as build-stage

# 16. Build the api with "-ldflags" aka linker flags to reduce binary size
# -s = disable symbol table
# -w = disable DWARF generation
RUN GOOS=linux go build -ldflags "-s -w" -o main ./cmd/web

# 17. Extend the base stage and create a new stage named prod
FROM base as prod

# 18. Copy only the files we want from a few stages into the prod stage
COPY --from=build-stage /cmd/web main

# 19. Create a new group and user, recursively change directory ownership, then allow the binary to be executed
RUN addgroup gopher && adduser -D -G gopher gopher \
    && chown -R gopher:gopher /api && \
    chmod +x ./main

# 20. Change to a non-root user
USER gopher

# 21. Provide meta data about the port the container must expose
EXPOSE 5000

# 22. Define how Docker should test the container to check that it is still working
HEALTHCHECK CMD [ "wget", "-q", "0.0.0.0:5000" ]

# 23. Provide the default command for the production container
CMD ["./main"]
