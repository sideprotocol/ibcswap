DOCKER := $(shell which docker)
TEST_CONTAINERS=$(shell docker ps --filter "label=ibc-test" -a -q)

cleanup-ibc-test-containers:
	for id in $(TEST_CONTAINERS) ; do \
		$(DOCKER) stop $$id ; \
		$(DOCKER) rm $$id ; \
	done

e2e-test: cleanup-ibc-test-containers
	./scripts/run-e2e.sh $(entrypoint) $(test)

e2e-test-interchainswap-basic: cleanup-ibc-test-containers
	./scripts/run-e2e.sh TestInterchainswapTestSuite TestBasicMsgPacket

e2e-test-atomicswap-basic: cleanup-ibc-test-containers
	./scripts/run-e2e.sh TestAtomicSwapTestSuite TestMakeSwap

compatibility-tests:
	./scripts/run-compatibility-tests.sh $(release_branch)

.PHONY: cleanup-ibc-test-containers e2e-test compatibility-tests
