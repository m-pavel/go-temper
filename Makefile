TEMPER=HID-TEMPerHUM
CGO_CFLAGS="-I${PWD}/${TEMPER}"
CGO_LDFLAGS="-L${PWD}/${TEMPER}"
GF=CGO_CFLAGS=${CGO_CFLAGS} CGO_LDFLAGS=${CGO_LDFLAGS} LD_LIBRARY_PATH=${PWD}/${TEMPER}

all: temper test build

deps:
	${GF} go get -v -d ./...
test: deps
	${GF} go test -v $$(go list ./... | grep -v /vendor/)

build: deps
	${GF} go build -o temper-influx-cli ./influx
	${GF} go build -o temper-mqtt-cli ./mqtt

temper:
	git clone https://github.com/m-pavel/HID-TEMPerHUM
	cd ${TEMPER} && make

temper-clean:
	rm -rf ./${TEMPER}

clean: temper-clean

