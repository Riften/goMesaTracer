.PHONY: lib tracer install

NAME ?= CGO
BUILD_DIR ?= build
PREFIX ?= ${HOME}/.local
INCLUDE_DIR ?= include
LIB_DIR ?= lib

lib:
	mkdir -p ${BUILD_DIR}
	go build -buildmode c-shared -o ${BUILD_DIR}/libMesaTracer.so github.com/Riften/goMesaTracer
	sed -i "s/extern void cgoAddTrace(int p0, char\* p1);/extern void cgoAddTrace(int p0, const char\* p1);/g" ${BUILD_DIR}/libMesaTracer.h

tracer:
	mkdir -p ${BUILD_DIR}
	go build -o ${BUILD_DIR}/tracer github.com/Riften/goMesaTracer

install:
	mkdir -p ${PREFIX}/${INCLUDE_DIR}/${NAME}
	mkdir -p ${PREFIX}/${LIB_DIR}/${NAME}
	cp ${BUILD_DIR}/libMesaTracer.h ${PREFIX}/${INCLUDE_DIR}/${NAME}/
	cp ${BUILD_DIR}/libMesaTracer.so ${PREFIX}/${LIB_DIR}/${NAME}/

all: lib tracer
	echo "DONE"
