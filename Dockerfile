FROM golang:1.20-buster AS build
# buster based on debian

WORKDIR /app 

COPY . ./
# install dependencies
RUN go mod download 

RUN CGO_ENABLED=0 go build -o /bin/app

# deploy the application binary into a lean image
FROM gcr.io/distroless/static-debian11
# distroless image lightweight than alpine & debian
# FROM debian:buster-slim

COPY --from=build /bin/app /bin
COPY .env.prod /bin

EXPOSE 3000

USER nonroot:nonroot

ENTRYPOINT ["/bin/app", "/bin/.env.prod"]




