package huejack
import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"log"
	"io"
	"io/ioutil"
)
var handlerMap map[string]huestate

func init() {
	log.SetOutput(ioutil.Discard)
	handlerMap = make(map[string]huestate)
	upnpTemplateInit()
}

func SetLogger(w io.Writer) {
	log.SetOutput(w)
}

func ListenAndServe(addr string) {
	router := httprouter.New()
	router.GET(upnp_uri, upnpSetup(addr))

	router.GET("/api/:userId", getLightsList)
	router.PUT("/api/:userId/lights/:lightId/state", setLightState)
	router.GET("/api/:userId/lights/:lightId", getLightInfo)

	go upnpResponder(addr, upnp_uri)
	http.ListenAndServe(addr, requestLogger(router))
}

type Handler func(Request) (state bool)

func Handle(deviceName string, h Handler) {
	log.Println("[HANDLE]", deviceName)
	handlerMap[deviceName] = huestate{
		Handler:h,
		OnState:false,
	}
}

func requestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("[WEB]", r.RemoteAddr, r.Method, r.URL)
		h.ServeHTTP(w, r)
	})
}