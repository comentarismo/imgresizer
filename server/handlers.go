package server

import (
	"log"
	"net/http"
	"bytes"
	"html/template"
	"fmt"
	"strconv"
	"strings"
	"image"
//    "image/gif"
    "image/jpeg"
    "image/png"
	"time"

	"gopkg.in/redis.v3"
	resize "github.com/nfnt/resize"
	cache "github.com/pmylund/go-cache"
	"io"
	"github.com/disintegration/imaging"
)

var (
	templates = template.Must(template.ParseFiles(
		"upload.html",
	))
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
//	AllowOrigin(w, r)

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

	widthStr := r.Form["width"]
	log.Println("width: ", widthStr)
	if len(widthStr) == 0 {
		log.Println("404 not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	width, err := strconv.ParseUint(widthStr[0],10,32)
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


	log.Println("Verify if is available on cache")

	cached, found := Cache.Get(url[0]+widthStr[0]+heightStr[0])
	abyte, found := cached.(image.Image)
	if found {
		o := jpeg.Options{quality}
		w.Header().Set("Content-Type", "image/jpeg")
		jpeg.Encode(w, abyte, &o)
		return
	}
//	fmt.Println("unknown type",cached)

	log.Println("Doing HTTP GET")

	res, err := http.Get(strings.Trim(url[0]," "))
	if err != nil {
		log.Printf("http.Get -> %v", err)
		return
	}

//	log.Println(res)

	// decode jpeg into image.Image
	img, imgtype, err := image.Decode(res.Body)
	log.Println(imgtype)
	if err != nil {
		log.Println(err)
		return
	}

	// resize to width height using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

//	w.Header().Set("Content-Length", fmt.Sprint(res.ContentLength))
	w.Header().Set("Content-Type", res.Header.Get("Content-Type"))

	o := jpeg.Options{quality}

	jpeg.Encode(w, m, &o)
	Cache.Set(url[0]+widthStr[0]+heightStr[0], m, cache.DefaultExpiration)

	defer res.Body.Close()
//	defer res.Close()
	return
}

func AllowOrigin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	//TODO: add origin validation
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
}

// writeError renders the error in the HTTP response.
func writeError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	if err := templates.ExecuteTemplate(w, "error.html", err); err != nil {
	}
}



func RedisImgPostHandler(w http.ResponseWriter, r *http.Request) {
	//allow http requests
	AllowOrigin(w, r)

	r.ParseForm()  //Parse url parameters passed, then parse the response packet for the POST body (request body)
	log.Println(r.Form) // print information on server side.

	log.Println("RedisImgPostHandler")

	url := r.Form["url"]
	if len(url) == 0 {
		log.Println("RedisImgPostHandler width 404 not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Println("url: ", url[0])
	if len(url) == 0 {
		log.Println("RedisImgPostHandler url 404 not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	widthStr := r.Form["width"]
	log.Println("width: ", widthStr)
	if len(widthStr) == 0 {
		log.Println("RedisImgPostHandler width 404 not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	width, err := strconv.ParseInt(widthStr[0],10,32)
	if err != nil {
		log.Println("RedisImgPostHandler width ParseUint 404 not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	heightStr := r.Form["height"]
	log.Println("height: ", heightStr)
	if len(heightStr) == 0 {
		log.Println("RedisImgPostHandler height 404 not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	height, err := strconv.ParseInt(heightStr[0],10,32)
	if err != nil {
		log.Println("RedisImgPostHandler height ParseUint 404 not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

//	qualityStr := r.Form["quality"]
//	var quality int
//	quality = 30

//	log.Println("ImgPostHandler qualityStr: ", qualityStr)
//	if len(qualityStr) != 0 {
//		quality64, err := strconv.ParseInt(qualityStr[0], 10, 32)
//		if err != nil {
//			log.Println("ImgPostHandler qualityStr 404 not found")
//			w.WriteHeader(http.StatusNotFound)
//			return
//		}
//		quality = int(quality64)
//	}


	log.Println("RedisImgPostHandler Check form params --> ok")


	log.Println("RedisImgPostHandler Verify if is available on cache")

	cached, err := Client.Get("ImgPostHandler"+url[0]+widthStr[0]+heightStr[0]).Result()
	if err == redis.Nil {
		log.Println("RedisImgPostHandler "+url[0]+widthStr[0]+heightStr[0]+" does not exists on cache")
	} else if err != nil {
		panic(err)
		return
	} else {
		log.Println("RedisImgPostHandler Serve from cache")
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write([]byte(cached))
		return
	}
	//	fmt.Println("unknown type",cached)

	log.Println("Doing HTTP GET")

	res, err := http.Get(strings.Trim(url[0]," "))
	if err != nil {
		log.Printf("http.Get -> %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//	log.Println(res)

	img_,imgType, err := image.Decode(res.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Println(imgType)
	img := imaging.Thumbnail(img_, int(width), int(height), imaging.CatmullRom)

//	img, err := decodeImage(res.Body)
//	if err || img == nil {
//		log.Println(err)
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}


	// resize to width height using Lanczos resampling
	// and preserve aspect ratio
//	m := resize.Resize( img, resize.Lanczos3)
//
//	w.Header().Set("Content-Type", fmt.Sprint(res.ContentLength))
//
//	o := jpeg.Options{quality}
//
//	err = jpeg.Encode(w, m, &o)
//	if err != nil {
//		log.Println(err)
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
	log.Println("ImgPostHandler saving on redis cache")

	buf := new(bytes.Buffer)
	jpeg.Encode(buf, img, nil)
	send_s3 := buf.Bytes()

	Client.Set("ImgPostHandler"+url[0]+widthStr[0]+heightStr[0], string(send_s3), time.Hour * 24 * 2)

	w.Write([]byte(send_s3))
	w.Header().Set("Content-Type", "image/jpeg")

	defer r.Body.Close()
	//	defer res.Close()
	return
}

func decodeImage(res io.Reader) (image.Image, error) {
	// decode jpeg into image.Image
	config, imgFormat, err := image.DecodeConfig(res)
	if err != nil {
		log.Println("error when DecodeConfig into image.Image")
		log.Println(err)
//		w.WriteHeader(http.StatusNotFound)
		return nil,err
	}
	log.Println(imgFormat)
	log.Println(config)

	if imgFormat == "jpeg" {
		img, err := jpeg.Decode(res)
		if err != nil {
			log.Println("error when decoding jpeg into image.Image")
			log.Println(err)
			return nil,err
		}
		return img,nil
	}else if imgFormat == "png" {
//		image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
		img, err := png.Decode(res)
		if err != nil {
			log.Println("error when decoding png into image.Image")
			log.Println(err)
			return nil,err
		}
		return img,nil

	}else {
		return nil,nil
	}

//	log.Println(imgtype)
//	if err != nil {
//		log.Println("error when decoding jpeg into image.Image")
//		log.Println(err)
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}

	return nil,nil
}

func RedisImgGetHandler (w http.ResponseWriter, req *http.Request) {
	//allow http requests
	AllowOrigin(w, req)



//	operator := req.URL.Query().Get(":operator")
//
//	key := req.URL.Query().Get(":key")
//	value := req.URL.Query().Get(":value")
//
//	widthStr := req.URL.Query().Get(":width")
//	heightStr := req.URL.Query().Get(":height")
//
//	log.Println("width: ", widthStr)
//	if len(widthStr) == 0 {
//		log.Println("ImgPostHandler width 404 not found")
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//
//	log.Println("height: ", heightStr)
//	if len(heightStr) == 0 {
//		log.Println("ImgPostHandler height 404 not found")
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//
//	width, err := strconv.ParseUint(widthStr,10,32)
//	if err != nil {
//		log.Println("ImgPostHandler width ParseUint 404 not found")
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//
//	height, err := strconv.ParseUint(heightStr,10,32)
//	if err != nil {
//		log.Println("ImgPostHandler height ParseUint 404 not found")
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//
//	qualityStr := req.URL.Query().Get(":quality")
//	var quality int
//	quality = 30
//
//	log.Println("ImgPostHandler qualityStr: ", qualityStr)
//	if len(qualityStr) != 0 {
//		quality64, err := strconv.ParseInt(qualityStr, 10, 32)
//		if err != nil {
//			log.Println("ImgPostHandler qualityStr 404 not found")
//			w.WriteHeader(http.StatusNotFound)
//			return
//		}
//		quality = int(quality64)
//	}
//	log.Println("ImgPostHandler Check form params --> ok")


	//GET JSON WITH THE URL



}