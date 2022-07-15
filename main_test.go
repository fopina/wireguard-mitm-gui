package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
)

func assertError(t *testing.T, err error, message string) {
	if err.Error() != message {
		t.Fatalf("unexpected err: %v", err)
	}
}

func waitForIt(t *testing.T) {
	var err error

	for {
		// FIXME: how to properly wait for main() to get ready without modifying it...?
		// gief python monkeypatching power
		_, err = http.Get("http://127.0.0.1:9999/")
		if err == nil {
			return
		}
		if !strings.Contains(err.Error(), "connection refused") {
			t.Fatal(err)
		}
		t.Log("http server not ready yet...")
	}
}

func setupTest(t *testing.T) {
	*clientIp = "1.2.3.4"
	*bindAddress = "127.0.0.1:9999"
	*ipHeader = "x-real-ip"
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal("unable to get current dir")
	}
	*iptablesPath = path.Join(cwd, "fake_iptables.sh")
	*iptablesSavePath = path.Join(cwd, "fake_iptables_save.sh")
	// needed to reset HandleFunc registry across tests
	http.DefaultServeMux = http.NewServeMux()
	os.Setenv("MITMGUI_TESTING_TESTDB", path.Join(t.TempDir(), "test.db"))
}

func TestParseIp(t *testing.T) {
	setupTest(t)
	*clientIp = "1.2.3"

	err := runIt()
	assertError(t, err, "1.2.3 is not a valid IP")

	*clientIp = "1.2.3.4"
	go runIt()
	waitForIt(t)
	_, err = http.Get("http://127.0.0.1:9999/api/check")
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
}

func TestIptables(t *testing.T) {
	setupTest(t)
	s, err := readIptables()
	if err != nil {
		t.Fatalf("readIptables: %v", err)
	}
	if s != nil {
		t.Fatalf("iptables not empty: %v", s)
	}

	testConfig := Config{Ip: net.IPv4(12, 12, 12, 12), Port: 8080}
	err = writeIptables(&testConfig)
	if err != nil {
		t.Fatalf("writeIptables: %v", err)
	}

	s, err = readIptables()
	if err != nil {
		t.Fatalf("readIptables2: %v", err)
	}
	if !s.Equal(&testConfig) {
		t.Fatalf("readIptables2 output: %v", s)
	}

	err = clearIptables()
	if err != nil {
		t.Fatalf("clearIptables: %v", err)
	}

	s, err = readIptables()
	if err != nil {
		t.Fatalf("readIptables3: %v", err)
	}
	if s != nil {
		t.Fatalf("iptables not empty: %v", s)
	}
}

func TestAPIReadEmpty(t *testing.T) {
	setupTest(t)
	go runIt()
	waitForIt(t)

	r, err := http.Get("http://127.0.0.1:9999/api/config")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if r.StatusCode != 200 {
		t.Fatalf("StatusCode unexpected: %v", r.StatusCode)
	}
	rs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if string(rs) != `{"Config":null,"YourIP":"127.0.0.1"}` {
		t.Fatalf("response unexpected: %v", string(rs))
	}
}

func TestAPIReadSomething(t *testing.T) {
	setupTest(t)
	go runIt()
	waitForIt(t)

	testConfig := Config{Ip: net.IPv4(12, 12, 12, 12), Port: 8080}
	err := writeIptables(&testConfig)
	if err != nil {
		t.Fatalf("writeIptables: %v", err)
	}

	r, err := http.Get("http://127.0.0.1:9999/api/config")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if r.StatusCode != 200 {
		t.Fatalf("StatusCode unexpected: %v", r.StatusCode)
	}
	rs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if string(rs) != `{"Config":{"Ip":"12.12.12.12","Port":8080},"YourIP":"127.0.0.1"}` {
		t.Fatalf("response unexpected: %v", string(rs))
	}
}

func TestAPIReadIpHeader(t *testing.T) {
	setupTest(t)
	go runIt()
	waitForIt(t)

	req, err := http.NewRequest("GET", "http://127.0.0.1:9999/api/config", nil)
	if err != nil {
		t.Fatalf("req: %v", err)
	}
	req.Header.Add("x-real-ip", "1.1.1.1")
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if r.StatusCode != 200 {
		t.Fatalf("StatusCode unexpected: %v", r.StatusCode)
	}
	rs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if string(rs) != `{"Config":null,"YourIP":"1.1.1.1"}` {
		t.Fatalf("response unexpected: %v", string(rs))
	}

	req, err = http.NewRequest("GET", "http://127.0.0.1:9999/api/config", nil)
	if err != nil {
		t.Fatalf("req: %v", err)
	}
	req.Header.Add("x-real-ip", "1.1.1.1,2.2.2.2")
	r, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if r.StatusCode != 200 {
		t.Fatalf("StatusCode unexpected: %v", r.StatusCode)
	}
	rs, err = ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if string(rs) != `{"Config":null,"YourIP":"1.1.1.1"}` {
		t.Fatalf("response unexpected: %v", string(rs))
	}
}

func TestAPIUpdateFailNoHeader(t *testing.T) {
	setupTest(t)
	go runIt()
	waitForIt(t)

	r, err := http.Post("http://127.0.0.1:9999/api/config", "application/json", strings.NewReader(`{"Ip":"1.1.1.1","Port":8080}`))
	if err != nil {
		t.Fatalf("post failed: %v", err)
	}
	// should return 503 as it missed header
	if r.StatusCode != 503 {
		t.Fatalf("StatusCode unexpected: %v", r.StatusCode)
	}
	// data should not have changed!
	s, err := readIptables()
	if err != nil {
		t.Fatalf("readIptables: %v", err)
	}
	if s != nil {
		t.Fatalf("iptables not empty: %v", s)
	}
}

func TestAPIUpdate(t *testing.T) {
	setupTest(t)
	go runIt()
	waitForIt(t)

	req, err := http.NewRequest("POST", "http://127.0.0.1:9999/api/config", strings.NewReader(`{"Ip":"1.1.1.1","Port":8080}`))
	if err != nil {
		t.Fatalf("req: %v", err)
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("post failed: %v", err)
	}
	if r.StatusCode != 200 {
		t.Fatalf("StatusCode unexpected: %v", r.StatusCode)
	}
	// data should WAS changed!
	s, err := readIptables()
	if err != nil {
		t.Fatalf("readIptables: %v", err)
	}
	if !s.Equal(&Config{Ip: net.IPv4(1, 1, 1, 1), Port: 8080}) {
		t.Fatalf("readIptables unexpected: %v", s)
	}
}

func TestAPIDisable(t *testing.T) {
	setupTest(t)
	go runIt()
	waitForIt(t)

	testConfig := Config{Ip: net.IPv4(12, 12, 12, 12), Port: 8080}
	err := writeIptables(&testConfig)
	if err != nil {
		t.Fatalf("writeIptables: %v", err)
	}

	req, err := http.NewRequest("POST", "http://127.0.0.1:9999/api/config", strings.NewReader(`null`))
	if err != nil {
		t.Fatalf("req: %v", err)
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("post failed: %v", err)
	}
	if r.StatusCode != 200 {
		t.Fatalf("StatusCode unexpected: %v", r.StatusCode)
	}
	// data should WAS changed!
	s, err := readIptables()
	if err != nil {
		t.Fatalf("readIptables: %v", err)
	}
	if s != nil {
		t.Fatalf("readIptables unexpected: %v", s)
	}
}
