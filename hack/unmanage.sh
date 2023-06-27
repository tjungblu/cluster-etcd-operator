#!/bin/bash

set +x

oc patch clusterversion version --type json -p "$(cat unmanaged_patch.yaml)"
oc scale deployment -n openshift-etcd-operator --replicas=0 etcd-operator
