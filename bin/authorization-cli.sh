#!/bin/bash
echo "$@"
kubectl exec `kubectl get pod -l app=authorization | tail -n +2 | awk '$1 {print $1}' | head -1` -- app $@
