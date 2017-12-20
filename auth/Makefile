KUBECTL_GLAGS= --ignore-not-found=true
SUBPROJECTS += database

.PHONY: install uninstall subprojs
.DEFAULT: subprojs

install: build
	kubectl create -f auth.yaml
	kubectl create -f secrets.yaml

uninstall:
	kubectl delete $(KUBECTL_GLAGS) -f auth.yaml
	kubectl delete $(KUBECTL_GLAGS) -f secrets.yaml

build:
	eval $(shell minikube docker-env)
	docker build -t auth .

subprojs:
	@for proj in $(SUBPROJECTS); \
        do \
		$(MAKE) -C $$proj $(MAKECMDGOALS); \
	done
