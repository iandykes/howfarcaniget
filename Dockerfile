FROM golang:1.9 AS builder

ARG VERSION
ARG COMMIT
ARG BUILD_DATE

WORKDIR $GOPATH/src/github.com/iandykes/howfarcaniget
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo "-X 'main.version=${VERSION}' -X 'main.buildDate=${BUILD_DATE}' -X 'main.commitHash=${COMMIT}'" -o /app .

# Could use scratch, but would need certs for https calls to google
FROM alpine
COPY --from=builder /app ./
COPY --from=builder /static ./static/
COPY --from=builder /templates ./templates/

ENTRYPOINT ["./app"]

EXPOSE 80