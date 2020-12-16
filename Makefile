.PHONY: lib
lib:
	go build -buildmode c-shared -o build/libMesaTracer.so github.com/Riften/goMesaTracer