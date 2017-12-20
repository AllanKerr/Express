KUBECTL_GLAGS= --ignore-not-found=true

.PHONY: install uninstall

install:
	kubectl create -f ingress.yaml

uninstall:
	kubectl delete -f ingress.yaml
