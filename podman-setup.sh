#!/bin/sh
set -e

apk add --no-cache --update podman

cat > /etc/containers/containers.conf <<EOF
[containers]
netns="host"
userns="host"
ipcns="host"
utsns="host"
cgroupns="host"
cgroups="disabled"
log_driver = "k8s-file"
[engine]
cgroup_manager = "cgroupfs"
events_logger="file"
runtime="crun"
EOF

chmod 644 /etc/containers/containers.conf
sed -i -e 's|^#mount_program|mount_program|g' -e '/additionalimage.*/a "/var/lib/shared",' -e 's|^mountopt[[:space:]]*=.*$|mountopt = "nodev,fsync=0"|g' /etc/containers/storage.conf

mkdir -p /var/lib/shared/overlay-images
touch /var/lib/shared/overlay-images/images.lock

mkdir -p /var/lib/shared/overlay-layers
touch /var/lib/shared/overlay-layers/layers.lock

mkdir -p /var/lib/shared/vfs-images
touch /var/lib/shared/vfs-images/images.lock

mkdir -p /var/lib/shared/vfs-layers
touch /var/lib/shared/vfs-layers/layers.lock

sed -i 's/\/var\/lib\/containers\/storage/$EARTHLY_DOCKERD_DATA_ROOT/g' /etc/containers/storage.conf
