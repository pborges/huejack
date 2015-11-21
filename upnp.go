package huejack
import (
	"github.com/julienschmidt/httprouter"
	"text/template"
	"net/http"
	"log"
	"net"
	"strings"
	"bytes"
)

const (
	upnp_multicast_address = "239.255.255.250:1900"
	upnp_uri = "/upnp/setup.xml"
)

var responseTemplateText =
`HTTP/1.1 200 OK
CACHE-CONTROL: max-age=86400
EXT:
LOCATION: http://{{.}}
OPT: "http://schemas.upnp.org/upnp/1/0/"; ns=01
ST: urn:schemas-upnp-org:device:basic:1
USN: uuid:Socket-1_0-221438K0100073::urn:Belkin:device:**

`

var setupTemplateText =
`<?xml version="1.0"?>
<root xmlns="urn:schemas-upnp-org:device-1-0">
        <specVersion>
                <major>1</major>
                <minor>0</minor>
        </specVersion>
        <URLBase>http://{{.}}/</URLBase>
        <device>
			<deviceType>urn:schemas-upnp-org:device:Basic:1</deviceType>
			<friendlyName>huejack</friendlyName>
			<manufacturer>Royal Philips Electronics</manufacturer>
			<modelName>Philips hue bridge 2012</modelName>
			<modelNumber>929000226503</modelNumber>
			<UDN>uuid:f6543a06-800d-48ba-8d8f-bc2949eddc33</UDN>
        </device>
</root>`

type upnpData struct {
	Addr string
	Uri  string
}

var setupTemplate *template.Template
func upnpTemplateInit() {
	var err error
	setupTemplate, err = template.New("").Parse(setupTemplateText)
	if err != nil {
		log.Fatalln("upnpTemplateInit:", err)
	}
}

func upnpSetup(addr string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/xml")
		err := setupTemplate.Execute(w, addr)
		if err != nil {
			log.Fatalln("[WEB] upnpSetup:", err)
		}
	}
}


func upnpResponder(hostAddr string, endpoint string) {
	responseTemplate, err := template.New("").Parse(responseTemplateText)

	log.Println("[UPNP] listening...")
	addr, err := net.ResolveUDPAddr("udp", upnp_multicast_address)
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.ListenMulticastUDP("udp", nil, addr)
	l.SetReadBuffer(1024)

	for {
		b := make([]byte, 1024)
		n, src, err := l.ReadFromUDP(b)
		if err != nil {
			log.Fatal("[UPNP] ReadFromUDP failed:", err)
		}

		if strings.Contains(string(b[:n]), "MAN: \"ssdp:discover\"") {
			c, err := net.DialUDP("udp", nil, src)
			if err != nil {
				log.Fatal("[UPNP] DialUDP failed:", err)
			}

			log.Println("[UPNP] discovery request from", src)

			// For whatever reason I can't execute the template using c as the reader,
			// you HAVE to put it in a buffer first
			// possible timing issue?
			// don't believe me? try it
			b := &bytes.Buffer{}
			err = responseTemplate.Execute(b, hostAddr + endpoint)
			if err != nil {
				log.Fatal("[UPNP] execute template failed:", err)
			}
			c.Write(b.Bytes())
		}
	}
}