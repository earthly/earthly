Summary: Build automation tool for the container era
Name: earthly
Version: __earthly_version__
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
