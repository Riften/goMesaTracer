#include "CGO/libMesaTracer.h"
#include "stdio.h"
#include <unistd.h>

typedef const char cchar_t;

int main() {
    cgoAddTrace(1, __func__);
    // cgoAddTrace(2, __func__);
    // cgoAddTrace(3, __func__);
    // cgoAddTrace(4, __func__);
    // cgoAddTrace(5, __func__);
    // cgoAddTrace(6, __func__);
    // cgoAddTrace(7, __func__);
    // cgoAddTrace(8, __func__);
    cgoStopAndWait();
    return 0;
}