package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

const (
	API_NUM = 4
)

var (
	API_LIST = []string{
		"http://myexternalip.com/raw",
		"http://api.ipify.org/",
		"https://myexternalip.com/raw",
		"https://api.ipify.org/",
	}
)

func main() {
	args := os.Args
	success := false
	debug := false
	fresh := false
	exist := true

	for _, arg := range args {
		if arg == "-d" {
			debug = true
		} else if arg == "-f" {
			fresh = true
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		if debug {
			log.Fatal(err)
		}
	}

	fp := path.Join(home, ".myexternalip")

	_, err = os.Stat(fp)
	if os.IsNotExist(err) {
		exist = false
	} else if err != nil {
		if debug {
			log.Fatal(err)
		}
	}

	if exist && !fresh {
		file, err := os.OpenFile(fp, os.O_RDONLY, 0666)
		if err != nil {
			if debug {
				log.Fatal(err)
			}
		}
		ip, err := io.ReadAll(file)
		if err != nil {
			if debug {
				log.Fatal(err)
			}
		}

		fmt.Print(string(ip))
		file.Close()

	} else {
		ch := make(chan string, API_NUM)
		var ip string

		for _, api := range API_LIST {
			go getip(api, debug, &ch)
		}

		select {
		case <-time.After(time.Second * 3):

			success = false
			break

		case _ip := <-ch:

			ip = _ip
			success = true
			break

		}

		if success {

			fmt.Print(string(ip))

			if !exist {
				file, err := os.OpenFile(fp, os.O_CREATE|os.O_WRONLY, 0666)
				if err != nil {
					if debug {
						log.Fatal(err)
					}
				}

				_, err = file.WriteString(string(ip))
				if err != nil {
					if debug {
						log.Fatal(err)
					}
				}
				file.Close()

			} else {
				file, err := os.OpenFile(fp, os.O_TRUNC, 0666)
				if err != nil {
					if debug {
						log.Fatal(err)
					}
				}

				_, err = file.WriteString(string(ip))
				if err != nil {
					if debug {
						log.Fatal(err)
					}
				}

				file.Close()
			}

		} else {
			if debug {
				log.Fatal("can't get ip from all apis")
			}
		}
	}
}

func getip(api string, debug bool, ch *chan string) {
	resp, err := http.Get(api)
	if err != nil {
		if debug {
			log.Fatal(err)
		}
	}
	defer resp.Body.Close()
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		if debug {
			log.Fatal(err)
		}
	}
	*ch <- string(ip)
}
