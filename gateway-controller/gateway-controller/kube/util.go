 /***************************************************************************************
 *    Title: Example: Deploying Cassandra with Stateful Sets
 *    Author: Kubernetes
 *    Date: Sep. 15, 2017
 *    Code version: 1.0
 *    Availability: https://kubernetes.io/docs/tutorials/stateful-application/cassandra/
 *
 ***************************************************************************************/

package kube

func hashString(s string) int {
	h := 0
	for i := 0; i < len(s); i++ {
		h = 31*h + int(s[i])
	}
	return h
}