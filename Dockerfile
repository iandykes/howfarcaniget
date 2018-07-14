FROM golang:1.9 AS builder

ARG VERSION
ARG COMMIT

WORKDIR /howfarcaniget
COPY . ./
RUN VERSION=${VERSION} COMMIT=${COMMIT} ./build-linux.sh 

# Could use scratch, but would need certs for https calls to google
FROM alpine
COPY --from=builder /howfarcaniget/output /howfarcaniget/

ENTRYPOINT ["/howfarcaniget/app"]

EXPOSE 80