## Process this file with automake to produce Makefile.in

COMPILE_DIRS = CUnit doc Man Share

if ENABLE_EXAMPLES
EXAMPLE_COMPILE_DIR = Examples
endif 

COMPILE_DIRS += $(EXAMPLE_COMPILE_DIR)

SUBDIRS = ${COMPILE_DIRS}

distdir=${PACKAGE}-${VERSION}-${RELEASE}

WD=`pwd`
RPM_BUILD_ROOT=${WD}/BUILD
RPM_SOURCE_DIR=${WD}/SOURCES
#RPM_DEFINES=--define "_topdir ${WD}" --define "_prefix ${RPM_BUILD_ROOT}" --define "_unpackaged_files_terminate_build 0"
RPM_DEFINES=--define "_topdir ${WD}" --define "_unpackaged_files_terminate_build 0" --define "_prefix /usr/local"
RPMFLAGS=--buildroot ${RPM_BUILD_ROOT}
ZIP=zip
RPM=rpmbuild

PACKAGE_TMPSPEC=./SPECS/@PACKAGE@-@VERSION@-@RELEASE@.spec
SPECFILE=@PACKAGE@".spec"
CHANGELOG="ChangeLog"
SOURCE_GZIP_FILE=@PACKAGE@-@VERSION@-@RELEASE@.tar.gz
SOURCE_ZIP_FILE=@PACKAGE@-@VERSION@-@RELEASE@.zip
SOURCE_TAR_FILE=@PACKAGE@-@VERSION@-@RELEASE@.tar
TAR_EXCLUDES=--exclude '*/BUILD/*' --exclude '*/SOURCES/*' --exclude '*/RPMS/*' --exclude '*/SRPMS/*' \
	--exclude '*/SPECS/*' --exclude '*.tar' --exclude '*.zip' --exclude '*~' --exclude '*.gz' \
	--exclude '*.rpm' --exclude 'autom4te.cache'
ZIP_EXCLUDES=-x "*/BUILD/*" -x "*/SOURCES/*" -x "*/RPMS/*" -x "*/SRPMS/*" -x "*/SPECS/*" \
	-x "*.tar" -x "*.zip" -x "*~" -x "*.gz" -x "*.rpm" -x "autom4te.cache"

RPM_ATHLON=./RPMS/athlon/@PACKAGE@-@VERSION@-@RELEASE@.athlon.rpm
RPM_I386=./RPMS/i386/@PACKAGE@-@VERSION@-@RELEASE@.i386.rpm
RPM_I486=./RPMS/i486/@PACKAGE@-@VERSION@-@RELEASE@.i486.rpm
RPM_I586=./RPMS/i586/@PACKAGE@-@VERSION@-@RELEASE@.i586.rpm
RPM_I686=./RPMS/i686/@PACKAGE@-@VERSION@-@RELEASE@.i686.rpm
RPM_K6=./RPMS/athlon/@PACKAGE@-@VERSION@-@RELEASE@.k6.rpm
RPM_SRC=./SRPMS/@PACKAGE@-@VERSION@-@RELEASE@.src.rpm

EXTRA_DIST=cunit.pc.in
EXTRA_REMOVE=CUnit-Memory-Dump.lst ${SOURCE_TAR_FILE} ${SOURCE_ZIP_FILE} ${SOURCE_GZIP_FILE} \
	*.rpm *.o *.ncb *.plg *.opt *.clw *.ilk *.exe *.pdb *.lib

pkgconfigdir = $(libdir)/pkgconfig
pkgconfig_DATA = cunit.pc

RPMDIRS=BUILD SOURCES RPMS RPMS/i386 SRPMS SPECS

clean-all: clean
	rm -f ${SOURCE_GZIP_FILE} ${SOURCE_ZIP_FILE} ${SOURCE_TAR_FILE} *.rpm
	rm -rf ${RPMDIRS} `find . -name "Temp"` `find . -name "Debug"`
	for file in ${EXTRA_REMOVE}; do rm -f `find . -name "${file}"`; done	

distclean-all: distclean
	rm -f ${SOURCE_GZIP_FILE} ${SOURCE_ZIP_FILE} ${SOURCE_TAR_FILE} *.rpm
	rm -rf ${RPMDIRS}

rpm: rpm-package

zip: clean
	cd .. && \
	rm -f ${SOURCE_ZIP_FILE} && \
	${ZIP} -rT ${SOURCE_ZIP_FILE} @PACKAGE@-@VERSION@-@RELEASE@/* ${ZIP_EXCLUDES} && \
	cp -f ${SOURCE_ZIP_FILE} ./@PACKAGE@-@VERSION@-@RELEASE@/ && \
	cd @PACKAGE@-@VERSION@-@RELEASE@

tarball: clean
	set -x
	cd .. && \
	rm -f ${SOURCE_TAR_FILE} && \
	tar cvf ${SOURCE_TAR_FILE} ${TAR_EXCLUDES} @PACKAGE@-@VERSION@-@RELEASE@/* && \
	cp -f ${SOURCE_TAR_FILE} ./@PACKAGE@-@VERSION@-@RELEASE@/  && \
	rm -f ${SOURCE_TAR_FILE} && \
	cd @PACKAGE@-@VERSION@-@RELEASE@

gzip: clean tarball
	rm -f ${SOURCE_GZIP_FILE} && \
	gzip ${SOURCE_TAR_FILE}

rpm-package: gzip
	rm -rf ${RPMDIRS} && \
	mkdir -p ${RPMDIRS} && \
	cp ${SPECFILE} ${PACKAGE_TMPSPEC} && \
	echo "%changelog" >> ${PACKAGE_TMPSPEC} && \
	cat ${CHANGELOG} >> ${PACKAGE_TMPSPEC} && \
	export RPM_BUILD_ROOT=${RPM_BUILD_ROOT} && \
	cp ${SOURCE_GZIP_FILE} ${RPM_SOURCE_DIR} 
	${RPM} ${RPM_DEFINES} ${RPMFLAGS} -bb ${PACKAGE_TMPSPEC} && \
	( (cp -f ${RPM_ATHLON} .) || (cp -f ${RPM_I386} .) || (cp -f ${RPM_I486} .) \
	  || (cp -f ${RPM_I586} .) || (cp -f ${RPM_I686} .) \
	  || (cp -f ${RPM_K6} .) || (cp -f ${RPM_NOARCH} .) ) && \
	rm -rf ${RPMDIRS}
