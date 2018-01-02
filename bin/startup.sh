#!/bin/bash

minikube start --vm-driver=vmwarefusion --memory 5120 --cpus=4 -v 10

# "/Applications/VMware Fusion.app/Contents/Library/vmrun" -gu docker -gp tcuser runScriptInGuest /Users/$USER/.minikube/machines/minikube/minikube.vmx /bin/sh "mkdir /home/docker/.ssh; cp /Users/$USER/.minikube/machines/minikube/id_rsa.pub /home/docker/.ssh/authorized_keys; chown -R docker:staff /home/docker/.ssh"
