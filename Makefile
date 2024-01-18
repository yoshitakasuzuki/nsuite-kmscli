NAME := nsuite-kmscli

# secp256k1.RecoverPubkeyが別アーキテクチャだとコンパイルできないので
# dockerでプラットフォームを指定してビルドする
bin/linux_x86_64/${NAME}: *.go
	docker run --platform linux/amd64 -it --rm -v `pwd`:/app golang:1.20 bash -c " \
	cd /app && \
	GOOS=linux GOARCH=amd64 go build -o bin/linux_x86_64/nsuite-kmscli . \
	"

bin/darwin_arm64/${NAME}: *.go
	GOOS=darwin GOARCH=arm64 go build -o bin/darwin_arm64/nsuite-kmscli .


bin/${NAME}.linux_x86_64.zip: bin/linux_x86_64/${NAME}
	cd bin && \
	zip -j ${NAME}.linux_x86_64.zip linux_x86_64/${NAME}

bin/${NAME}.darwin_arm64.zip: bin/darwin_arm64/${NAME}
	cd bin && \
	zip -j ${NAME}.darwin_arm64.zip darwin_arm64/${NAME}

package: bin/${NAME}.darwin_arm64.zip bin/${NAME}.linux_x86_64.zip

