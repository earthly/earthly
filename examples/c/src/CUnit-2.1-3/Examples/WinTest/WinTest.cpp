// WinTest.cpp : Defines the entry point for the application. 
//

#include "stdafx.h"

extern "C" {
  #include "Win.h"
}

int APIENTRY WinMain(HINSTANCE hInstance,
                     HINSTANCE hPrevInstance,
                     LPSTR     lpCmdLine,
                     int       nCmdShow)
{
  CU_UNREFERENCED_PARAMETER(hInstance);
  CU_UNREFERENCED_PARAMETER(hPrevInstance);
  CU_UNREFERENCED_PARAMETER(lpCmdLine);
  CU_UNREFERENCED_PARAMETER(nCmdShow);

  // TODO: Place code here.
  CU_win_run_tests();
  return 0;
}



