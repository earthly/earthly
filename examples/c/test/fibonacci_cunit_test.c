#include "fibonacci_cunit_test.h"

int initialise_suite(void)
{
   return 0;
}

int cleanup_suite(void)
{
    return 0;
}

void test_fibonacci_1(void)
{
    CU_ASSERT_EQUAL(fibonacci(1), 1);
}

void test_fibonacci_2(void)
{
    CU_ASSERT_EQUAL(fibonacci(2), 1);
}

void test_fibonacci_3(void)
{
    CU_ASSERT_EQUAL(fibonacci(3), 2);
}

void test_fibonacci_30(void)
{
    CU_ASSERT_EQUAL(fibonacci(30), 832040);
}
