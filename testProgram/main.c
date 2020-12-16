#include "CGO/libMesaTracer.h"
#include "stdio.h"
#include <unistd.h>

int main() {
    cgoAddTrace(CGO_START);
    cgoAddTrace(CGO_START);
    cgoAddTrace(CGO_START);
    cgoAddTrace(CGO_START);
    cgoAddTrace(CGO_START);
    cgoAddTrace(CGO_START);
    cgoAddTrace(CGO_START);
    cgoAddTrace(CGO_END);
    cgoStopAndWait();
    return 0;
}