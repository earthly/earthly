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
#include "ExampleTests.h"

int main(int argc, char* argv[])
{
  CU_BOOL Run = CU_FALSE ;

  setvbuf(stdout, NULL, _IONBF, 0);

  if (argc > 1) {
    if (!strcmp("-i", argv[1])) {
      Run = CU_TRUE ;
      CU_set_error_action(CUEA_IGNORE);
    }
    else if (!strcmp("-f", argv[1])) {
      Run = CU_TRUE ;
      CU_set_error_action(CUEA_FAIL);
    }
    else if (!strcmp("-A", argv[1])) {
      Run = CU_TRUE ;
      CU_set_error_action(CUEA_ABORT);
    }
    else if (!strcmp("-e", argv[1])) {
      print_example_results();
    }
    else {
      printf("\nUsage:  CursesTest [option]\n\n"
               "        Options: -i  Run, ignoring framework errors [default].\n"
               "                 -f  Run, failing on framework error.\n"
               "                 -A  Run, aborting on framework error.\n"
               "                 -e  Print expected test results and exit.\n"
               "                 -h  Print this message.\n\n");
    }
  }
  else {
    Run = CU_TRUE;
    CU_set_error_action(CUEA_IGNORE);
  }

  if (CU_TRUE == Run) {
    if (CU_initialize_registry()) {
      printf("\nInitialization of Test Registry failed.");
    }
    else {
      AddTests();
      CU_curses_run_tests();
      CU_cleanup_registry();
    }
  }

  return 0;
}

