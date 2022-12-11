TARGET_FILE:=${shell head -n1 go.mod | sed -r 's/.*\/(.*)/\1/g' }
BUILD_DIR=.build

.PHONY: target build run-app clean mk-build-dir build-deps build-app build-all

target: build-app

echo:
	@echo $(TARGET_FILE)

clean:
	rm -rf $(TARGET_FILE) $(BUILD_DIR)

mk-build-dir:
	@mkdir -p ${BUILD_DIR}

build-deps:
	@go get -d -v ./...

build-app: clean build-deps test
	go build -o $(TARGET_FILE) cmd/app/main.go

run-app:
	@go build -o $(TARGET_FILE) cmd/app/main.go && go run .

build-all: clean mk-build-dir build-deps test
	GOOS=linux go build -o $(TARGET_FILE) cmd/app/main.go && zip -9 $(TARGET_FILE)-linux64.zip $(TARGET_FILE) && rm $(TARGET_FILE)
	GOOS=windows go build -o $(TARGET_FILE) cmd/app/main.go && zip -9 $(TARGET_FILE)-win64.zip $(TARGET_FILE) && rm $(TARGET_FILE)
	GOOS=darwin go build -o $(TARGET_FILE) cmd/app/main.go && zip -9 $(TARGET_FILE)-osx64.zip $(TARGET_FILE) && rm $(TARGET_FILE)
	mv *.zip ${BUILD_DIR}