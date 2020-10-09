#include <stdio.h>
#include "fibonacci.h"

int main(int argc, char **argv) 
{
    if (argc <= 1)
    {
        print_usage();
        return(1);
    }
    int n = atoi(argv[1]);
    if (n < 0)
    {
        print_usage();
        return(2);
    }
    int fib = fibonacci(n);
    printf("fibonacci(%d) = %d\n", n, fib);
    return(0);
}