KUBECTL_GLAGS= --ignore-not-found=true

.PHONY: install uninstall subprojs
.DEFAULT: subprojs

install: build
uninstall:

build:
	eval $(shell minikube docker-env)
	docker build -t allankerr/authorization .

subprojs:
	@for proj in $(SUBPROJECTS); \
        do \
		$(MAKE) -C $$proj $(MAKECMDGOALS); \
	done
