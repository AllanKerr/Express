KUBECTL_GLAGS= --ignore-not-found=true
SUBPROJECTS += database

.PHONY: install uninstall subprojs
.DEFAULT: subprojs

install: build
	kubectl create -f auth.yaml

uninstall:
	kubectl delete $(KUBECTL_GLAGS) -f auth.yaml

build:
	eval $(shell minikube docker-env)
	docker build -t auth .

subprojs:
	@for proj in $(SUBPROJECTS); \
        do \
		$(MAKE) -C $$proj $(MAKECMDGOALS); \
	done
