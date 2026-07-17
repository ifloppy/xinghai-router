# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS router-build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/xinghai-router ./cmd/router

FROM alpine:3.22 AS router
RUN addgroup -S router && adduser -S -G router router
COPY --from=router-build /out/xinghai-router /usr/local/bin/xinghai-router
USER router
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/xinghai-router"]

FROM node:22-alpine AS web-dependencies
WORKDIR /src/web
COPY web/package.json web/package-lock.json ./
RUN npm ci

FROM web-dependencies AS web-build
COPY web ./
RUN npm run build

FROM node:22-alpine AS web
WORKDIR /app
ENV NODE_ENV=production
ENV HOST=0.0.0.0
ENV PORT=3000
COPY --from=web-build /src/web/.output ./.output
RUN addgroup -S web && adduser -S -G web web
USER web
EXPOSE 3000
CMD ["node", ".output/server/index.mjs"]
