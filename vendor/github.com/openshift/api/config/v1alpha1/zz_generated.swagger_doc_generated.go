package v1alpha1

// This file contains a collection of methods that can be used from go-restful to
// generate Swagger API documentation for its models. Please read this PR for more
// information on the implementation: https://github.com/emicklei/go-restful/pull/215
//
// TODOs are ignored from the parser (e.g. TODO(andronat):... || TODO:...) if and only if
// they are on one line! For multiple line or blocks that you want to ignore use ---.
// Any context after a --- is ignored.
//
// Those methods can be generated by using hack/update-swagger-docs.sh

// AUTO-GENERATED FUNCTIONS START HERE
var map_Backup = map[string]string{
	"":         "\n\nBackup provides configuration for performing backups of the openshift cluster.\n\nCompatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.",
	"metadata": "metadata is the standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata",
	"spec":     "spec holds user settable values for configuration",
	"status":   "status holds observed values from the cluster. They may not be overridden.",
}

func (Backup) SwaggerDoc() map[string]string {
	return map_Backup
}

var map_BackupList = map[string]string{
	"":         "BackupList is a collection of items\n\nCompatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.",
	"metadata": "metadata is the standard list's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata",
}

func (BackupList) SwaggerDoc() map[string]string {
	return map_BackupList
}

var map_EtcdBackupSpec = map[string]string{
	"":                "EtcdBackupSpec provides configuration for automated etcd backups to the cluster-etcd-operator",
	"schedule":        "Schedule defines the recurring backup schedule in Cron format every 2 hours: 0 */2 * * * every day at 3am: 0 3 * * * Setting to an empty string \"\" means disabling scheduled backups Default: \"\" See if the upstream CronJob has CEL validation already defined somewhere.",
	"timeZone":        "The time zone name for the given schedule, see https://en.wikipedia.org/wiki/List_of_tz_database_time_zones. If not specified, this will default to the time zone of the kube-controller-manager process. See https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/#time-zones",
	"retentionPolicy": "RetentionPolicy defines the retention policy for retaining and deleting existing backups.",
	"pvcName":         "PVCName specifies the name of the PersistentVolumeClaim which binds a PersistentVolume where the etcd backup files would be saved",
}

func (EtcdBackupSpec) SwaggerDoc() map[string]string {
	return map_EtcdBackupSpec
}

var map_RetentionCountConfig = map[string]string{
	"":                   "RetentionCountConfig specifies the configuration of the retention policy on the number of backups",
	"maxNumberOfBackups": "MaxNumberOfBackups defines the maximum number of backups to retain. If the number of successful backups matches retentionCount the oldest backup will be removed before a new backup is initiated. The count here is for the total number of backups",
}

func (RetentionCountConfig) SwaggerDoc() map[string]string {
	return map_RetentionCountConfig
}

var map_RetentionPolicy = map[string]string{
	"":               "RetentionPolicy defines the retention policy for retaining and deleting existing backups. This struct is a discriminated union that allows users to select the type of retention policy from the supported types.",
	"retentionType":  "RetentionType sets the type of retention policy. The currently supported and valid values are \"retentionCount\" Currently, the only valid policies are retention by count (RetentionCount) and by size (RetentionSize). More policies or types may be added in the future. which would need to be explicitly set as the discriminant or can we set that as the default value.",
	"retentionCount": "RetentionCount configures the retention policy based on the number of backups",
	"retentionSize":  "RetentionSize configures the retention policy based on the size of backups",
}

func (RetentionPolicy) SwaggerDoc() map[string]string {
	return map_RetentionPolicy
}

var map_RetentionSizeConfig = map[string]string{
	"":                   "RetentionSizeConfig specifies the configuration of the retention policy on the total size of backups",
	"maxSizeOfBackupsMb": "MaxSizeOfBackupsMb defines the total size in Mb of backups to retain. If the current total size backups exceeds MaxSizeOfBackupsMb then the oldest backup will be removed before a new backup is initiated.",
}

func (RetentionSizeConfig) SwaggerDoc() map[string]string {
	return map_RetentionSizeConfig
}

var map_GatherConfig = map[string]string{
	"":                  "gatherConfig provides data gathering configuration options.",
	"dataPolicy":        "dataPolicy allows user to enable additional global obfuscation of the IP addresses and base domain in the Insights archive data. Valid values are \"None\" and \"ObfuscateNetworking\". When set to None the data is not obfuscated. When set to ObfuscateNetworking the IP addresses and the cluster domain name are obfuscated. When omitted, this means no opinion and the platform is left to choose a reasonable default, which is subject to change over time. The current default is None.",
	"disabledGatherers": "disabledGatherers is a list of gatherers to be excluded from the gathering. All the gatherers can be disabled by providing \"all\" value. If all the gatherers are disabled, the Insights operator does not gather any data. The particular gatherers IDs can be found at https://github.com/openshift/insights-operator/blob/master/docs/gathered-data.md. Run the following command to get the names of last active gatherers: \"oc get insightsoperators.operator.openshift.io cluster -o json | jq '.status.gatherStatus.gatherers[].name'\" An example of disabling gatherers looks like this: `disabledGatherers: [\"clusterconfig/machine_configs\", \"workloads/workload_info\"]`",
}

func (GatherConfig) SwaggerDoc() map[string]string {
	return map_GatherConfig
}

var map_InsightsDataGather = map[string]string{
	"":         "\n\nInsightsDataGather provides data gather configuration options for the the Insights Operator.\n\nCompatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.",
	"metadata": "metadata is the standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata",
	"spec":     "spec holds user settable values for configuration",
	"status":   "status holds observed values from the cluster. They may not be overridden.",
}

func (InsightsDataGather) SwaggerDoc() map[string]string {
	return map_InsightsDataGather
}

var map_InsightsDataGatherList = map[string]string{
	"":         "InsightsDataGatherList is a collection of items\n\nCompatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.",
	"metadata": "metadata is the standard list's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata",
}

func (InsightsDataGatherList) SwaggerDoc() map[string]string {
	return map_InsightsDataGatherList
}

var map_InsightsDataGatherSpec = map[string]string{
	"gatherConfig": "gatherConfig spec attribute includes all the configuration options related to gathering of the Insights data and its uploading to the ingress.",
}

func (InsightsDataGatherSpec) SwaggerDoc() map[string]string {
	return map_InsightsDataGatherSpec
}

// AUTO-GENERATED FUNCTIONS END HERE
