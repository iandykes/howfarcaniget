FROM golang as builder
ADD . /go/src/github.com/iandykes/howfarcaniget
RUN cd /go/src/github.com/iandykes/howfarcaniget && go get && go build

FROM alpine
COPY --from=builder /go/src/github.com/iandykes/howfarcaniget/howfarcaniget /app/

ENV LOG_LEVEL Debug
ENV INCLUDE_DEBUG_HANDLERS 1

ENTRYPOINT /app/howfarcaniget
EXPOSE 80