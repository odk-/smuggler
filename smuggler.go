package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type pack struct {
	reader *io.PipeReader
	writer *io.PipeWriter
}

type packs struct {
	sync.RWMutex
	p map[string]*pack
}

type Err struct {
	Msg string
}

func (e Err) Error() string {
	return fmt.Sprint(e.Msg)
}

var packList = packs{p: make(map[string]*pack)}

//register new pack to smuggle
func smuggle(w http.ResponseWriter, r *http.Request) {
	pid := getPackID()
	packList.Lock()
	reader, writer := io.Pipe()
	packList.p[pid] = &pack{reader: reader, writer: writer}
	packList.Unlock()
	w.Write([]byte(fmt.Sprintf("To send data use:\ncurl -X POST --data-binary \"@file\" https://%s/pack?packID=%s\nTo retrieve them use:\ncurl https://%s/pack?packID=%s\n", r.Host, pid, r.Host, pid)))
}

//check if someone want to send or get package
func routePack(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		packSender(w, r)
	case "POST":
		packReceiver(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

//send package to client, will block until sending side shows up
func packSender(w http.ResponseWriter, r *http.Request) {
	pkt, err := getPack(r.URL.Query().Get("packID"))
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	defer pkt.reader.Close()
	defer delPack(r.URL.Query().Get("packID"))
	io.Copy(w, pkt.reader)
}

//get package from client, will block until receiving side shows up
func packReceiver(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	pkt, err := getPack(r.URL.Query().Get("packID"))
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	defer pkt.writer.Close()
	io.Copy(pkt.writer, r.Body)
}

func getPackID() string {
	uuid := make([]byte, 16)
	n, _ := rand.Read(uuid)
	if n != len(uuid) {
		return ""
	}
	return hex.EncodeToString(uuid)
}

//get one pack from list
func getPack(id string) (*pack, error) {
	packList.RLock()
	defer packList.RUnlock()
	if pc, ok := packList.p[id]; ok {
		return pc, nil
	} else {
		return nil, Err{Msg: "ERROR: Pack not found\n"}
	}
}

//del one pack from list
func delPack(id string) {
	packList.Lock()
	delete(packList.p, id)
	packList.Unlock()
}

//TODO periodic clenup of registered but not used packs
func main() {
	mx := http.NewServeMux()

	mx.HandleFunc("/smuggle", smuggle)
	mx.HandleFunc("/pack", routePack)

	http.ListenAndServe(":8080", mx)
}
