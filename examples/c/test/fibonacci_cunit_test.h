#ifndef FIBONACCI_CUNIT_TEST_H
#define FIBONACCI_CUNIT_TEST_H

#include "CUnit/Automated.h"
#include "CUnit/Basic.h"
#include "fibonacci.h"


int initialise_suite(void);


int cleanup_suite(void);

void test_fibonacci_1(void);

void test_fibonacci_2(void);

void test_fibonacci_3(void);

void test_fibonacci_30(void);

#endif
