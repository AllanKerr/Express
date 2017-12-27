# Teardown Command Design

The teardown command is designed to use the same transaction model as the deploy command to maximize code reuse.

## Transactions

The teardown command uses the [`AutoscalerTransaction`](https://github.com/AllanKerr/Services/blob/master/gateway-controller/gateway-controller/kube/transaction_autoscaler.go), [`DeploymentTransaction`](https://github.com/AllanKerr/Services/blob/master/gateway-controller/gateway-controller/kube/transaction_deployment.go), [`IngressTransaction`](https://github.com/AllanKerr/Services/blob/master/gateway-controller/gateway-controller/kube/transaction_ingress.go), and [`ServiceTransaction`](https://github.com/AllanKerr/Services/blob/master/gateway-controller/gateway-controller/kube/transaction_service.go). These transactions implement the [`Transaction interface`](https://github.com/AllanKerr/Services/blob/master/gateway-controller/gateway-controller/kube/transaction.go) to allow transactions to be executed and rolled back. However, only `Rollback` is used by the teardown command.

A transaction is created for each of the four Kubernetes objects that were created by the `deploy` command. These transactions replicate the transactions that were created during the `deploy` command. This allows `Rollback` to be called on each transaction without calling `Execute` because the transactions have already been executed during deploy. As a result, teardown is designed to rollback the transactions executed by the `deploy` command.

The only information required to rebuild the deploy transactions is the deploy name. Because of this, no transaction state needs to be saved.
