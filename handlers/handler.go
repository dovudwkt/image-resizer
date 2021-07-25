package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	service "github.com/dovudwkt/playground/server/services"
)

func UserHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("error reading body: ", err.Error())
	}

	var u User

	err = json.Unmarshal(body, &u)
	if err != nil {
		fmt.Println("unmarshal body: ", err.Error())
	}

	fmt.Printf("user object: %+v", u)

	u.ID = rand.Int()

	json.NewEncoder(w).Encode(u)
}

type User struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// ----------------------------------------

type ImageHTTPHandler struct {
	Service service.Service
}

func (h ImageHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	buffer := make([]byte, req.ContentLength)
	_, err := io.ReadFull(req.Body, buffer)
	if err != nil {
		log.Fatalln(err)
	}

	width, err := strconv.ParseUint(req.URL.Query().Get("w"), 10, 32)
	if err != nil {
		log.Fatalln(err)
	}
	height, err := strconv.ParseUint(req.URL.Query().Get("h"), 10, 32)
	if err != nil {
		log.Fatalln(err)
	}

	img, _, err := image.Decode(bytes.NewReader(buffer))
	if err != nil {
		log.Fatalln(err)
	}

	resizedImg, err := h.Service.ResizeImage(&img, service.ResizeConfig{
		W: uint(height), H: uint(width), Interp: service.NearestNeighbor,
	})
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

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

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

	if _, err := w.Write(buffer.Bytes()); err != nil {
		return errors.New("unable to write image")
	}

	return nil
}

// ------------------------------------

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
