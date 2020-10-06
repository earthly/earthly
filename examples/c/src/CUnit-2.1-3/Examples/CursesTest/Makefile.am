## Process this file with automake to produce Makefile.in

exampledir=${datadir}/@PACKAGE@/Examples/Curses

if ENABLE_EXAMPLES
if ENABLE_CURSES

example_PROGRAMS = CursesTest

CursesTest_SOURCES = CursesTest.c
CursesTest_LDADD = -lcunit -lcunitexamples
CursesTest_LDFLAGS = -L.. -lm
INCLUDES = -I..

if ENABLE_CURSES
CursesTest_LDADD += -l$(CURSES_LIB)
endif
example_DATA = README

endif
endif
