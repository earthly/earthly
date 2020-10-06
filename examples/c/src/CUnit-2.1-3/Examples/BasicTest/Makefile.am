## Process this file with automake to produce Makefile.in

exampledir=${datadir}/@PACKAGE@/Examples/Basic

if ENABLE_EXAMPLES
if ENABLE_BASIC

example_PROGRAMS = BasicTest

BasicTest_SOURCES = BasicTest.c
BasicTest_LDADD = -lcunit -lcunitexamples
BasicTest_LDFLAGS = -L.. -lm
INCLUDES = -I..

if ENABLE_CURSES
BasicTest_LDADD += -l$(CURSES_LIB)
endif

example_DATA = README

endif
endif