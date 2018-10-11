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
	"strconv"
	"time"

	"github.com/aquilax/go-perlin"
	"github.com/gorilla/websocket"
)

/* adapted from examples at https://github.com/gorilla/websocket/blob/master/examples/echo/server.go
and https://gist.github.com/ismasan/3fb75381cd2deb6bfa9c */

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{}

/*
message constants to send over control channels
*/
const (
	CLOSE int = 0
	WRITE int = 1
	DONE  int = 2
)

/*
client holds the state of a websocket connection
*/
type client struct {
	id           string
	ws           *websocket.Conn
	ctrl         *controller
	toWriter     chan int
	toController chan int
}

/*
reader is the function that accepts input from the browser side of the socket
*/
func (c *client) reader() {
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading: %v", err)
			c.toWriter <- CLOSE
			return
		}
		log.Printf("Got %v from client %v", string(msg), c.id)
	}
}

/*
writer is the function that sends images to the browser when told to
*/
func (c *client) writer() {
	for {
		select {
		case event := <-c.toWriter:
			if event == CLOSE {
				c.ctrl.leavingClients <- c
				return
			}
			sw, err := c.ws.NextWriter(websocket.BinaryMessage)
			if err != nil {
				log.Printf("Error getting next writer: %v", err)
				c.ctrl.leavingClients <- c
				return
			}
			err = png.Encode(sw, c.ctrl.img)
			sw.Close()
			if err != nil {
				log.Printf("Error encoding image: %v", err)
			}
			c.toController <- DONE
		}
	}
}

/*
close cleans up a client connection/channels
*/
func (c *client) close() {
	log.Printf("closing client %v", c.id)
	close(c.toController)
	close(c.toWriter)
	c.ws.Close()
}

/*
controller is the object that coordinates the work between the clients
*/
type controller struct {
	img            *image.RGBA
	wp             *perlin.Perlin
	hp             *perlin.Perlin
	bg             color.RGBA
	fg             color.RGBA
	clients        map[string]*client
	newClients     chan *client
	leavingClients chan *client
	count          int
}

/*
newController instatiates a new controller object
*/
func newController(width, height int, bg, fg color.RGBA) *controller {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.ZP, draw.Src)
	rs := rand.NewSource(time.Now().UTC().UnixNano())
	wp := perlin.NewPerlinRandSource(2., 2., 3, rs)
	hp := perlin.NewPerlinRandSource(2., 2., 3, rs)
	clients := make(map[string]*client)
	newClients := make(chan *client)
	leavingClients := make(chan *client)
	return &controller{img, wp, hp, bg, fg, clients, newClients, leavingClients, 0}
}

/*
wsConnect accepts new websocket connections for clients
*/
func (ctrl *controller) wsConnect(w http.ResponseWriter, r *http.Request) {
	log.Printf("New client")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade failed: %v", err)
		return
	}
	toWriter := make(chan int)
	toController := make(chan int)
	c := &client{"", ws, ctrl, toWriter, toController}
	go c.reader()
	log.Printf("notifying ctrl of new client")
	ctrl.newClients <- c
	c.writer()
}

/*
run is the main logic loop for handling clients and generating images
*/
func (ctrl *controller) run() {
	width := ctrl.img.Rect.Max.X
	height := ctrl.img.Rect.Max.Y
	minWait := 250 * time.Millisecond
	i, j := 0., 0.
	for {
		select {
		case c := <-ctrl.newClients:
			c.id = strconv.Itoa(ctrl.count)
			log.Printf("id: %v from %v", c.id, ctrl.count)
			ctrl.clients[c.id] = c
			ctrl.count++
			log.Printf("added client %v", c.id)
		case c := <-ctrl.leavingClients:
			delete(ctrl.clients, c.id)
			c.close()
			log.Printf("removed client %v", c.id)
		default:
			start := time.Now()
			for id, c := range ctrl.clients {
				select {
				case c.toWriter <- WRITE:
				case <-time.After(1 * time.Second):
					log.Printf("Client %v took too long to accept message", id)
				}
			}
			for id, c := range ctrl.clients {
				select {
				case <-c.toController:
				case <-time.After(1 * time.Second):
					log.Printf("Client %v took too long to process img", id)
				}
			}
			ni := ctrl.wp.Noise1D(i)
			nj := ctrl.hp.Noise1D(j)
			px := int(math.Abs(ni) * float64(width))
			py := int(math.Abs(nj) * float64(height))
			for dy := -1; dy < 2; dy++ {
				for dx := -1; dx < 2; dx++ {
					x := px + dx
					y := py + dy
					if x > -1 && dx < width && y > -1 && y < height {
						ctrl.img.Set(x, y, ctrl.fg)
					}
				}
			}
			end := time.Now()
			total := end.Sub(start)
			wait := minWait - total
			if wait > 0 {
				time.Sleep(wait)
			}
			i += .01
			j += .01
		}
	}
}

/*
serveClient handles requests for the HTML/javascript
*/
func serveClient(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/data")
}

func main() {
	width, height := 512, 480
	green := color.RGBA{124, 252, 0, 255}
	brown := color.RGBA{210, 105, 30, 255}
	ctrl := newController(width, height, green, brown)
	go ctrl.run()
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/data", ctrl.wsConnect)
	http.HandleFunc("/", serveClient)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

/*
homeTemplate contains the HTML/javascript
*/
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
