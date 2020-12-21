.PHONY: lib
.PHONY: install
.PHONY: tracer
.PHONY: inject
lib:
	go build -buildmode c-shared -o build/libMesaTracer.so github.com/Riften/goMesaTracer
tracer:
	go build -o build/tracer github.com/Riften/goMesaTracer
install:
	sudo mkdir -p /usr/local/include/CGO
	sudo cp build/libMesaTracer.h /usr/local/include/CGO/
	chmod +x build/libMesaTracer.so
	sudo cp build/libMesaTracer.so /usr/local/lib/x86_64-linux-gnu/
	sudo rm -f /usr/local/lib/libMesaTracer.so
	sudo ln -s /usr/local/lib/x86_64-linux-gnu/libMesaTracer.so /usr/local/lib/libMesaTracer.so
all: lib tracer install
	echo "DONE"
