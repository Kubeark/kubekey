/*
 Copyright 2022 The KubeSphere Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package templates

import (
	"text/template"

	"github.com/lithammer/dedent"
)

var RookOperator = template.Must(template.New("operator.yaml").Parse(
	dedent.Dedent(`
#################################################################################################################
# The deployment for the rook operator
# Contains the common settings for most Kubernetes deployments.
# For example, to create the rook-ceph cluster:
#   kubectl create -f crds.yaml -f common.yaml -f operator.yaml
#   kubectl create -f cluster.yaml
#
# Also see other operator sample files for variations of operator.yaml:
# - operator-openshift.yaml: Common settings for running in OpenShift
###############################################################################################################

# Rook Ceph Operator Config ConfigMap
# Use this ConfigMap to override Rook-Ceph Operator configurations.
# NOTE! Precedence will be given to this config if the same Env Var config also exists in the
#       Operator Deployment.
# To move a configuration(s) from the Operator Deployment to this ConfigMap, add the config
# here. It is recommended to then remove it from the Deployment to eliminate any future confusion.
kind: ConfigMap
apiVersion: v1
metadata:
  name: rook-ceph-operator-config
  # should be in the namespace of the operator
  namespace: rook-ceph # namespace:operator
data:
  # The logging level for the operator: ERROR | WARNING | INFO | DEBUG
  ROOK_LOG_LEVEL: "INFO"

  # Enable the CSI driver.
  # To run the non-default version of the CSI driver, see the override-able image properties in operator.yaml
  ROOK_CSI_ENABLE_CEPHFS: "true"
  # Enable the default version of the CSI RBD driver. To start another version of the CSI driver, see image properties below.
  ROOK_CSI_ENABLE_RBD: "true"
  # Enable the CSI NFS driver. To start another version of the CSI driver, see image properties below.
  ROOK_CSI_ENABLE_NFS: "false"
  ROOK_CSI_ENABLE_GRPC_METRICS: "false"

  # Set to true to enable Ceph CSI pvc encryption support.
  CSI_ENABLE_ENCRYPTION: "false"

  # Set to true to enable host networking for CSI CephFS and RBD nodeplugins. This may be necessary
  # in some network configurations where the SDN does not provide access to an external cluster or
  # there is significant drop in read/write performance.
  # CSI_ENABLE_HOST_NETWORK: "true"

  # Set logging level for cephCSI containers maintained by the cephCSI.
  # Supported values from 0 to 5. 0 for general useful logs, 5 for trace level verbosity.
  # CSI_LOG_LEVEL: "0"

  # Set logging level for Kubernetes-csi sidecar containers.
  # Supported values from 0 to 5. 0 for general useful logs (the default), 5 for trace level verbosity.
  # CSI_SIDECAR_LOG_LEVEL: "0"

  # Set replicas for csi provisioner deployment.
  CSI_PROVISIONER_REPLICAS: "2"

  # OMAP generator will generate the omap mapping between the PV name and the RBD image.
  # CSI_ENABLE_OMAP_GENERATOR need to be enabled when we are using rbd mirroring feature.
  # By default OMAP generator sidecar is deployed with CSI provisioner pod, to disable
  # it set it to false.
  # CSI_ENABLE_OMAP_GENERATOR: "false"

  # set to false to disable deployment of snapshotter container in CephFS provisioner pod.
  CSI_ENABLE_CEPHFS_SNAPSHOTTER: "true"

  # set to false to disable deployment of snapshotter container in RBD provisioner pod.
  CSI_ENABLE_RBD_SNAPSHOTTER: "true"

  # Enable cephfs kernel driver instead of ceph-fuse.
  # If you disable the kernel client, your application may be disrupted during upgrade.
  # See the upgrade guide: https://rook.io/docs/rook/latest/ceph-upgrade.html
  # NOTE! cephfs quota is not supported in kernel version < 4.17
  CSI_FORCE_CEPHFS_KERNEL_CLIENT: "true"

  # (Optional) policy for modifying a volume's ownership or permissions when the RBD PVC is being mounted.
  # supported values are documented at https://kubernetes-csi.github.io/docs/support-fsgroup.html
  CSI_RBD_FSGROUPPOLICY: "ReadWriteOnceWithFSType"

  # (Optional) policy for modifying a volume's ownership or permissions when the CephFS PVC is being mounted.
  # supported values are documented at https://kubernetes-csi.github.io/docs/support-fsgroup.html
  CSI_CEPHFS_FSGROUPPOLICY: "ReadWriteOnceWithFSType"

  # (Optional) policy for modifying a volume's ownership or permissions when the NFS PVC is being mounted.
  # supported values are documented at https://kubernetes-csi.github.io/docs/support-fsgroup.html
  CSI_NFS_FSGROUPPOLICY: "ReadWriteOnceWithFSType"

  # (Optional) Allow starting unsupported ceph-csi image
  ROOK_CSI_ALLOW_UNSUPPORTED_VERSION: "false"

  # (Optional) control the host mount of /etc/selinux for csi plugin pods.
  CSI_PLUGIN_ENABLE_SELINUX_HOST_MOUNT: "false"

  # The default version of CSI supported by Rook will be started. To change the version
  # of the CSI driver to something other than what is officially supported, change
  # these images to the desired release of the CSI driver.
  # ROOK_CSI_CEPH_IMAGE: "quay.io/cephcsi/cephcsi:v3.6.2"
  # ROOK_CSI_REGISTRAR_IMAGE: "registry.k8s.io/sig-storage/csi-node-driver-registrar:v2.5.1"
  # ROOK_CSI_RESIZER_IMAGE: "registry.k8s.io/sig-storage/csi-resizer:v1.4.0"
  # ROOK_CSI_PROVISIONER_IMAGE: "registry.k8s.io/sig-storage/csi-provisioner:v3.1.0"
  # ROOK_CSI_SNAPSHOTTER_IMAGE: "registry.k8s.io/sig-storage/csi-snapshotter:v6.0.1"
  # ROOK_CSI_ATTACHER_IMAGE: "registry.k8s.io/sig-storage/csi-attacher:v3.4.0"
  # ROOK_CSI_NFS_IMAGE: "registry.k8s.io/sig-storage/nfsplugin:v4.0.0"

  # (Optional) set user created priorityclassName for csi plugin pods.
  CSI_PLUGIN_PRIORITY_CLASSNAME: "system-node-critical"

  # (Optional) set user created priorityclassName for csi provisioner pods.
  CSI_PROVISIONER_PRIORITY_CLASSNAME: "system-cluster-critical"

  # CSI CephFS plugin daemonset update strategy, supported values are OnDelete and RollingUpdate.
  # Default value is RollingUpdate.
  # CSI_CEPHFS_PLUGIN_UPDATE_STRATEGY: "OnDelete"
  # CSI RBD plugin daemonset update strategy, supported values are OnDelete and RollingUpdate.
  # Default value is RollingUpdate.
  # CSI_RBD_PLUGIN_UPDATE_STRATEGY: "OnDelete"

  # CSI NFS plugin daemonset update strategy, supported values are OnDelete and RollingUpdate.
  # Default value is RollingUpdate.
  # CSI_NFS_PLUGIN_UPDATE_STRATEGY: "OnDelete"

  # kubelet directory path, if kubelet configured to use other than /var/lib/kubelet path.
  # ROOK_CSI_KUBELET_DIR_PATH: "/var/lib/kubelet"

  # Labels to add to the CSI CephFS Deployments and DaemonSets Pods.
  # ROOK_CSI_CEPHFS_POD_LABELS: "key1=value1,key2=value2"
  # Labels to add to the CSI RBD Deployments and DaemonSets Pods.
  # ROOK_CSI_RBD_POD_LABELS: "key1=value1,key2=value2"
  # Labels to add to the CSI NFS Deployments and DaemonSets Pods.
  # ROOK_CSI_NFS_POD_LABELS: "key1=value1,key2=value2"

  # (Optional) CephCSI provisioner NodeAffinity (applied to both CephFS and RBD provisioner).
  # CSI_PROVISIONER_NODE_AFFINITY: "role=storage-node; storage=rook, ceph"
  # (Optional) CephCSI provisioner tolerations list(applied to both CephFS and RBD provisioner).
  # Put here list of taints you want to tolerate in YAML format.
  # CSI provisioner would be best to start on the same nodes as other ceph daemons.
  # CSI_PROVISIONER_TOLERATIONS: |
  #   - effect: NoSchedule
  #     key: node-role.kubernetes.io/control-plane
  #     operator: Exists
  #   - effect: NoExecute
  #     key: node-role.kubernetes.io/etcd
  #     operator: Exists
  # (Optional) CephCSI plugin NodeAffinity (applied to both CephFS and RBD plugin).
  # CSI_PLUGIN_NODE_AFFINITY: "role=storage-node; storage=rook, ceph"
  # (Optional) CephCSI plugin tolerations list(applied to both CephFS and RBD plugin).
  # Put here list of taints you want to tolerate in YAML format.
  # CSI plugins need to be started on all the nodes where the clients need to mount the storage.
  # CSI_PLUGIN_TOLERATIONS: |
  #   - effect: NoSchedule
  #     key: node-role.kubernetes.io/control-plane
  #     operator: Exists
  #   - effect: NoExecute
  #     key: node-role.kubernetes.io/etcd
  #     operator: Exists

  # (Optional) CephCSI RBD provisioner NodeAffinity (if specified, overrides CSI_PROVISIONER_NODE_AFFINITY).
  # CSI_RBD_PROVISIONER_NODE_AFFINITY: "role=rbd-node"
  # (Optional) CephCSI RBD provisioner tolerations list(if specified, overrides CSI_PROVISIONER_TOLERATIONS).
  # Put here list of taints you want to tolerate in YAML format.
  # CSI provisioner would be best to start on the same nodes as other ceph daemons.
  # CSI_RBD_PROVISIONER_TOLERATIONS: |
  #   - key: node.rook.io/rbd
  #     operator: Exists
  # (Optional) CephCSI RBD plugin NodeAffinity (if specified, overrides CSI_PLUGIN_NODE_AFFINITY).
  # CSI_RBD_PLUGIN_NODE_AFFINITY: "role=rbd-node"
  # (Optional) CephCSI RBD plugin tolerations list(if specified, overrides CSI_PLUGIN_TOLERATIONS).
  # Put here list of taints you want to tolerate in YAML format.
  # CSI plugins need to be started on all the nodes where the clients need to mount the storage.
  # CSI_RBD_PLUGIN_TOLERATIONS: |
  #   - key: node.rook.io/rbd
  #     operator: Exists

  # (Optional) CephCSI CephFS provisioner NodeAffinity (if specified, overrides CSI_PROVISIONER_NODE_AFFINITY).
  # CSI_CEPHFS_PROVISIONER_NODE_AFFINITY: "role=cephfs-node"
  # (Optional) CephCSI CephFS provisioner tolerations list(if specified, overrides CSI_PROVISIONER_TOLERATIONS).
  # Put here list of taints you want to tolerate in YAML format.
  # CSI provisioner would be best to start on the same nodes as other ceph daemons.
  # CSI_CEPHFS_PROVISIONER_TOLERATIONS: |
  #   - key: node.rook.io/cephfs
  #     operator: Exists
  # (Optional) CephCSI CephFS plugin NodeAffinity (if specified, overrides CSI_PLUGIN_NODE_AFFINITY).
  # CSI_CEPHFS_PLUGIN_NODE_AFFINITY: "role=cephfs-node"
  # NOTE: Support for defining NodeAffinity for operators other than "In" and "Exists" requires the user to input a
  # valid v1.NodeAffinity JSON or YAML string. For example, the following is valid YAML v1.NodeAffinity:
  # CSI_CEPHFS_PLUGIN_NODE_AFFINITY: |
  #   requiredDuringSchedulingIgnoredDuringExecution:
  #     nodeSelectorTerms:
  #       - matchExpressions:
  #         - key: myKey
  #           operator: DoesNotExist
  # (Optional) CephCSI CephFS plugin tolerations list(if specified, overrides CSI_PLUGIN_TOLERATIONS).
  # Put here list of taints you want to tolerate in YAML format.
  # CSI plugins need to be started on all the nodes where the clients need to mount the storage.
  # CSI_CEPHFS_PLUGIN_TOLERATIONS: |
  #   - key: node.rook.io/cephfs
  #     operator: Exists

  # (Optional) CephCSI NFS provisioner NodeAffinity (overrides CSI_PROVISIONER_NODE_AFFINITY).
  # CSI_NFS_PROVISIONER_NODE_AFFINITY: "role=nfs-node"
  # (Optional) CephCSI NFS provisioner tolerations list (overrides CSI_PROVISIONER_TOLERATIONS).
  # Put here list of taints you want to tolerate in YAML format.
  # CSI provisioner would be best to start on the same nodes as other ceph daemons.
  # CSI_NFS_PROVISIONER_TOLERATIONS: |
  #   - key: node.rook.io/nfs
  #     operator: Exists
  # (Optional) CephCSI NFS plugin NodeAffinity (overrides CSI_PLUGIN_NODE_AFFINITY).
  # CSI_NFS_PLUGIN_NODE_AFFINITY: "role=nfs-node"
  # (Optional) CephCSI NFS plugin tolerations list (overrides CSI_PLUGIN_TOLERATIONS).
  # Put here list of taints you want to tolerate in YAML format.
  # CSI plugins need to be started on all the nodes where the clients need to mount the storage.
  # CSI_NFS_PLUGIN_TOLERATIONS: |
  #   - key: node.rook.io/nfs
  #     operator: Exists

  # (Optional) CEPH CSI RBD provisioner resource requirement list, Put here list of resource
  # requests and limits you want to apply for provisioner pod
  #CSI_RBD_PROVISIONER_RESOURCE: |
  #  - name : csi-provisioner
  #    resource:
  #      requests:
  #        memory: 128Mi
  #        cpu: 100m
  #      limits:
  #        memory: 256Mi
  #        cpu: 200m
  #  - name : csi-resizer
  #    resource:
  #      requests:
  #        memory: 128Mi
  #        cpu: 100m
  #      limits:
  #        memory: 256Mi
  #        cpu: 200m
  #  - name : csi-attacher
  #    resource:
  #      requests:
  #        memory: 128Mi
  #        cpu: 100m
  #      limits:
  #        memory: 256Mi
  #        cpu: 200m
  #  - name : csi-snapshotter
  #    resource:
  #      requests:
  #        memory: 128Mi
  #        cpu: 100m
  #      limits:
  #        memory: 256Mi
  #        cpu: 200m
  #  - name : csi-rbdplugin
  #    resource:
  #      requests:
  #        memory: 512Mi
  #        cpu: 250m
  #      limits:
  #        memory: 1Gi
  #        cpu: 500m
  #  - name : csi-omap-generator
  #    resource:
  #      requests:
  #        memory: 512Mi
  #        cpu: 250m
  #      limits:
  #        memory: 1Gi
  #        cpu: 500m
  #  - name : liveness-prometheus
  #    resource:
  #      requests:
  #        memory: 128Mi
  #        cpu: 50m
  #      limits:
  #        memory: 256Mi
  #        cpu: 100m
  # (Optional) CEPH CSI RBD plugin resource requirement list, Put here list of resource
  # requests and limits you want to apply for plugin pod
  #CSI_RBD_PLUGIN_RESOURCE: |
  #  - name : driver-registrar
  #    resource:
  #      requests:
  #        memory: 128Mi
  #        cpu: 50m
  #      limits:
  #        memory: 256Mi
  #        cpu: 100m
  #  - name : csi-rbdplugin
  #    resource:
  #      requests:
  #        memory: 512Mi
  #        cpu: 250m
  #      limits:
  #        memory: 1Gi
  #        cpu: 500m
  #  - name : liveness-prometheus
  #    resource:
  #      requests:
  #        memory: 128Mi
  #        cpu: 50m
  #      limits:
  #        memory: 256Mi
  #        cpu: 100m
  # (Optional) CEPH CSI CephFS provisioner resource requirement list, Put here list of resource
  # requests and limits you want to apply for provisioner pod
  #CSI_CEPHFS_PROVISIONER_RESOURCE: |
  #  - name : csi-provisioner
  #    resource:
  #      requests:
  #        memory: 128Mi
  #        cpu: 100m
  #      limits:
  #        memory: 256Mi
  #        cpu: 200m
  #  - name : csi-resizer
  #    resource:
  #      requests:
  #        memory: 128Mi
  #        cpu: 100m
  #      limits:
  #        memory: 256Mi
  #        cpu: 200m
  #  - name : csi-attacher
  #    resource:
  #      requests:
  #        memory: 128Mi
  #        cpu: 100m
  #      limits:
  #        memory: 256Mi
  #        cpu: 200m
  #  - name : csi-snapshotter
  #    resource:
  #      requests:
  #        memory: 128Mi
  #        cpu: 100m
  #      limits:
  #        memory: 256Mi
  #        cpu: 200m
  #  - name : csi-cephfsplugin
  #    resource:
  #      requests:
  #        memory: 512Mi
  #        cpu: 250m
  #      limits:
  #        memory: 1Gi
  #        cpu: 500m
  #  - name : liveness-prometheus
  #    resource:
  #      requests:
  #        memory: 128Mi
  #        cpu: 50m
  #      limits:
  #        memory: 256Mi
  #        cpu: 100m
  # (Optional) CEPH CSI CephFS plugin resource requirement list, Put here list of resource
  # requests and limits you want to apply for plugin pod
  #CSI_CEPHFS_PLUGIN_RESOURCE: |
  #  - name : driver-registrar
  #    resource:
  #      requests:
  #        memory: 128Mi
  #        cpu: 50m
  #      limits:
  #        memory: 256Mi
  #        cpu: 100m
  #  - name : csi-cephfsplugin
  #    resource:
  #      requests:
  #        memory: 512Mi
  #        cpu: 250m
  #      limits:
  #        memory: 1Gi
  #        cpu: 500m
  #  - name : liveness-prometheus
  #    resource:
  #      requests:
  #        memory: 128Mi
  #        cpu: 50m
  #      limits:
  #        memory: 256Mi
  #        cpu: 100m

  # (Optional) CEPH CSI NFS provisioner resource requirement list, Put here list of resource
  # requests and limits you want to apply for provisioner pod
  # CSI_NFS_PROVISIONER_RESOURCE: |
  #  - name : csi-provisioner
  #    resource:
  #      requests:
  #        memory: 128Mi
  #        cpu: 100m
  #      limits:
  #        memory: 256Mi
  #        cpu: 200m
  #  - name : csi-nfsplugin
  #    resource:
  #      requests:
  #        memory: 512Mi
  #        cpu: 250m
  #      limits:
  #        memory: 1Gi
  #        cpu: 500m
  # (Optional) CEPH CSI NFS plugin resource requirement list, Put here list of resource
  # requests and limits you want to apply for plugin pod
  # CSI_NFS_PLUGIN_RESOURCE: |
  #  - name : driver-registrar
  #    resource:
  #      requests:
  #        memory: 128Mi
  #        cpu: 50m
  #      limits:
  #        memory: 256Mi
  #        cpu: 100m
  #  - name : csi-nfsplugin
  #    resource:
  #      requests:
  #        memory: 512Mi
  #        cpu: 250m
  #      limits:
  #        memory: 1Gi
  #        cpu: 500m

  # Configure CSI Ceph FS grpc and liveness metrics port
  # Set to true to enable Ceph CSI liveness container.
  CSI_ENABLE_LIVENESS: "false"
  # CSI_CEPHFS_GRPC_METRICS_PORT: "9091"
  # CSI_CEPHFS_LIVENESS_METRICS_PORT: "9081"
  # Configure CSI RBD grpc and liveness metrics port
  # CSI_RBD_GRPC_METRICS_PORT: "9090"
  # CSI_RBD_LIVENESS_METRICS_PORT: "9080"
  # CSIADDONS_PORT: "9070"

  # Whether the OBC provisioner should watch on the operator namespace or not, if not the namespace of the cluster will be used
  ROOK_OBC_WATCH_OPERATOR_NAMESPACE: "true"

  # Whether to start the discovery daemon to watch for raw storage devices on nodes in the cluster.
  # This daemon does not need to run if you are only going to create your OSDs based on StorageClassDeviceSets with PVCs.
  ROOK_ENABLE_DISCOVERY_DAEMON: "false"
  # The timeout value (in seconds) of Ceph commands. It should be >= 1. If this variable is not set or is an invalid value, it's default to 15.
  ROOK_CEPH_COMMANDS_TIMEOUT_SECONDS: "15"
  # Enable the volume replication controller.
  # Before enabling, ensure the Volume Replication CRDs are created.
  # See https://rook.io/docs/rook/latest/ceph-csi-drivers.html#rbd-mirroring
  CSI_ENABLE_VOLUME_REPLICATION: "false"
  # CSI_VOLUME_REPLICATION_IMAGE: "quay.io/csiaddons/volumereplication-operator:v0.3.0"
  # Enable the csi addons sidecar.
  CSI_ENABLE_CSIADDONS: "false"
  # ROOK_CSIADDONS_IMAGE: "quay.io/csiaddons/k8s-sidecar:v0.4.0"
  # The CSI GRPC timeout value (in seconds). It should be >= 120. If this variable is not set or is an invalid value, it's default to 150.
  CSI_GRPC_TIMEOUT_SECONDS: "150"
---
# OLM: BEGIN OPERATOR DEPLOYMENT
apiVersion: apps/v1
kind: Deployment
metadata:
  name: rook-ceph-operator
  namespace: rook-ceph # namespace:operator
  labels:
    operator: rook
    storage-backend: ceph
    app.kubernetes.io/name: rook-ceph
    app.kubernetes.io/instance: rook-ceph
    app.kubernetes.io/component: rook-ceph-operator
    app.kubernetes.io/part-of: rook-ceph-operator
spec:
  selector:
    matchLabels:
      app: rook-ceph-operator
  strategy:
    type: Recreate
  replicas: 1
  template:
    metadata:
      labels:
        app: rook-ceph-operator
    spec:
      serviceAccountName: rook-ceph-system
      containers:
        - name: rook-ceph-operator
          image: rook/ceph:v1.9.9
          args: ["ceph", "operator"]
          securityContext:
            runAsNonRoot: true
            runAsUser: 2016
            runAsGroup: 2016
          volumeMounts:
            - mountPath: /var/lib/rook
              name: rook-config
            - mountPath: /etc/ceph
              name: default-config-dir
            - mountPath: /etc/webhook
              name: webhook-cert
          ports:
            - containerPort: 9443
              name: https-webhook
              protocol: TCP
          env:
            # If the operator should only watch for cluster CRDs in the same namespace, set this to "true".
            # If this is not set to true, the operator will watch for cluster CRDs in all namespaces.
            - name: ROOK_CURRENT_NAMESPACE_ONLY
              value: "false"
            # Rook Discover toleration. Will tolerate all taints with all keys.
            # Choose between NoSchedule, PreferNoSchedule and NoExecute:
            # - name: DISCOVER_TOLERATION
            #   value: "NoSchedule"
            # (Optional) Rook Discover toleration key. Set this to the key of the taint you want to tolerate
            # - name: DISCOVER_TOLERATION_KEY
            #   value: "<KeyOfTheTaintToTolerate>"
            # (Optional) Rook Discover tolerations list. Put here list of taints you want to tolerate in YAML format.
            # - name: DISCOVER_TOLERATIONS
            #   value: |
            #     - effect: NoSchedule
            #       key: node-role.kubernetes.io/control-plane
            #       operator: Exists
            #     - effect: NoExecute
            #       key: node-role.kubernetes.io/etcd
            #       operator: Exists
            # (Optional) Rook Discover priority class name to set on the pod(s)
            # - name: DISCOVER_PRIORITY_CLASS_NAME
            #   value: "<PriorityClassName>"
            # (Optional) Discover Agent NodeAffinity.
            # - name: DISCOVER_AGENT_NODE_AFFINITY
            #   value: "role=storage-node; storage=rook, ceph"
            # (Optional) Discover Agent Pod Labels.
            # - name: DISCOVER_AGENT_POD_LABELS
            #   value: "key1=value1,key2=value2"

            # The duration between discovering devices in the rook-discover daemonset.
            - name: ROOK_DISCOVER_DEVICES_INTERVAL
              value: "60m"

            # Whether to start pods as privileged that mount a host path, which includes the Ceph mon and osd pods.
            # Set this to true if SELinux is enabled (e.g. OpenShift) to workaround the anyuid issues.
            # For more details see https://github.com/rook/rook/issues/1314#issuecomment-355799641
            - name: ROOK_HOSTPATH_REQUIRES_PRIVILEGED
              value: "false"

            # In some situations SELinux relabelling breaks (times out) on large filesystems, and doesn't work with cephfs ReadWriteMany volumes (last relabel wins).
            # Disable it here if you have similar issues.
            # For more details see https://github.com/rook/rook/issues/2417
            - name: ROOK_ENABLE_SELINUX_RELABELING
              value: "true"

            # In large volumes it will take some time to chown all the files. Disable it here if you have performance issues.
            # For more details see https://github.com/rook/rook/issues/2254
            - name: ROOK_ENABLE_FSGROUP
              value: "true"

            # Disable automatic orchestration when new devices are discovered
            - name: ROOK_DISABLE_DEVICE_HOTPLUG
              value: "false"

            # Provide customised regex as the values using comma. For eg. regex for rbd based volume, value will be like "(?i)rbd[0-9]+".
            # In case of more than one regex, use comma to separate between them.
            # Default regex will be "(?i)dm-[0-9]+,(?i)rbd[0-9]+,(?i)nbd[0-9]+"
            # Add regex expression after putting a comma to blacklist a disk
            # If value is empty, the default regex will be used.
            - name: DISCOVER_DAEMON_UDEV_BLACKLIST
              value: "(?i)dm-[0-9]+,(?i)rbd[0-9]+,(?i)nbd[0-9]+"

            # Time to wait until the node controller will move Rook pods to other
            # nodes after detecting an unreachable node.
            # Pods affected by this setting are:
            # mgr, rbd, mds, rgw, nfs, PVC based mons and osds, and ceph toolbox
            # The value used in this variable replaces the default value of 300 secs
            # added automatically by k8s as Toleration for
            # <node.kubernetes.io/unreachable>
            # The total amount of time to reschedule Rook pods in healthy nodes
            # before detecting a <not ready node> condition will be the sum of:
            #  --> node-monitor-grace-period: 40 seconds (k8s kube-controller-manager flag)
            #  --> ROOK_UNREACHABLE_NODE_TOLERATION_SECONDS: 5 seconds
            - name: ROOK_UNREACHABLE_NODE_TOLERATION_SECONDS
              value: "5"

            - name: ROOK_DISABLE_ADMISSION_CONTROLLER
              value: "false"

            # The name of the node to pass with the downward API
            - name: NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            # The pod name to pass with the downward API
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            # The pod namespace to pass with the downward API
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
          # Recommended resource requests and limits, if desired
          #resources:
          #  limits:
          #    cpu: 500m
          #    memory: 512Mi
          #  requests:
          #    cpu: 100m
          #    memory: 128Mi

          #  Uncomment it to run lib bucket provisioner in multithreaded mode
          #- name: LIB_BUCKET_PROVISIONER_THREADS
          #  value: "5"

      # Uncomment it to run rook operator on the host network
      #hostNetwork: true
      volumes:
        - name: rook-config
          emptyDir: {}
        - name: default-config-dir
          emptyDir: {}
        - name: webhook-cert
          emptyDir: {}
# OLM: END OPERATOR DEPLOYMENT
  `)))
