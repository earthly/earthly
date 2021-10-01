/*
 * This example test suite is based on the code on the CUnit project's website.
 * See there for more information: http://cunit.sourceforge.net/example.html
 */

#include <stdio.h>
#include "CUnit/Basic.h"
#include "../src/fibonacci.h"

void testFIBONACCI(void)
{
  CU_ASSERT(fibonacci(0) == 0);
  CU_ASSERT(fibonacci(1) == 1);
  CU_ASSERT(fibonacci(3) == 2);
  CU_ASSERT(fibonacci(15) == 610);
  CU_ASSERT(fibonacci(23) == 28657);
}

int main()
{
  CU_pSuite pSuite = NULL;

  if (CUE_SUCCESS != CU_initialize_registry())
    return CU_get_error();

  pSuite = CU_add_suite("TestSuite", NULL, NULL);
  if (NULL == pSuite) {
    CU_cleanup_registry();
    return CU_get_error();
  }

  if ((NULL == CU_add_test(pSuite, "test of fibonacci()", testFIBONACCI)))
  {
    CU_cleanup_registry();
    return CU_get_error();
  }
  
  CU_basic_set_mode(CU_BRM_VERBOSE);
  CU_basic_run_tests();
  CU_cleanup_registry();
  return CU_get_error();
}
