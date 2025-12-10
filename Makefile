.PHONY: clear build run

EXECUTE := mini-gateway.exe


clear:
	@echo clear ${EXECUTE} ...
	@TASKLIST | FINDSTR ${EXECUTE} && TASKKILL /F /IM ${EXECUTE} 
	@IF EXIST ${EXECUTE} del ${EXECUTE}


build:
	@echo build ${EXECUTE} ...
	@go build -ldflags "-w -s" -trimpath

run: build
	@echo run ${EXECUTE} ...
	.\${EXECUTE}

test:
	@go test ./tests -v --count=1

cert:
	@echo generate cert ...
	@openssl req -x509 -newkey rsa:4096 -nodes -out cert.pem -keyout key.pem -days 365