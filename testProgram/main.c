#include "CGO/libMesaTracer.h"

int main() {
    cgoAddString("test");
    cgoTestEnum(CGO_START);
    cgoTestEnum(CGO_END);
    return 0;
}