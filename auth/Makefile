KUBECTL_GLAGS= --ignore-not-found=true
SUBPROJECTS += database

.PHONY: install uninstall subprojs
.DEFAULT: subprojs

install:
	kubectl create -f hydra-hydravolume-persistentvolumeclaim.yaml
	kubectl create -f postgresd-service.yaml
	kubectl create -f postgresd-deployment.yaml
	kubectl create -f hydra-service.yaml
	kubectl create -f hydra-deployment.yaml

uninstall:
	kubectl delete $(KUBECTL_GLAGS) -f hydra-deployment.yaml
	kubectl delete $(KUBECTL_GLAGS) -f hydra-service.yaml
	kubectl delete $(KUBECTL_GLAGS) -f postgresd-deployment.yaml
	kubectl delete $(KUBECTL_GLAGS) -f postgresd-service.yaml
	kubectl delete $(KUBECTL_GLAGS) -f hydra-hydravolume-persistentvolumeclaim.yaml

subprojs:
	@for proj in $(SUBPROJECTS); \
        do \
		$(MAKE) -C $$proj $(MAKECMDGOALS); \
	done
