package main

import (
	"flag"
	"html/template"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/aquilax/go-perlin"
	"github.com/gorilla/websocket"
)

/* adapted from examples at https://github.com/gorilla/websocket/blob/master/examples/echo/server.go
and https://gist.github.com/ismasan/3fb75381cd2deb6bfa9c */

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{}

func reader(ws *websocket.Conn) {
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading: %v", err)
			break
		}
		log.Printf("Got %v", string(msg))
	}
}

func wsConnect(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade failed: %v", err)
		return
	}
	defer ws.Close()
	go reader(ws)
	width, height := 512, 480
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	green := color.RGBA{124, 252, 0, 255}
	brown := color.RGBA{210, 105, 30, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{green}, image.ZP, draw.Src)
	rs := rand.NewSource(time.Now().UTC().UnixNano())
	wp := perlin.NewPerlinRandSource(2., 2., 3, rs)
	hp := perlin.NewPerlinRandSource(2., 2., 3, rs)
	i, j := 0., 0.
	for {
		ni := wp.Noise1D(i)
		nj := hp.Noise1D(j)
		px := int(math.Abs(ni) * float64(width))
		py := int(math.Abs(nj) * float64(height))
		for dy := -1; dy < 2; dy++ {
			for dx := -1; dx < 2; dx++ {
				x := px + dx
				y := py + dy
				if x > -1 && dx < width && y > -1 && y < height {
					img.Set(x, y, brown)
				}
			}
		}
		time.Sleep(time.Millisecond * 250)
		sw, err := ws.NextWriter(websocket.BinaryMessage)
		if err != nil {
			log.Printf("Error getting next writer: %v", err)
			return
		}
		err = png.Encode(sw, img)
		sw.Close()
		if err != nil {
			log.Printf("Error encoding image: %v", err)
			return
		}
		i += .01
		j += .01
	}
}

func serveClient(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/data")
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/data", wsConnect)
	http.HandleFunc("/", serveClient)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function onMessage(evt) {
            var img = document.getElementById("dispimg")
            if (img == null){
                img = document.createElement("img")
                img.setAttribute("id", "dispimg")
            }
            var urlObject = URL.createObjectURL(evt.data);
            img.src = urlObject;
            document.body.appendChild(img);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
