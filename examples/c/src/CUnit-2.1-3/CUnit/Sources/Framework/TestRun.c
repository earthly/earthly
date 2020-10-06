/*
 *  CUnit - A Unit testing framework library for C.
 *  Copyright (C) 2001       Anil Kumar
 *  Copyright (C) 2004-2006  Anil Kumar, Jerry St.Clair
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

/*
 *  Implementation of Test Run Interface.
 *
 *  Aug 2001      Initial implementaion (AK)
 *
 *  19/Aug/2001   Added initial registry/Suite/test framework implementation. (AK)
 *
 *  24/Aug/2001   Changed Data structure from SLL to DLL for all linked lists. (AK)
 *
 *  25/Nov/2001   Added notification for Suite Initialization failure condition. (AK)
 *
 *  5-Aug-2004    New interface, doxygen comments, moved add_failure on suite
 *                initialization so called even if a callback is not registered,
 *                moved CU_assertImplementation into TestRun.c, consolidated
 *                all run summary info out of CU_TestRegistry into TestRun.c,
 *                revised counting and reporting of run stats to cleanly
 *                differentiate suite, test, and assertion failures. (JDS)
 *
 *  1-Sep-2004    Modified CU_assertImplementation() and run_single_test() for
 *                setjmp/longjmp mechanism of aborting test runs, add asserts in
 *                CU_assertImplementation() to trap use outside a registered
 *                test function during an active test run. (JDS)
 *
 *  22-Sep-2004   Initial implementation of internal unit tests, added nFailureRecords
 *                to CU_Run_Summary, added CU_get_n_failure_records(), removed
 *                requirement for registry to be initialized in order to run
 *                CU_run_suite() and CU_run_test(). (JDS)
 *
 *  30-Apr-2005   Added callback for suite cleanup function failure,
 *                updated unit tests. (JDS)
 *
 *  23-Apr-2006   Added testing for suite/test deactivation, changing functions.
 *                Moved doxygen comments for public functions into header.
 *                Added type marker to CU_FailureRecord.
 *                Added support for tracking inactive suites/tests. (JDS)
 *
 *  02-May-2006   Added internationalization hooks.  (JDS)
 *
 *  02-Jun-2006   Added support for elapsed time.  Added handlers for suite
 *                start and complete events.  Reworked test run routines to
 *                better support these features, suite/test activation. (JDS)
 *
 *  16-Avr-2007   Added setup and teardown functions. (CJN)
 *
 */

/** @file
 *  Test run management functions (implementation).
 */
/** @addtogroup Framework
 @{
*/

#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include <stdio.h>
#include <setjmp.h>
#include <time.h>

#include "CUnit.h"
#include "MyMem.h"
#include "TestDB.h"
#include "TestRun.h"
#include "Util.h"
#include "CUnit_intl.h"

/*=================================================================
 *  Global/Static Definitions
 *=================================================================*/
static CU_BOOL   f_bTestIsRunning = CU_FALSE; /**< Flag for whether a test run is in progress */
static CU_pSuite f_pCurSuite = NULL;          /**< Pointer to the suite currently being run. */
static CU_pTest  f_pCurTest  = NULL;          /**< Pointer to the test currently being run. */

/** CU_RunSummary to hold results of each test run. */
static CU_RunSummary f_run_summary = {"", 0, 0, 0, 0, 0, 0, 0, 0, 0, 0};

/** CU_pFailureRecord to hold head of failure record list of each test run. */
static CU_pFailureRecord f_failure_list = NULL;

/** CU_pFailureRecord to hold head of failure record list of each test run. */
static CU_pFailureRecord f_last_failure = NULL;

/** Flag for whether inactive suites/tests are treated as failures. */
static CU_BOOL f_failure_on_inactive = CU_TRUE;

/** Variable for storage of start time for test run. */
static clock_t f_start_time;


/** Pointer to the function to be called before running a suite. */
static CU_SuiteStartMessageHandler          f_pSuiteStartMessageHandler = NULL;

/** Pointer to the function to be called before running a test. */
static CU_TestStartMessageHandler           f_pTestStartMessageHandler = NULL;

/** Pointer to the function to be called after running a test. */
static CU_TestCompleteMessageHandler        f_pTestCompleteMessageHandler = NULL;

/** Pointer to the function to be called after running a suite. */
static CU_SuiteCompleteMessageHandler       f_pSuiteCompleteMessageHandler = NULL;

/** Pointer to the function to be called when all tests have been run. */
static CU_AllTestsCompleteMessageHandler    f_pAllTestsCompleteMessageHandler = NULL;

/** Pointer to the function to be called if a suite initialization function returns an error. */
static CU_SuiteInitFailureMessageHandler    f_pSuiteInitFailureMessageHandler = NULL;

/** Pointer to the function to be called if a suite cleanup function returns an error. */
static CU_SuiteCleanupFailureMessageHandler f_pSuiteCleanupFailureMessageHandler = NULL;

/*=================================================================
 * Private function forward declarations
 *=================================================================*/
static void         clear_previous_results(CU_pRunSummary pRunSummary, CU_pFailureRecord* ppFailure);
static void         cleanup_failure_list(CU_pFailureRecord* ppFailure);
static CU_ErrorCode run_single_suite(CU_pSuite pSuite, CU_pRunSummary pRunSummary);
static CU_ErrorCode run_single_test(CU_pTest pTest, CU_pRunSummary pRunSummary);
static void         add_failure(CU_pFailureRecord* ppFailure,
                                CU_pRunSummary pRunSummary,
                                CU_FailureType type,
                                unsigned int uiLineNumber,
                                const char *szCondition,
                                const char *szFileName,
                                CU_pSuite pSuite,
                                CU_pTest pTest);

/*=================================================================
 *  Public Interface functions
 *=================================================================*/
CU_BOOL CU_assertImplementation(CU_BOOL bValue,
                                unsigned int uiLine,
                                const char *strCondition,
                                const char *strFile,
                                const char *strFunction,
                                CU_BOOL bFatal)
{
  /* not used in current implementation - stop compiler warning */
  CU_UNREFERENCED_PARAMETER(strFunction);

  /* these should always be non-NULL (i.e. a test run is in progress) */
  assert(NULL != f_pCurSuite);
  assert(NULL != f_pCurTest);

  ++f_run_summary.nAsserts;
  if (CU_FALSE == bValue) {
    ++f_run_summary.nAssertsFailed;
    add_failure(&f_failure_list, &f_run_summary, CUF_AssertFailed,
                uiLine, strCondition, strFile, f_pCurSuite, f_pCurTest);

    if ((CU_TRUE == bFatal) && (NULL != f_pCurTest->pJumpBuf)) {
      longjmp(*(f_pCurTest->pJumpBuf), 1);
    }
  }

  return bValue;
}

/*------------------------------------------------------------------------*/
void CU_set_suite_start_handler(CU_SuiteStartMessageHandler pSuiteStartHandler)
{
  f_pSuiteStartMessageHandler = pSuiteStartHandler;
}

/*------------------------------------------------------------------------*/
void CU_set_test_start_handler(CU_TestStartMessageHandler pTestStartHandler)
{
  f_pTestStartMessageHandler = pTestStartHandler;
}

/*------------------------------------------------------------------------*/
void CU_set_test_complete_handler(CU_TestCompleteMessageHandler pTestCompleteHandler)
{
  f_pTestCompleteMessageHandler = pTestCompleteHandler;
}

/*------------------------------------------------------------------------*/
void CU_set_suite_complete_handler(CU_SuiteCompleteMessageHandler pSuiteCompleteHandler)
{
  f_pSuiteCompleteMessageHandler = pSuiteCompleteHandler;
}

/*------------------------------------------------------------------------*/
void CU_set_all_test_complete_handler(CU_AllTestsCompleteMessageHandler pAllTestsCompleteHandler)
{
  f_pAllTestsCompleteMessageHandler = pAllTestsCompleteHandler;
}

/*------------------------------------------------------------------------*/
void CU_set_suite_init_failure_handler(CU_SuiteInitFailureMessageHandler pSuiteInitFailureHandler)
{
  f_pSuiteInitFailureMessageHandler = pSuiteInitFailureHandler;
}

/*------------------------------------------------------------------------*/
void CU_set_suite_cleanup_failure_handler(CU_SuiteCleanupFailureMessageHandler pSuiteCleanupFailureHandler)
{
  f_pSuiteCleanupFailureMessageHandler = pSuiteCleanupFailureHandler;
}

/*------------------------------------------------------------------------*/
CU_SuiteStartMessageHandler CU_get_suite_start_handler(void)
{
  return f_pSuiteStartMessageHandler;
}

/*------------------------------------------------------------------------*/
CU_TestStartMessageHandler CU_get_test_start_handler(void)
{
  return f_pTestStartMessageHandler;
}

/*------------------------------------------------------------------------*/
CU_TestCompleteMessageHandler CU_get_test_complete_handler(void)
{
  return f_pTestCompleteMessageHandler;
}

/*------------------------------------------------------------------------*/
CU_SuiteCompleteMessageHandler CU_get_suite_complete_handler(void)
{
  return f_pSuiteCompleteMessageHandler;
}

/*------------------------------------------------------------------------*/
CU_AllTestsCompleteMessageHandler CU_get_all_test_complete_handler(void)
{
  return f_pAllTestsCompleteMessageHandler;
}

/*------------------------------------------------------------------------*/
CU_SuiteInitFailureMessageHandler CU_get_suite_init_failure_handler(void)
{
  return f_pSuiteInitFailureMessageHandler;
}

/*------------------------------------------------------------------------*/
CU_SuiteCleanupFailureMessageHandler CU_get_suite_cleanup_failure_handler(void)
{
  return f_pSuiteCleanupFailureMessageHandler;
}

/*------------------------------------------------------------------------*/
unsigned int CU_get_number_of_suites_run(void)
{
  return f_run_summary.nSuitesRun;
}

/*------------------------------------------------------------------------*/
unsigned int CU_get_number_of_suites_failed(void)
{
  return f_run_summary.nSuitesFailed;
}

/*------------------------------------------------------------------------*/
unsigned int CU_get_number_of_suites_inactive(void)
{
  return f_run_summary.nSuitesInactive;
}

/*------------------------------------------------------------------------*/
unsigned int CU_get_number_of_tests_run(void)
{
  return f_run_summary.nTestsRun;
}

/*------------------------------------------------------------------------*/
unsigned int CU_get_number_of_tests_failed(void)
{
  return f_run_summary.nTestsFailed;
}

/*------------------------------------------------------------------------*/
unsigned int CU_get_number_of_tests_inactive(void)
{
  return f_run_summary.nTestsInactive;
}

/*------------------------------------------------------------------------*/
unsigned int CU_get_number_of_asserts(void)
{
  return f_run_summary.nAsserts;
}

/*------------------------------------------------------------------------*/
unsigned int CU_get_number_of_successes(void)
{
  return (f_run_summary.nAsserts - f_run_summary.nAssertsFailed);
}

/*------------------------------------------------------------------------*/
unsigned int CU_get_number_of_failures(void)
{
  return f_run_summary.nAssertsFailed;
}

/*------------------------------------------------------------------------*/
unsigned int CU_get_number_of_failure_records(void)
{
  return f_run_summary.nFailureRecords;
}

/*------------------------------------------------------------------------*/
double CU_get_elapsed_time(void)
{
  if (CU_TRUE == f_bTestIsRunning) {
    return ((double)clock() - (double)f_start_time)/(double)CLOCKS_PER_SEC;
  }
  else {
    return f_run_summary.ElapsedTime;
  }
}

/*------------------------------------------------------------------------*/
CU_pFailureRecord CU_get_failure_list(void)
{
  return f_failure_list;
}

/*------------------------------------------------------------------------*/
CU_pRunSummary CU_get_run_summary(void)
{
  return &f_run_summary;
}

/*------------------------------------------------------------------------*/
CU_ErrorCode CU_run_all_tests(void)
{
  CU_pTestRegistry pRegistry = CU_get_registry();
  CU_pSuite pSuite = NULL;
  CU_ErrorCode result = CUE_SUCCESS;
  CU_ErrorCode result2;

  /* Clear results from the previous run */
  clear_previous_results(&f_run_summary, &f_failure_list);

  if (NULL == pRegistry) {
    result = CUE_NOREGISTRY;
  }
  else {
    /* test run is starting - set flag */
    f_bTestIsRunning = CU_TRUE;
    f_start_time = clock();

    pSuite = pRegistry->pSuite;
    while ((NULL != pSuite) && ((CUE_SUCCESS == result) || (CU_get_error_action() == CUEA_IGNORE))) {
      result2 = run_single_suite(pSuite, &f_run_summary);
      result = (CUE_SUCCESS == result) ? result2 : result;  /* result = 1st error encountered */
      pSuite = pSuite->pNext;
    }

    /* test run is complete - clear flag */
    f_bTestIsRunning = CU_FALSE;
    f_run_summary.ElapsedTime = ((double)clock() - (double)f_start_time)/(double)CLOCKS_PER_SEC;

    if (NULL != f_pAllTestsCompleteMessageHandler) {
     (*f_pAllTestsCompleteMessageHandler)(f_failure_list);
    }
  }

  CU_set_error(result);
  return result;
}

/*------------------------------------------------------------------------*/
CU_ErrorCode CU_run_suite(CU_pSuite pSuite)
{
  CU_ErrorCode result = CUE_SUCCESS;

  /* Clear results from the previous run */
  clear_previous_results(&f_run_summary, &f_failure_list);

  if (NULL == pSuite) {
    result = CUE_NOSUITE;
  }
  else {
    /* test run is starting - set flag */
    f_bTestIsRunning = CU_TRUE;
    f_start_time = clock();

    result = run_single_suite(pSuite, &f_run_summary);

    /* test run is complete - clear flag */
    f_bTestIsRunning = CU_FALSE;
    f_run_summary.ElapsedTime = ((double)clock() - (double)f_start_time)/(double)CLOCKS_PER_SEC;

    /* run handler for overall completion, if any */
    if (NULL != f_pAllTestsCompleteMessageHandler) {
      (*f_pAllTestsCompleteMessageHandler)(f_failure_list);
    }
  }

  CU_set_error(result);
  return result;
}

/*------------------------------------------------------------------------*/
CU_ErrorCode CU_run_test(CU_pSuite pSuite, CU_pTest pTest)
{
  CU_ErrorCode result = CUE_SUCCESS;
  CU_ErrorCode result2;

  /* Clear results from the previous run */
  clear_previous_results(&f_run_summary, &f_failure_list);

  if (NULL == pSuite) {
    result = CUE_NOSUITE;
  }
  else if (NULL == pTest) {
    result = CUE_NOTEST;
  }
  else if (CU_FALSE == pSuite->fActive) {
    f_run_summary.nSuitesInactive++;
    if (CU_FALSE != f_failure_on_inactive) {
      add_failure(&f_failure_list, &f_run_summary, CUF_SuiteInactive,
                  0, _("Suite inactive"), _("CUnit System"), pSuite, NULL);
    }
    result = CUE_SUITE_INACTIVE;
  }
  else if ((NULL == pTest->pName) || (NULL == CU_get_test_by_name(pTest->pName, pSuite))) {
    result = CUE_TEST_NOT_IN_SUITE;
  }
  else {
    /* test run is starting - set flag */
    f_bTestIsRunning = CU_TRUE;
    f_start_time = clock();

    f_pCurTest = NULL;
    f_pCurSuite = pSuite;

    pSuite->uiNumberOfTestsFailed = 0;
    pSuite->uiNumberOfTestsSuccess = 0;

    /* run handler for suite start, if any */
    if (NULL != f_pSuiteStartMessageHandler) {
      (*f_pSuiteStartMessageHandler)(pSuite);
    }

    /* run the suite initialization function, if any */
    if ((NULL != pSuite->pInitializeFunc) && (0 != (*pSuite->pInitializeFunc)())) {
      /* init function had an error - call handler, if any */
      if (NULL != f_pSuiteInitFailureMessageHandler) {
        (*f_pSuiteInitFailureMessageHandler)(pSuite);
      }
      f_run_summary.nSuitesFailed++;
      add_failure(&f_failure_list, &f_run_summary, CUF_SuiteInitFailed, 0,
                  _("Suite Initialization failed - Suite Skipped"),
                  _("CUnit System"), pSuite, NULL);
      result = CUE_SINIT_FAILED;
    }
    /* reach here if no suite initialization, or if it succeeded */
    else {
      result2 = run_single_test(pTest, &f_run_summary);
      result = (CUE_SUCCESS == result) ? result2 : result;

      /* run the suite cleanup function, if any */
      if ((NULL != pSuite->pCleanupFunc) && (0 != (*pSuite->pCleanupFunc)())) {
        /* cleanup function had an error - call handler, if any */
        if (NULL != f_pSuiteCleanupFailureMessageHandler) {
          (*f_pSuiteCleanupFailureMessageHandler)(pSuite);
        }
        f_run_summary.nSuitesFailed++;
        add_failure(&f_failure_list, &f_run_summary, CUF_SuiteCleanupFailed,
                    0, _("Suite cleanup failed."), _("CUnit System"), pSuite, NULL);
        result = (CUE_SUCCESS == result) ? CUE_SCLEAN_FAILED : result;
      }
    }

    /* run handler for suite completion, if any */
    if (NULL != f_pSuiteCompleteMessageHandler) {
      (*f_pSuiteCompleteMessageHandler)(pSuite, NULL);
    }

    /* test run is complete - clear flag */
    f_bTestIsRunning = CU_FALSE;
    f_run_summary.ElapsedTime = ((double)clock() - (double)f_start_time)/(double)CLOCKS_PER_SEC;

    /* run handler for overall completion, if any */
    if (NULL != f_pAllTestsCompleteMessageHandler) {
      (*f_pAllTestsCompleteMessageHandler)(f_failure_list);
    }

    f_pCurSuite = NULL;
  }

  CU_set_error(result);
  return result;
}

/*------------------------------------------------------------------------*/
void CU_clear_previous_results(void)
{
  clear_previous_results(&f_run_summary, &f_failure_list);
}

/*------------------------------------------------------------------------*/
CU_pSuite CU_get_current_suite(void)
{
  return f_pCurSuite;
}

/*------------------------------------------------------------------------*/
CU_pTest CU_get_current_test(void)
{
  return f_pCurTest;
}

/*------------------------------------------------------------------------*/
CU_BOOL CU_is_test_running(void)
{
  return f_bTestIsRunning;
}

/*------------------------------------------------------------------------*/
CU_EXPORT void CU_set_fail_on_inactive(CU_BOOL new_inactive)
{
  f_failure_on_inactive = new_inactive;
}

/*------------------------------------------------------------------------*/
CU_EXPORT CU_BOOL CU_get_fail_on_inactive(void)
{
  return f_failure_on_inactive;
}

/*------------------------------------------------------------------------*/
CU_EXPORT void CU_print_run_results(FILE *file)
{
  char *summary_string;

  assert(NULL != file);
  summary_string = CU_get_run_results_string();
  if (NULL != summary_string) {
    fprintf(file, "%s", summary_string);
    CU_FREE(summary_string);
  }
  else {
    fprintf(file, _("An error occurred printing the run results."));
  }
}

/*------------------------------------------------------------------------*/
CU_EXPORT char * CU_get_run_results_string(void)

{
  CU_pRunSummary pRunSummary = &f_run_summary;
  CU_pTestRegistry pRegistry = CU_get_registry();
  size_t width[9];
  size_t len;
  char *result;

  assert(NULL != pRunSummary);
  assert(NULL != pRegistry);

  width[0] = strlen(_("Run Summary:"));
  width[1] = CU_MAX(6,
                    CU_MAX(strlen(_("Type")),
                           CU_MAX(strlen(_("suites")),
                                  CU_MAX(strlen(_("tests")),
                                         strlen(_("asserts")))))) + 1;
  width[2] = CU_MAX(6,
                    CU_MAX(strlen(_("Total")),
                           CU_MAX(CU_number_width(pRegistry->uiNumberOfSuites),
                                  CU_MAX(CU_number_width(pRegistry->uiNumberOfTests),
                                         CU_number_width(pRunSummary->nAsserts))))) + 1;
  width[3] = CU_MAX(6,
                    CU_MAX(strlen(_("Ran")),
                           CU_MAX(CU_number_width(pRunSummary->nSuitesRun),
                                  CU_MAX(CU_number_width(pRunSummary->nTestsRun),
                                         CU_number_width(pRunSummary->nAsserts))))) + 1;
  width[4] = CU_MAX(6,
                    CU_MAX(strlen(_("Passed")),
                           CU_MAX(strlen(_("n/a")),
                                  CU_MAX(CU_number_width(pRunSummary->nTestsRun - pRunSummary->nTestsFailed),
                                         CU_number_width(pRunSummary->nAsserts - pRunSummary->nAssertsFailed))))) + 1;
  width[5] = CU_MAX(6,
                    CU_MAX(strlen(_("Failed")),
                           CU_MAX(CU_number_width(pRunSummary->nSuitesFailed),
                                  CU_MAX(CU_number_width(pRunSummary->nTestsFailed),
                                         CU_number_width(pRunSummary->nAssertsFailed))))) + 1;
  width[6] = CU_MAX(6,
                    CU_MAX(strlen(_("Inactive")),
                           CU_MAX(CU_number_width(pRunSummary->nSuitesInactive),
                                  CU_MAX(CU_number_width(pRunSummary->nTestsInactive),
                                         strlen(_("n/a")))))) + 1;

  width[7] = strlen(_("Elapsed time = "));
  width[8] = strlen(_(" seconds"));

  len = 13 + 4*(width[0] + width[1] + width[2] + width[3] + width[4] + width[5] + width[6]) + width[7] + width[8] + 1;
  result = (char *)CU_MALLOC(len);

  if (NULL != result) {
    snprintf(result, len, "%*s%*s%*s%*s%*s%*s%*s\n"   /* if you change this, be sure  */
                          "%*s%*s%*u%*u%*s%*u%*u\n"   /* to change the calculation of */
                          "%*s%*s%*u%*u%*u%*u%*u\n"   /* len above!                   */
                          "%*s%*s%*u%*u%*u%*u%*s\n\n"
                          "%*s%8.3f%*s",
            width[0], _("Run Summary:"),
            width[1], _("Type"),
            width[2], _("Total"),
            width[3], _("Ran"),
            width[4], _("Passed"),
            width[5], _("Failed"),
            width[6], _("Inactive"),
            width[0], " ",
            width[1], _("suites"),
            width[2], pRegistry->uiNumberOfSuites,
            width[3], pRunSummary->nSuitesRun,
            width[4], _("n/a"),
            width[5], pRunSummary->nSuitesFailed,
            width[6], pRunSummary->nSuitesInactive,
            width[0], " ",
            width[1], _("tests"),
            width[2], pRegistry->uiNumberOfTests,
            width[3], pRunSummary->nTestsRun,
            width[4], pRunSummary->nTestsRun - pRunSummary->nTestsFailed,
            width[5], pRunSummary->nTestsFailed,
            width[6], pRunSummary->nTestsInactive,
            width[0], " ",
            width[1], _("asserts"),
            width[2], pRunSummary->nAsserts,
            width[3], pRunSummary->nAsserts,
            width[4], pRunSummary->nAsserts - pRunSummary->nAssertsFailed,
            width[5], pRunSummary->nAssertsFailed,
            width[6], _("n/a"),
            width[7], _("Elapsed time = "), CU_get_elapsed_time(),  /* makes sure time is updated */
            width[8], _(" seconds")
            );
     result[len-1] = '\0';
  }
  return result;
}

/*=================================================================
 *  Static Function Definitions
 *=================================================================*/
/**
 *  Records a runtime failure.
 *  This function is called whenever a runtime failure occurs.
 *  This includes user assertion failures, suite initialization and
 *  cleanup failures, and inactive suites/tests when set as failures.
 *  This function records the details of the failure in a new
 *  failure record in the linked list of runtime failures.
 *
 *  @param ppFailure    Pointer to head of linked list of failure
 *                      records to append with new failure record.
 *                      If it points to a NULL pointer, it will be set
 *                      to point to the new failure record.
 *  @param pRunSummary  Pointer to CU_RunSummary keeping track of failure records
 *                      (ignored if NULL).
 *  @param type         Type of failure.
 *  @param uiLineNumber Line number of the failure, if applicable.
 *  @param szCondition  Description of failure condition
 *  @param szFileName   Name of file, if applicable
 *  @param pSuite       The suite being run at time of failure
 *  @param pTest        The test being run at time of failure
 */
static void add_failure(CU_pFailureRecord* ppFailure,
                        CU_pRunSummary pRunSummary,
                        CU_FailureType type,
                        unsigned int uiLineNumber,
                        const char *szCondition,
                        const char *szFileName,
                        CU_pSuite pSuite,
                        CU_pTest pTest)
{
  CU_pFailureRecord pFailureNew = NULL;
  CU_pFailureRecord pTemp = NULL;

  assert(NULL != ppFailure);

  pFailureNew = (CU_pFailureRecord)CU_MALLOC(sizeof(CU_FailureRecord));

  if (NULL == pFailureNew) {
    return;
  }

  pFailureNew->strFileName = NULL;
  pFailureNew->strCondition = NULL;
  if (NULL != szFileName) {
    pFailureNew->strFileName = (char*)CU_MALLOC(strlen(szFileName) + 1);
    if(NULL == pFailureNew->strFileName) {
      CU_FREE(pFailureNew);
      return;
    }
    strcpy(pFailureNew->strFileName, szFileName);
  }

  if (NULL != szCondition) {
    pFailureNew->strCondition = (char*)CU_MALLOC(strlen(szCondition) + 1);
    if (NULL == pFailureNew->strCondition) {
      if(NULL != pFailureNew->strFileName) {
        CU_FREE(pFailureNew->strFileName);
      }
      CU_FREE(pFailureNew);
      return;
    }
    strcpy(pFailureNew->strCondition, szCondition);
  }

  pFailureNew->type = type;
  pFailureNew->uiLineNumber = uiLineNumber;
  pFailureNew->pTest = pTest;
  pFailureNew->pSuite = pSuite;
  pFailureNew->pNext = NULL;
  pFailureNew->pPrev = NULL;

  pTemp = *ppFailure;
  if (NULL != pTemp) {
    while (NULL != pTemp->pNext) {
      pTemp = pTemp->pNext;
    }
    pTemp->pNext = pFailureNew;
    pFailureNew->pPrev = pTemp;
  }
  else {
    *ppFailure = pFailureNew;
  }

  if (NULL != pRunSummary) {
    ++(pRunSummary->nFailureRecords);
  }
  f_last_failure = pFailureNew;
}

/*
 *  Local function for result set initialization/cleanup.
 */
/*------------------------------------------------------------------------*/
/**
 *  Initializes the run summary information in the specified structure.
 *  Resets the run counts to zero, and calls cleanup_failure_list() if
 *  failures were recorded by the last test run.  Calling this function
 *  multiple times, while inefficient, will not cause an error condition.
 *
 *  @param pRunSummary CU_RunSummary to initialize (non-NULL).
 *  @param ppFailure   The failure record to clean (non-NULL).
 *  @see CU_clear_previous_results()
 */
static void clear_previous_results(CU_pRunSummary pRunSummary, CU_pFailureRecord* ppFailure)
{
  assert(NULL != pRunSummary);
  assert(NULL != ppFailure);

  pRunSummary->nSuitesRun = 0;
  pRunSummary->nSuitesFailed = 0;
  pRunSummary->nSuitesInactive = 0;
  pRunSummary->nTestsRun = 0;
  pRunSummary->nTestsFailed = 0;
  pRunSummary->nTestsInactive = 0;
  pRunSummary->nAsserts = 0;
  pRunSummary->nAssertsFailed = 0;
  pRunSummary->nFailureRecords = 0;
  pRunSummary->ElapsedTime = 0.0;

  if (NULL != *ppFailure) {
    cleanup_failure_list(ppFailure);
  }

  f_last_failure = NULL;
}

/*------------------------------------------------------------------------*/
/**
 *  Frees all memory allocated for the linked list of test failure
 *  records.  pFailure is reset to NULL after its list is cleaned up.
 *
 *  @param ppFailure Pointer to head of linked list of
 *                   CU_pFailureRecords to clean.
 *  @see CU_clear_previous_results()
 */
static void cleanup_failure_list(CU_pFailureRecord* ppFailure)
{
  CU_pFailureRecord pCurFailure = NULL;
  CU_pFailureRecord pNextFailure = NULL;

  pCurFailure = *ppFailure;

  while (NULL != pCurFailure) {

    if (NULL != pCurFailure->strCondition) {
      CU_FREE(pCurFailure->strCondition);
    }

    if (NULL != pCurFailure->strFileName) {
      CU_FREE(pCurFailure->strFileName);
    }

    pNextFailure = pCurFailure->pNext;
    CU_FREE(pCurFailure);
    pCurFailure = pNextFailure;
  }

  *ppFailure = NULL;
}

/*------------------------------------------------------------------------*/
/**
 *  Runs all tests in a specified suite.
 *  Internal function to run all tests in a suite.  The suite need
 *  not be registered in the test registry to be run.  Only
 *  suites having their fActive flags set CU_TRUE will actually be
 *  run.  If the CUnit framework is in an error condition after
 *  running a test, no additional tests are run.
 *
 *  @param pSuite The suite containing the test (non-NULL).
 *  @param pRunSummary The CU_RunSummary to receive the results (non-NULL).
 *  @return A CU_ErrorCode indicating the status of the run.
 *  @see CU_run_suite() for public interface function.
 *  @see CU_run_all_tests() for running all suites.
 */
static CU_ErrorCode run_single_suite(CU_pSuite pSuite, CU_pRunSummary pRunSummary)
{
  CU_pTest pTest = NULL;
  unsigned int nStartFailures;
  /* keep track of the last failure BEFORE running the test */
  CU_pFailureRecord pLastFailure = f_last_failure;
  CU_ErrorCode result = CUE_SUCCESS;
  CU_ErrorCode result2;

  assert(NULL != pSuite);
  assert(NULL != pRunSummary);

  nStartFailures = pRunSummary->nFailureRecords;

  f_pCurTest = NULL;
  f_pCurSuite = pSuite;

  /* run handler for suite start, if any */
  if (NULL != f_pSuiteStartMessageHandler) {
    (*f_pSuiteStartMessageHandler)(pSuite);
  }

  /* run suite if it's active */
  if (CU_FALSE != pSuite->fActive) {

    /* run the suite initialization function, if any */
    if ((NULL != pSuite->pInitializeFunc) && (0 != (*pSuite->pInitializeFunc)())) {
      /* init function had an error - call handler, if any */
      if (NULL != f_pSuiteInitFailureMessageHandler) {
        (*f_pSuiteInitFailureMessageHandler)(pSuite);
      }
      pRunSummary->nSuitesFailed++;
      add_failure(&f_failure_list, &f_run_summary, CUF_SuiteInitFailed, 0,
                  _("Suite Initialization failed - Suite Skipped"),
                  _("CUnit System"), pSuite, NULL);
      result = CUE_SINIT_FAILED;
    }

    /* reach here if no suite initialization, or if it succeeded */
    else {
      pTest = pSuite->pTest;
      while ((NULL != pTest) && ((CUE_SUCCESS == result) || (CU_get_error_action() == CUEA_IGNORE))) {
        if (CU_FALSE != pTest->fActive) {
          result2 = run_single_test(pTest, pRunSummary);
          result = (CUE_SUCCESS == result) ? result2 : result;
        }
        else {
          f_run_summary.nTestsInactive++;
          if (CU_FALSE != f_failure_on_inactive) {
            add_failure(&f_failure_list, &f_run_summary, CUF_TestInactive,
                        0, _("Test inactive"), _("CUnit System"), pSuite, pTest);
            result = CUE_TEST_INACTIVE;
          }
        }
        pTest = pTest->pNext;

        if (CUE_SUCCESS == result) {
          pSuite->uiNumberOfTestsFailed++;
        }
        else {
          pSuite->uiNumberOfTestsSuccess++;
        }
      }
      pRunSummary->nSuitesRun++;

      /* call the suite cleanup function, if any */
      if ((NULL != pSuite->pCleanupFunc) && (0 != (*pSuite->pCleanupFunc)())) {
        if (NULL != f_pSuiteCleanupFailureMessageHandler) {
          (*f_pSuiteCleanupFailureMessageHandler)(pSuite);
        }
        pRunSummary->nSuitesFailed++;
        add_failure(&f_failure_list, &f_run_summary, CUF_SuiteCleanupFailed,
                    0, _("Suite cleanup failed."), _("CUnit System"), pSuite, NULL);
        result = (CUE_SUCCESS == result) ? CUE_SCLEAN_FAILED : result;
      }
    }
  }

  /* otherwise record inactive suite and failure if appropriate */
  else {
    f_run_summary.nSuitesInactive++;
    if (CU_FALSE != f_failure_on_inactive) {
      add_failure(&f_failure_list, &f_run_summary, CUF_SuiteInactive,
                  0, _("Suite inactive"), _("CUnit System"), pSuite, NULL);
      result = CUE_SUITE_INACTIVE;
    }
  }

  /* if additional failures have occurred... */
  if (pRunSummary->nFailureRecords > nStartFailures) {
    if (NULL != pLastFailure) {
      pLastFailure = pLastFailure->pNext;  /* was a previous failure, so go to next one */
    }
    else {
      pLastFailure = f_failure_list;       /* no previous failure - go to 1st one */
    }
  }
  else {
    pLastFailure = NULL;                   /* no additional failure - set to NULL */
  }

  /* run handler for suite completion, if any */
  if (NULL != f_pSuiteCompleteMessageHandler) {
    (*f_pSuiteCompleteMessageHandler)(pSuite, pLastFailure);
  }

  f_pCurSuite = NULL;
  return result;
}

/*------------------------------------------------------------------------*/
/**
 *  Runs a specific test.
 *  Internal function to run a test case.  This includes calling
 *  any handler to be run before executing the test, running the
 *  test's function (if any), and calling any handler to be run
 *  after executing a test.  Suite initialization and cleanup functions
 *  are not called by this function.  A current suite must be set and
 *  active (checked by assertion).
 *
 *  @param pTest The test to be run (non-NULL).
 *  @param pRunSummary The CU_RunSummary to receive the results (non-NULL).
 *  @return A CU_ErrorCode indicating the status of the run.
 *  @see CU_run_test() for public interface function.
 *  @see CU_run_all_tests() for running all suites.
 */
static CU_ErrorCode run_single_test(CU_pTest pTest, CU_pRunSummary pRunSummary)
{
  volatile unsigned int nStartFailures;
  /* keep track of the last failure BEFORE running the test */
  volatile CU_pFailureRecord pLastFailure = f_last_failure;
  jmp_buf buf;
  CU_ErrorCode result = CUE_SUCCESS;

  assert(NULL != f_pCurSuite);
  assert(CU_FALSE != f_pCurSuite->fActive);
  assert(NULL != pTest);
  assert(NULL != pRunSummary);

  nStartFailures = pRunSummary->nFailureRecords;

  f_pCurTest = pTest;

  if (NULL != f_pTestStartMessageHandler) {
    (*f_pTestStartMessageHandler)(f_pCurTest, f_pCurSuite);
  }

  /* run test if it is active */
  if (CU_FALSE != pTest->fActive) {

    if (NULL != f_pCurSuite->pSetUpFunc) {
      (*f_pCurSuite->pSetUpFunc)();
    }

    /* set jmp_buf and run test */
    pTest->pJumpBuf = &buf;
    if (0 == setjmp(buf)) {
      if (NULL != pTest->pTestFunc) {
        (*pTest->pTestFunc)();
      }
    }

    if (NULL != f_pCurSuite->pTearDownFunc) {
       (*f_pCurSuite->pTearDownFunc)();
    }

    pRunSummary->nTestsRun++;
  }
  else {
    f_run_summary.nTestsInactive++;
    if (CU_FALSE != f_failure_on_inactive) {
      add_failure(&f_failure_list, &f_run_summary, CUF_TestInactive,
                  0, _("Test inactive"), _("CUnit System"), f_pCurSuite, f_pCurTest);
    }
    result = CUE_TEST_INACTIVE;
  }

  /* if additional failures have occurred... */
  if (pRunSummary->nFailureRecords > nStartFailures) {
    pRunSummary->nTestsFailed++;
    if (NULL != pLastFailure) {
      pLastFailure = pLastFailure->pNext;  /* was a previous failure, so go to next one */
    }
    else {
      pLastFailure = f_failure_list;       /* no previous failure - go to 1st one */
    }
  }
  else {
    pLastFailure = NULL;                   /* no additional failure - set to NULL */
  }

  if (NULL != f_pTestCompleteMessageHandler) {
    (*f_pTestCompleteMessageHandler)(f_pCurTest, f_pCurSuite, pLastFailure);
  }

  pTest->pJumpBuf = NULL;
  f_pCurTest = NULL;

  return result;
}

/** @} */

#ifdef CUNIT_BUILD_TESTS
#include "test_cunit.h"

/** Types of framework events tracked by test system. */
typedef enum TET {
  SUITE_START = 1,
  TEST_START,
  TEST_COMPLETE,
  SUITE_COMPLETE,
  ALL_TESTS_COMPLETE,
  SUITE_INIT_FAILED,
  SUITE_CLEANUP_FAILED
} TestEventType;

/** Test event structure for recording details of a framework event. */
typedef struct TE {
  TestEventType     type;
  CU_pSuite         pSuite;
  CU_pTest          pTest;
  CU_pFailureRecord pFailure;
  struct TE *       pNext;
} TestEvent, * pTestEvent;

static int f_nTestEvents = 0;
static pTestEvent f_pFirstEvent = NULL;

/** Creates & stores a test event record having the specified details. */
static void add_test_event(TestEventType type, CU_pSuite psuite,
                           CU_pTest ptest, CU_pFailureRecord pfailure)
{
  pTestEvent pNewEvent = (pTestEvent)malloc(sizeof(TestEvent));
  pTestEvent pNextEvent = f_pFirstEvent;

  if (NULL == pNewEvent) {
    fprintf(stderr, "Memory allocation failed in add_test_event().");
    exit(1);
  }

  pNewEvent->type = type;
  pNewEvent->pSuite = psuite;
  pNewEvent->pTest = ptest;
  pNewEvent->pFailure = pfailure;
  pNewEvent->pNext = NULL;

  if (pNextEvent) {
    while (pNextEvent->pNext) {
      pNextEvent = pNextEvent->pNext;
    }
    pNextEvent->pNext = pNewEvent;
  }
  else {
    f_pFirstEvent = pNewEvent;
  }
  ++f_nTestEvents;
}

/** Deallocates all test event data. */
static void clear_test_events(void)
{
  pTestEvent pCurrentEvent = f_pFirstEvent;
  pTestEvent pNextEvent = NULL;

  while (pCurrentEvent) {
    pNextEvent = pCurrentEvent->pNext;
    free(pCurrentEvent);
    pCurrentEvent = pNextEvent;
  }

  f_pFirstEvent = NULL;
  f_nTestEvents = 0;
}

static void suite_start_handler(const CU_pSuite pSuite)
{
  TEST(CU_is_test_running());
  TEST(pSuite == CU_get_current_suite());
  TEST(NULL == CU_get_current_test());

  add_test_event(SUITE_START, pSuite, NULL, NULL);
}

static void test_start_handler(const CU_pTest pTest, const CU_pSuite pSuite)
{
  TEST(CU_is_test_running());
  TEST(pSuite == CU_get_current_suite());
  TEST(pTest == CU_get_current_test());

  add_test_event(TEST_START, pSuite, pTest, NULL);
}

static void test_complete_handler(const CU_pTest pTest, const CU_pSuite pSuite,
                                  const CU_pFailureRecord pFailure)
{
  TEST(CU_is_test_running());
  TEST(pSuite == CU_get_current_suite());
  TEST(pTest == CU_get_current_test());

  add_test_event(TEST_COMPLETE, pSuite, pTest, pFailure);
}

static void suite_complete_handler(const CU_pSuite pSuite,
                                   const CU_pFailureRecord pFailure)
{
  TEST(CU_is_test_running());
  TEST(pSuite == CU_get_current_suite());
  TEST(NULL == CU_get_current_test());

  add_test_event(SUITE_COMPLETE, pSuite, NULL, pFailure);
}

static void test_all_complete_handler(const CU_pFailureRecord pFailure)
{
  TEST(!CU_is_test_running());

  add_test_event(ALL_TESTS_COMPLETE, NULL, NULL, pFailure);
}

static void suite_init_failure_handler(const CU_pSuite pSuite)
{
  TEST(CU_is_test_running());
  TEST(pSuite == CU_get_current_suite());

  add_test_event(SUITE_INIT_FAILED, pSuite, NULL, NULL);
}

static void suite_cleanup_failure_handler(const CU_pSuite pSuite)
{
  TEST(CU_is_test_running());
  TEST(pSuite == CU_get_current_suite());

  add_test_event(SUITE_CLEANUP_FAILED, pSuite, NULL, NULL);
}

/**
 *  Centralize test result testing - we're going to do it a lot!
 *  This is messy since we want to report the calling location upon failure.
 *
 *  Via calling test functions tests:
 *      CU_get_number_of_suites_run()
 *      CU_get_number_of_suites_failed()
 *      CU_get_number_of_tests_run()
 *      CU_get_number_of_tests_failed()
 *      CU_get_number_of_asserts()
 *      CU_get_number_of_successes()
 *      CU_get_number_of_failures()
 *      CU_get_number_of_failure_records()
 *      CU_get_run_summary()
 */
static void do_test_results(unsigned int nSuitesRun,
                            unsigned int nSuitesFailed,
                            unsigned int nSuitesInactive,
                            unsigned int nTestsRun,
                            unsigned int nTestsFailed,
                            unsigned int nTestsInactive,
                            unsigned int nAsserts,
                            unsigned int nSuccesses,
                            unsigned int nFailures,
                            unsigned int nFailureRecords,
                            const char *file,
                            unsigned int line)
{
  char msg[500];
  CU_pRunSummary pRunSummary = NULL;

  if (nSuitesRun == CU_get_number_of_suites_run()) {
    PASS();
  } else {
    snprintf(msg, 499, "%u == CU_get_number_of_suites_run() (called from %s:%u)",
                       nSuitesRun, file, line);
    msg[499] = '\0';
    FAIL(msg);
  }

  if (nSuitesInactive == CU_get_number_of_suites_inactive()) {
    PASS();
  } else {
    snprintf(msg, 499, "%u == CU_get_number_of_suites_inactive() (called from %s:%u)",
                       nSuitesInactive, file, line);
    msg[499] = '\0';
    FAIL(msg);
  }

  if (nSuitesFailed == CU_get_number_of_suites_failed()) {
    PASS();
  } else {
    snprintf(msg, 499, "%u == CU_get_number_of_suites_failed() (called from %s:%u)",
                       nSuitesFailed, file, line);
    msg[499] = '\0';
    FAIL(msg);
  }

  if (nTestsRun == CU_get_number_of_tests_run()) {
    PASS();
  } else {
    snprintf(msg, 499, "%u == CU_get_number_of_tests_run() (called from %s:%u)",
                       nTestsRun, file, line);
    msg[499] = '\0';
    FAIL(msg);
  }

  if (nTestsFailed == CU_get_number_of_tests_failed()) {
    PASS();
  } else {
    snprintf(msg, 499, "%u == CU_get_number_of_tests_failed() (called from %s:%u)",
                       nTestsFailed, file, line);
    msg[499] = '\0';
    FAIL(msg);
  }

  if (nTestsInactive == CU_get_number_of_tests_inactive()) {
    PASS();
  } else {
    snprintf(msg, 499, "%u == CU_get_number_of_tests_inactive() (called from %s:%u)",
                       nTestsInactive, file, line);
    msg[499] = '\0';
    FAIL(msg);
  }

  if (nAsserts == CU_get_number_of_asserts()) {
    PASS();
  } else {
    snprintf(msg, 499, "%u == CU_get_number_of_asserts() (called from %s:%u)",
                       nAsserts, file, line);
    msg[499] = '\0';
    FAIL(msg);
  }

  if (nSuccesses == CU_get_number_of_successes()) {
    PASS();
  } else {
    snprintf(msg, 499, "%u == CU_get_number_of_successes() (called from %s:%u)",
                       nSuccesses, file, line);
    msg[499] = '\0';
    FAIL(msg);
  }

  if (nFailures == CU_get_number_of_failures()) {
    PASS();
  } else {
    snprintf(msg, 499, "%u == CU_get_number_of_failures() (called from %s:%u)",
                       nFailures, file, line);
    msg[499] = '\0';
    FAIL(msg);
  }

  if (nFailureRecords == CU_get_number_of_failure_records()) {
    PASS();
  } else {
    snprintf(msg, 499, "%u == CU_get_number_of_failure_records() (called from %s:%u)",
                       nFailureRecords, file, line);
    msg[499] = '\0';
    FAIL(msg);
  }

  pRunSummary = CU_get_run_summary();

  if (pRunSummary->nSuitesRun == CU_get_number_of_suites_run()) {
    PASS();
  } else {
    snprintf(msg, 499, "%u == CU_get_number_of_suites_run() (called from %s:%u)",
                       pRunSummary->nSuitesRun, file, line);
    msg[499] = '\0';
    FAIL(msg);
  }

  if (pRunSummary->nSuitesFailed == CU_get_number_of_suites_failed()) {
    PASS();
  } else {
    snprintf(msg, 499, "%u == CU_get_number_of_suites_failed() (called from %s:%u)",
                       pRunSummary->nSuitesFailed, file, line);
    msg[499] = '\0';
    FAIL(msg);
  }

  if (pRunSummary->nTestsRun == CU_get_number_of_tests_run()) {
    PASS();
  } else {
    snprintf(msg, 499, "%u == CU_get_number_of_tests_run() (called from %s:%u)",
                       pRunSummary->nTestsRun, file, line);
    msg[499] = '\0';
    FAIL(msg);
  }

  if (pRunSummary->nTestsFailed == CU_get_number_of_tests_failed()) {
    PASS();
  } else {
    snprintf(msg, 499, "%u == CU_get_number_of_tests_failed() (called from %s:%u)",
                       pRunSummary->nTestsFailed, file, line);
    msg[499] = '\0';
    FAIL(msg);
  }

  if (pRunSummary->nAsserts == CU_get_number_of_asserts()) {
    PASS();
  } else {
    snprintf(msg, 499, "%u == CU_get_number_of_asserts() (called from %s:%u)",
                       pRunSummary->nAsserts, file, line);
    msg[499] = '\0';
    FAIL(msg);
  }

  if (pRunSummary->nAssertsFailed == CU_get_number_of_failures()) {
    PASS();
  } else {
    snprintf(msg, 499, "%u == CU_get_number_of_failures() (called from %s:%u)",
                       pRunSummary->nAssertsFailed, file, line);
    msg[499] = '\0';
    FAIL(msg);
  }

  if (pRunSummary->nFailureRecords == CU_get_number_of_failure_records()) {
    PASS();
  } else {
    snprintf(msg, 499, "%u == CU_get_number_of_failure_records() (called from %s:%u)",
                       pRunSummary->nFailureRecords, file, line);
    msg[499] = '\0';
    FAIL(msg);
  }
}

#define test_results(nSuitesRun, nSuitesFailed, nSuitesInactive, nTestsRun, nTestsFailed,    \
                     nTestsInactive, nAsserts, nSuccesses, nFailures, nFailureRecords)      \
        do_test_results(nSuitesRun, nSuitesFailed, nSuitesInactive, nTestsRun, nTestsFailed, \
                        nTestsInactive, nAsserts, nSuccesses, nFailures, nFailureRecords,   \
                        __FILE__, __LINE__)

static void test_succeed(void) { CU_TEST(CU_TRUE); }
static void test_fail(void) { CU_TEST(CU_FALSE); }
static int suite_succeed(void) { return 0; }
static int suite_fail(void) { return 1; }

static CU_BOOL SetUp_Passed;

static void test_succeed_if_setup(void) { CU_TEST(SetUp_Passed); }
static void test_fail_if_not_setup(void) { CU_TEST(SetUp_Passed); }

static void suite_setup(void) { SetUp_Passed = CU_TRUE; }
static void suite_teardown(void) { SetUp_Passed = CU_FALSE; }


/*-------------------------------------------------*/
/* tests:
 *      CU_set_suite_start_handler()
 *      CU_set_test_start_handler()
 *      CU_set_test_complete_handler()
 *      CU_set_suite_complete_handler()
 *      CU_set_all_test_complete_handler()
 *      CU_set_suite_init_failure_handler()
 *      CU_set_suite_cleanup_failure_handler()
 *      CU_get_suite_start_handler()
 *      CU_get_test_start_handler()
 *      CU_get_test_complete_handler()
 *      CU_get_suite_complete_handler()
 *      CU_get_all_test_complete_handler()
 *      CU_get_suite_init_failure_handler()
 *      CU_get_suite_cleanup_failure_handler()
 *      CU_is_test_running()
 *  via handlers tests:
 *      CU_get_current_suite()
 *      CU_get_current_test()
 */
static void test_message_handlers(void)
{
  CU_pSuite pSuite1 = NULL;
  CU_pSuite pSuite2 = NULL;
  CU_pSuite pSuite3 = NULL;
  CU_pTest  pTest1 = NULL;
  CU_pTest  pTest2 = NULL;
  CU_pTest  pTest3 = NULL;
  CU_pTest  pTest4 = NULL;
  CU_pTest  pTest5 = NULL;
  pTestEvent pEvent = NULL;

  TEST(!CU_is_test_running());

  /* handlers should be NULL on startup */
  TEST(NULL == CU_get_suite_start_handler());
  TEST(NULL == CU_get_test_start_handler());
  TEST(NULL == CU_get_test_complete_handler());
  TEST(NULL == CU_get_suite_complete_handler());
  TEST(NULL == CU_get_all_test_complete_handler());
  TEST(NULL == CU_get_suite_init_failure_handler());
  TEST(NULL == CU_get_suite_cleanup_failure_handler());

  /* register some suites and tests */
  CU_initialize_registry();
  pSuite1 = CU_add_suite("suite1", NULL, NULL);
  pTest1 = CU_add_test(pSuite1, "test1", test_succeed);
  pTest2 = CU_add_test(pSuite1, "test2", test_fail);
  pTest3 = CU_add_test(pSuite1, "test3", test_succeed);
  pSuite2 = CU_add_suite("suite2", suite_fail, NULL);
  pTest4 = CU_add_test(pSuite2, "test4", test_succeed);
  pSuite3 = CU_add_suite("suite3", suite_succeed, suite_fail);
  pTest5 = CU_add_test(pSuite3, "test5", test_fail);

  TEST_FATAL(CUE_SUCCESS == CU_get_error());

  /* first run tests without handlers set */
  clear_test_events();
  CU_run_all_tests();

  TEST(0 == f_nTestEvents);
  TEST(NULL == f_pFirstEvent);
  test_results(2,2,0,4,2,0,4,2,2,4);

  /* set handlers to local functions */
  CU_set_suite_start_handler(&suite_start_handler);
  CU_set_test_start_handler(&test_start_handler);
  CU_set_test_complete_handler(&test_complete_handler);
  CU_set_suite_complete_handler(&suite_complete_handler);
  CU_set_all_test_complete_handler(&test_all_complete_handler);
  CU_set_suite_init_failure_handler(&suite_init_failure_handler);
  CU_set_suite_cleanup_failure_handler(&suite_cleanup_failure_handler);

  /* confirm handlers set properly */
  TEST(suite_start_handler == CU_get_suite_start_handler());
  TEST(test_start_handler == CU_get_test_start_handler());
  TEST(test_complete_handler == CU_get_test_complete_handler());
  TEST(suite_complete_handler == CU_get_suite_complete_handler());
  TEST(test_all_complete_handler == CU_get_all_test_complete_handler());
  TEST(suite_init_failure_handler == CU_get_suite_init_failure_handler());
  TEST(suite_cleanup_failure_handler == CU_get_suite_cleanup_failure_handler());

  /* run tests again with handlers set */
  clear_test_events();
  CU_run_all_tests();

  TEST(17 == f_nTestEvents);
  if (17 == f_nTestEvents) {
    pEvent = f_pFirstEvent;
    TEST(SUITE_START == pEvent->type);
    TEST(pSuite1 == pEvent->pSuite);
    TEST(NULL == pEvent->pTest);
    TEST(NULL == pEvent->pFailure);

    pEvent = pEvent->pNext;
    TEST(TEST_START == pEvent->type);
    TEST(pSuite1 == pEvent->pSuite);
    TEST(pTest1 == pEvent->pTest);
    TEST(NULL == pEvent->pFailure);

    pEvent = pEvent->pNext;
    TEST(TEST_COMPLETE == pEvent->type);
    TEST(pSuite1 == pEvent->pSuite);
    TEST(pTest1 == pEvent->pTest);
    TEST(NULL == pEvent->pFailure);

    pEvent = pEvent->pNext;
    TEST(TEST_START == pEvent->type);
    TEST(pSuite1 == pEvent->pSuite);
    TEST(pTest2 == pEvent->pTest);
    TEST(NULL == pEvent->pFailure);

    pEvent = pEvent->pNext;
    TEST(TEST_COMPLETE == pEvent->type);
    TEST(pSuite1 == pEvent->pSuite);
    TEST(pTest2 == pEvent->pTest);
    TEST(NULL != pEvent->pFailure);

    pEvent = pEvent->pNext;
    TEST(TEST_START == pEvent->type);
    TEST(pSuite1 == pEvent->pSuite);
    TEST(pTest3 == pEvent->pTest);
    TEST(NULL == pEvent->pFailure);

    pEvent = pEvent->pNext;
    TEST(TEST_COMPLETE == pEvent->type);
    TEST(pSuite1 == pEvent->pSuite);
    TEST(pTest3 == pEvent->pTest);
    TEST(NULL == pEvent->pFailure);

    pEvent = pEvent->pNext;
    TEST(SUITE_COMPLETE == pEvent->type);
    TEST(pSuite1 == pEvent->pSuite);
    TEST(NULL == pEvent->pTest);
    TEST(NULL != pEvent->pFailure);

    pEvent = pEvent->pNext;
    TEST(SUITE_START == pEvent->type);
    TEST(pSuite2 == pEvent->pSuite);
    TEST(NULL == pEvent->pTest);
    TEST(NULL == pEvent->pFailure);

    pEvent = pEvent->pNext;
    TEST(SUITE_INIT_FAILED == pEvent->type);
    TEST(pSuite2 == pEvent->pSuite);
    TEST(NULL == pEvent->pTest);
    TEST(NULL == pEvent->pFailure);

    pEvent = pEvent->pNext;
    TEST(SUITE_COMPLETE == pEvent->type);
    TEST(pSuite2 == pEvent->pSuite);
    TEST(NULL == pEvent->pTest);
    TEST(NULL != pEvent->pFailure);

    pEvent = pEvent->pNext;
    TEST(SUITE_START == pEvent->type);
    TEST(pSuite3 == pEvent->pSuite);
    TEST(NULL == pEvent->pTest);
    TEST(NULL == pEvent->pFailure);

    pEvent = pEvent->pNext;
    TEST(TEST_START == pEvent->type);
    TEST(pSuite3 == pEvent->pSuite);
    TEST(pTest5 == pEvent->pTest);
    TEST(NULL == pEvent->pFailure);

    pEvent = pEvent->pNext;
    TEST(TEST_COMPLETE == pEvent->type);
    TEST(pSuite3 == pEvent->pSuite);
    TEST(pTest5 == pEvent->pTest);
    TEST(NULL != pEvent->pFailure);

    pEvent = pEvent->pNext;
    TEST(SUITE_CLEANUP_FAILED == pEvent->type);
    TEST(pSuite3 == pEvent->pSuite);
    TEST(NULL == pEvent->pTest);
    TEST(NULL == pEvent->pFailure);

    pEvent = pEvent->pNext;
    TEST(SUITE_COMPLETE == pEvent->type);
    TEST(pSuite3 == pEvent->pSuite);
    TEST(NULL == pEvent->pTest);
    TEST(NULL != pEvent->pFailure);

    pEvent = pEvent->pNext;
    TEST(ALL_TESTS_COMPLETE == pEvent->type);
    TEST(NULL == pEvent->pSuite);
    TEST(NULL == pEvent->pTest);
    TEST(NULL != pEvent->pFailure);
    if (4 == CU_get_number_of_failure_records()) {
      TEST(NULL != pEvent->pFailure->pNext);
      TEST(NULL != pEvent->pFailure->pNext->pNext);
      TEST(NULL != pEvent->pFailure->pNext->pNext->pNext);
      TEST(NULL == pEvent->pFailure->pNext->pNext->pNext->pNext);
    }
    TEST(pEvent->pFailure == CU_get_failure_list());
  }

  test_results(2,2,0,4,2,0,4,2,2,4);

  /* clear handlers and run again */
  CU_set_suite_start_handler(NULL);
  CU_set_test_start_handler(NULL);
  CU_set_test_complete_handler(NULL);
  CU_set_suite_complete_handler(NULL);
  CU_set_all_test_complete_handler(NULL);
  CU_set_suite_init_failure_handler(NULL);
  CU_set_suite_cleanup_failure_handler(NULL);

  TEST(NULL == CU_get_suite_start_handler());
  TEST(NULL == CU_get_test_start_handler());
  TEST(NULL == CU_get_test_complete_handler());
  TEST(NULL == CU_get_suite_complete_handler());
  TEST(NULL == CU_get_all_test_complete_handler());
  TEST(NULL == CU_get_suite_init_failure_handler());
  TEST(NULL == CU_get_suite_cleanup_failure_handler());

  clear_test_events();
  CU_run_all_tests();

  TEST(0 == f_nTestEvents);
  TEST(NULL == f_pFirstEvent);
  test_results(2,2,0,4,2,0,4,2,2,4);

  CU_cleanup_registry();
  clear_test_events();
}

static CU_BOOL f_exit_called = CU_FALSE;

/* intercept exit for testing of CUEA_ABORT action */
void test_exit(int status)
{
  CU_UNREFERENCED_PARAMETER(status);  /* not used */
  f_exit_called = CU_TRUE;
}


/*-------------------------------------------------*/
static void test_CU_fail_on_inactive(void)
{
  CU_pSuite pSuite1 = NULL;
  CU_pSuite pSuite2 = NULL;
  CU_pTest pTest1 = NULL;
  CU_pTest pTest2 = NULL;
  CU_pTest pTest3 = NULL;
  CU_pTest pTest4 = NULL;

  CU_set_error_action(CUEA_IGNORE);
  CU_initialize_registry();

  /* register some suites and tests */
  CU_initialize_registry();
  pSuite1 = CU_add_suite("suite1", NULL, NULL);
  pTest1 = CU_add_test(pSuite1, "test1", test_succeed);
  pTest2 = CU_add_test(pSuite1, "test2", test_fail);
  pSuite2 = CU_add_suite("suite2", suite_fail, NULL);
  pTest3 = CU_add_test(pSuite2, "test3", test_succeed);
  pTest4 = CU_add_test(pSuite2, "test4", test_succeed);

  /* test initial conditions */
  TEST(CU_TRUE == CU_get_fail_on_inactive());
  TEST(CU_TRUE == pSuite1->fActive);
  TEST(CU_TRUE == pSuite2->fActive);
  TEST(CU_TRUE == pTest1->fActive);
  TEST(CU_TRUE == pTest2->fActive);
  TEST(CU_TRUE == pTest3->fActive);
  TEST(CU_TRUE == pTest4->fActive);

  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CU_TRUE == CU_get_fail_on_inactive());
  TEST(CUE_SINIT_FAILED == CU_run_all_tests());   /* all suites/tests active */
  test_results(1,1,0,2,1,0,2,1,1,2);
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CU_FALSE == CU_get_fail_on_inactive());
  TEST(CUE_SINIT_FAILED == CU_run_all_tests());
  test_results(1,1,0,2,1,0,2,1,1,2);

  CU_set_suite_active(pSuite1, CU_FALSE);
  CU_set_suite_active(pSuite2, CU_FALSE);
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_SUITE_INACTIVE == CU_run_all_tests());   /* all suites inactive */
  test_results(0,0,2,0,0,0,0,0,0,2);
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SUCCESS == CU_run_all_tests());
  test_results(0,0,2,0,0,0,0,0,0,0);
  CU_set_suite_active(pSuite1, CU_TRUE);
  CU_set_suite_active(pSuite2, CU_TRUE);

  CU_set_suite_active(pSuite2, CU_FALSE);
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_SUITE_INACTIVE == CU_run_all_tests());   /* some suites inactive */
  test_results(1,0,1,2,1,0,2,1,1,2);
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SUCCESS == CU_run_all_tests());
  test_results(1,0,1,2,1,0,2,1,1,1);
  CU_set_suite_active(pSuite2, CU_TRUE);

  CU_set_test_active(pTest1, CU_FALSE);
  CU_set_test_active(pTest2, CU_FALSE);
  CU_set_test_active(pTest3, CU_FALSE);
  CU_set_test_active(pTest4, CU_FALSE);
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_TEST_INACTIVE == CU_run_all_tests());    /* all tests inactive */
  test_results(1,1,0,0,0,2,0,0,0,3);
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SINIT_FAILED == CU_run_all_tests());
  test_results(1,1,0,0,0,2,0,0,0,1);
  CU_set_test_active(pTest1, CU_TRUE);
  CU_set_test_active(pTest2, CU_TRUE);
  CU_set_test_active(pTest3, CU_TRUE);
  CU_set_test_active(pTest4, CU_TRUE);

  CU_set_test_active(pTest2, CU_FALSE);
  CU_set_test_active(pTest4, CU_FALSE);
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_TEST_INACTIVE == CU_run_all_tests());    /* some tests inactive */
  test_results(1,1,0,1,0,1,1,1,0,2);
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SINIT_FAILED == CU_run_all_tests());
  test_results(1,1,0,1,0,1,1,1,0,1);
  CU_set_test_active(pTest2, CU_TRUE);
  CU_set_test_active(pTest4, CU_TRUE);

  CU_set_suite_active(pSuite2, CU_FALSE);
  CU_set_test_active(pTest1, CU_FALSE);
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_TEST_INACTIVE == CU_run_all_tests());    /* some suites & tests inactive */
  test_results(1,0,1,1,1,1,1,0,1,3);
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SUCCESS == CU_run_all_tests());
  test_results(1,0,1,1,1,1,1,0,1,1);
  CU_set_suite_active(pSuite2, CU_TRUE);
  CU_set_test_active(pTest1, CU_TRUE);

  /* clean up */
  CU_cleanup_registry();
}

/*-------------------------------------------------*/
static void test_CU_run_all_tests(void)
{
  CU_pSuite pSuite1 = NULL;
  CU_pSuite pSuite2 = NULL;
  CU_pSuite pSuite3 = NULL;
  CU_pSuite pSuite4 = NULL;
  CU_pTest pTest1 = NULL;
  CU_pTest pTest2 = NULL;
  CU_pTest pTest3 = NULL;
  CU_pTest pTest4 = NULL;
  CU_pTest pTest5 = NULL;
  CU_pTest pTest6 = NULL;
  CU_pTest pTest7 = NULL;
  CU_pTest pTest8 = NULL;
  CU_pTest pTest9 = NULL;
  CU_pTest pTest10 = NULL;

  /* error - uninitialized registry  (CUEA_IGNORE) */
  CU_cleanup_registry();
  CU_set_error_action(CUEA_IGNORE);

  TEST(CUE_NOREGISTRY == CU_run_all_tests());
  TEST(CUE_NOREGISTRY == CU_get_error());

  /* error - uninitialized registry  (CUEA_FAIL) */
  CU_cleanup_registry();
  CU_set_error_action(CUEA_FAIL);

  TEST(CUE_NOREGISTRY == CU_run_all_tests());
  TEST(CUE_NOREGISTRY == CU_get_error());

  /* error - uninitialized registry  (CUEA_ABORT) */
  CU_cleanup_registry();
  CU_set_error_action(CUEA_ABORT);

  f_exit_called = CU_FALSE;
  CU_run_all_tests();
  TEST(CU_TRUE == f_exit_called);
  f_exit_called = CU_FALSE;

  /* run with no suites or tests registered */
  CU_initialize_registry();

  CU_set_error_action(CUEA_IGNORE);
  TEST(CUE_SUCCESS == CU_run_all_tests());
  test_results(0,0,0,0,0,0,0,0,0,0);

  /* register some suites and tests */
  CU_initialize_registry();
  pSuite1 = CU_add_suite("suite1", NULL, NULL);
  pTest1 = CU_add_test(pSuite1, "test1", test_succeed);
  pTest2 = CU_add_test(pSuite1, "test2", test_fail);
  pTest3 = CU_add_test(pSuite1, "test1", test_succeed); /* duplicate test name OK */
  pTest4 = CU_add_test(pSuite1, "test4", test_fail);
  pTest5 = CU_add_test(pSuite1, "test1", test_succeed); /* duplicate test name OK */
  pSuite2 = CU_add_suite("suite2", suite_fail, NULL);
  pTest6 = CU_add_test(pSuite2, "test6", test_succeed);
  pTest7 = CU_add_test(pSuite2, "test7", test_succeed);
  pSuite3 = CU_add_suite("suite1", NULL, NULL);         /* duplicate suite name OK */
  pTest8 = CU_add_test(pSuite3, "test8", test_fail);
  pTest9 = CU_add_test(pSuite3, "test9", test_succeed);
  pSuite4 = CU_add_suite("suite4", NULL, suite_fail);
  pTest10 = CU_add_test(pSuite4, "test10", test_succeed);

  TEST_FATAL(4 == CU_get_registry()->uiNumberOfSuites);
  TEST_FATAL(10 == CU_get_registry()->uiNumberOfTests);

  /* run all tests (CUEA_IGNORE) */
  CU_set_error_action(CUEA_IGNORE);
  TEST(CUE_SINIT_FAILED == CU_run_all_tests());     /* all suites/tests active */
  test_results(3,2,0,8,3,0,8,5,3,5);

  CU_set_suite_active(pSuite1, CU_FALSE);
  CU_set_suite_active(pSuite2, CU_FALSE);
  CU_set_suite_active(pSuite3, CU_FALSE);
  CU_set_suite_active(pSuite4, CU_FALSE);
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SUCCESS == CU_run_all_tests());          /* suites inactive */
  test_results(0,0,4,0,0,0,0,0,0,0);
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_SUITE_INACTIVE == CU_run_all_tests());
  test_results(0,0,4,0,0,0,0,0,0,4);

  CU_set_suite_active(pSuite1, CU_FALSE);
  CU_set_suite_active(pSuite2, CU_TRUE);
  CU_set_suite_active(pSuite3, CU_TRUE);
  CU_set_suite_active(pSuite4, CU_FALSE);
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SINIT_FAILED == CU_run_all_tests());   /* some suites inactive */
  test_results(1,1,2,2,1,0,2,1,1,2);
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_SUITE_INACTIVE == CU_run_all_tests());
  test_results(1,1,2,2,1,0,2,1,1,4);

  CU_set_suite_active(pSuite1, CU_TRUE);
  CU_set_suite_active(pSuite2, CU_TRUE);
  CU_set_suite_active(pSuite3, CU_TRUE);
  CU_set_suite_active(pSuite4, CU_TRUE);

  CU_set_test_active(pTest1, CU_FALSE);
  CU_set_test_active(pTest2, CU_FALSE);
  CU_set_test_active(pTest3, CU_FALSE);
  CU_set_test_active(pTest4, CU_FALSE);
  CU_set_test_active(pTest5, CU_FALSE);
  CU_set_test_active(pTest6, CU_FALSE);
  CU_set_test_active(pTest7, CU_FALSE);
  CU_set_test_active(pTest8, CU_FALSE);
  CU_set_test_active(pTest9, CU_FALSE);
  CU_set_test_active(pTest10, CU_FALSE);
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SINIT_FAILED == CU_run_all_tests());   /* no tests active */
  test_results(3,2,0,0,0,8,0,0,0,2);
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_TEST_INACTIVE == CU_run_all_tests());
  test_results(3,2,0,0,0,8,0,0,0,10);

  CU_set_test_active(pTest1, CU_TRUE);
  CU_set_test_active(pTest2, CU_FALSE);
  CU_set_test_active(pTest3, CU_TRUE);
  CU_set_test_active(pTest4, CU_FALSE);
  CU_set_test_active(pTest5, CU_TRUE);
  CU_set_test_active(pTest6, CU_FALSE);
  CU_set_test_active(pTest7, CU_TRUE);
  CU_set_test_active(pTest8, CU_FALSE);
  CU_set_test_active(pTest9, CU_TRUE);
  CU_set_test_active(pTest10, CU_FALSE);
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SINIT_FAILED == CU_run_all_tests());   /* some tests active */
  test_results(3,2,0,4,0,4,4,4,0,2);
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_TEST_INACTIVE == CU_run_all_tests());
  test_results(3,2,0,4,0,4,4,4,0,6);

  CU_set_test_active(pTest1, CU_TRUE);
  CU_set_test_active(pTest2, CU_TRUE);
  CU_set_test_active(pTest3, CU_TRUE);
  CU_set_test_active(pTest4, CU_TRUE);
  CU_set_test_active(pTest5, CU_TRUE);
  CU_set_test_active(pTest6, CU_TRUE);
  CU_set_test_active(pTest7, CU_TRUE);
  CU_set_test_active(pTest8, CU_TRUE);
  CU_set_test_active(pTest9, CU_TRUE);
  CU_set_test_active(pTest10, CU_TRUE);

  CU_set_suite_initfunc(pSuite1, &suite_fail);
  TEST(CUE_SINIT_FAILED == CU_run_all_tests());   /* change a suite init function */
  CU_set_suite_initfunc(pSuite1, NULL);
  test_results(2,3,0,3,1,0,3,2,1,4);

  CU_set_suite_cleanupfunc(pSuite4, NULL);
  TEST(CUE_SINIT_FAILED == CU_run_all_tests());   /* change a suite cleanup function */
  CU_set_suite_cleanupfunc(pSuite4, &suite_fail);
  test_results(3,1,0,8,3,0,8,5,3,4);

  CU_set_test_func(pTest2, &test_succeed);
  CU_set_test_func(pTest4, &test_succeed);
  CU_set_test_func(pTest8, &test_succeed);
  TEST(CUE_SINIT_FAILED == CU_run_all_tests());   /* change a test function */
  CU_set_test_func(pTest2, &test_fail);
  CU_set_test_func(pTest4, &test_fail);
  CU_set_test_func(pTest8, &test_fail);
  test_results(3,2,0,8,0,0,8,8,0,2);

  /* run all tests (CUEA_FAIL) */
  CU_set_error_action(CUEA_FAIL);
  TEST(CUE_SINIT_FAILED == CU_run_all_tests()); /* all suites active */
  test_results(1,1,0,5,2,0,5,3,2,3);

  CU_set_suite_active(pSuite1, CU_TRUE);
  CU_set_suite_active(pSuite2, CU_FALSE);
  CU_set_suite_active(pSuite3, CU_FALSE);
  CU_set_suite_active(pSuite4, CU_TRUE);
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SCLEAN_FAILED == CU_run_all_tests()); /* some suites inactive */
  test_results(2,1,2,6,2,0,6,4,2,3);
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_SUITE_INACTIVE == CU_run_all_tests());
  test_results(1,0,1,5,2,0,5,3,2,3);

  CU_set_suite_active(pSuite1, CU_TRUE);
  CU_set_suite_active(pSuite2, CU_TRUE);
  CU_set_suite_active(pSuite3, CU_TRUE);
  CU_set_suite_active(pSuite4, CU_TRUE);

  CU_set_test_active(pTest1, CU_FALSE);
  CU_set_test_active(pTest2, CU_FALSE);
  CU_set_test_active(pTest3, CU_FALSE);
  CU_set_test_active(pTest4, CU_FALSE);
  CU_set_test_active(pTest5, CU_FALSE);
  CU_set_test_active(pTest6, CU_FALSE);
  CU_set_test_active(pTest7, CU_FALSE);
  CU_set_test_active(pTest8, CU_FALSE);
  CU_set_test_active(pTest9, CU_FALSE);
  CU_set_test_active(pTest10, CU_FALSE);
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SINIT_FAILED == CU_run_all_tests());   /* no tests active */
  test_results(1,1,0,0,0,5,0,0,0,1);
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_TEST_INACTIVE == CU_run_all_tests());
  test_results(1,0,0,0,0,1,0,0,0,1);

  CU_set_test_active(pTest1, CU_FALSE);
  CU_set_test_active(pTest2, CU_TRUE);
  CU_set_test_active(pTest3, CU_FALSE);
  CU_set_test_active(pTest4, CU_TRUE);
  CU_set_test_active(pTest5, CU_FALSE);
  CU_set_test_active(pTest6, CU_TRUE);
  CU_set_test_active(pTest7, CU_FALSE);
  CU_set_test_active(pTest8, CU_TRUE);
  CU_set_test_active(pTest9, CU_FALSE);
  CU_set_test_active(pTest10, CU_TRUE);
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SINIT_FAILED == CU_run_all_tests());   /* some tests active */
  test_results(1,1,0,2,2,3,2,0,2,3);
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_TEST_INACTIVE == CU_run_all_tests());
  test_results(1,0,0,0,0,1,0,0,0,1);

  CU_set_test_active(pTest1, CU_TRUE);
  CU_set_test_active(pTest2, CU_TRUE);
  CU_set_test_active(pTest3, CU_TRUE);
  CU_set_test_active(pTest4, CU_TRUE);
  CU_set_test_active(pTest5, CU_TRUE);
  CU_set_test_active(pTest6, CU_TRUE);
  CU_set_test_active(pTest7, CU_TRUE);
  CU_set_test_active(pTest8, CU_TRUE);
  CU_set_test_active(pTest9, CU_TRUE);
  CU_set_test_active(pTest10, CU_TRUE);

  CU_set_suite_initfunc(pSuite2, NULL);
  TEST(CUE_SCLEAN_FAILED == CU_run_all_tests());   /* change a suite init function */
  CU_set_suite_initfunc(pSuite2, &suite_fail);
  test_results(4,1,0,10,3,0,10,7,3,4);

  CU_set_suite_cleanupfunc(pSuite1, &suite_fail);
  TEST(CUE_SCLEAN_FAILED == CU_run_all_tests());   /* change a suite cleanup function */
  CU_set_suite_cleanupfunc(pSuite1, NULL);
  test_results(1,1,0,5,2,0,5,3,2,3);

  CU_set_test_func(pTest1, &test_fail);
  CU_set_test_func(pTest3, &test_fail);
  CU_set_test_func(pTest5, &test_fail);
  CU_set_test_func(pTest9, &test_fail);
  CU_set_test_func(pTest10, &test_fail);
  TEST(CUE_SINIT_FAILED == CU_run_all_tests());   /* change a test function */
  CU_set_test_func(pTest1, &test_succeed);
  CU_set_test_func(pTest3, &test_succeed);
  CU_set_test_func(pTest5, &test_succeed);
  CU_set_test_func(pTest9, &test_succeed);
  CU_set_test_func(pTest10, &test_succeed);
  test_results(1,1,0,5,5,0,5,0,5,6);

  /* run all tests (CUEA_ABORT) */
  f_exit_called = CU_FALSE;
  CU_set_error_action(CUEA_ABORT);
  CU_set_suite_active(pSuite1, CU_TRUE);
  CU_set_suite_active(pSuite2, CU_TRUE);
  CU_set_suite_active(pSuite3, CU_TRUE);
  CU_set_suite_active(pSuite4, CU_TRUE);
  TEST(CUE_SINIT_FAILED == CU_run_all_tests()); /* all suites active */
  TEST(CU_TRUE == f_exit_called);
  test_results(1,1,0,5,2,0,5,3,2,3);

  CU_set_suite_active(pSuite1, CU_FALSE);
  CU_set_suite_active(pSuite2, CU_FALSE);
  CU_set_suite_active(pSuite3, CU_FALSE);
  CU_set_suite_active(pSuite4, CU_FALSE);
  f_exit_called = CU_FALSE;
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SUCCESS == CU_run_all_tests());         /* no suites active, so no abort() */
  TEST(CU_FALSE == f_exit_called);
  test_results(0,0,4,0,0,0,0,0,0,0);
  f_exit_called = CU_FALSE;
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_SUITE_INACTIVE == CU_run_all_tests());
  TEST(CU_TRUE == f_exit_called);
  test_results(0,0,1,0,0,0,0,0,0,1);

  CU_set_suite_active(pSuite1, CU_TRUE);
  CU_set_suite_active(pSuite2, CU_FALSE);
  CU_set_suite_active(pSuite3, CU_TRUE);
  CU_set_suite_active(pSuite4, CU_TRUE);
  f_exit_called = CU_FALSE;
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SCLEAN_FAILED == CU_run_all_tests()); /* some suites active */
  TEST(CU_TRUE == f_exit_called);
  test_results(3,1,1,8,3,0,8,5,3,4);
  f_exit_called = CU_FALSE;
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_SUITE_INACTIVE == CU_run_all_tests());
  TEST(CU_TRUE == f_exit_called);
  test_results(1,0,1,5,2,0,5,3,2,3);

  CU_set_suite_active(pSuite1, CU_TRUE);
  CU_set_suite_active(pSuite2, CU_TRUE);
  CU_set_suite_active(pSuite3, CU_TRUE);
  CU_set_suite_active(pSuite4, CU_TRUE);

  CU_set_test_active(pTest1, CU_FALSE);
  CU_set_test_active(pTest2, CU_FALSE);
  CU_set_test_active(pTest3, CU_FALSE);
  CU_set_test_active(pTest4, CU_FALSE);
  CU_set_test_active(pTest5, CU_FALSE);
  CU_set_test_active(pTest6, CU_FALSE);
  CU_set_test_active(pTest7, CU_FALSE);
  CU_set_test_active(pTest8, CU_FALSE);
  CU_set_test_active(pTest9, CU_FALSE);
  CU_set_test_active(pTest10, CU_FALSE);
  f_exit_called = CU_FALSE;
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SINIT_FAILED == CU_run_all_tests());   /* no tests active */
  TEST(CU_TRUE == f_exit_called);
  test_results(1,1,0,0,0,5,0,0,0,1);
  f_exit_called = CU_FALSE;
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_TEST_INACTIVE == CU_run_all_tests());
  TEST(CU_TRUE == f_exit_called);
  test_results(1,0,0,0,0,1,0,0,0,1);

  CU_set_test_active(pTest1, CU_FALSE);
  CU_set_test_active(pTest2, CU_TRUE);
  CU_set_test_active(pTest3, CU_FALSE);
  CU_set_test_active(pTest4, CU_TRUE);
  CU_set_test_active(pTest5, CU_FALSE);
  CU_set_test_active(pTest6, CU_TRUE);
  CU_set_test_active(pTest7, CU_FALSE);
  CU_set_test_active(pTest8, CU_TRUE);
  CU_set_test_active(pTest9, CU_FALSE);
  CU_set_test_active(pTest10, CU_TRUE);
  f_exit_called = CU_FALSE;
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SINIT_FAILED == CU_run_all_tests());   /* some tests active */
  TEST(CU_TRUE == f_exit_called);
  test_results(1,1,0,2,2,3,2,0,2,3);
  f_exit_called = CU_FALSE;
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_TEST_INACTIVE == CU_run_all_tests());
  TEST(CU_TRUE == f_exit_called);
  test_results(1,0,0,0,0,1,0,0,0,1);

  CU_set_test_active(pTest1, CU_TRUE);
  CU_set_test_active(pTest2, CU_TRUE);
  CU_set_test_active(pTest3, CU_TRUE);
  CU_set_test_active(pTest4, CU_TRUE);
  CU_set_test_active(pTest5, CU_TRUE);
  CU_set_test_active(pTest6, CU_TRUE);
  CU_set_test_active(pTest7, CU_TRUE);
  CU_set_test_active(pTest8, CU_TRUE);
  CU_set_test_active(pTest9, CU_TRUE);
  CU_set_test_active(pTest10, CU_TRUE);

  f_exit_called = CU_FALSE;
  CU_set_suite_initfunc(pSuite1, &suite_fail);
  TEST(CUE_SINIT_FAILED == CU_run_all_tests());   /* change a suite init function */
  CU_set_suite_initfunc(pSuite1, NULL);
  TEST(CU_TRUE == f_exit_called);
  test_results(0,1,0,0,0,0,0,0,0,1);

  f_exit_called = CU_FALSE;
  CU_set_suite_cleanupfunc(pSuite1, &suite_fail);
  TEST(CUE_SCLEAN_FAILED == CU_run_all_tests());   /* change a suite cleanup function */
  CU_set_suite_cleanupfunc(pSuite1, NULL);
  TEST(CU_TRUE == f_exit_called);
  test_results(1,1,0,5,2,0,5,3,2,3);

  f_exit_called = CU_FALSE;
  CU_set_test_func(pTest1, &test_fail);
  CU_set_test_func(pTest3, &test_fail);
  CU_set_test_func(pTest5, &test_fail);
  CU_set_test_func(pTest9, &test_fail);
  CU_set_test_func(pTest10, &test_fail);
  TEST(CUE_SINIT_FAILED == CU_run_all_tests());   /* change a test function */
  CU_set_test_func(pTest1, &test_succeed);
  CU_set_test_func(pTest3, &test_succeed);
  CU_set_test_func(pTest5, &test_succeed);
  CU_set_test_func(pTest9, &test_succeed);
  CU_set_test_func(pTest10, &test_succeed);
  TEST(CU_TRUE == f_exit_called);
  test_results(1,1,0,5,5,0,5,0,5,6);

  /* clean up after testing */
  CU_set_error_action(CUEA_IGNORE);
  CU_cleanup_registry();
}


/*-------------------------------------------------*/
static void test_CU_run_suite(void)
{
  CU_pSuite pSuite1 = NULL;
  CU_pSuite pSuite2 = NULL;
  CU_pSuite pSuite3 = NULL;
  CU_pSuite pSuite4 = NULL;
  CU_pSuite pSuite5 = NULL;
  CU_pSuite pSuite6 = NULL;
  CU_pTest pTest1 = NULL;
  CU_pTest pTest2 = NULL;
  CU_pTest pTest3 = NULL;
  CU_pTest pTest4 = NULL;
  CU_pTest pTest5 = NULL;
  CU_pTest pTest6 = NULL;
  CU_pTest pTest7 = NULL;
  CU_pTest pTest8 = NULL;
  CU_pTest pTest9 = NULL;
  CU_pTest pTest10 = NULL;
  CU_pTest pTest11 = NULL;

  /* error - NULL suite (CUEA_IGNORE) */
  CU_set_error_action(CUEA_IGNORE);

  TEST(CUE_NOSUITE == CU_run_suite(NULL));
  TEST(CUE_NOSUITE == CU_get_error());

  /* error - NULL suite (CUEA_FAIL) */
  CU_set_error_action(CUEA_FAIL);

  TEST(CUE_NOSUITE == CU_run_suite(NULL));
  TEST(CUE_NOSUITE == CU_get_error());

  /* error - NULL suite (CUEA_ABORT) */
  CU_set_error_action(CUEA_ABORT);

  f_exit_called = CU_FALSE;
  CU_run_suite(NULL);
  TEST(CU_TRUE == f_exit_called);
  f_exit_called = CU_FALSE;

  /* register some suites and tests */
  CU_initialize_registry();
  pSuite1 = CU_add_suite("suite1", NULL, NULL);
  pTest1 = CU_add_test(pSuite1, "test1", test_succeed);
  pTest2 = CU_add_test(pSuite1, "test2", test_fail);
  pTest3 = CU_add_test(pSuite1, "test3", test_succeed);
  pTest4 = CU_add_test(pSuite1, "test4", test_fail);
  pTest5 = CU_add_test(pSuite1, "test5", test_succeed);
  pSuite2 = CU_add_suite("suite1", suite_fail, NULL);   /* duplicate suite name OK */
  pTest6 = CU_add_test(pSuite2, "test6", test_succeed);
  pTest7 = CU_add_test(pSuite2, "test7", test_succeed);
  pSuite3 = CU_add_suite("suite3", NULL, suite_fail);
  pTest8 = CU_add_test(pSuite3, "test8", test_fail);
  pTest9 = CU_add_test(pSuite3, "test8", test_succeed); /* duplicate test name OK */
  pSuite4 = CU_add_suite("suite4", NULL, NULL);
  pSuite5 = CU_add_suite_with_setup_and_teardown("suite5", NULL, NULL, suite_setup, suite_teardown);
  pTest10 = CU_add_test(pSuite5, "test10", test_succeed_if_setup);
  pSuite6 = CU_add_suite("suite6", NULL, NULL);
  pTest11 = CU_add_test(pSuite6, "test11", test_fail_if_not_setup);

  TEST_FATAL(6 == CU_get_registry()->uiNumberOfSuites);
  TEST_FATAL(11 == CU_get_registry()->uiNumberOfTests);

  /* run each suite (CUEA_IGNORE) */
  CU_set_error_action(CUEA_IGNORE);

  TEST(CUE_SUCCESS == CU_run_suite(pSuite1));   /* suites/tests active */
  test_results(1,0,0,5,2,0,5,3,2,2);

  TEST(CUE_SINIT_FAILED == CU_run_suite(pSuite2));
  test_results(0,1,0,0,0,0,0,0,0,1);

  TEST(CUE_SCLEAN_FAILED == CU_run_suite(pSuite3));
  test_results(1,1,0,2,1,0,2,1,1,2);

  TEST(CUE_SUCCESS == CU_run_suite(pSuite4));
  test_results(1,0,0,0,0,0,0,0,0,0);

  TEST(CUE_SUCCESS == CU_run_suite(pSuite5));
  test_results(1,0,0,1,0,0,1,1,0,0);

  TEST(CUE_SUCCESS == CU_run_suite(pSuite6));
  test_results(1,0,0,1,1,0,1,0,1,1);

  CU_set_suite_active(pSuite3, CU_FALSE);
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SUCCESS == CU_run_suite(pSuite3));   /* suite inactive */
  test_results(0,0,1,0,0,0,0,0,0,0);
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_SUITE_INACTIVE == CU_run_suite(pSuite3));
  test_results(0,0,1,0,0,0,0,0,0,1);
  CU_set_suite_active(pSuite3, CU_TRUE);

  CU_set_test_active(pTest1, CU_FALSE);
  CU_set_test_active(pTest2, CU_FALSE);
  CU_set_test_active(pTest3, CU_FALSE);
  CU_set_test_active(pTest4, CU_FALSE);
  CU_set_test_active(pTest5, CU_FALSE);
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SUCCESS == CU_run_suite(pSuite1));   /* all tests inactive */
  test_results(1,0,0,0,0,5,0,0,0,0);
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_TEST_INACTIVE == CU_run_suite(pSuite1));
  test_results(1,0,0,0,0,5,0,0,0,5);

  CU_set_test_active(pTest1, CU_TRUE);
  CU_set_test_active(pTest2, CU_FALSE);
  CU_set_test_active(pTest3, CU_TRUE);
  CU_set_test_active(pTest4, CU_FALSE);
  CU_set_test_active(pTest5, CU_TRUE);
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SUCCESS == CU_run_suite(pSuite1));   /* some tests inactive */
  test_results(1,0,0,3,0,2,3,3,0,0);
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_TEST_INACTIVE == CU_run_suite(pSuite1));
  test_results(1,0,0,3,0,2,3,3,0,2);
  CU_set_test_active(pTest2, CU_TRUE);
  CU_set_test_active(pTest4, CU_TRUE);

  CU_set_suite_initfunc(pSuite1, &suite_fail);
  TEST(CUE_SINIT_FAILED == CU_run_suite(pSuite1));   /* change a suite init function */
  CU_set_suite_initfunc(pSuite1, NULL);
  test_results(0,1,0,0,0,0,0,0,0,1);

  CU_set_suite_cleanupfunc(pSuite1, &suite_fail);
  TEST(CUE_SCLEAN_FAILED == CU_run_suite(pSuite1));   /* change a suite cleanup function */
  CU_set_suite_cleanupfunc(pSuite1, NULL);
  test_results(1,1,0,5,2,0,5,3,2,3);

  CU_set_test_func(pTest1, &test_fail);
  CU_set_test_func(pTest3, &test_fail);
  CU_set_test_func(pTest5, &test_fail);
  TEST(CUE_SUCCESS == CU_run_suite(pSuite1));   /* change a test function */
  CU_set_test_func(pTest1, &test_succeed);
  CU_set_test_func(pTest3, &test_succeed);
  CU_set_test_func(pTest5, &test_succeed);
  test_results(1,0,0,5,5,0,5,0,5,5);

  /* run each suite (CUEA_FAIL) */
  CU_set_error_action(CUEA_FAIL);

  TEST(CUE_SUCCESS == CU_run_suite(pSuite1));   /* suite active */
  test_results(1,0,0,5,2,0,5,3,2,2);

  TEST(CUE_SINIT_FAILED == CU_run_suite(pSuite2));
  test_results(0,1,0,0,0,0,0,0,0,1);

  TEST(CUE_SCLEAN_FAILED == CU_run_suite(pSuite3));
  test_results(1,1,0,2,1,0,2,1,1,2);

  TEST(CUE_SUCCESS == CU_run_suite(pSuite4));
  test_results(1,0,0,0,0,0,0,0,0,0);

  CU_set_suite_active(pSuite1, CU_FALSE);
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SUCCESS == CU_run_suite(pSuite1));         /* suite inactive */
  test_results(0,0,1,0,0,0,0,0,0,0);
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_SUITE_INACTIVE == CU_run_suite(pSuite1));
  test_results(0,0,1,0,0,0,0,0,0,1);
  CU_set_suite_active(pSuite1, CU_TRUE);

  CU_set_test_active(pTest8, CU_FALSE);
  CU_set_test_active(pTest9, CU_FALSE);
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SCLEAN_FAILED == CU_run_suite(pSuite3));   /* all tests inactive */
  test_results(1,1,0,0,0,2,0,0,0,1);
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_TEST_INACTIVE == CU_run_suite(pSuite3));
  test_results(1,1,0,0,0,1,0,0,0,2);

  CU_set_test_active(pTest8, CU_TRUE);
  CU_set_test_active(pTest9, CU_FALSE);
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SCLEAN_FAILED == CU_run_suite(pSuite3));   /* some tests inactive */
  test_results(1,1,0,1,1,1,1,0,1,2);
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_TEST_INACTIVE == CU_run_suite(pSuite3));
  test_results(1,1,0,1,1,1,1,0,1,3);
  CU_set_test_active(pTest9, CU_TRUE);

  CU_set_suite_initfunc(pSuite2, NULL);
  TEST(CUE_SUCCESS == CU_run_suite(pSuite2));         /* change a suite init function */
  CU_set_suite_initfunc(pSuite2, &suite_fail);
  test_results(1,0,0,2,0,0,2,2,0,0);

  CU_set_suite_cleanupfunc(pSuite1, &suite_fail);
  TEST(CUE_SCLEAN_FAILED == CU_run_suite(pSuite1));   /* change a suite cleanup function */
  CU_set_suite_cleanupfunc(pSuite1, NULL);
  test_results(1,1,0,5,2,0,5,3,2,3);

  CU_set_test_func(pTest2, &test_succeed);
  CU_set_test_func(pTest4, &test_succeed);
  TEST(CUE_SUCCESS == CU_run_suite(pSuite1));   /* change a test function */
  CU_set_test_func(pTest2, &test_fail);
  CU_set_test_func(pTest4, &test_fail);
  test_results(1,0,0,5,0,0,5,5,0,0);

  /* run each suite (CUEA_ABORT) */
  CU_set_error_action(CUEA_ABORT);

  f_exit_called = CU_FALSE;
  TEST(CUE_SUCCESS == CU_run_suite(pSuite1));   /* suite active */
  TEST(CU_FALSE == f_exit_called);
  test_results(1,0,0,5,2,0,5,3,2,2);

  f_exit_called = CU_FALSE;
  TEST(CUE_SINIT_FAILED == CU_run_suite(pSuite2));
  TEST(CU_TRUE == f_exit_called);
  f_exit_called = CU_FALSE;
  test_results(0,1,0,0,0,0,0,0,0,1);

  f_exit_called = CU_FALSE;
  TEST(CUE_SCLEAN_FAILED == CU_run_suite(pSuite3));
  TEST(CU_TRUE == f_exit_called);
  test_results(1,1,0,2,1,0,2,1,1,2);

  f_exit_called = CU_FALSE;
  TEST(CUE_SUCCESS == CU_run_suite(pSuite4));
  TEST(CU_FALSE == f_exit_called);
  test_results(1,0,0,0,0,0,0,0,0,0);

  CU_set_suite_active(pSuite2, CU_FALSE);
  f_exit_called = CU_FALSE;
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SUCCESS == CU_run_suite(pSuite2));         /* suite inactive, but not a failure */
  TEST(CU_FALSE == f_exit_called);
  test_results(0,0,1,0,0,0,0,0,0,0);
  f_exit_called = CU_FALSE;
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_SUITE_INACTIVE == CU_run_suite(pSuite2));
  TEST(CU_TRUE == f_exit_called);
  test_results(0,0,1,0,0,0,0,0,0,1);
  CU_set_suite_active(pSuite2, CU_TRUE);

  CU_set_test_active(pTest8, CU_FALSE);
  CU_set_test_active(pTest9, CU_FALSE);
  f_exit_called = CU_FALSE;
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SCLEAN_FAILED == CU_run_suite(pSuite3));   /* all tests inactive */
  TEST(CU_TRUE == f_exit_called);
  test_results(1,1,0,0,0,2,0,0,0,1);
  f_exit_called = CU_FALSE;
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_TEST_INACTIVE == CU_run_suite(pSuite3));
  TEST(CU_TRUE == f_exit_called);
  test_results(1,1,0,0,0,1,0,0,0,2);

  CU_set_test_active(pTest8, CU_FALSE);
  CU_set_test_active(pTest9, CU_TRUE);
  f_exit_called = CU_FALSE;
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SCLEAN_FAILED == CU_run_suite(pSuite3));   /* some tests inactive */
  TEST(CU_TRUE == f_exit_called);
  test_results(1,1,0,1,0,1,1,1,0,1);
  f_exit_called = CU_FALSE;
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_TEST_INACTIVE == CU_run_suite(pSuite3));
  TEST(CU_TRUE == f_exit_called);
  test_results(1,1,0,0,0,1,0,0,0,2);
  CU_set_test_active(pTest8, CU_TRUE);

  f_exit_called = CU_FALSE;
  CU_set_suite_initfunc(pSuite1, &suite_fail);
  TEST(CUE_SINIT_FAILED == CU_run_suite(pSuite1));    /* change a suite init function */
  CU_set_suite_initfunc(pSuite1, NULL);
  TEST(CU_TRUE == f_exit_called);
  test_results(0,1,0,0,0,0,0,0,0,1);

  f_exit_called = CU_FALSE;
  CU_set_suite_cleanupfunc(pSuite1, &suite_fail);
  TEST(CUE_SCLEAN_FAILED == CU_run_suite(pSuite1));   /* change a suite cleanup function */
  CU_set_suite_cleanupfunc(pSuite1, NULL);
  TEST(CU_TRUE == f_exit_called);
  test_results(1,1,0,5,2,0,5,3,2,3);

  f_exit_called = CU_FALSE;
  CU_set_test_func(pTest8, &test_succeed);
  CU_set_test_func(pTest9, &test_fail);
  TEST(CUE_SCLEAN_FAILED == CU_run_suite(pSuite3));   /* change a test function */
  CU_set_test_func(pTest8, &test_fail);
  CU_set_test_func(pTest9, &test_succeed);
  TEST(CU_TRUE == f_exit_called);
  test_results(1,1,0,2,1,0,2,1,1,2);

  /* clean up after testing */
  CU_set_error_action(CUEA_IGNORE);
  CU_cleanup_registry();
}

/*-------------------------------------------------*/
static void test_CU_run_test(void)
{
  CU_pSuite pSuite1 = NULL;
  CU_pSuite pSuite2 = NULL;
  CU_pSuite pSuite3 = NULL;
  CU_pTest pTest1 = NULL;
  CU_pTest pTest2 = NULL;
  CU_pTest pTest3 = NULL;
  CU_pTest pTest4 = NULL;
  CU_pTest pTest5 = NULL;
  CU_pTest pTest6 = NULL;
  CU_pTest pTest7 = NULL;
  CU_pTest pTest8 = NULL;
  CU_pTest pTest9 = NULL;

  /* register some suites and tests */
  CU_initialize_registry();
  pSuite1 = CU_add_suite("suite1", NULL, NULL);
  pTest1 = CU_add_test(pSuite1, "test1", test_succeed);
  pTest2 = CU_add_test(pSuite1, "test2", test_fail);
  pTest3 = CU_add_test(pSuite1, "test3", test_succeed);
  pTest4 = CU_add_test(pSuite1, "test4", test_fail);
  pTest5 = CU_add_test(pSuite1, "test5", test_succeed);
  pSuite2 = CU_add_suite("suite2", suite_fail, NULL);
  pTest6 = CU_add_test(pSuite2, "test6", test_succeed);
  pTest7 = CU_add_test(pSuite2, "test7", test_succeed);
  pSuite3 = CU_add_suite("suite2", NULL, suite_fail);   /* duplicate suite name OK */
  pTest8 = CU_add_test(pSuite3, "test8", test_fail);
  pTest9 = CU_add_test(pSuite3, "test8", test_succeed); /* duplicate test name OK */

  TEST_FATAL(3 == CU_get_registry()->uiNumberOfSuites);
  TEST_FATAL(9 == CU_get_registry()->uiNumberOfTests);

  /* error - NULL suite (CUEA_IGNORE) */
  CU_set_error_action(CUEA_IGNORE);

  TEST(CUE_NOSUITE == CU_run_test(NULL, pTest1));
  TEST(CUE_NOSUITE == CU_get_error());

  /* error - NULL suite (CUEA_FAIL) */
  CU_set_error_action(CUEA_FAIL);

  TEST(CUE_NOSUITE == CU_run_test(NULL, pTest1));
  TEST(CUE_NOSUITE == CU_get_error());

  /* error - NULL test (CUEA_ABORT) */
  CU_set_error_action(CUEA_ABORT);

  f_exit_called = CU_FALSE;
  CU_run_test(NULL, pTest1);
  TEST(CU_TRUE == f_exit_called);
  f_exit_called = CU_FALSE;

  /* error - NULL test (CUEA_IGNORE) */
  CU_set_error_action(CUEA_IGNORE);

  TEST(CUE_NOTEST == CU_run_test(pSuite1, NULL));
  TEST(CUE_NOTEST == CU_get_error());

  /* error - NULL test (CUEA_FAIL) */
  CU_set_error_action(CUEA_FAIL);

  TEST(CUE_NOTEST == CU_run_test(pSuite1, NULL));
  TEST(CUE_NOTEST == CU_get_error());

  /* error - NULL test (CUEA_ABORT) */
  CU_set_error_action(CUEA_ABORT);

  f_exit_called = CU_FALSE;
  CU_run_test(pSuite1, NULL);
  TEST(CU_TRUE == f_exit_called);
  f_exit_called = CU_FALSE;

  /* error - test not in suite (CUEA_IGNORE) */
  CU_set_error_action(CUEA_IGNORE);

  TEST(CUE_TEST_NOT_IN_SUITE == CU_run_test(pSuite3, pTest1));
  TEST(CUE_TEST_NOT_IN_SUITE == CU_get_error());

  /* error - NULL test (CUEA_FAIL) */
  CU_set_error_action(CUEA_FAIL);

  TEST(CUE_TEST_NOT_IN_SUITE == CU_run_test(pSuite3, pTest1));
  TEST(CUE_TEST_NOT_IN_SUITE == CU_get_error());

  /* error - NULL test (CUEA_ABORT) */
  CU_set_error_action(CUEA_ABORT);

  f_exit_called = CU_FALSE;
  CU_run_test(pSuite3, pTest1);
  TEST(CU_TRUE == f_exit_called);
  f_exit_called = CU_FALSE;

  /* run each test (CUEA_IGNORE) */
  CU_set_error_action(CUEA_IGNORE);

  TEST(CUE_SUCCESS == CU_run_test(pSuite1, pTest1));  /* all suite/tests active */
  test_results(0,0,0,1,0,0,1,1,0,0);

  TEST(CUE_SUCCESS == CU_run_test(pSuite1, pTest2));
  test_results(0,0,0,1,1,0,1,0,1,1);

  TEST(CUE_SUCCESS == CU_run_test(pSuite1, pTest3));
  test_results(0,0,0,1,0,0,1,1,0,0);

  TEST(CUE_SUCCESS == CU_run_test(pSuite1, pTest4));
  test_results(0,0,0,1,1,0,1,0,1,1);

  TEST(CUE_SUCCESS == CU_run_test(pSuite1, pTest5));
  test_results(0,0,0,1,0,0,1,1,0,0);

  TEST(CUE_SINIT_FAILED == CU_run_test(pSuite2, pTest6));
  test_results(0,1,0,0,0,0,0,0,0,1);

  TEST(CUE_SINIT_FAILED == CU_run_test(pSuite2, pTest7));
  test_results(0,1,0,0,0,0,0,0,0,1);

  TEST(CUE_SCLEAN_FAILED == CU_run_test(pSuite3, pTest8));
  test_results(0,1,0,1,1,0,1,0,1,2);

  TEST(CUE_SCLEAN_FAILED == CU_run_test(pSuite3, pTest9));
  test_results(0,1,0,1,0,0,1,1,0,1);

  CU_set_suite_active(pSuite1, CU_FALSE);
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SUITE_INACTIVE == CU_run_test(pSuite1, pTest1));  /* suite inactive */
  test_results(0,0,1,0,0,0,0,0,0,0);
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_SUITE_INACTIVE == CU_run_test(pSuite1, pTest1));
  test_results(0,0,1,0,0,0,0,0,0,1);
  CU_set_suite_active(pSuite1, CU_TRUE);

  CU_set_test_active(pTest1, CU_FALSE);
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_TEST_INACTIVE == CU_run_test(pSuite1, pTest1));   /* test inactive */
  test_results(0,0,0,0,0,1,0,0,0,0);
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_TEST_INACTIVE == CU_run_test(pSuite1, pTest1));
  test_results(0,0,0,0,1,1,0,0,0,1);
  CU_set_test_active(pTest1, CU_TRUE);

  CU_set_suite_initfunc(pSuite1, &suite_fail);
  TEST(CUE_SINIT_FAILED == CU_run_test(pSuite1, pTest1));    /* change a suite init function */
  CU_set_suite_initfunc(pSuite1, NULL);
  test_results(0,1,0,0,0,0,0,0,0,1);

  CU_set_suite_cleanupfunc(pSuite1, &suite_fail);
  TEST(CUE_SCLEAN_FAILED == CU_run_test(pSuite1, pTest1));   /* change a suite cleanup function */
  CU_set_suite_cleanupfunc(pSuite1, NULL);
  test_results(0,1,0,1,0,0,1,1,0,1);

  CU_set_test_func(pTest8, &test_succeed);
  TEST(CUE_SCLEAN_FAILED == CU_run_test(pSuite3, pTest8));   /* change a test function */
  CU_set_test_func(pTest8, &test_fail);
  test_results(0,1,0,1,0,0,1,1,0,1);

  /* run each test (CUEA_FAIL) */
  CU_set_error_action(CUEA_FAIL);

  TEST(CUE_SUCCESS == CU_run_test(pSuite1, pTest1));  /* suite/test active */
  test_results(0,0,0,1,0,0,1,1,0,0);

  TEST(CUE_SUCCESS == CU_run_test(pSuite1, pTest2));
  test_results(0,0,0,1,1,0,1,0,1,1);

  TEST(CUE_SUCCESS == CU_run_test(pSuite1, pTest3));
  test_results(0,0,0,1,0,0,1,1,0,0);

  TEST(CUE_SUCCESS == CU_run_test(pSuite1, pTest4));
  test_results(0,0,0,1,1,0,1,0,1,1);

  TEST(CUE_SUCCESS == CU_run_test(pSuite1, pTest5));
  test_results(0,0,0,1,0,0,1,1,0,0);

  TEST(CUE_SINIT_FAILED == CU_run_test(pSuite2, pTest6));
  test_results(0,1,0,0,0,0,0,0,0,1);

  TEST(CUE_SINIT_FAILED == CU_run_test(pSuite2, pTest7));
  test_results(0,1,0,0,0,0,0,0,0,1);

  TEST(CUE_SCLEAN_FAILED == CU_run_test(pSuite3, pTest8));
  test_results(0,1,0,1,1,0,1,0,1,2);

  TEST(CUE_SCLEAN_FAILED == CU_run_test(pSuite3, pTest9));
  test_results(0,1,0,1,0,0,1,1,0,1);

  CU_set_suite_active(pSuite2, CU_FALSE);
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SUITE_INACTIVE == CU_run_test(pSuite2, pTest7));   /* suite inactive */
  test_results(0,0,1,0,0,0,0,0,0,0);
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_SUITE_INACTIVE == CU_run_test(pSuite2, pTest7));
  test_results(0,0,1,0,0,0,0,0,0,1);
  CU_set_suite_active(pSuite2, CU_TRUE);

  CU_set_test_active(pTest7, CU_FALSE);
  CU_set_fail_on_inactive(CU_FALSE);
  TEST(CUE_SINIT_FAILED == CU_run_test(pSuite2, pTest7));     /* test inactive */
  test_results(0,1,0,0,0,0,0,0,0,1);
  CU_set_fail_on_inactive(CU_TRUE);
  TEST(CUE_SINIT_FAILED == CU_run_test(pSuite2, pTest7));
  test_results(0,1,0,0,0,0,0,0,0,1);
  CU_set_test_active(pTest7, CU_TRUE);

  CU_set_suite_initfunc(pSuite2, NULL);
  TEST(CUE_SUCCESS == CU_run_test(pSuite2, pTest6));          /* change a suite init function */
  CU_set_suite_initfunc(pSuite2, &suite_fail);
  test_results(0,0,0,1,0,0,1,1,0,0);

  CU_set_suite_cleanupfunc(pSuite3, NULL);
  TEST(CUE_SUCCESS == CU_run_test(pSuite3, pTest8));          /* change a suite cleanup function */
  CU_set_suite_cleanupfunc(pSuite3, &suite_fail);
  test_results(0,0,0,1,1,0,1,0,1,1);

  CU_set_test_func(pTest8, &test_succeed);
  TEST(CUE_SCLEAN_FAILED == CU_run_test(pSuite3, pTest8));   /* change a test function */
  CU_set_test_func(pTest8, &test_fail);
  test_results(0,1,0,1,0,0,1,1,0,1);

  /* run each test (CUEA_ABORT) */
  CU_set_error_action(CUEA_ABORT);

  f_exit_called = CU_FALSE;
  TEST(CUE_SUCCESS == CU_run_test(pSuite1, pTest1));
  TEST(CU_FALSE == f_exit_called);
  test_results(0,0,0,1,0,0,1,1,0,0);

  f_exit_called = CU_FALSE;
  TEST(CUE_SUCCESS == CU_run_test(pSuite1, pTest2));
  TEST(CU_FALSE == f_exit_called);
  test_results(0,0,0,1,1,0,1,0,1,1);

  f_exit_called = CU_FALSE;
  TEST(CUE_SUCCESS == CU_run_test(pSuite1, pTest3));
  TEST(CU_FALSE == f_exit_called);
  test_results(0,0,0,1,0,0,1,1,0,0);

  f_exit_called = CU_FALSE;
  TEST(CUE_SUCCESS == CU_run_test(pSuite1, pTest4));
  TEST(CU_FALSE == f_exit_called);
  test_results(0,0,0,1,1,0,1,0,1,1);

  f_exit_called = CU_FALSE;
  TEST(CUE_SUCCESS == CU_run_test(pSuite1, pTest5));
  TEST(CU_FALSE == f_exit_called);
  test_results(0,0,0,1,0,0,1,1,0,0);

  f_exit_called = CU_FALSE;
  TEST(CUE_SINIT_FAILED == CU_run_test(pSuite2, pTest6));
  TEST(CU_TRUE == f_exit_called);
  test_results(0,1,0,0,0,0,0,0,0,1);

  f_exit_called = CU_FALSE;
  TEST(CUE_SINIT_FAILED == CU_run_test(pSuite2, pTest7));
  TEST(CU_TRUE == f_exit_called);
  test_results(0,1,0,0,0,0,0,0,0,1);

  f_exit_called = CU_FALSE;
  TEST(CUE_SCLEAN_FAILED == CU_run_test(pSuite3, pTest8));
  TEST(CU_TRUE == f_exit_called);
  test_results(0,1,0,1,1,0,1,0,1,2);

  f_exit_called = CU_FALSE;
  TEST(CUE_SCLEAN_FAILED == CU_run_test(pSuite3, pTest9));
  TEST(CU_TRUE == f_exit_called);
  test_results(0,1,0,1,0,0,1,1,0,1);

  CU_set_suite_active(pSuite2, CU_FALSE);
  CU_set_fail_on_inactive(CU_FALSE);
  f_exit_called = CU_FALSE;
  TEST(CUE_SUITE_INACTIVE == CU_run_test(pSuite2, pTest6));   /* suite inactive */
  TEST(CU_TRUE == f_exit_called);
  test_results(0,0,1,0,0,0,0,0,0,0);
  CU_set_fail_on_inactive(CU_TRUE);
  f_exit_called = CU_FALSE;
  TEST(CUE_SUITE_INACTIVE == CU_run_test(pSuite2, pTest6));
  TEST(CU_TRUE == f_exit_called);
  test_results(0,0,1,0,0,0,0,0,0,1);
  CU_set_suite_active(pSuite2, CU_TRUE);

  CU_set_test_active(pTest6, CU_FALSE);
  CU_set_fail_on_inactive(CU_FALSE);
  f_exit_called = CU_FALSE;
  TEST(CUE_SINIT_FAILED == CU_run_test(pSuite2, pTest6));     /* test inactive */
  TEST(CU_TRUE == f_exit_called);
  test_results(0,1,0,0,0,0,0,0,0,1);
  CU_set_fail_on_inactive(CU_TRUE);
  f_exit_called = CU_FALSE;
  TEST(CUE_SINIT_FAILED == CU_run_test(pSuite2, pTest6));
  TEST(CU_TRUE == f_exit_called);
  test_results(0,1,0,0,0,0,0,0,0,1);
  CU_set_test_active(pTest6, CU_TRUE);

  f_exit_called = CU_FALSE;
  CU_set_suite_initfunc(pSuite2, NULL);
  TEST(CUE_SUCCESS == CU_run_test(pSuite2, pTest6));          /* change a suite init function */
  CU_set_suite_initfunc(pSuite2, &suite_fail);
  TEST(CU_FALSE == f_exit_called);
  test_results(0,0,0,1,0,0,1,1,0,0);

  f_exit_called = CU_FALSE;
  CU_set_suite_cleanupfunc(pSuite1, &suite_fail);
  TEST(CUE_SCLEAN_FAILED == CU_run_test(pSuite1, pTest1));    /* change a suite cleanup function */
  CU_set_suite_cleanupfunc(pSuite1, NULL);
  TEST(CU_TRUE == f_exit_called);
  test_results(0,1,0,1,0,0,1,1,0,1);

  f_exit_called = CU_FALSE;
  CU_set_test_func(pTest8, &test_succeed);
  TEST(CUE_SCLEAN_FAILED == CU_run_test(pSuite3, pTest8));    /* change a test function */
  CU_set_test_func(pTest8, &test_fail);
  TEST(CU_TRUE == f_exit_called);
  test_results(0,1,0,1,0,0,1,1,0,1);

  /* clean up after testing */
  CU_set_error_action(CUEA_IGNORE);
  CU_cleanup_registry();
}

/*-------------------------------------------------*/
/*  tests CU_assertImplementation()
 *        CU_get_failure_list()
 *        CU_clear_previous_results()
 */
static void test_CU_assertImplementation(void)
{
  CU_Test dummy_test;
  CU_Suite dummy_suite;
  CU_pFailureRecord pFailure1 = NULL;
  CU_pFailureRecord pFailure2 = NULL;
  CU_pFailureRecord pFailure3 = NULL;
  CU_pFailureRecord pFailure4 = NULL;
  CU_pFailureRecord pFailure5 = NULL;
  CU_pFailureRecord pFailure6 = NULL;

  CU_clear_previous_results();

  TEST(NULL == CU_get_failure_list());
  TEST(0 == CU_get_number_of_asserts());
  TEST(0 == CU_get_number_of_failures());
  TEST(0 == CU_get_number_of_failure_records());

  /* fool CU_assertImplementation into thinking test run is in progress */
  f_pCurTest = &dummy_test;
  f_pCurSuite = &dummy_suite;

  /* asserted value is CU_TRUE*/
  TEST(CU_TRUE == CU_assertImplementation(CU_TRUE, 100, "Nothing happened 0.", "dummy0.c", "dummy_func0", CU_FALSE));

  TEST(NULL == CU_get_failure_list());
  TEST(1 == CU_get_number_of_asserts());
  TEST(0 == CU_get_number_of_failures());
  TEST(0 == CU_get_number_of_failure_records());

  TEST(CU_TRUE == CU_assertImplementation(CU_TRUE, 101, "Nothing happened 1.", "dummy1.c", "dummy_func1", CU_FALSE));

  TEST(NULL == CU_get_failure_list());
  TEST(2 == CU_get_number_of_asserts());
  TEST(0 == CU_get_number_of_failures());
  TEST(0 == CU_get_number_of_failure_records());

  /* asserted value is CU_FALSE */
  TEST(CU_FALSE == CU_assertImplementation(CU_FALSE, 102, "Something happened 2.", "dummy2.c", "dummy_func2", CU_FALSE));

  TEST(NULL != CU_get_failure_list());
  TEST(3 == CU_get_number_of_asserts());
  TEST(1 == CU_get_number_of_failures());
  TEST(1 == CU_get_number_of_failure_records());

  TEST(CU_FALSE == CU_assertImplementation(CU_FALSE, 103, "Something happened 3.", "dummy3.c", "dummy_func3", CU_FALSE));

  TEST(NULL != CU_get_failure_list());
  TEST(4 == CU_get_number_of_asserts());
  TEST(2 == CU_get_number_of_failures());
  TEST(2 == CU_get_number_of_failure_records());

  TEST(CU_FALSE == CU_assertImplementation(CU_FALSE, 104, "Something happened 4.", "dummy4.c", "dummy_func4", CU_FALSE));

  TEST(NULL != CU_get_failure_list());
  TEST(5 == CU_get_number_of_asserts());
  TEST(3 == CU_get_number_of_failures());
  TEST(3 == CU_get_number_of_failure_records());

  if (3 == CU_get_number_of_failure_records()) {
    pFailure1 = CU_get_failure_list();
    TEST(102 == pFailure1->uiLineNumber);
    TEST(!strcmp("dummy2.c", pFailure1->strFileName));
    TEST(!strcmp("Something happened 2.", pFailure1->strCondition));
    TEST(&dummy_test == pFailure1->pTest);
    TEST(&dummy_suite == pFailure1->pSuite);
    TEST(NULL != pFailure1->pNext);
    TEST(NULL == pFailure1->pPrev);

    pFailure2 = pFailure1->pNext;
    TEST(103 == pFailure2->uiLineNumber);
    TEST(!strcmp("dummy3.c", pFailure2->strFileName));
    TEST(!strcmp("Something happened 3.", pFailure2->strCondition));
    TEST(&dummy_test == pFailure2->pTest);
    TEST(&dummy_suite == pFailure2->pSuite);
    TEST(NULL != pFailure2->pNext);
    TEST(pFailure1 == pFailure2->pPrev);

    pFailure3 = pFailure2->pNext;
    TEST(104 == pFailure3->uiLineNumber);
    TEST(!strcmp("dummy4.c", pFailure3->strFileName));
    TEST(!strcmp("Something happened 4.", pFailure3->strCondition));
    TEST(&dummy_test == pFailure3->pTest);
    TEST(&dummy_suite == pFailure3->pSuite);
    TEST(NULL == pFailure3->pNext);
    TEST(pFailure2 == pFailure3->pPrev);
  }
  else
    FAIL("Unexpected number of failure records.");

  /* confirm destruction of failure records */
  pFailure4 = pFailure1;
  pFailure5 = pFailure2;
  pFailure6 = pFailure3;
  TEST(0 != test_cunit_get_n_memevents(pFailure4));
  TEST(test_cunit_get_n_allocations(pFailure4) != test_cunit_get_n_deallocations(pFailure4));
  TEST(0 != test_cunit_get_n_memevents(pFailure5));
  TEST(test_cunit_get_n_allocations(pFailure5) != test_cunit_get_n_deallocations(pFailure5));
  TEST(0 != test_cunit_get_n_memevents(pFailure6));
  TEST(test_cunit_get_n_allocations(pFailure6) != test_cunit_get_n_deallocations(pFailure6));

  CU_clear_previous_results();
  TEST(0 != test_cunit_get_n_memevents(pFailure4));
  TEST(test_cunit_get_n_allocations(pFailure4) == test_cunit_get_n_deallocations(pFailure4));
  TEST(0 != test_cunit_get_n_memevents(pFailure5));
  TEST(test_cunit_get_n_allocations(pFailure5) == test_cunit_get_n_deallocations(pFailure5));
  TEST(0 != test_cunit_get_n_memevents(pFailure6));
  TEST(test_cunit_get_n_allocations(pFailure6) == test_cunit_get_n_deallocations(pFailure6));
  TEST(0 == CU_get_number_of_asserts());
  TEST(0 == CU_get_number_of_successes());
  TEST(0 == CU_get_number_of_failures());
  TEST(0 == CU_get_number_of_failure_records());

  f_pCurTest = NULL;
  f_pCurSuite = NULL;
}

/*-------------------------------------------------*/
static void test_add_failure(void)
{
  CU_Test test1;
  CU_Suite suite1;
  CU_pFailureRecord pFailure1 = NULL;
  CU_pFailureRecord pFailure2 = NULL;
  CU_pFailureRecord pFailure3 = NULL;
  CU_pFailureRecord pFailure4 = NULL;
  CU_RunSummary run_summary = {"", 0, 0, 0, 0, 0, 0, 0, 0, 0, 0};

  /* test under memory exhaustion */
  test_cunit_deactivate_malloc();
  add_failure(&pFailure1, &run_summary, CUF_AssertFailed, 100, "condition 0", "file0.c", &suite1, &test1);
  TEST(NULL == pFailure1);
  TEST(0 == run_summary.nFailureRecords);
  test_cunit_activate_malloc();

  /* normal operation */
  add_failure(&pFailure1, &run_summary, CUF_AssertFailed, 101, "condition 1", "file1.c", &suite1, &test1);
  TEST(1 == run_summary.nFailureRecords);
  if (TEST(NULL != pFailure1)) {
    TEST(101 == pFailure1->uiLineNumber);
    TEST(!strcmp("condition 1", pFailure1->strCondition));
    TEST(!strcmp("file1.c", pFailure1->strFileName));
    TEST(&test1 == pFailure1->pTest);
    TEST(&suite1 == pFailure1->pSuite);
    TEST(NULL == pFailure1->pNext);
    TEST(NULL == pFailure1->pPrev);
    TEST(pFailure1 == f_last_failure);
    TEST(0 != test_cunit_get_n_memevents(pFailure1));
    TEST(test_cunit_get_n_allocations(pFailure1) != test_cunit_get_n_deallocations(pFailure1));
  }

  add_failure(&pFailure1, &run_summary, CUF_AssertFailed, 102, "condition 2", "file2.c", NULL, &test1);
  TEST(2 == run_summary.nFailureRecords);
  if (TEST(NULL != pFailure1)) {
    TEST(101 == pFailure1->uiLineNumber);
    TEST(!strcmp("condition 1", pFailure1->strCondition));
    TEST(!strcmp("file1.c", pFailure1->strFileName));
    TEST(&test1 == pFailure1->pTest);
    TEST(&suite1 == pFailure1->pSuite);
    TEST(NULL != pFailure1->pNext);
    TEST(NULL == pFailure1->pPrev);
    TEST(pFailure1 != f_last_failure);
    TEST(0 != test_cunit_get_n_memevents(pFailure1));
    TEST(test_cunit_get_n_allocations(pFailure1) != test_cunit_get_n_deallocations(pFailure1));

    if (TEST(NULL != (pFailure2 = pFailure1->pNext))) {
      TEST(102 == pFailure2->uiLineNumber);
      TEST(!strcmp("condition 2", pFailure2->strCondition));
      TEST(!strcmp("file2.c", pFailure2->strFileName));
      TEST(&test1 == pFailure2->pTest);
      TEST(NULL == pFailure2->pSuite);
      TEST(NULL == pFailure2->pNext);
      TEST(pFailure1 == pFailure2->pPrev);
      TEST(pFailure2 == f_last_failure);
      TEST(0 != test_cunit_get_n_memevents(pFailure2));
      TEST(test_cunit_get_n_allocations(pFailure2) != test_cunit_get_n_deallocations(pFailure2));
    }
  }

  pFailure3 = pFailure1;
  pFailure4 = pFailure2;
  clear_previous_results(&run_summary, &pFailure1);

  TEST(0 == run_summary.nFailureRecords);
  TEST(0 != test_cunit_get_n_memevents(pFailure3));
  TEST(test_cunit_get_n_allocations(pFailure3) == test_cunit_get_n_deallocations(pFailure3));
  TEST(0 != test_cunit_get_n_memevents(pFailure4));
  TEST(test_cunit_get_n_allocations(pFailure4) == test_cunit_get_n_deallocations(pFailure4));
}

/*-------------------------------------------------*/
void test_cunit_TestRun(void)
{
  test_cunit_start_tests("TestRun.c");

  test_message_handlers();
  test_CU_fail_on_inactive();
  test_CU_run_all_tests();
  test_CU_run_suite();
  test_CU_run_test();
  test_CU_assertImplementation();
  test_add_failure();

  test_cunit_end_tests();
}

#endif    /* CUNIT_BUILD_TESTS */
