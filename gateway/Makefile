APP=nginx
KUBECTL_GLAGS= --ignore-not-found=true

.PHONY: install uninstall

install:
	kubectl create -f config.yaml
	kubectl create -f service.yaml
	kubectl create -f deployment.yaml

uninstall:
	kubectl delete $(KUBECTL_GLAGS) service -l "app=$(APP)"
	kubectl delete $(KUBECTL_GLAGS) deploy -l "app=$(APP)"
	kubectl delete $(KUBECTL_GLAGS) configmap -l "app=$(APP)"	
