KUBECTL_GLAGS= --ignore-not-found=true
SUBPROJECTS += database

.PHONY: install uninstall subprojs
.DEFAULT: subprojs

install: build

uninstall:
	rm ./services

build:
	go build -o ./services ./gateway-controller

subprojs:
	@for proj in $(SUBPROJECTS); \
        do \
		$(MAKE) -C $$proj $(MAKECMDGOALS); \
	done
