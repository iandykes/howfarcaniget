FROM golang as builder
ADD . /go/src/github.com/iandykes/howfarcaniget
RUN go install github.com/iandykes/howfarcaniget

FROM alpine
COPY --from=builder /go/bin/howfarcaniget /app/

ENV LOG_LEVEL Debug
ENV INCLUDE_DEBUG_HANDLERS 1

ENTRYPOINT /app/howfarcaniget
EXPOSE 80