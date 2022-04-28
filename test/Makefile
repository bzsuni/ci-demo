include ../Makefile.defs

WHEREABOUTS_IMAGE=ghcr.io/k8snetworkplumbingwg/whereabouts:latest-amd64

all: test

test:
ifeq (${shell (which ginkgo &>/dev/null) && echo EXIST },EXIST)
	@echo "use ginkgo from local machine"
	$(QUIET)ginkgo -vv -r ./
else
	@echo "use ginkgo from vendor"
	$(QUIET)./ginkgo.sh
endif

.PHONY: kind-init
kind-init:
	@echo "Install kind cluster..."
	@kube_proxy_mode=$(KUBE_PROXY_MODE) ip_family=$(IP_FAMILY) singe_node=$(SINGE_NODE) kind_image_tag=$(KIND_IMAGE_TAG) disable_default_cni=$(DISABLE_DEFAULT_CNI) p2ctl -t ./yaml/kind/kind-config.tmpl > ./yaml/kind/kind-config.yaml
	@cat ./yaml/kind/kind-config.yaml
	kind create cluster --config ./yaml/kind/kind-config.yaml --name $(CLUSTER_NAME)
	kubectl apply -f ./yaml/cni/cni-install.yaml
ifneq ($(INSTALL_MACVLAN),)
	@echo "Install macvlan + whereabouts"
	@bash ./hack/install-macvlan.sh
endif
	kubectl get nodes


#============ e2e ====================
.PHONY: e2e-test
e2e-test:
	@echo "Run E2E"
	#ginkgo -mod=mod -progress -reportPassed --slowSpecThreshold=60 e2e/*
	bash ./ginkgo.sh
