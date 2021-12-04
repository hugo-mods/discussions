FROM golang:1.17-alpine AS build

# RUN apk add --no-cache make

WORKDIR /go/src/godisc/

COPY . .
RUN go build -o /bin/godisc .

FROM scratch
COPY --from=build /bin/godisc ./godisc
ENTRYPOINT ["./godisc"]