# cert-scanner

## SSL Expiration Watcher

### Overview

Creates a Kubernetes [cronjob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) resource in the monitoring namespace that runs on a schedule and checks all SSL certificates for hostnames listed in a [hosts.txt](go/hosts.txt) file, to ensure that all have expiration dates at least [daysWarning](charts/cert-scanner/values.yaml) days in the future.

Any certificates with fewer than [daysWarning](charts/cert-scanner/values.yaml) days remaining will cause the job to log detailed errors in JSON format.

### Usage

Edit a list of plain text hostnames in the [hosts.txt](go/hosts.txt) file on the main branch to control which external hosts are being checked during each scheduled pass.

### Runtime Kubernetes Details

The past several job executions are visible as pods in the monitoring namespace, which end up either in a Completed or an Error state, depending on whether any SSL expiration warnings were found. An alarm can then watch for a steady stream of Completed job events coming from the monitoring namespace, and be triggered if they stop recurring for any reason.

### Helm Configuration Details

See [charts/cert-scanner/values.yaml](charts/cert-scanner/values.yaml) for all exposed helm chart values and their current defaults, as described below.

| Parameter| Description | Default |
|---:|:---|:---:|
|schedule | cron schedule for recurring checks | 0 * * * * |
| successfulJobsHistoryLimit | how many recent pods in the Completed state to retain | 5 |
| failedJobsHistoryLimit | how many recent pods in Error state to retain | 5 |
| daysWarning| how many days warning for SSL expirations | 21 |
