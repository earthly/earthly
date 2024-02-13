#include <stdio.h>
#include <string.h>
#include "fibonacci.h"

void usage() {
  fprintf(stderr, "Usage: ./fibonacci <num>\n");
}

int main(int argc, char** argv) {
  int n;

  if (argc != 2) {
    usage();
    return 1;
  }

  if (sscanf(argv[1], "%d", &n) < 0 || n < 0) {
    fprintf(stderr, "Could not read a positive integer from the input\n");
    usage();
    return 1;
  }

  printf("fib(%d) = %d\n", n, fibonacci(n));
  return 0;
}
