KUBECTL_GLAGS= --ignore-not-found=true

.PHONY: install uninstall

install:
	kubectl create -f service.yaml
	kubectl create -f statefulset.yaml

uninstall:
	kubectl delete $(KUBECTL_GLAGS) -f statefulset.yaml
	kubectl delete $(KUBECTL_GLAGS) -f service.yaml
