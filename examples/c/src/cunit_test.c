#include<stdio.h>
#include "sum.h"

void test_sum(void)
{

CU_ASSERT(sum(10,15) == 25);
CU_ASSERT(sum(3,4) == 7);
CU_ASSERT(sum(6,9) == 15);

}