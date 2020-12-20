.PHONY: lib
.PHONY: install
.PHONY: tracer
.PHONY: inject
lib:
	echo "== Build shared lib libMesaTracer =="
	go build -buildmode c-shared -o build/libMesaTracer.so github.com/Riften/goMesaTracer
tracer:
	echo "== Build cmd tool tracer =="
	go build -o build/tracer github.com/Riften/goMesaTracer
install:
	echo "== Install header and lib file to system folders =="
	sudo mkdir -p /usr/local/include/CGO
	sudo cp build/libMesaTracer.h /usr/local/include/CGO/
	chmod +x build/libMesaTracer.so
	sudo cp build/libMesaTracer.so /usr/local/lib/x86_64-linux-gnu/
	sudo rm -f /usr/local/lib/libMesaTracer.so
	sudo ln -s /usr/local/lib/x86_64-linux-gnu/libMesaTracer.so /usr/local/lib/libMesaTracer.so
inject:
	echo "== Inject header files to glmark2 and mesa =="
	./inject_glmark2.sh
	./inject_mesa.sh
all: lib tracer install inject
	echo "DONE"
