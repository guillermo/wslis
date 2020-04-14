


ssl.key:
	openssl genrsa -out ssl.key 2048
	openssl ecparam -genkey -name secp384r1 -out ssl.key



ssl.crt: ssl.key
	openssl req -new -x509 -sha256 -key ssl.key -out ssl.crt -days 3650

install-crt: ssl.crt
	sudo cp ssl.crt /etc/ssl/certs

test: install-crt
	go test


.PHONY: test
