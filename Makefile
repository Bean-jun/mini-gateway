.PHONY: clear build run

EXECUTE := g-server.exe


clear:
	@echo clear ${EXECUTE} ...
	@TASKLIST | FINDSTR ${EXECUTE} && TASKKILL /F /IM ${EXECUTE} 
	@IF EXIST ${EXECUTE} del ${EXECUTE}


build:
	@echo build ${EXECUTE} ...
	@go build

run: build
	@echo run ${EXECUTE} ...
	.\${EXECUTE}
