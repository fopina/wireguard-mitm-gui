package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/fopina/wireguard-mitm-gui/assets"
	flag "github.com/spf13/pflag"
)

// top level for easier testing
var bindAddress = flag.StringP("bind", "b", "127.0.0.1:8081", "address:port to bind webserver")
var iptablesPath = flag.StringP("iptables-bin", "i", "/usr/sbin/iptables", "Path to iptables")
var iptablesSavePath = flag.StringP("iptables-save-bin", "s", "/sbin/iptables-save", "Path to iptables-save")
var clientIp = flag.StringP("client-ip", "c", "192.168.0.222", "Client IP that should be redirected")

var iptablesParser = regexp.MustCompile(`-A PREROUTING -s (\d+\.\d+\.\d+\.\d+)/32 -p tcp -m tcp --dport (\d+) -j DNAT --to-destination (\d+\.\d+\.\d+\.\d+):(\d+)`)

type Config struct {
	Ip   net.IP
	Port int
}

func (c *Config) IpPort() string {
	return fmt.Sprintf("%s:%d", c.Ip, c.Port)
}

func (c *Config) Equal(other *Config) bool {
	if c == other {
		// same reference, no need to compare anything...
		return true
	}
	if (other == nil) || (c == nil) {
		// if one is nil (and other is not), not equal
		return false
	}
	return c.Ip.Equal(other.Ip) && (c.Port == other.Port)
}

func writeIptables(c *Config) error {
	// iptables -t nat -A PREROUTING -p tcp -s $sip --dport $dport -j DNAT --to-destination $dest
	for _, port := range []string{"80", "443", "8080"} {
		cmd := exec.Command(*iptablesPath, "-t", "nat", "-A", "PREROUTING", "-p", "tcp", "-s", *clientIp, "--dport", port, "-j", "DNAT", "--to-destination", c.IpPort())
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func parseIptables(process func(line []string) (bool, error)) error {
	cmd := exec.Command(*iptablesSavePath)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	for scanner.Scan() {
		data := iptablesParser.FindStringSubmatch(scanner.Text())
		if data != nil {
			if data[1] == *clientIp {
				done, err := process(data)
				if err != nil {
					return err
				}
				if done {
					break
				}
			}
		}
	}
	return nil
}

func clearIptables() error {
	err := parseIptables(func(data []string) (bool, error) {
		cmd := exec.Command(*iptablesPath, "-t", "nat", "-D", "PREROUTING", "-p", "tcp", "-s", data[1], "--dport", data[2], "-j", "DNAT", "--to-destination", fmt.Sprintf("%s:%s", data[3], data[4]))
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return false, err
		}
		return false, nil
	})
	return err
}

func readIptables() (*Config, error) {
	var c *Config
	err := parseIptables(func(data []string) (bool, error) {
		port, err := strconv.Atoi(data[4])
		if err != nil {
			return false, err
		}
		c = &Config{Ip: net.ParseIP(data[3]), Port: port}
		return true, nil
	})
	return c, err
}

func runIt() error {
	if net.ParseIP(*clientIp) == nil {
		return fmt.Errorf("%v is not a valid IP", *clientIp)
	}
	http.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			config, err := readIptables()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			c, err := json.Marshal(map[string]interface{}{"Config": config, "YourIP": strings.Split(r.RemoteAddr, ":")[0]})
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			fmt.Fprint(w, string(c))
		case "POST":
			// JSON payload already protects from CSRF but checking header does not hurt
			val, ok := r.Header["X-Requested-With"]
			if ok && (val[0] == "XMLHttpRequest") {
				// all good
			} else {
				http.Error(w, "", 503)
				return
			}

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			var c Config
			err = json.Unmarshal(body, &c)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			clearIptables()
			writeIptables(&c)
		default:
			http.Error(w, "Method not allowed", 405)
		}
	})

	http.Handle("/", http.FileServer(assets.Assets))
	log.Println("Listening on http://" + *bindAddress)
	return http.ListenAndServe(*bindAddress, nil)
}

func main() {
	flag.Parse()
	err := runIt()
	if err != nil {
		log.Fatal(err)
	}
}
