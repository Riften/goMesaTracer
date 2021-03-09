# goMesaTracer
A tracer and analyzer for glmark2 and mesa writen in golang.

## Build Shared Library
**Linux Only**

| Variable | Explanation |
|---|---|
| BUILD_DIR | The build dest folder, default is `build` |
| PREFIX | The prefix folder to install CGO, default is `$HOME/.local` |
| INCLUDE_DIR | The header files will be copied to this folder, default is `$PREFIX/include/CGO` |
| LIB_DIR | The lib files will be copied to this folder, default is `$PREFIX/lib/CGO` |

```shell
make lib
```

## Install Shared Library
**Linux Only**

```shell
make install
```

## Use Shared Library
```cpp
#include <CGO/libMesaTracer.h>

// When you need to add a trace
cgoAddTrace(int counter, char* funcName);

// After the last trace
cgoStopAndWait();
```

`CgoType` is defined in `main.go` as macro definition.

## Build & Run Sample Progarm
Build
```shell
cd testProgram
make
# Executable file would output to testProgram/test
```

Run
```bash
./test
```

## NOTE

Maybe you should configure `C_INCLUDE_PATH`, `CPP_INCLUDE_PATH`, `LD_LIBRARY_PATH` to use the header files and the lib files.

## Analyze Tool
**TODO**