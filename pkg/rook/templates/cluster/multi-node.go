/*
 Copyright 2021 The KubeSphere Authors.

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

var RookCluster = template.Must(template.New("multi-node.yaml").Parse(
	dedent.Dedent(`
apiVersion: ceph.rook.io/v1
kind: CephCluster
metadata:
  name: rook-ceph
  namespace: rook-ceph
spec:
  cephVersion:
    # The container image used to launch the Ceph daemon pods (mon, mgr, osd, mds, rgw).
    # v14 is nautilus, v15 is octopus, and v16 is pacific.
    # RECOMMENDATION: In production, use a specific version tag instead of the general v14 flag, which pulls the latest release and could result in different
    # versions running within the cluster. See tags available at https://hub.docker.com/r/ceph/ceph/tags/.
    # If you want to be more precise, you can always use a timestamp tag such quay.io/ceph/ceph:v16.2.6-20210918
    # This tag might not contain a new Ceph version, just security fixes from the underlying operating system, which will reduce vulnerabilities
    image: quay.io/ceph/ceph:v16.2.6
    # Whether to allow unsupported versions of Ceph. Currently nautilus, octopus, and pacific are supported.
    # Future versions such as pacific would require this to be set to true.
    # Do not set to true in production.
    allowUnsupported: false
  # The path on the host where configuration files will be persisted. Must be specified.
  # Important: if you reinstall the cluster, make sure you delete this directory from each host or else the mons will fail to start on the new cluster.
  # In Minikube, the '/data' directory is configured to persist across reboots. Use "/data/rook" in Minikube environment.
  dataDirHostPath: /var/lib/rook
  # Whether or not upgrade should continue even if a check fails
  # This means Ceph's status could be degraded and we don't recommend upgrading but you might decide otherwise
  # Use at your OWN risk
  # To understand Rook's upgrade process of Ceph, read https://rook.io/docs/rook/master/ceph-upgrade.html#ceph-version-upgrades
  skipUpgradeChecks: false
  # Whether or not continue if PGs are not clean during an upgrade
  continueUpgradeAfterChecksEvenIfNotHealthy: false
  # WaitTimeoutForHealthyOSDInMinutes defines the time (in minutes) the operator would wait before an OSD can be stopped for upgrade or restart.
  # If the timeout exceeds and OSD is not ok to stop, then the operator would skip upgrade for the current OSD and proceed with the next one
  # if continueUpgradeAfterChecksEvenIfNotHealthy is false. If continueUpgradeAfterChecksEvenIfNotHealthy is true, then opertor would
  # continue with the upgrade of an OSD even if its not ok to stop after the timeout. This timeout won't be applied if skipUpgradeChecks is true.
  # The default wait timeout is 10 minutes.
  waitTimeoutForHealthyOSDInMinutes: 10
  mon:
    # Set the number of mons to be started. Generally recommended to be 3.
    # For highest availability, an odd number of mons should be specified.
    count: 7
    # The mons should be on unique nodes. For production, at least 3 nodes are recommended for this reason.
    # Mons should only be allowed on the same node for test environments where data loss is acceptable.
    allowMultiplePerNode: false
  mgr:
    # When higher availability of the mgr is needed, increase the count to 2.
    # In that case, one mgr will be active and one in standby. When Ceph updates which
    # mgr is active, Rook will update the mgr services to match the active mgr.
    count: 2
    modules:
      # Several monitoring should not need to be included in this list. The "dashboard" and "monitoring" monitoring
      # are already enabled by other settings in the cluster CR.
      - name: pg_autoscaler
        enabled: true
  # enable the ceph dashboard for viewing cluster status
  dashboard:
    enabled: true
    # serve the dashboard under a subpath (useful when you are accessing the dashboard via a reverse proxy)
    # urlPrefix: /ceph-dashboard
    # serve the dashboard at the given port.
    # port: 8443
    # serve the dashboard using SSL
    ssl: true
  # enable prometheus alerting for cluster
  monitoring:
    # requires Prometheus to be pre-installed
    enabled: false
    # namespace to deploy prometheusRule in. If empty, namespace of the cluster will be used.
    # Recommended:
    # If you have a single rook-ceph cluster, set the rulesNamespace to the same namespace as the cluster or keep it empty.
    # If you have multiple rook-ceph clusters in the same k8s cluster, choose the same namespace (ideally, namespace with prometheus
    # deployed) to set rulesNamespace for all the clusters. Otherwise, you will get duplicate alerts with multiple alert definitions.
    rulesNamespace: rook-ceph
  network:
    # enable host networking
    #provider: host
    # enable the Multus network provider
    #provider: multus
    #selectors:
      # The selector keys are required to be public and cluster.
      # Based on the configuration, the operator will do the following:
      #   1. if only the public selector key is specified both public_network and cluster_network Ceph settings will listen on that interface
      #   2. if both public and cluster selector keys are specified the first one will point to 'public_network' flag and the second one to 'cluster_network'
      #
      # In order to work, each selector value must match a NetworkAttachmentDefinition object in Multus
      #
      #public: public-conf --> NetworkAttachmentDefinition object name in Multus
      #cluster: cluster-conf --> NetworkAttachmentDefinition object name in Multus
    # Provide internet protocol version. IPv6, IPv4 or empty string are valid options. Empty string would mean IPv4
    #ipFamily: "IPv6"
    # Ceph daemons to listen on both IPv4 and Ipv6 networks
    #dualStack: false
  # enable the crash collector for ceph daemon crash collection
  crashCollector:
    disable: false
    # Uncomment daysToRetain to prune ceph crash entries older than the
    # specified number of days.
    #daysToRetain: 30
  # enable log collector, daemons will log on files and rotate
  # logCollector:
  #   enabled: true
  #   periodicity: 24h # SUFFIX may be 'h' for hours or 'd' for days.
  # automate [data cleanup process](https://github.com/rook/rook/blob/master/Documentation/ceph-teardown.md#delete-the-data-on-hosts) in cluster destruction.
  cleanupPolicy:
    # Since cluster cleanup is destructive to data, confirmation is required.
    # To destroy all Rook data on hosts during uninstall, confirmation must be set to "yes-really-destroy-data".
    # This value should only be set when the cluster is about to be deleted. After the confirmation is set,
    # Rook will immediately stop configuring the cluster and only wait for the delete command.
    # If the empty string is set, Rook will not destroy any data on hosts during uninstall.
    confirmation: ""
    # sanitizeDisks represents settings for sanitizing OSD disks on cluster deletion
    sanitizeDisks:
      # method indicates if the entire disk should be sanitized or simply ceph's metadata
      # in both case, re-install is possible
      # possible choices are 'complete' or 'quick' (default)
      method: quick
      # dataSource indicate where to get random bytes from to write on the disk
      # possible choices are 'zero' (default) or 'random'
      # using random sources will consume entropy from the system and will take much more time then the zero source
      dataSource: zero
      # iteration overwrite N times instead of the default (1)
      # takes an integer value
      iteration: 1
    # allowUninstallWithVolumes defines how the uninstall should be performed
    # If set to true, cephCluster deletion does not wait for the PVs to be deleted.
    allowUninstallWithVolumes: false
  # To control where various services will be scheduled by kubernetes, use the placement configuration sections below.
  # The example under 'all' would have all services scheduled on kubernetes nodes labeled with 'role=storage-node' and
  # tolerate taints with a key of 'storage-node'.
#  placement:
#    all:
#      nodeAffinity:
#        requiredDuringSchedulingIgnoredDuringExecution:
#          nodeSelectorTerms:
#          - matchExpressions:
#            - key: role
#              operator: In
#              values:
#              - storage-node
#      podAffinity:
#      podAntiAffinity:
#      topologySpreadConstraints:
#      tolerations:
#      - key: storage-node
#        operator: Exists
# The above placement information can also be specified for mon, osd, and mgr components
#    mon:
# Monitor deployments may contain an anti-affinity rule for avoiding monitor
# collocation on the same node. This is a required rule when host network is used
# or when AllowMultiplePerNode is false. Otherwise this anti-affinity rule is a
# preferred rule with weight: 50.
#    osd:
#    mgr:
#    cleanup:
  annotations:
#    all:
#    mon:
#    osd:
#    cleanup:
#    prepareosd:
# If no mgr annotations are set, prometheus scrape annotations will be set by default.
#    mgr:
  labels:
#    all:
#    mon:
#    osd:
#    cleanup:
#    mgr:
#    prepareosd:
# monitoring is a list of key-value pairs. It is injected into all the monitoring resources created by operator.
# These labels can be passed as LabelSelector to Prometheus
#    monitoring:
#    crashcollector:
  resources:
# The requests and limits set here, allow the mgr pod to use half of one CPU core and 1 gigabyte of memory
#    mgr:
#      limits:
#        cpu: "500m"
#        memory: "1024Mi"
#      requests:
#        cpu: "500m"
#        memory: "1024Mi"
# The above example requests/limits can also be added to the other components
#    mon:
#    osd:
# For OSD it also is a possible to specify requests/limits based on device class
#    osd-hdd:
#    osd-ssd:
#    osd-nvme:
#    prepareosd:
#    mgr-sidecar:
#    crashcollector:
#    logcollector:
#    cleanup:
  # The option to automatically remove OSDs that are out and are safe to destroy.
  removeOSDsIfOutAndSafeToRemove: false
#  priorityClassNames:
#    all: rook-ceph-default-priority-class
#    mon: rook-ceph-mon-priority-class
#    osd: rook-ceph-osd-priority-class
#    mgr: rook-ceph-mgr-priority-class
  storage: # cluster level storage configuration and selection
    useAllNodes: true
    useAllDevices: true
    #deviceFilter:
    config:
      # crushRoot: "custom-root" # specify a non-default root label for the CRUSH map
      # metadataDevice: "md0" # specify a non-rotational storage so ceph-volume will use it as block db device of bluestore.
      # databaseSizeMB: "1024" # uncomment if the disks are smaller than 100 GB
      # journalSizeMB: "1024"  # uncomment if the disks are 20 GB or smaller
      # osdsPerDevice: "1" # this value can be overridden at the node or device level
      # encryptedDevice: "true" # the default value for this option is "false"
# Individual nodes and their config can be specified as well, but 'useAllNodes' above must be set to false. Then, only the named
# nodes below will be used as storage resources.  Each node's 'name' field should match their 'kubernetes.io/hostname' label.
    # nodes:
    #   - name: "172.17.4.201"
    #     devices: # specific devices to use for storage can be specified for each node
    #       - name: "sdb"
    #       - name: "nvme01" # multiple osds can be created on high performance devices
    #         config:
    #           osdsPerDevice: "5"
    #       - name: "/dev/disk/by-id/ata-ST4000DM004-XXXX" # devices can be specified using full udev paths
    #     config: # configuration can be specified at the node level which overrides the cluster level config
    #   - name: "172.17.4.301"
    #     deviceFilter: "^sd."
    # when onlyApplyOSDPlacement is false, will merge both placement.All() and placement.osd
    onlyApplyOSDPlacement: false
  # The section for configuring management of daemon disruptions during upgrade or fencing.
  disruptionManagement:
    # If true, the operator will create and manage PodDisruptionBudgets for OSD, Mon, RGW, and MDS daemons. OSD PDBs are managed dynamically
    # via the strategy outlined in the [design](https://github.com/rook/rook/blob/master/design/ceph/ceph-managed-disruptionbudgets.md). The operator will
    # block eviction of OSDs by default and unblock them safely when drains are detected.
    managePodBudgets: true
    # A duration in minutes that determines how long an entire failureDomain like region/zone/host will be held in noout (in addition to the
    # default DOWN/OUT interval) when it is draining. This is only relevant when  managePodBudgets is true. The default value is 30 minutes.
    osdMaintenanceTimeout: 30
    # A duration in minutes that the operator will wait for the placement groups to become healthy (active+clean) after a drain was completed and OSDs came back up.
    # Operator will continue with the next drain if the timeout exceeds. It only works if managePodBudgets is true.
    # No values or 0 means that the operator will wait until the placement groups are healthy before unblocking the next drain.
    pgHealthCheckTimeout: 0
    # If true, the operator will create and manage MachineDisruptionBudgets to ensure OSDs are only fenced when the cluster is healthy.
    # Only available on OpenShift.
    manageMachineDisruptionBudgets: false
    # Namespace in which to watch for the MachineDisruptionBudgets.
    machineDisruptionBudgetNamespace: openshift-machine-api

  # healthChecks
  # Valid values for daemons are 'mon', 'osd', 'status'
  healthCheck:
    daemonHealth:
      mon:
        disabled: false
        interval: 45s
      osd:
        disabled: false
        interval: 60s
      status:
        disabled: false
        interval: 60s
    # Change pod liveness probe, it works for all mon,mgr,osd daemons
    livenessProbe:
      mon:
        disabled: false
      mgr:
        disabled: false
      osd:
        disabled: false
`)))
