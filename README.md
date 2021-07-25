# image-resizer

## Installing dependencies: 
`go get github.com/nfnt/resize`

## Starting server
From the project directory, run `go run server.go` to start the server. By default the port is set to `3001`.

## Endpoints

### POST `/images/resize?w=xxx&h=xxx`

The endpoint that accepts an image and returns the resized version. Set `w` (width) and `h` (height) query parameters to desired values.

