/*
 *  CUnit - A Unit testing framework library for C.
 *  Copyright (C) 2001        Anil Kumar
 *  Copyright (C) 2004, 2005  Anil Kumar, Jerry St.Clair
 *
 *  This library is free software; you can redistribute it and/or
 *  modify it under the terms of the GNU Library General Public
 *  License as published by the Free Software Foundation; either
 *  version 2 of the License, or (at your option) any later version.
 *
 *  This library is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 *  Library General Public License for more details.
 *
 *  You should have received a copy of the GNU Library General Public
 *  License along with this library; if not, write to the Free Software
 *  Foundation, Inc., 59 Temple Place, Suite 330, Boston, MA  02111-1307  USA
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>

#include "CUnit.h"
#include "ExampleTests.h"

/* WARNING - MAINTENANCE NIGHTMARE AHEAD
 *
 * If you change any of the tests & suites below, you also need
 * to keep track of changes in the result statistics and reflect
 * any changes in the result report counts in print_example_results().
 *
 * Yes, this could have been designed better using a more
 * automated mechanism.  No, it was not done that way.
 */

/* Suite initialization/cleanup functions */
static int suite_success_init(void) { return 0; }
static int suite_success_clean(void) { return 0; }

static int suite_failure_init(void) { return 1;}
static int suite_failure_clean(void) { return 1; }

static void testSuccess1(void) { CU_ASSERT(1); }
static void testSuccess2(void) { CU_ASSERT(2); }
static void testSuccess3(void) { CU_ASSERT(3); }

static void testSuiteFailure1(void) { CU_ASSERT(0); }
static void testSuiteFailure2(void) { CU_ASSERT(2); }

static void testFailure1(void) { CU_ASSERT(0); }
static void testFailure2(void) { CU_ASSERT(0); }
static void testFailure3(void) { CU_ASSERT(0); }

/* Test functions for CUnit assertions */
static void testSimpleAssert(void)
{
  CU_ASSERT(1);
  CU_ASSERT(!0);
  CU_TEST(1);

  CU_ASSERT(0);
  CU_ASSERT(!1);
  CU_TEST(0);
}

static void testFail(void)
{
  CU_FAIL("This is a failure.");
  CU_FAIL("This is another failure.");
}

static void testAssertTrue(void)
{
  CU_ASSERT_TRUE(CU_TRUE);
  CU_ASSERT_TRUE(!CU_FALSE);

  CU_ASSERT_TRUE(!CU_TRUE);
  CU_ASSERT_TRUE(CU_FALSE);
}

static void testAssertFalse(void)
{
  CU_ASSERT_FALSE(CU_FALSE);
  CU_ASSERT_FALSE(!CU_TRUE);

  CU_ASSERT_FALSE(!CU_FALSE);
  CU_ASSERT_FALSE(CU_TRUE);
}

static void testAssertEqual(void)
{
  CU_ASSERT_EQUAL(10, 10);
  CU_ASSERT_EQUAL(0, 0);
  CU_ASSERT_EQUAL(0, -0);
  CU_ASSERT_EQUAL(-12, -12);

  CU_ASSERT_EQUAL(10, 11);
  CU_ASSERT_EQUAL(0, 1);
  CU_ASSERT_EQUAL(0, -1);
  CU_ASSERT_EQUAL(-12, 12);
}

static void testAssertNotEqual(void)
{
  CU_ASSERT_NOT_EQUAL(10, 11);
  CU_ASSERT_NOT_EQUAL(0, -1);
  CU_ASSERT_NOT_EQUAL(-12, -11);

  CU_ASSERT_NOT_EQUAL(10, 10);
  CU_ASSERT_NOT_EQUAL(0, -0);
  CU_ASSERT_NOT_EQUAL(0, 0);
  CU_ASSERT_NOT_EQUAL(-12, -12);
}

static void testAssertPtrEqual(void)
{
  CU_ASSERT_PTR_EQUAL((void*)0x100, (void*)0x100);

  CU_ASSERT_PTR_EQUAL((void*)0x100, (void*)0x101);
}

static void testAssertPtrNotEqual(void)
{
  CU_ASSERT_PTR_NOT_EQUAL((void*)0x100, (void*)0x101);

  CU_ASSERT_PTR_NOT_EQUAL((void*)0x100, (void*)0x100);
}

static void testAssertPtrNull(void)
{
  CU_ASSERT_PTR_NULL((void*)(NULL));
  CU_ASSERT_PTR_NULL((void*)(0x0));

  CU_ASSERT_PTR_NULL((void*)0x23);
}

static void testAssertPtrNotNull(void)
{
  CU_ASSERT_PTR_NOT_NULL((void*)0x100);

  CU_ASSERT_PTR_NOT_NULL(NULL);
  CU_ASSERT_PTR_NOT_NULL((void*)0x0);
}

static void testAssertStringEqual(void)
{
  char str1[] = "test";
  char str2[] = "test";
  char str3[] = "suite";

  CU_ASSERT_STRING_EQUAL(str1, str2);

  CU_ASSERT_STRING_EQUAL(str1, str3);
  CU_ASSERT_STRING_EQUAL(str3, str2);
}

static void testAssertStringNotEqual(void)
{
  char str1[] = "test";
  char str2[] = "test";
  char str3[] = "suite";

  CU_ASSERT_STRING_NOT_EQUAL(str1, str3);
  CU_ASSERT_STRING_NOT_EQUAL(str3, str2);

  CU_ASSERT_STRING_NOT_EQUAL(str1, str2);
}

static void testAssertNStringEqual(void)
{
  char str1[] = "test";
  char str2[] = "testgfsg";
  char str3[] = "tesgfsg";

  CU_ASSERT_NSTRING_EQUAL(str1, str2, strlen(str1));
  CU_ASSERT_NSTRING_EQUAL(str1, str1, strlen(str1));
  CU_ASSERT_NSTRING_EQUAL(str1, str1, strlen(str1) + 1);

  CU_ASSERT_NSTRING_EQUAL(str2, str3, 4);
  CU_ASSERT_NSTRING_EQUAL(str1, str3, strlen(str1));
}

static void testAssertNStringNotEqual(void)
{
  char str1[] = "test";
  char str2[] = "tevt";
  char str3[] = "testgfsg";

  CU_ASSERT_NSTRING_NOT_EQUAL(str1, str2, 3);
  CU_ASSERT_NSTRING_NOT_EQUAL(str1, str3, strlen(str1) + 1);

  CU_ASSERT_NSTRING_NOT_EQUAL(str1, str2, 2);
  CU_ASSERT_NSTRING_NOT_EQUAL(str2, str3, 2);
}

static void testAssertDoubleEqual(void)
{
  CU_ASSERT_DOUBLE_EQUAL(10, 10.0001, 0.0001);
  CU_ASSERT_DOUBLE_EQUAL(10, 10.0001, -0.0001);
  CU_ASSERT_DOUBLE_EQUAL(-10, -10.0001, 0.0001);
  CU_ASSERT_DOUBLE_EQUAL(-10, -10.0001, -0.0001);

  CU_ASSERT_DOUBLE_EQUAL(10, 10.0001, 0.00001);
  CU_ASSERT_DOUBLE_EQUAL(10, 10.0001, -0.00001);
  CU_ASSERT_DOUBLE_EQUAL(-10, -10.0001, 0.00001);
  CU_ASSERT_DOUBLE_EQUAL(-10, -10.0001, -0.00001);
}

static void testAssertDoubleNotEqual(void)
{
  CU_ASSERT_DOUBLE_NOT_EQUAL(10, 10.001, 0.0001);
  CU_ASSERT_DOUBLE_NOT_EQUAL(10, 10.001, -0.0001);
  CU_ASSERT_DOUBLE_NOT_EQUAL(-10, -10.001, 0.0001);
  CU_ASSERT_DOUBLE_NOT_EQUAL(-10, -10.001, -0.0001);

  CU_ASSERT_DOUBLE_NOT_EQUAL(10, 10.001, 0.01);
  CU_ASSERT_DOUBLE_NOT_EQUAL(10, 10.001, -0.01);
  CU_ASSERT_DOUBLE_NOT_EQUAL(-10, -10.001, 0.01);
  CU_ASSERT_DOUBLE_NOT_EQUAL(-10, -10.001, -0.01);
}

static void testAbort(void)
{
  CU_TEST_FATAL(CU_TRUE);
  CU_TEST_FATAL(CU_FALSE);
  fprintf(stderr, "\nFatal assertion failed to abort test in testAbortIndirect1\n");
  exit(1);
}

static void testAbortIndirect(void)
{
  testAbort();
  fprintf(stderr, "\nFatal assertion failed to abort test in testAbortIndirect2\n");
  exit(1);
}

static void testFatal(void)
{
  CU_TEST_FATAL(CU_TRUE);
  testAbortIndirect();
  fprintf(stderr, "\nFatal assertion failed to abort test in testFatal\n");
  exit(1);
}

static CU_TestInfo tests_success[] = {
  { "testSuccess1", testSuccess1 },
  { "testSuccess2", testSuccess2 },
  { "testSuccess3", testSuccess3 },
	CU_TEST_INFO_NULL,
};

static CU_TestInfo tests_failure[] = {
  { "testFailure1", testFailure1 },
  { "testFailure2", testFailure2 },
  { "testFailure3", testFailure3 },
	CU_TEST_INFO_NULL,
};

static CU_TestInfo tests_suitefailure[] = {
  { "testSuiteFailure1", testSuiteFailure1 },
  { "testSuiteFailure2", testSuiteFailure2 },
	CU_TEST_INFO_NULL,
};

static CU_TestInfo tests_simple[] = {
  { "testSimpleAssert", testSimpleAssert },
  { "testFail", testFail },
	CU_TEST_INFO_NULL,
};

static CU_TestInfo tests_bool[] = {
  { "testAssertTrue", testAssertTrue },
  { "testAssertFalse", testAssertFalse },
	CU_TEST_INFO_NULL,
};

static CU_TestInfo tests_equal[] = {
  { "testAssertEqual", testAssertEqual },
  { "testAssertNotEqual", testAssertNotEqual },
	CU_TEST_INFO_NULL,
};

static CU_TestInfo tests_ptr[] = {
  { "testAssertPtrEqual", testAssertPtrEqual },
  { "testAssertPtrNotEqual", testAssertPtrNotEqual },
	CU_TEST_INFO_NULL,
};

static CU_TestInfo tests_null[] = {
  { "testAssertPtrNull", testAssertPtrNull },
  { "testAssertPtrNotNull", testAssertPtrNotNull },
	CU_TEST_INFO_NULL,
};

static CU_TestInfo tests_string[] = {
  { "testAssertStringEqual", testAssertStringEqual },
  { "testAssertStringNotEqual", testAssertStringNotEqual },
	CU_TEST_INFO_NULL,
};

static CU_TestInfo tests_nstring[] = {
  { "testAssertNStringEqual", testAssertNStringEqual },
  { "testAssertNStringNotEqual", testAssertNStringNotEqual },
	CU_TEST_INFO_NULL,
};

static CU_TestInfo tests_double[] = {
  { "testAssertDoubleEqual", testAssertDoubleEqual },
  { "testAssertDoubleNotEqual", testAssertDoubleNotEqual },
	CU_TEST_INFO_NULL,
};

static CU_TestInfo tests_fatal[] = {
  { "testFatal", testFatal },
	CU_TEST_INFO_NULL,
};

static CU_SuiteInfo suites[] = {
  { "suite_success_both",  suite_success_init, suite_success_clean, NULL, NULL, tests_success},
  { "suite_success_init",  suite_success_init, NULL,                NULL, NULL, tests_success},
  { "suite_success_clean", NULL,               suite_success_clean, NULL, NULL, tests_success},
  { "test_failure",        NULL,               NULL,                NULL, NULL, tests_failure},
  { "suite_failure_both",  suite_failure_init, suite_failure_clean, NULL, NULL, tests_suitefailure}, /* tests should not run */
  { "suite_failure_init",  suite_failure_init, NULL,                NULL, NULL, tests_suitefailure}, /* tests should not run */
  { "suite_success_but_failure_clean", NULL,   suite_failure_clean, NULL, NULL, tests_suitefailure}, /* tests will run, suite counted as running, but suite tagged as a failure */
  { "TestSimpleAssert",    NULL,               NULL,                NULL, NULL, tests_simple},
  { "TestBooleanAssert",   NULL,               NULL,                NULL, NULL, tests_bool},
  { "TestEqualityAssert",  NULL,               NULL,                NULL, NULL, tests_equal},
  { "TestPointerAssert",   NULL,               NULL,                NULL, NULL, tests_ptr},
  { "TestNullnessAssert",  NULL,               NULL,                NULL, NULL, tests_null},
  { "TestStringAssert",    NULL,               NULL,                NULL, NULL, tests_string},
  { "TestNStringAssert",   NULL,               NULL,                NULL, NULL, tests_nstring},
  { "TestDoubleAssert",    NULL,               NULL,                NULL, NULL, tests_double},
  { "TestFatal",           NULL,               NULL,                NULL, NULL, tests_fatal},
	CU_SUITE_INFO_NULL,
};

void AddTests(void)
{
  assert(NULL != CU_get_registry());
  assert(!CU_is_test_running());

	/* Register suites. */
	if (CU_register_suites(suites) != CUE_SUCCESS) {
		fprintf(stderr, "suite registration failed - %s\n",
			CU_get_error_msg());
		exit(EXIT_FAILURE);
	}

/* implementation without shortcut registration
  CU_pSuite pSuite;

  pSuite = CU_add_suite("suite_success_both", suite_success_init, suite_success_clean);
  CU_add_test(pSuite, "testSuccess1", testSuccess1);
  CU_add_test(pSuite, "testSuccess2", testSuccess2);
  CU_add_test(pSuite, "testSuccess3", testSuccess3);

  pSuite = CU_add_suite("suite_success_init", suite_success_init, NULL);
  CU_add_test(pSuite, "testSuccess1", testSuccess1);
  CU_add_test(pSuite, "testSuccess2", testSuccess2);
  CU_add_test(pSuite, "testSuccess3", testSuccess3);

  pSuite = CU_add_suite("suite_success_clean", NULL, suite_success_clean);
  CU_add_test(pSuite, "testSuccess1", testSuccess1);
  CU_add_test(pSuite, "testSuccess2", testSuccess2);
  CU_add_test(pSuite, "testSuccess3", testSuccess3);

  pSuite = CU_add_suite("test_failure", NULL, NULL);
  CU_add_test(pSuite, "testFailure1", testFailure1);
  CU_add_test(pSuite, "testFailure2", testFailure2);
  CU_add_test(pSuite, "testFailure3", testFailure3);

  / * tests should not run * /
  pSuite = CU_add_suite("suite_failure_both", suite_failure_init, suite_failure_clean);
  CU_add_test(pSuite, "testSuiteFailure1", testSuiteFailure1);
  CU_add_test(pSuite, "testSuiteFailure2", testSuiteFailure2);

  / * tests should not run * /
  pSuite = CU_add_suite("suite_failure_init", suite_failure_init, NULL);
  CU_add_test(pSuite, "testSuiteFailure1", testSuiteFailure1);
  CU_add_test(pSuite, "testSuiteFailure2", testSuiteFailure2);

  / * tests will run, suite counted as running, but suite tagged as a failure * /
  pSuite = CU_add_suite("suite_success_but_failure_clean", NULL, suite_failure_clean);
  CU_add_test(pSuite, "testSuiteFailure1", testSuiteFailure1);
  CU_add_test(pSuite, "testSuiteFailure2", testSuiteFailure2);

  pSuite = CU_add_suite("TestSimpleAssert", NULL, NULL);
  CU_add_test(pSuite, "testSimpleAssert", testSimpleAssert);
  CU_add_test(pSuite, "testFail", testFail);

  pSuite = CU_add_suite("TestBooleanAssert", NULL, NULL);
  CU_add_test(pSuite, "testAssertTrue", testAssertTrue);
  CU_add_test(pSuite, "testAssertFalse", testAssertFalse);

  pSuite = CU_add_suite("TestEqualityAssert", NULL, NULL);
  CU_add_test(pSuite, "testAssertEqual", testAssertEqual);
  CU_add_test(pSuite, "testAssertNotEqual", testAssertNotEqual);

  pSuite = CU_add_suite("TestPointerAssert", NULL, NULL);
  CU_add_test(pSuite, "testAssertPtrEqual", testAssertPtrEqual);
  CU_add_test(pSuite, "testAssertPtrNotEqual", testAssertPtrNotEqual);

  pSuite = CU_add_suite("TestNullnessAssert", NULL, NULL);
  CU_add_test(pSuite, "testAssertPtrNull", testAssertPtrNull);
  CU_add_test(pSuite, "testAssertPtrNotNull", testAssertPtrNotNull);

  pSuite = CU_add_suite("TestStringAssert", NULL, NULL);
  CU_add_test(pSuite, "testAssertStringEqual", testAssertStringEqual);
  CU_add_test(pSuite, "testAssertStringNotEqual", testAssertStringNotEqual);

  pSuite = CU_add_suite("TestNStringAssert", NULL, NULL);
  CU_add_test(pSuite, "testAssertNStringEqual", testAssertNStringEqual);
  CU_add_test(pSuite, "testAssertNStringNotEqual", testAssertNStringNotEqual);

  pSuite = CU_add_suite("TestDoubleAssert", NULL, NULL);
  CU_add_test(pSuite, "testAssertDoubleEqual", testAssertDoubleEqual);
  CU_add_test(pSuite, "testAssertDoubleNotEqual", testAssertDoubleNotEqual);

  pSuite = CU_add_suite("TestFatal", NULL, NULL);
  CU_add_test(pSuite, "testFatal", testFatal);
*/
}

void print_example_results(void)
{
  fprintf(stdout, "\n\nExpected Test Results:"
                  "\n\n  Error Handling  Type      # Run   # Pass   # Fail"
                  "\n\n  ignore errors   suites%9u%9u%9u"
                    "\n                  tests %9u%9u%9u"
                    "\n                  asserts%8u%9u%9u"
                  "\n\n  stop on error   suites%9u%9u%9u"
                    "\n                  tests %9u%9u%9u"
                    "\n                  asserts%8u%9u%9u\n\n",
                  14, 14, 3,
                  31, 10, 21,
                  89, 47, 42,
                  4, 4, 1,
                  12, 9, 3,
                  12, 9, 3);
}
