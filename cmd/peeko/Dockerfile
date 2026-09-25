FROM golang:1.27-alpine AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /peeko ./cmd/peeko

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /peeko /peeko
EXPOSE 8080
ENTRYPOINT ["/peeko"]
