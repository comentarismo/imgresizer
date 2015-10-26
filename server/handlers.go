package server

import (
	"log"
	"net/http"
	"bytes"
	"html/template"
	"fmt"
	"image/jpeg"
	resize "github.com/nfnt/resize"

	"strconv"
)

var (
	templates = template.Must(template.ParseFiles(
		"upload.html",
	))
//	imageQuality = jpeg.Options{95}
)


func ImgHandler(w http.ResponseWriter, r *http.Request){
	b := &bytes.Buffer{}
	if err := templates.ExecuteTemplate(b, "upload.html", nil); err != nil {
		writeError(w, r, err)
		return
	}
	b.WriteTo(w)
	return
}

func ImgPostHandler(w http.ResponseWriter, r *http.Request) {
	//allow http requests
	AllowOrigin(w, r)

	r.ParseForm()  //Parse url parameters passed, then parse the response packet for the POST body (request body)
	fmt.Println(r.Form) // print information on server side.

	log.Println("ImgPostHandler")

	url := r.Form["url"]
	log.Println("url: ", url[0])
	if len(url) == 0 {
		log.Println("404 not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	withStr := r.Form["with"]
	log.Println("with: ", withStr)
	if len(withStr) == 0 {
		log.Println("404 not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	with, err := strconv.ParseUint(withStr[0],10,32)
	if err != nil {
		log.Println("404 not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	heightStr := r.Form["height"]
	log.Println("height: ", heightStr)
	if len(heightStr) == 0 {
		log.Println("404 not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	height, err := strconv.ParseUint(heightStr[0],10,32)
	if err != nil {
		log.Println("404 not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	qualityStr := r.Form["quality"]
	var quality int
	quality = 30

	log.Println("qualityStr: ", qualityStr)
	if len(qualityStr) != 0 {
		quality64, err := strconv.ParseInt(qualityStr[0], 10, 32)
		if err != nil {
			log.Println("404 not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		quality = int(quality64)
	}


	log.Println("Check form params --> ok")

	log.Println("Doing HTTP GET")

	res, err := http.Get(url[0])
	if err != nil {
		log.Fatalf("http.Get -> %v", err)
		return
	}

	log.Println(res)

	// decode jpeg into image.Image
	img, err := jpeg.Decode(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(uint(with), uint(height), img, resize.Lanczos3)

	w.Header().Set("Content-Length", fmt.Sprint(res.ContentLength))
	w.Header().Set("Content-Type", res.Header.Get("Content-Type"))

	o := jpeg.Options{quality}

	jpeg.Encode(w, m, &o)

	res.Body.Close()
	return
}

func AllowOrigin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	//	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	//TODO: add origin validation
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	//	w.Header().Set("Access-Control-Allow-Origin","http://192.168.0.22:3001,http://")
}

// writeError renders the error in the HTTP response.
func writeError(w http.ResponseWriter, r *http.Request, err error) {
//	w.Errorf("Error: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
	if err := templates.ExecuteTemplate(w, "error.html", err); err != nil {
//		w.Errorf("templates.ExecuteTemplate: %v", err)
	}
}
