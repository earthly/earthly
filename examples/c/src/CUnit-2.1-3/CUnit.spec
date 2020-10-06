Summary:  A unit testing framework for 'C'
Name:     CUnit
Version:  2.1
Release:  3
Source:   http://www.sourceforge.net/projects/cunit/CUnit-2.1-3.tar.gz
Group:    Development/Tools
License:  GPL
URL:      http://cunit.sourceforge.net
Packager: Jerry St. Clair <jds2@users.sourceforge.net>

%description
CUnit is a unit testing framework for C.
This package installs the CUnit static library,
headers, and documentation files.

%prep
echo "Preparing for Installation."
%setup -q -n CUnit-2.1-3

%build
echo "Preparing for Building."
./configure --prefix=%{_prefix} --enable-automated --enable-basic --enable-console --enable-curses --enable-examples --enable-test && \
make

%install
echo "Preparing for Make install."
make DESTDIR=$RPM_BUILD_ROOT install

%clean

%files
%defattr(-,root,root)

########### Include Files
%{_prefix}/include/CUnit/Automated.h
%{_prefix}/include/CUnit/Basic.h
%{_prefix}/include/CUnit/Console.h
%{_prefix}/include/CUnit/CUError.h
%{_prefix}/include/CUnit/CUnit.h
%{_prefix}/include/CUnit/CUnit_intl.h
%{_prefix}/include/CUnit/CUCurses.h
%{_prefix}/include/CUnit/MyMem.h
%{_prefix}/include/CUnit/TestDB.h
%{_prefix}/include/CUnit/TestRun.h
%{_prefix}/include/CUnit/Util.h

########## Library Files
%{_prefix}/lib/libcunit.a
%{_prefix}/lib/libcunit.so.1.0.1

########## doc Files
%{_prefix}/doc/CUnit/CUnit_doc.css
%{_prefix}/doc/CUnit/error_handling.html
%{_prefix}/doc/CUnit/fdl.html
%{_prefix}/doc/CUnit/index.html
%{_prefix}/doc/CUnit/introduction.html
%{_prefix}/doc/CUnit/managing_tests.html
%{_prefix}/doc/CUnit/running_tests.html
%{_prefix}/doc/CUnit/test_registry.html
%{_prefix}/doc/CUnit/writing_tests.html
%{_prefix}/doc/CUnit/headers/Automated.h
%{_prefix}/doc/CUnit/headers/Basic.h
%{_prefix}/doc/CUnit/headers/Console.h
%{_prefix}/doc/CUnit/headers/CUError.h
%{_prefix}/doc/CUnit/headers/CUnit.h
%{_prefix}/doc/CUnit/headers/CUnit_intl.h
%{_prefix}/doc/CUnit/headers/CUCurses.h
%{_prefix}/doc/CUnit/headers/MyMem.h
%{_prefix}/doc/CUnit/headers/TestDB.h
%{_prefix}/doc/CUnit/headers/TestRun.h
%{_prefix}/doc/CUnit/headers/Util.h
%{_prefix}/doc/CUnit/headers/Win.h

########## Manpage Files
%{_prefix}/man/man3/CUnit.3*

########## Share information and Example Files
%{_prefix}/share/CUnit/Examples/Automated/README
%{_prefix}/share/CUnit/Examples/Automated/AutomatedTest
%{_prefix}/share/CUnit/Examples/Basic/README
%{_prefix}/share/CUnit/Examples/Basic/BasicTest
%{_prefix}/share/CUnit/Examples/Console/README
%{_prefix}/share/CUnit/Examples/Console/ConsoleTest
%{_prefix}/share/CUnit/Examples/Curses/README
%{_prefix}/share/CUnit/Examples/Curses/CursesTest
%{_prefix}/share/CUnit/Test/test_cunit
%{_prefix}/share/CUnit/CUnit-List.dtd
%{_prefix}/share/CUnit/CUnit-List.xsl
%{_prefix}/share/CUnit/CUnit-Run.dtd
%{_prefix}/share/CUnit/CUnit-Run.xsl
%{_prefix}/share/CUnit/Memory-Dump.dtd
%{_prefix}/share/CUnit/Memory-Dump.xsl

# Add the change log in ChangeLog file located under source home directory.
# The same file is used internally to populate the change log for the RPM creation.
