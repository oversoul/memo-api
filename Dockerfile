FROM golang:1.23 as build

COPY . /src

WORKDIR /src

RUN make build

FROM scratch

COPY --from=build /src/.env .
COPY --from=build /src/go-app .

RUN mkdir -p /public/images

EXPOSE 8000

CMD ["./go-app"]
