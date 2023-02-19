# syntax=docker/dockerfile:1
FROM --platform=${BUILDPLATFORM} golang:1.20.1 as compiler
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go install .

FROM --platform=${BUILDPLATFORM} gcr.io/distroless/static
COPY --from=compiler /go/bin/nginx-auth /usr/local/bin/
EXPOSE 8080
USER 65532
ENTRYPOINT [ "nginx-auth" ]
