package test
//
//
//
//import (
//	"testing"
//	"log"
//
//	"github.com/drewolson/testflight"
//	"imgfling/server"
//	"github.com/stretchr/testify/assert"
//	"encoding/json"
//)
//
//type FormGif struct {
//	Gifs []string `json:"gifs"`
//}
//
//func Test_GifGenerator(t *testing.T) {
//
//	testflight.WithServer(server.InitRouting(), func(r *testflight.Requester) {
//		arrayGifs := []string{"","",""}
//
//		gifs := FormGif{Gifs:arrayGifs}
//
//		jsonBytes, err := json.Marshal(&gifs)
//		if err != nil {
//			log.Println("Error: GifHandler, ", err)
//			//w.Write([]byte(`[]`))
//			return
//		}
//
//
//		response := r.Post("/gif/", testflight.JSON, string(jsonBytes));
//
//		log.Println(len(response.Body))
//		assert.Equal(t, 200, response.StatusCode)
//		assert.True(t,len(response.Body) > 0)
//	})
//}