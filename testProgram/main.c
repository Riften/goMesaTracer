#include "CGO/libMesaTracer.h"
#include "stdio.h"
#include <unistd.h>

int main() {
    //cgoAddString("test");
    cgoTestEnum(CGO_START);
    cgoTestEnum(CGO_END);
    sleep(3);
    return 0;
}