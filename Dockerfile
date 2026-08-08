FROM golang:1.26.5-alpine as build
ADD ./ ./
RUN go build main.go

FROM alpine:latest
COPY --from=build /go/main /app/main
WORKDIR /app
CMD [ "./main" ]
