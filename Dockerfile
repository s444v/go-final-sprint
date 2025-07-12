FROM ubuntu:latest
WORKDIR /app

COPY server .
COPY web ./web

ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db
ENV TODO_PASSWORD=12345

EXPOSE 7540

CMD ["./server"]