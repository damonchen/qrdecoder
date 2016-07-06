# qrdecoder
qrcode decoder with web server(wrapper).

## build

in your `GOPATH`, run the command

```bash
git clone https://github.com/damonchen/qrdecoder
go get github.com/golang/glog
cd qrdecoder
go build
```

## run

download the qrcode decode component from `http://sourceforge.jp/projects/qrcode/files/`, and then unzip it.

suppose your unzip location is 

`/home/damonchen/qrcode`

so the bin path is `/home/damonchen/qrcode/bin` .

run the command below

```bash

mkdir conf
cat <<EOF
{
    "port": ":8090",
    "cmdline": "/home/damonchen/qrcode/bin/qrcode"
}
EOF > conf/app.json

./qrdecorder -config=conf/app.json
```

then you could use 

```bash
curl -F "uploadFile=@test-image.png" http://localhost:8090/
```

to decode your qrcode image through net. 


