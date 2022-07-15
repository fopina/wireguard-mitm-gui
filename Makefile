dev:
	@go run -tags=dev main.go -i ./fake_iptables.sh -s ./fake_iptables_save.sh

test:
	@go test -cover ./...

testv:
	@go test -v ./...

gen:
	@go generate ./...
