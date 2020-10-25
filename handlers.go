package statsviz

import (
	"log"
	"net/http"
	"strings"
)

// Index responds to a request for /debug/statsviz with the statsviz HTML page
// which shows a live visualization of the statistics sent by the application
// over the websocket handler Ws.
var Index = IndexAtRoot(defaultRoot)

// IndexAtRoot returns an index statsviz handler rooted at root. It's useful if
// you desire your server to responds with the statsviz HTML page at a
// path that is different than /debug/statsviz.
func IndexAtRoot(root string) http.Handler {
	return http.StripPrefix(strings.TrimRight(root, "/")+"/", http.FileServer(assets))
}

// Ws upgrades the HTTP server connection to the WebSocket protocol and sends
// application statistics every second.
//
// If the upgrade fails, an HTTP error response is sent to the client.
func Ws(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ws: Upgrade error:", err)
		return
	}
	defer ws.Close()

	// Explicitly ignore this error. We don't want to spam standard output
	// each time the other end of the websocket connection closes.
	_ = sendStats(ws)
}
