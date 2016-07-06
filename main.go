package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/golang/glog"
)

var (
	port    string
	cmdline string
	config  string
)

func initConfig(configFile string) error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		glog.Errorln(err)
		return err
	}

	configData := make(map[string]string)
	err = json.Unmarshal(data, &configData)
	if err != nil {
		glog.Errorln(err)
		return err
	}

	configPort, ok := configData["port"]
	if ok {
		port = configPort
	}

	configCmdline, ok := configData["cmdline"]
	if ok {
		cmdline = configCmdline
	}
	return nil
}

func init() {
	flag.StringVar(&port, "port", ":8080", "listen port")
	flag.StringVar(&config, "config", "", "json config file")
	flag.Parse()

	if config != "" {
		err := initConfig(config)
		if err != nil {
			glog.Fatal(err)
			return
		}
	}

	if cmdline == "" {
		path, err := exec.LookPath("qrcode")
		if err != nil {
			glog.Fatal("set qrcode to your path")
			return
		}
		cmdline = path
	}
}

func decode(fileName string) (string, error) {
	cmd := exec.Command(cmdline, fileName)
	out, err := cmd.Output()
	if err != nil {
		glog.Errorln(err)
		return "", err
	}

	return string(out), nil
}

type ResponseData struct {
	Status bool   `json:"status"`
	ErrMsg string `json:"errMsg"`
	Data   string `json:"data"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	const (
		maxMemory = 32 << 10
	)

	data := &ResponseData{
		Status: false,
		ErrMsg: "Unsupport",
	}

	defer func() {
		v, err := json.Marshal(data)
		if err != nil {
			glog.Fatalln(err)
		}

		w.Write(v)
	}()

	if r.Method == "POST" {
		err := r.ParseMultipartForm(maxMemory)
		if err != nil {
			glog.Errorln(err)
			return
		}

		file, handler, err := r.FormFile("uploadFile")
		if err != nil {
			glog.Errorln(err)
			return
		}
		defer file.Close()

		fileName := os.TempDir() + "/" + handler.Filename

		f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			glog.Errorln(err)
			return
		}
		defer f.Close()
		defer func() {
			os.Remove(fileName)
		}()
		io.Copy(f, file)

		result, err := decode(fileName)
		if err != nil {
			glog.Errorln(err)
			return
		}

		data.Status = true
		data.Data = result
		data.ErrMsg = ""
	}

}

func main() {
	glog.V(2).Infoln("will listen port", port)

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		glog.Fatalln(err)
	}
}
