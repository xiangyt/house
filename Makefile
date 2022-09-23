.PHONY: release debug

release:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=mod -o ./release/house .
	docker build -t xiangyt/house:0.0.1 .
	docker push xiangyt/house:0.0.1

debug:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=mod -o ./release/house .
	docker build -t xiangyt/house:mac .
	docker save -o ./release/house-mac.tar xiangyt/house:mac