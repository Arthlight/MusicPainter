# :art: :musical_note: MusicPainter
Music Painter is a simple app that lets you connect your [Spotify](https://www.spotify.com) account and analyses the music you are currently listening to. The recommended way of use is letting the app run in the background while listening to your favourite playlists. Music Painter will then go along and paint a customized picture according to the metadata of all the songs you choose to listen to. (☞ ͡° ͜ʖ ͡°)☞

The Frontend for this webapp can be found [here](https://github.com/Arthlight/MusicPainter-Frontend).

You can visit the live version [here](http://161.35.173.232:8081/).

# Demo :movie_camera:
![](Demo/Demo.gif)

# Structure :open_file_folder:
```bash
├── api
│   └── v1
├── go.mod
├── go.sum
├── handler
│   └── socket.go
├── main.go
├── models
│   ├── colors.go
│   ├── events.go
│   ├── frontend.go
│   ├── socket.go
│   └── spotify.go
└── spotify
    └── spotify.go
```

# Frameworks & APIs :hammer_and_pick:
- Websocket: [Gorilla](https://github.com/gorilla/websocket)
- HTTP Router: [Chi](https://github.com/go-chi/chi)
- Dotenv: [godotenv](https://github.com/joho/godotenv)
- Deployment: [Docker](https://www.docker.com)
- Spotify API: [Spotify For Developers](https://developer.spotify.com)

# Installation :gear:
**Requirement**:
- [Golang](https://golang.org) >=1.14

**With Goland**:
- [Goland](https://www.jetbrains.com/go/) >=2020.1
- ```$ git clone https://github.com/Arthlight/MusicPainter.git```
- Open the project with Goland and click on the attached Terminal on the bottom.
- ```$ go run main.go```
- Open http://localhost:4000/

**Without Goland**:
- ```$ git clone https://github.com/Arthlight/MusicPainter.git```
- ```$ cd MusicPainter```
- ```$ go run main.go```
- Open http://localhost:4000/

Please refer to the corresponding installation guide for the frontend [here](https://github.com/Arthlight/MusicPainter-Frontend).
