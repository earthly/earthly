## Process this file with automake to produce Makefile.in

EXAMPLE_DIRS = . 

if ENABLE_BASIC
BASIC_EXAMPLE_DIRS = BasicTest
endif

if ENABLE_CONSOLE
CONSOLE_EXAMPLE_DIRS = ConsoleTest
endif

if ENABLE_CURSES
CURSES_EXAMPLE_DIRS = CursesTest
endif

if ENABLE_AUTOMATED
AUTOMATED_EXAMPLE_DIRS = AutomatedTest
endif

EXAMPLE_DIRS += $(BASIC_EXAMPLE_DIRS) $(CONSOLE_EXAMPLE_DIRS) $(CURSES_EXAMPLE_DIRS) $(AUTOMATED_EXAMPLE_DIRS)

SUBDIRS = ${EXAMPLE_DIRS}

if ENABLE_EXAMPLES
noinst_LIBRARIES = libcunitexamples.a
libcunitexamples_a_SOURCES = ExampleTests.c
endif

