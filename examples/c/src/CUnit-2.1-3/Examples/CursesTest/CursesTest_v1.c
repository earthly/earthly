/*
 *  CUnit - A Unit testing framework library for C.
 *  Copyright (C) 2001  Anil Kumar
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

#include "CUCurses.h"

int success_init(void) { return 0; }
int success_clean(void) { return 0; }

void testSuccess1(void) { ASSERT(1); }
void testSuccess2(void) { ASSERT(1); }
void testSuccess3(void) { ASSERT(1); }

void testfailure1(void) { ASSERT(0); }
void testfailure2(void) { ASSERT(2); }
void testfailure3(void) { ASSERT(3); }

int group_failure_init(void) { return 1;}
int group_failure_clean(void) { return 1; }

void testGroupFailure1(void) { ASSERT(0); }
void testGroupFailure2(void) { ASSERT(2); }

int init(void) { return 0; }
int clean(void) { return 0; }

void test1(void)
{
	ASSERT((char *)2 != "THis is positive test.");
	ASSERT((char *)2 == "THis is negative test. test 1");
}

void test2(void)
{
	ASSERT((char *)2 != "THis is positive test.");
	ASSERT((char *)3 == "THis is negative test. test 2");
}

void testSuccessSimpleAssert(void)
{
	ASSERT(1) ;
	ASSERT(!0) ;
}

void testSuccessAssertTrue(void)
{
	ASSERT_TRUE(TRUE) ;
	ASSERT_TRUE(!FALSE) ;
}

void testSuccessAssertFalse(void)
{
	ASSERT_FALSE(FALSE) ;
	ASSERT_FALSE(!TRUE) ;
}

void testSuccessAssertEqual(void)
{
	ASSERT_EQUAL(10, 10) ;
	ASSERT_EQUAL(0, 0) ;
	ASSERT_EQUAL(0, -0) ;
	ASSERT_EQUAL(-12, -12) ;
}

void testSuccessAssertNotEqual(void)
{
	ASSERT_NOT_EQUAL(10, 11) ;
	ASSERT_NOT_EQUAL(0, -1) ;
	ASSERT_NOT_EQUAL(-12, -11) ;
}

void testSuccessAssertPtrEqual(void)
{
	ASSERT_PTR_EQUAL((void*)0x100, (void*)0x101) ;
}

void testSuccessAssertPtrNotEqual(void)
{
	ASSERT_PTR_NOT_EQUAL((void*)0x100, (void*)0x100) ;
}

void testSuccessAssertPtrNull(void)
{
	ASSERT_PTR_NULL((void*)0x23) ;
}

void testSuccessAssertPtrNotNull(void)
{
	ASSERT_PTR_NOT_NULL(NULL) ;
	ASSERT_PTR_NOT_NULL(0x0) ;
}

void testSuccessAssertStringEqual(void)
{
	char str1[] = "test" ;
	char str2[] = "test" ;

	ASSERT_STRING_EQUAL(str1, str2) ;
}

void testSuccessAssertStringNotEqual(void)
{
	char str1[] = "test" ;
	char str2[] = "testtsg" ;

	ASSERT_STRING_NOT_EQUAL(str1, str2) ;
}

void testSuccessAssertNStringEqual(void)
{
	char str1[] = "test" ;
	char str2[] = "testgfsg" ;

	ASSERT_NSTRING_EQUAL(str1, str2, strlen(str1)) ;
	ASSERT_NSTRING_EQUAL(str1, str1, strlen(str1)) ;
	ASSERT_NSTRING_EQUAL(str1, str1, strlen(str1) + 1) ;
}

void testSuccessAssertNStringNotEqual(void)
{
	char str1[] = "test" ;
	char str2[] = "teet" ;
	char str3[] = "testgfsg" ;

	ASSERT_NSTRING_NOT_EQUAL(str1, str2, 3) ;
	ASSERT_NSTRING_NOT_EQUAL(str1, str3, strlen(str1) + 1) ;
}

void testSuccessAssertDoubleEqual(void)
{
	ASSERT_DOUBLE_EQUAL(10, 10.0001, 0.0001) ;
	ASSERT_DOUBLE_EQUAL(10, 10.0001, -0.0001) ;
	ASSERT_DOUBLE_EQUAL(-10, -10.0001, 0.0001) ;
	ASSERT_DOUBLE_EQUAL(-10, -10.0001, -0.0001) ;
}

void testSuccessAssertDoubleNotEqual(void)
{
	ASSERT_DOUBLE_NOT_EQUAL(10, 10.001, 0.0001) ;
	ASSERT_DOUBLE_NOT_EQUAL(10, 10.001, -0.0001) ;
	ASSERT_DOUBLE_NOT_EQUAL(-10, -10.001, 0.0001) ;
	ASSERT_DOUBLE_NOT_EQUAL(-10, -10.001, -0.0001) ;
}

void AddTests(void)
{
	PTestGroup pGroup = NULL;
	PTestCase pTest = NULL;

	pGroup = add_test_group("Sucess", success_init, success_clean);
	pTest = add_test_case(pGroup, "testSuccess1", testSuccess1);
	pTest = add_test_case(pGroup, "testSuccess2", testSuccess2);
	pTest = add_test_case(pGroup, "testSuccess3", testSuccess3);

	pGroup = add_test_group("failure", NULL, NULL);
	pTest = add_test_case(pGroup, "testfailure1", testfailure1);
	pTest = add_test_case(pGroup, "testfailure2", testfailure2);
	pTest = add_test_case(pGroup, "testfailure3", testfailure3);

	pGroup = add_test_group("group_failure", group_failure_init, group_failure_clean);
	pTest = add_test_case(pGroup, "testGroupFailure1", testGroupFailure1);
	pTest = add_test_case(pGroup, "testGroupFailure2", testGroupFailure2);
}

void AddAssertTests(void)
{
	PTestGroup pGroup = NULL;
	PTestCase pTest = NULL;

	pGroup = add_test_group("TestSimpleAssert", NULL, NULL);
	pTest = add_test_case(pGroup, "testSuccessSimpleAssert", testSuccessSimpleAssert);

	pGroup = add_test_group("TestBooleanAssert", NULL, NULL);
	pTest = add_test_case(pGroup, "testSuccessAssertTrue", testSuccessAssertTrue);
	pTest = add_test_case(pGroup, "testSuccessAssertFalse", testSuccessAssertFalse);

	pGroup = add_test_group("TestEqualityAssert", NULL, NULL);
	pTest = add_test_case(pGroup, "testSuccessAssertEqual", testSuccessAssertEqual);
	pTest = add_test_case(pGroup, "testSuccessAssertNotEqual", testSuccessAssertNotEqual);

	pGroup = add_test_group("TestPointerAssert", NULL, NULL);
	pTest = add_test_case(pGroup, "testSuccessAssertPtrEqual", testSuccessAssertPtrEqual);
	pTest = add_test_case(pGroup, "testSuccessAssertPtrNotEqual", testSuccessAssertPtrNotEqual);

	pGroup = add_test_group("TestNullnessAssert", NULL, NULL);
	pTest = add_test_case(pGroup, "testSuccessAssertPtrNull", testSuccessAssertPtrNull);
	pTest = add_test_case(pGroup, "testSuccessAssertPtrNotNull", testSuccessAssertPtrNotNull);

	pGroup = add_test_group("TestStringAssert", NULL, NULL);
	pTest = add_test_case(pGroup, "testSuccessAssertStringEqual", testSuccessAssertStringEqual);
	pTest = add_test_case(pGroup, "testSuccessAssertStringNotEqual", testSuccessAssertStringNotEqual);

	pGroup = add_test_group("TestNStringAssert", NULL, NULL);
	pTest = add_test_case(pGroup, "testSuccessAssertNStringEqual", testSuccessAssertNStringEqual);
	pTest = add_test_case(pGroup, "testSuccessAssertNStringNotEqual", testSuccessAssertNStringNotEqual);

	pGroup = add_test_group("TestDoubleAssert", NULL, NULL);
	pTest = add_test_case(pGroup, "testSuccessAssertDoubleEqual", testSuccessAssertDoubleEqual);
	pTest = add_test_case(pGroup, "testSuccessAssertDoubleNotEqual", testSuccessAssertDoubleNotEqual);
}


int main(int argc, char* argv[])
{
	setvbuf(stdout, NULL, _IONBF, 0);

	if (argc > 1) {
		BOOL Run = FALSE ;
		if (initialize_registry()) {
			printf("\nInitialize of test Registry failed.");
		}

		if (!strcmp("--test", argv[1])) {
			Run = TRUE ;
			AddTests();
		}
		else if (!strcmp("--atest", argv[1])) {
			Run = TRUE ;
			AddAssertTests();
		}
		else if (!strcmp("--alltest", argv[1])) {
			Run = TRUE ;
			AddTests();
			AddAssertTests();
		}

		if (TRUE == Run) {
			curses_run_tests();
		}

		cleanup_registry();
	}

	return 0;
}
