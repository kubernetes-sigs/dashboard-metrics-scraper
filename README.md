# dashboard-metrics-scraper [WIP]

Small binary to scrape and store a small window of metrics from the Metrics Server in Kubernetes.

## Command-Line Arguments
| Flag  | Description  | Default  |
|---|---|---|
| kubeconfig  | Absolute path to the kubeconfig file  | `~/.kube`  |
| db-file  | What file to use as a SQLite3 database.  |  `:memory:` |
| refresh-interval | Frequency (in seconds) to update the metrics database.  | `5` |
| max-window | Window of time you wish to retain records (in minutes). | `15` |

## Community, discussion, contribution, and support

Learn how to engage with the Kubernetes community on the [community page](http://kubernetes.io/community/).

You can reach the maintainers of this project at:

- [Slack](http://slack.k8s.io/)
- [Mailing List](https://groups.google.com/forum/#!forum/kubernetes-dev)

### Code of conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](code-of-conduct.md).

[owners]: https://git.k8s.io/community/contributors/guide/owners.md
[Creative Commons 4.0]: https://git.k8s.io/website/LICENSE
