#include <iostream>

#include "fib.h"

int main(int argc, char *argv[]) {
	std::cout << "Hello World!" << std::endl;
	for( int i = 0; i < 5; i++ ) {
		std::cout << "fib(" << i << ") = " << fib(i) << std::endl;
	}
	return 0;
}
