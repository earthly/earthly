#include <stdio.h>

#include "sum.h"

int main(int argc, char *argv[]) {
	int n1 = 100;
    int n2 = 200;
    int s = sum(n1, n2);
    printf("The sum is: %d\n", s);
	return 0;
}