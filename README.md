# dashboard-metrics-scraper [WIP]

Small binary to scrape and store a small window of metrics from the Metrics Server in Kubernetes.

## Command-Line Arguments
| Flag  | Description  | Default  |
|---|---|---|
| kubeconfig  | The path to the kubeconfig used to connect to the Kubernetes API server and the Kubelets (defaults to in-cluster config)  |  |
| db-file  | What file to use as a SQLite3 database.  |  `/tmp/metrics.db` |
| metric-resolution | The resolution at which dashboard-metrics-scraper will poll metrics.  | `1m` |
| metric-duration | The duration after which metrics are purged from the database. | `15m` |
| namespace | The namespace to use for all metric calls. When provided, skip node metrics. (defaults to cluster level metrics). | |

## Community, discussion, contribution, and support

Learn how to engage with the Kubernetes community on the [community page](http://kubernetes.io/community/).

You can reach the maintainers of this project at:

- [Slack](http://slack.k8s.io/)
- [Mailing List](https://groups.google.com/forum/#!forum/kubernetes-dev)

### Code of conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](code-of-conduct.md).

[owners]: https://git.k8s.io/community/contributors/guide/owners.md
[Creative Commons 4.0]: https://git.k8s.io/website/LICENSE
