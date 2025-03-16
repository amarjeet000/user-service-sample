FROM golang:1.23-alpine AS base

WORKDIR /service
COPY . .

FROM base AS build

RUN  go build -o userservice .

FROM scratch AS img

COPY --from=build /service/userservice /service

EXPOSE 3030
ENTRYPOINT ["/service"]
