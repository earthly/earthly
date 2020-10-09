#include <stdio.h>
#include <string.h>
#include "fibonacci_cunit_test.h"

/*
 * Setup and run tests.
 *
 * @return CUE_SUCCESS if successful, else a CUnit error code if
 * any problems arise.
 */

int main()
{
    if (CUE_SUCCESS != CU_initialize_registry())
    {
        return CU_get_error();
    }

    CU_pSuite suite = 
        CU_add_suite("Fibonacci Suite", initialise_suite, cleanup_suite);
    if (NULL == suite) 
    {
        CU_cleanup_registry();
        return CU_get_error();
    }

    if ((NULL == CU_add_test(suite, "test_fibonacci_1", test_fibonacci_1)) ||
        (NULL == CU_add_test(suite, "test_fibonacci_2", test_fibonacci_2)) ||
        (NULL == CU_add_test(suite, "test_fibonacci_3", test_fibonacci_3)) ||
        (NULL == CU_add_test(suite, "test_fibonacci_30", test_fibonacci_30)))
    {
        CU_cleanup_registry();
        return CU_get_error();
    }

    CU_basic_set_mode(CU_BRM_VERBOSE);
    CU_basic_run_tests();

    CU_list_tests_to_file();

    CU_automated_run_tests();

    CU_cleanup_registry();
    return CU_get_error();
}
