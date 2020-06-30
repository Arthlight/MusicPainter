FROM golang:1.14

RUN mkdir app

WORKDIR /app

ADD . /app

RUN go build

EXPOSE 4000

CMD ["./Spotify-Visualizer"]