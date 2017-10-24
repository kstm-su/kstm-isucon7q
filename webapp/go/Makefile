GOPATH := ${PWD}
export GOPATH

deploy:
		sudo systemctl stop isubata.golang
		go build isubata
		sudo systemctl start isubata.golang

build:
		go build -v isubata

vet:
		go vet ./src/isubata/...
