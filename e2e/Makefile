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

e2e-test-interchainswap-pool: cleanup-ibc-test-containers
	./scripts/run-e2e.sh TestInterchainswapTestSuite TestPoolStatus
	

e2e-test-interchainswap-single: cleanup-ibc-test-containers
	./scripts/run-e2e.sh TestInterchainswapTestSuite TestSingleDepositStatus

e2e-test-interchainswap-double: cleanup-ibc-test-containers
	./scripts/run-e2e.sh TestInterchainswapTestSuite TestDoubleDepositStatus


e2e-test-interchainswap-timeout: cleanup-ibc-test-containers
	./scripts/run-e2e.sh TestInterchainswapTestSuite TestRefundMsgOnTimeoutPacket


e2e-test-atomic-swap-base: cleanup-ibc-test-containers
	./scripts/run-e2e.sh TestAtomicSwapTestSuite TestAtomicSwap_HappyPath


e2e-test-atomic-swap: cleanup-ibc-test-containers
	./scripts/run-e2e-atomic-swap.sh

compatibility-tests:
	./scripts/run-compatibility-tests.sh $(release_branch)

.PHONY: cleanup-ibc-test-containers e2e-test compatibility-tests
