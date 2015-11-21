package huejack
import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
	"log"
	"strconv"
)

type light struct {
	State            struct {
						 On        bool `json:"on"`
						 Bri       int `json:"bri"`
						 Hue       int `json:"hue"`
						 Sat       int `json:"sat"`
						 Effect    string `json:"effect"`
						 Ct        int `json:"ct"`
						 Alert     string `json:"alert"`
						 Colormode string `json:"colormode"`
						 Reachable bool `json:"reachable"`
						 XY        []float64 `json:"xy"`
					 } `json:"state"`
	Type             string `json:"type"`
	Name             string `json:"name"`
	ModelId          string `json:"modelid"`
	ManufacturerName string `json:"manufacturername"`
	UniqueId         string `json:"uniqueid"`
	SwVersion        string `json:"swversion"`
	PointSymbol      struct {
						 One   string `json:"1"`
						 Two   string `json:"2"`
						 Three string `json:"3"`
						 Four  string `json:"4"`
						 Five  string `json:"5"`
						 Six   string `json:"6"`
						 Seven string `json:"7"`
						 Eight string `json:"8"`
					 } `json:"pointsymbol"`
}

type lights struct {
	Lights map[string]light `json:"lights"`
}

func initLight(name string) light {
	l := light{
		Type:"Extended color light",
		ModelId:"LCT001",
		SwVersion:"65003148",
		ManufacturerName:"Philips",
		Name:name,
		UniqueId:name,
	}
	l.State.Reachable = true
	l.State.XY = []float64{0.4255, 0.3998}  // this seems to be voodoo, if it is nil the echo says it could not turn on/off the device, useful...
	return l
}

func enumerateLights() lights {
	lightList := lights{}
	lightList.Lights = make(map[string]light)
	for name, hstate := range handlerMap {
		l := initLight(name)
		l.State.On = hstate.OnState
		lightList.Lights[l.UniqueId] = l
	}
	return lightList
}

func getLightsList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(enumerateLights())
	if err != nil {
		log.Fatalln("[WEB] Error encoding json", err)
	}
}

func setLightState(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	req := make(map[string]bool)
	json.NewDecoder(r.Body).Decode(&req)

	l := initLight(p.ByName("lightId"))

	log.Println("[DEVICE]", p.ByName("userId"), "requested state:", req["on"])
	state := false;
	if hstate, ok := handlerMap[p.ByName("lightId")]; ok {
		state = hstate.Handler(Request{
			UserId:p.ByName("userId"),
			RequestedOnState:req["on"],
			RemoteAddr:r.RemoteAddr,
		})
		log.Println("[DEVICE] handler replied with state:", state)
		hstate.OnState = state
		handlerMap[p.ByName("lightId")] = hstate
	}

	// this is very ugly...
	w.Write([]byte("[{\"success\":{\"/lights/" + l.UniqueId + "/state/on\":" + strconv.FormatBool(state) + "}}]"))
}

func getLightInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	l := initLight(p.ByName("lightId"))

	if hstate, ok := handlerMap[p.ByName("lightId")]; ok {
		if hstate.OnState {
			l.State.On = true
		}
	}

	err := json.NewEncoder(w).Encode(l)
	if err != nil {
		log.Fatalln("[WEB] Error encoding json", err)
	}
}