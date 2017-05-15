package server

import (
	"bytes"
	"fmt"
	"html/template"
	"image"
	"log"
	"net/http"
	"strconv"
	"strings"
	//    "image/gif"
	"image/jpeg"
	"image/png"
	"time"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	resize "github.com/nfnt/resize"
	cache "github.com/pmylund/go-cache"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"gopkg.in/redis.v3"
	"io"
)

var (
	templates = template.Must(template.ParseFiles(
		"upload.html",
	))
)

func MemeHandler(w http.ResponseWriter, r *http.Request) {
	//allow http requests
	AllowOrigin(w, r)

	r.ParseForm()       //Parse url parameters passed, then parse the response packet for the POST body (request body)
	log.Println(r.Form) // print information on server side.

	log.Println("MemeHandler")

	url := r.Form["url"]
	if len(url) == 0 {
		log.Println("MemeHandler url 404 not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Println("url: ", url[0])
	if len(url) == 0 {
		log.Println("MemeHandler url 404 not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	widthStr := r.Form["width"]
	log.Println("width: ", widthStr)
	if len(widthStr) == 0 {
		log.Println("MemeHandler width 404 not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	width, err := strconv.ParseInt(widthStr[0], 10, 32)
	if err != nil {
		log.Println("MemeHandler width ParseUint 404 not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	heightStr := r.Form["height"]
	log.Println("height: ", heightStr)
	if len(heightStr) == 0 {
		log.Println("MemeHandler height 404 not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	height, err := strconv.ParseInt(heightStr[0], 10, 32)
	if err != nil {
		log.Println("MemeHandler height ParseUint 404 not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	qualityStr := r.Form["quality"]
	var quality int
	quality = 30

	log.Println("ImgPostHandler qualityStr: ", qualityStr)
	if len(qualityStr) != 0 {
		quality64, err := strconv.ParseInt(qualityStr[0], 10, 32)
		if err != nil {
			log.Println("ImgPostHandler qualityStr 404 not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		quality = int(quality64)
	}

	log.Println("MemeHandler Check form params --> ok", url, width, height, quality)

	cached, isCacheValid := GetFromCache("MemeHandler" + url[0] + widthStr[0] + heightStr[0])
	if isCacheValid {
		log.Println("MemeHandler, Return from cache ->")
		w.Write([]byte(cached))
		return
	}

	log.Println("MemeHandler HTTP GET")

	res, err := http.Get(strings.Trim(url[0], " "))
	if err != nil {
		log.Printf("MemeHandler, http.Get -> %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//	log.Println(res)

	img_, imgType, err := image.Decode(res.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Println(imgType)
	//img := imaging.Thumbnail(img_, int(width), int(height), imaging.CatmullRom)

	dc := gg.NewContextForImage(img_)
	//dc.SetRGB(1, 1, 1)
	//dc.Clear()
	//font := "/Library/Fonts/Impact.ttf"
	//font := "arial.ttf"
	f, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	font := truetype.NewFace(f, &truetype.Options{
		Size:    60,
		DPI:     72,
		Hinting: font.HintingNone,
	})

	dc.SetFontFace(font)

	dc.SetRGB(0, 0, 0)
	s := "ONE DOES NOT SIMPLY"
	sheader := "Qualquer Header text"
	smiddle := "Test Middle text"
	sfooter := "Test Footer text"
	n := 8 // "stroke" size
	const S =1024
	for dy := -n; dy <= n; dy++ {
		for dx := -n; dx <= n; dx++ {
			if dx*dx+dy*dy >= n*n {
				// give it rounded corners
				continue
			}
			x := float64(width) + float64(dx)
			y := float64(height) + float64(dy)
			xsheader := S/2 + float64(dx)
			ysheader := S/12 + float64(dy)
			xsmiddle := S/2 + float64(dx)
			ysmiddle := S/2.5 + float64(dy)
			xsfooter := S/2 + float64(dx)
			ysfooter := S/1.4 + float64(dy)
			dc.DrawStringAnchored(s, x, y, 0.5, 0.5)
			dc.DrawStringAnchored(sheader, xsheader, ysheader, 0.5, 0.5)
			dc.DrawStringAnchored(smiddle, xsmiddle, ysmiddle, 0.5, 0.5)
			dc.DrawStringAnchored(sfooter, xsfooter, ysfooter, 0.5, 0.5)
		}
	}	

	dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored(s, float64(width), float64(height), 0.5, 0.5)
	dc.DrawStringAnchored(sheader, S/2, S/12, 0.5, 0.5)
	dc.DrawStringAnchored(smiddle, S/2, S/2.5, 0.5, 0.5)
	dc.DrawStringAnchored(sfooter, S/2, S/1.4, 0.5, 0.5)

	buf := new(bytes.Buffer)

	dc.EncodePNG(buf)

	send_s3 := buf.Bytes()

	log.Println("MemeHandler saving on redis cache")
	SetToCache("MemeHandler"+url[0]+widthStr[0]+heightStr[0], string(send_s3), time.Hour*24*2)

	w.Header().Set("Content-Type", "image/png")
	w.Write([]byte(send_s3))

	defer r.Body.Close()
	//	defer res.Close()
	return

}

func ImgHandler(w http.ResponseWriter, r *http.Request) {
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

	r.ParseForm()       //Parse url parameters passed, then parse the response packet for the POST body (request body)
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
	width, err := strconv.ParseUint(widthStr[0], 10, 32)
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
	height, err := strconv.ParseUint(heightStr[0], 10, 32)
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

	cached, found := Cache.Get(url[0] + widthStr[0] + heightStr[0])
	abyte, found := cached.(image.Image)
	if found {
		o := jpeg.Options{quality}
		w.Header().Set("Content-Type", "image/jpeg")
		jpeg.Encode(w, abyte, &o)
		return
	}
	//	fmt.Println("unknown type",cached)

	log.Println("Doing HTTP GET")

	res, err := http.Get(strings.Trim(url[0], " "))
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

	r.ParseForm()       //Parse url parameters passed, then parse the response packet for the POST body (request body)
	log.Println(r.Form) // print information on server side.

	log.Println("RedisImgPostHandler")

	url := r.Form["url"]
	if len(url) == 0 {
		log.Println("RedisImgPostHandler url 404 not found")
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
	width, err := strconv.ParseInt(widthStr[0], 10, 32)
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
	height, err := strconv.ParseInt(heightStr[0], 10, 32)
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

	cached, err := Client.Get("ImgPostHandler" + url[0] + widthStr[0] + heightStr[0]).Result()
	if err == redis.Nil {
		log.Println("RedisImgPostHandler " + url[0] + widthStr[0] + heightStr[0] + " does not exists on cache")
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

	res, err := http.Get(strings.Trim(url[0], " "))
	if err != nil {
		log.Printf("http.Get -> %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//	log.Println(res)

	img_, imgType, err := image.Decode(res.Body)
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

	Client.Set("ImgPostHandler"+url[0]+widthStr[0]+heightStr[0], string(send_s3), time.Hour*24*2)

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
		return nil, err
	}
	log.Println(imgFormat)
	log.Println(config)

	if imgFormat == "jpeg" {
		img, err := jpeg.Decode(res)
		if err != nil {
			log.Println("error when decoding jpeg into image.Image")
			log.Println(err)
			return nil, err
		}
		return img, nil
	} else if imgFormat == "png" {
		//		image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
		img, err := png.Decode(res)
		if err != nil {
			log.Println("error when decoding png into image.Image")
			log.Println(err)
			return nil, err
		}
		return img, nil

	} else {
		return nil, nil
	}

	//	log.Println(imgtype)
	//	if err != nil {
	//		log.Println("error when decoding jpeg into image.Image")
	//		log.Println(err)
	//		w.WriteHeader(http.StatusNotFound)
	//		return
	//	}

	return nil, nil
}

func GifPostHandler(w http.ResponseWriter, r *http.Request) {
	//allow http requests
	AllowOrigin(w, r)

	r.ParseForm()       //Parse url parameters passed, then parse the response packet for the POST body (request body)
	fmt.Println(r.Form) // print information on server side.

	log.Println("GifPostHandler")

	//w.Header().Set("Content-Length", fmt.Sprint(res.ContentLength))
	//w.Header().Set("Content-Type", res.Header.Get("Content-Type"))

	//jpeg.Encode(w, m, jpeg.Options{30})

	//res.Body.Close()
	w.Write([]byte(`[]`))
	return
}
