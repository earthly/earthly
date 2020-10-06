## Process this file with automake to produce Makefile.in

exampledir=${datadir}/@PACKAGE@/Examples/Automated

if ENABLE_EXAMPLES
if ENABLE_AUTOMATED

example_PROGRAMS = AutomatedTest

AutomatedTest_SOURCES = AutomatedTest.c
AutomatedTest_LDADD = -lcunit -lcunitexamples
AutomatedTest_LDFLAGS = -L.. -lm
INCLUDES = -I..

if ENABLE_CURSES
AutomatedTest_LDADD += -l$(CURSES_LIB)
endif

example_DATA = README

endif
endif

