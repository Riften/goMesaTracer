# goMesaTracer
A tracer and analyzer for glmark2 and mesa writen in golang.

## Build Shared Library
**Linux Only**

```bash
mkdir build
make lib
```
The built shared library and header file would be in `build` folder. 

## Install Shared Library
**Linux Only**

```bash
sudo make install
```

Simply copy the header file to
```bash
/usr/local/include/CGO/libMesaTracer.h
```

Copy the lib file to
```bash
/usr/local/lib/x86_64-linux-gnu/libMesaTracer.so
```

Make a soft link of lib to
```bash
/usr/local/lib/libMesaTracer.so
```

## Use Shared Library
```cpp
#include "CGO/libMesaTracer.h"

// When you need to add a trace
cgoAddTrace(<CgoType>);

// After the last trace
cgoStopAndWait();
```

`CgoType` is defined in `main.go` as macro definition.

## Build & Run Sample Progarm
Build
```bash
cd testProgram
make
# Executable file would output to /build/test
```

Run
```bash
./test
# or
LD_LIBRARY_PATH=/usr/local/lib ./test
# if the shared lib can not be found
```

## Analyze Tool
**TODO**