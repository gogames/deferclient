package deferstats

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"testing"
)

func TestClient(t *testing.T) {

	dps := NewClient("token")

	var resbody = make(chan []byte)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		resbody <- body

	})

	// set listener
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error("http not listening")
	}

	dps.statsUrl = "http://" + l.Addr().String() + "/"

	go http.Serve(l, mux)

	// capture our stats && ship
	dps.capture()

	// ensure we have some stats
	var ds DeferStats
	err = json.Unmarshal(<-resbody, &ds)
	if err != nil {
		t.Error(err)
	}

	mem, _ := strconv.Atoi(ds.Mem)
	if mem <= 0 {
		t.Error("mems not gt 0")
	}

	grs, _ := strconv.Atoi(ds.GoRoutines)
	if grs <= 0 {
		t.Error("gr not gt 0")
	}

	gc, _ := strconv.Atoi(ds.GC)
	if gc <= 0 {
		t.Error("gc not gt 0")
	}

	cgos, _ := strconv.ParseInt(ds.Cgos, 10, 64)
	if cgos <= 0 {
		t.Error("cgos not gt 0")
	}

}
