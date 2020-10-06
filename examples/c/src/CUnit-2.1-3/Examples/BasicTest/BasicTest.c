/*
 *  CUnit - A Unit testing framework library for C.
 *  Copyright (C) 2004  Anil Kumar, Jerry St.Clair
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

#include "Basic.h"
#include "ExampleTests.h"

int main(int argc, char* argv[])
{
  CU_BasicRunMode mode = CU_BRM_VERBOSE;
  CU_ErrorAction error_action = CUEA_IGNORE;
  int i;

  setvbuf(stdout, NULL, _IONBF, 0);

  for (i=1 ; i<argc ; i++) {
    if (!strcmp("-i", argv[i])) {
      error_action = CUEA_IGNORE;
    }
    else if (!strcmp("-f", argv[i])) {
      error_action = CUEA_FAIL;
    }
    else if (!strcmp("-A", argv[i])) {
      error_action = CUEA_ABORT;
    }
    else if (!strcmp("-s", argv[i])) {
      mode = CU_BRM_SILENT;
    }
    else if (!strcmp("-n", argv[i])) {
      mode = CU_BRM_NORMAL;
    }
    else if (!strcmp("-v", argv[i])) {
      mode = CU_BRM_VERBOSE;
    }
    else if (!strcmp("-e", argv[i])) {
      print_example_results();
      return 0;
    }
    else {
      printf("\nUsage:  BasicTest [options]\n\n"
               "Options:   -i   ignore framework errors [default].\n"
               "           -f   fail on framework error.\n"
               "           -A   abort on framework error.\n\n"
               "           -s   silent mode - no output to screen.\n"
               "           -n   normal mode - standard output to screen.\n"
               "           -v   verbose mode - max output to screen [default].\n\n"
               "           -e   print expected test results and exit.\n"
               "           -h   print this message and exit.\n\n");
      return 0;
    }
  }

  if (CU_initialize_registry()) {
    printf("\nInitialization of Test Registry failed.");
  }
  else {
    AddTests();
    CU_basic_set_mode(mode);
    CU_set_error_action(error_action);
    printf("\nTests completed with return value %d.\n", CU_basic_run_tests());
    CU_cleanup_registry();
  }

  return 0;
}
