KUBECTL_GLAGS= --ignore-not-found=true
SUBPROJECTS += database

.PHONY: install uninstall subprojs
.DEFAULT: subprojs

install: $(GOPATH)/src/express-controller build

uninstall:
	rm ./services

$(GOPATH)/src/express-controller:
	ln -s "$(shell pwd)/express-controller" "$(GOPATH)/src/express-controller"

build:
	go build -o ./services ./express-controller

subprojs:
	@for proj in $(SUBPROJECTS); \
        do \
		$(MAKE) -C $$proj $(MAKECMDGOALS); \
	done
