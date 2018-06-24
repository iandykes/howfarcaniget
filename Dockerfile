FROM alpine:3.7
COPY howfarcaniget.exe /app/
WORKDIR /app
ENTRYPOINT ["howfarcaniget.exe"]
EXPOSE 80