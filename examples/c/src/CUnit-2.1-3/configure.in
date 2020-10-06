dnl Applied patch to handle AC_ARG_ENABLE targets (Patch 1250604) - Provided by Andrew Kephart <andrew.kephart@alereon.com

AC_INIT(CUnit, 2.1-3)
AM_INIT_AUTOMAKE(CUnit, 2.1-3)

dnl Package version information
PACKAGE=CUnit
VERSION_MAJOR=2
VERSION_MINOR=1
VERSION=${VERSION_MAJOR}.${VERSION_MINOR}
RELEASE=3

dnl Shared library version information
dnl
dnl Current  = interface version (increment when change user interface to lib)
dnl Revision = revision number (increment when change implementation, reset to 0 for new interface)
dnl Age      = number of previous interfaces this one is compatible with (i.e. is a superset of)
dnl
dnl Suffix1  = [Current - Age] (on Linux, anyway, soname suffix is [C-A].A.R)
dnl
LIBTOOL_CURRENT=1
LIBTOOL_REVISION=1
LIBTOOL_AGE=0
LIBTOOL_SUFFIX1=1
LIBTOOL_SUFFIX=${LIBTOOL_SUFFIX1}.${LIBTOOL_AGE}.${LIBTOOL_REVISION}

TOPDIR=`pwd`

CFLAGS="$CFLAGS -DRELEASE=@RELEASE@ -Wall -W -pedantic -Wshadow -ansi -I${PWD}/CUnit/Headers -std=c99"
LDFLAGS="$LDFLAGS -L${PWD}/CUnit/Sources"

AC_ARG_ENABLE(debug,
  [AS_HELP_STRING([--enable-debug],[enable debugging [default=no]])],
  [cu_do_debug=$enableval],
  [cu_do_debug="no"])
if test x"$cu_do_debug" = xyes ; then
	echo "++++++++++ Enabling debug mode compilation."
	CFLAGS="$CFLAGS -D_DEBUG -g3 -Wall"
	LDFLAGS="$LDFLAGS -g"
	ENABLE_DEBUG="TRUE"
else
	echo "---------- Disabling debug mode compilation."
	ENABLE_DEBUG="FALSE"
fi


AC_ARG_ENABLE(automated,
  [AS_HELP_STRING([--enable-automated],[compile CUnit automated interface [default=yes]])],
  [cu_do_automated=$enableval],
  [cu_do_automated="yes"])
if test x"$cu_do_automated" = xyes ; then
	echo "++++++++++ Enabling automated(XML) mode compilation"
	BUILD_AUTOMATED="TRUE"
else
	echo "---------- Disabling automated(XML) mode compilation"
	BUILD_AUTOMATED="FALSE"
fi
AM_CONDITIONAL(ENABLE_AUTOMATED, test x"$cu_do_automated" = xyes)


AC_ARG_ENABLE(basic,
  [AS_HELP_STRING([--enable-basic],[compile CUnit basic interface [default=yes]])],
  [cu_do_basic=$enableval],
  [cu_do_basic="yes"])
if test x"$cu_do_basic" = xyes ; then
	echo "++++++++++ Enabling basic mode compilation"
	BUILD_BASIC="TRUE"
else
	echo "---------- Disabling basic mode compilation"
	BUILD_BASIC="FALSE"
fi
AM_CONDITIONAL(ENABLE_BASIC, test x"$cu_do_basic" = xyes)


AC_ARG_ENABLE(console,
  [AS_HELP_STRING([--enable-console],[compile CUnit console interface [default=yes]])],
  [cu_do_console=$enableval],
  [cu_do_console="yes"])
if test x"$cu_do_console" = xyes ; then
	echo "++++++++++ Enabling console mode compilation"
	BUILD_CONSOLE="TRUE"
else
	echo "---------- Disabling console mode compilation"
	BUILD_CONSOLE="FALSE"
fi
AM_CONDITIONAL(ENABLE_CONSOLE, test x"$cu_do_console" = xyes)


AC_ARG_ENABLE(curses,
  [AS_HELP_STRING([--enable-curses],[compile CUnit curses interface [default=no]])],
  [cu_do_curses=$enableval],
  [cu_do_curses="no"])
if test x"$cu_do_curses" = xyes ; then
	echo "++++++++++ Enabling curses mode compilation"
else
	echo "---------- Disabling curses mode compilation"
	BUILD_CURSES="FALSE"
fi


AC_ARG_ENABLE(examples,
  [AS_HELP_STRING([--enable-examples],[compile CUnit example programs [default=no]])],
  [cu_do_examples=$enableval],
  [cu_do_examples="no"])
if test x"$cu_do_examples" = xyes ; then
	echo "++++++++++ Enabling examples mode compilation"
	BUILD_EXAMPLES="TRUE"
else
	echo "---------- Disabling examples mode compilation"
	BUILD_EXAMPLES="FALSE"
fi
AM_CONDITIONAL(ENABLE_EXAMPLES, test x"$cu_do_examples" = xyes)


AC_ARG_ENABLE(test,
  [AS_HELP_STRING([--enable-test],[compile CUnit internal test program [default=no]])],
  [cu_do_test=$enableval],
  [cu_do_test="no"])
if test x"$cu_do_test" = xyes ; then
	echo "++++++++++ Enabling test mode compilation"
	BUILD_TEST="TRUE"
else
	echo "---------- Disabling test mode compilation"
	BUILD_TEST="FALSE"
fi
AM_CONDITIONAL(ENABLE_TEST, test x"$cu_do_test" = xyes)


AC_ARG_ENABLE(memtrace,
  [AS_HELP_STRING([--enable-memtrace],[enable CUnit internal memory tracking [default=no]])],
  [cu_do_memtrace=$enableval],
  [cu_do_memtrace="no"])
if  test x"$cu_do_memtrace" = xyes ; then
	echo "++++++++++ Enabling memtrace functionality at compile time"
	CFLAGS="$CFLAGS -DMEMTRACE"
	ENABLE_MEMTRACE="TRUE"
else
	echo "---------- Disabling memtrace functionality at compile time"
	ENABLE_MEMTRACE="FALSE"
fi


AC_ARG_ENABLE(deprecated,
  [AS_HELP_STRING([--enable-deprecated],[enable use of deprecated v1.1 names [default=no]])],
  [cu_do_deprecated=$enableval],
  [cu_do_deprecated="no"])
if test x"$cu_do_deprecated" = xyes ; then
	echo "++++++++++ Enabling use of deprecated v1.1 names"
	CFLAGS="$CFLAGS -DUSE_DEPRECATED_CUNIT_NAMES"
	ENABLE_DEPRECATED="TRUE"
else
	echo "---------- Disabling use of deprecated v1.1 names"
	ENABLE_DEPRECATED="FALSE"
fi


dnl Check for programs.
AC_PROG_CC
AC_PROG_INSTALL
AC_PROG_MAKE_SET
AC_PROG_LIBTOOL

dnl Check for libraries
AC_CHECK_LIB(c, main)

dnl TODO: We should provide a --with-curses=PREFIX option to allow user to point to curses lib
if test x"$cu_do_curses" = xyes ; then
	AC_CHECK_LIB([ncurses], [initscr], [cu_have_curses="ncurses"], [cu_have_curses="no"])
	if test x"$cu_have_curses" = xno ; then
		AC_CHECK_LIB([curses], [initscr], [cu_have_curses="curses"], [cu_have_curses="no"])
	fi
	if test x"$cu_have_curses" = xno ; then
		echo "---------- Disabling curses mode compilation (lib not found)"
		cu_do_curses="no"
		BUILD_CURSES=FALSE
	else
		BUILD_CURSES=TRUE
	fi
	CURSES_LIB=$cu_have_curses
fi
AM_CONDITIONAL(ENABLE_CURSES, test x"$cu_do_curses" = xyes)

dnl Check for standard functions
AC_CHECK_FUNC(malloc, [], [echo malloc not found; exit 1])
AC_CHECK_FUNC(free, [], [echo free not found; exit 1])
AC_CHECK_FUNC(calloc, [], [echo calloc not found; exit 1])
AC_CHECK_FUNC(realloc, [], [echo realloc not found; exit 1])

AC_CHECK_FUNC(strcpy, [], [echo strcpy not found; exit 1])
AC_CHECK_FUNC(strerror, [], [echo strerror not found; exit 1])

AC_CHECK_FUNC(fopen, [], [echo fopen not found; exit 1])
AC_CHECK_FUNC(fclose, [], [echo fclose not found; exit 1])
AC_CHECK_FUNC(fprintf, [], [echo fprintf not found; exit 1])
AC_CHECK_FUNC(snprintf, [], [echo snprintf not found; exit 1])
AC_CHECK_FUNC(setvbuf, [], [echo setvbuf not found; exit 1])

AC_CHECK_FUNC(time, [], [echo time not found; exit 1])
AC_CHECK_FUNC(ctime, [], [echo ctime not found; exit 1])

dnl Check for header files.
AC_HEADER_STDC
AC_CHECK_HEADERS(assert.h)
AC_CHECK_HEADERS(ctype.h)
AC_CHECK_HEADERS(math.h)
AC_CHECK_HEADERS(file.h)
AC_CHECK_HEADERS(stdio.h)
AC_CHECK_HEADERS(stdlib.h)
AC_CHECK_HEADERS(string.h)

AM_CONFIG_HEADER(config.h)

dnl Check for typedefs, structures, and compiler characteristics.
AC_C_CONST

dnl Set environment variables.
AC_SUBST(PACKAGE)
AC_SUBST(VERSION)
AC_SUBST(VERSION_MAJOR)
AC_SUBST(VERSION_MINOR)
AC_SUBST(RELEASE)
AC_SUBST(LIBTOOL_CURRENT)
AC_SUBST(LIBTOOL_REVISION)
AC_SUBST(LIBTOOL_AGE)
AC_SUBST(LIBTOOL_SUFFIX1)
AC_SUBST(LIBTOOL_SUFFIX)
AC_SUBST(ENABLE_DEBUG)
AC_SUBST(ENABLE_MEMTRACE)
AC_SUBST(ENABLE_DEPRECATED)
AC_SUBST(BUILD_AUTOMATED)
AC_SUBST(BUILD_BASIC)
AC_SUBST(BUILD_CONSOLE)
AC_SUBST(BUILD_CURSES)
AC_SUBST(CURSES_LIB)
AC_SUBST(BUILD_EXAMPLES)
AC_SUBST(BUILD_TEST)

dnl Configure Jamrules for user environment
AC_CONFIG_FILES([Jamrules:Jamrules.in])

dnl Configure Makefile set.
AC_CONFIG_FILES( Makefile \
		CUnit.spec \
		cunit.pc \
		CUnit/Makefile \
		CUnit/Headers/Makefile \
		CUnit/Headers/CUnit.h \
		CUnit/Sources/Makefile \
		CUnit/Sources/Framework/Makefile \
		CUnit/Sources/Automated/Makefile \
		CUnit/Sources/Basic/Makefile \
		CUnit/Sources/Console/Makefile \
		CUnit/Sources/Curses/Makefile \
		CUnit/Sources/Test/Makefile \
		doc/Makefile \
		doc/headers/Makefile \
		Examples/Makefile \
		Examples/AutomatedTest/Makefile \
		Examples/BasicTest/Makefile \
		Examples/ConsoleTest/Makefile \
		Examples/CursesTest/Makefile \
		Man/Makefile \
		Man/man3/Makefile \
		Share/Makefile )

AC_OUTPUT

