#!/bin/bash

# $ curl -O -L https://releases-art-rhcos.svc.ci.openshift.org/art/storage/releases/rhcos-4.7/47.83.202101161239-0/x86_64/rhcos-47.83.202101161239-0-qemu.x86_64.qcow2.gz
# $ cp rhcos-47.83.202101161239-0-qemu.x86_64.qcow2.gz /tmp
# $ sudo gunzip /tmp/rhcos-47.83.202101161239-0-qemu.x86_64.qcow2.gz

IGNITION_CONFIG="/var/lib/libvirt/images/aio.ign"
sudo cp "$1" "${IGNITION_CONFIG}"
sudo chown qemu:qemu "${IGNITION_CONFIG}"
sudo restorecon "${IGNITION_CONFIG}"

#RHCOS_IMAGE="/tmp/rhcos-47.83.202101161239-0-qemu.x86_64.qcow2"
RHCOS_IMAGE="/var/lib/libvirt/images/ushift.qcow2"
VM_NAME="aio-test"
OS_VARIANT="fedora33"
#OS_VARIANT="rhel8.1"
RAM_MB="4096"
DISK_GB="20"

virt-install \
    --connect qemu:///system \
    -n "${VM_NAME}" \
    -r "${RAM_MB}" \
    --os-variant="${OS_VARIANT}" \
    --import \
    --network=network:test-net,mac=52:54:00:ee:42:e1 \
    --graphics=none \
    --disk "size=${DISK_GB},backing_store=${RHCOS_IMAGE}" \
#    --qemu-commandline="-fw_cfg name=opt/com.coreos/config,file=${IGNITION_CONFIG}"
