# Voice Recorder Web Application

#### Overview
This is a full-stack voice recording application built with Vue.js (frontend) and Go (backend). The application allows users to record audio, stream it to the backend, and save audio chunks in real-time.

#### Prerequisites
Node.js (v14+)
Vue CLI
Go (v1.23+)
WebSocket-compatible browser

#### Frontend Setup
Navigate to client directory
- install all the dependencies
    ```shell
      npm install
    ```

- Run development server:
    ```shell
      npm run dev
    ```

#### Backend Setup
Navigate to server directory

- Install Go dependencies:
    ```shell
        go mod tidy
        go get github.com/gorilla/websocket
        go get github.com/google/uuid
    ```

- Run the Go server:
    ```shell
        go run main.go
    ```

#### Key Features

* Real-time audio streaming
* WebSocket communication
* 30-second audio chunk saving
* Pause/Resume/Stop recording controls

#### Notes

Ensure microphone permissions are granted
WebSocket server runs on localhost:8080
Audio chunks are saved in a recordings/ directory