# C++ example

In this example, we will walk through a simple C++ app that uses CMake.

Let's assume our code is structured as follows:
```
.
├── Earthfile
└── src
    ├── CMakeLists.txt
    ├── fib.cpp
    ├── fib.h
    └── main.cpp
```

Our program will be split between two different `cpp` files; the `main.cpp` file:

```C++
#include <iostream>

#include "fib.h"

int main(int argc, char *argv[]) {
	for( int i = 0; i < 5; i++ ) {
		std::cout << "fib(" << i << ") = " << fib(i) << std::endl;
	}
	return 0;
}
```

and a file containing our fibonacci function:

```C++
#include "fib.h"

int fib(int n)
{
	if( n <= 0 ) {
		return 0;
	}
	if( n == 1 ) {
		return 1;
	}
	return fib(n-1) + fib(n-2);
}
```

We will use CMake to manage the build process of the c++ code, with the following CMakeList.txt file:

```
cmake_minimum_required(VERSION 2.8.9)
project (fibonacci)
add_executable(fibonacci main.cpp fib.cpp)
```

CMake caches object files under CMakeFiles which allows CMake to only recompile objects when the corresponding
source code changes. We will use a [mount-based cache](https://docs.earthly.dev/docs/guides/advanced-local-caching) to cache these temporary
files to allow for faster builds on a local machine. Here's a sample `Earthfile`:

```Dockerfile
# Earthfile
VERSION 0.7
FROM ubuntu:20.10

# configure apt to be noninteractive
ENV DEBIAN_FRONTEND noninteractive
ENV DEBCONF_NONINTERACTIVE_SEEN true

# install dependencies
RUN apt-get update && apt-get install -y build-essential cmake

WORKDIR /code

code:
  COPY src src

build:
  FROM +code
  RUN cmake src
  # cache cmake temp files to prevent rebuilding .o files
  # when the .cpp files don't change
  RUN --mount=type=cache,target=/code/CMakeFiles make
  SAVE ARTIFACT fibonacci AS LOCAL fibonacci

docker:
  COPY +build/fibonacci /bin/fibonacci
  ENTRYPOINT ["/bin/fibonacci"]
  SAVE IMAGE --push earthly/examples:cpp
```

If you run `earthly +build` for the first time you should see:

```
...
+build | Scanning dependencies of target fibonacci
+build | [ 33%] Building CXX object CMakeFiles/fibonacci.dir/main.cpp.o
+build | [ 66%] Building CXX object CMakeFiles/fibonacci.dir/fib.cpp.o
+build | [100%] Linking CXX executable fibonacci
+build | [100%] Built target fibonacci
...
```

However on the next run since the object files were cached you should only see

```
...
+build | Scanning dependencies of target fibonacci
+build | [100%] Linking CXX executable fibonacci
+build | [100%] Built target fibonacci
...
```

If you need to force a full rebuild, you can run earthly `--no-cache +build` to trigger a clean build; however
this will also rebuild the entire base docker images.

And finally, the fibonacci program can be run via docker:

```
~/workspace/earthly/examples/cpp ❯ docker run --rm earthly/examples:cpp
fib(0) = 0
fib(1) = 1
fib(2) = 1
fib(3) = 2
fib(4) = 3
```
