package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"strconv"

	service "github.com/dovudwkt/playground/server/services"
)

// ImageHTTPHandler accepts an image and query parameters to resize image.
//  Queary parameters: 'w' - Width and 'h' - Height
type ImageHTTPHandler struct {
	Service service.Service
}

func (h ImageHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// get url query parameters
	width, err := strconv.ParseUint(req.URL.Query().Get("w"), 10, 32)
	if err != nil {
		log.Fatalln(err)
	}
	height, err := strconv.ParseUint(req.URL.Query().Get("h"), 10, 32)
	if err != nil {
		log.Fatalln(err)
	}

	// read request body to the buffer
	buffer := make([]byte, req.ContentLength)
	_, err = io.ReadFull(req.Body, buffer)
	if err != nil {
		log.Fatalln(err)
	}

	// decode image from buffer
	img, _, err := image.Decode(bytes.NewReader(buffer))
	if err != nil {
		log.Fatalln(err)
	}

	// call a service to resize the image providing configs
	resizedImg, err := h.Service.ResizeImage(&img, service.ResizeConfig{
		W: uint(height), H: uint(width), Interp: service.NearestNeighbor,
	})
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	// return the image in response writer
	err = writeImage(w, resizedImg)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
}

func writeImage(w http.ResponseWriter, img *image.Image) error {
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		return errors.New("unable to encode image")
	}

	// set response writer headers
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

	if _, err := w.Write(buffer.Bytes()); err != nil {
		return errors.New("unable to write image")
	}

	return nil
}

// -------------------------------------------------

type ImageFromURLHTTPHandler struct {
	Service service.Service
}

func (h ImageFromURLHTTPHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	client := http.Client{}

	reqImg, err := client.Get("https://images.unsplash.com/photo-1541963463532-d68292c34b19?ixid=MnwxMjA3fDB8MHxzZWFyY2h8Mnx8Ym9va3xlbnwwfHwwfHw%3D&ixlib=rb-1.2.1&w=1000&q=80)")
	if err != nil {
		fmt.Fprintf(res, "Error %d", err)
		return
	}

	buffer := make([]byte, reqImg.ContentLength)
	io.ReadFull(reqImg.Body, buffer)

	res.Header().Set("Content-Length", fmt.Sprint(reqImg.ContentLength))
	res.Header().Set("Content-Type", reqImg.Header.Get("Content-Type"))

	res.Write(buffer)
	req.Body.Close()
}
