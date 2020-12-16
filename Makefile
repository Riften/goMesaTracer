.PHONY: lib
.PHONY: install
lib:
	go build -buildmode c-shared -o build/libMesaTracer.so github.com/Riften/goMesaTracer
install : lib
	mkdir -p /usr/local/include/CGO
	cp build/libMesaTracer.h /usr/local/include/CGO/
	chmod +x build/libMesaTracer.so
	cp build/libMesaTracer.so /usr/local/lib/x86_64-linux-gnu/
