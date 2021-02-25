#/bin/sh
set -e

mkdir -p rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}

cat > earthly.spec << EOF
Summary: Build automation tool for the container era
Name: earthly
Version: 1.0.0
Release: 1
License: Business Source License
URL: http://earthly.dev
Group: System
Packager: Earthly team
Requires: bash
BuildRoot: /work/rpmbuild/

%description
Build automation tool for the container era

%install
mkdir -p %{buildroot}/usr/bin/
cp /usr/local/bin/earthly %{buildroot}/usr/bin/earthly

%files
/usr/bin/earthly

%changelog
* Thu Feb 25 2021 alex <alex@earthly.dev>
- initial poc

%clean
echo maybe-rm-rf $RPM_BUILD_ROOT

EOF

rpmbuild --target x86_64 -bb earthly.spec
find / | grep -w .rpm$
