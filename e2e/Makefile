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

e2e-test-interchainswap-make-pool: cleanup-ibc-test-containers
	./scripts/run-e2e.sh TestInterchainswapTestSuite TestMakePoolMsgPacket

e2e-test-interchainswap-pool: cleanup-ibc-test-containers
	./scripts/run-e2e.sh TestInterchainswapTestSuite TestSwapStatus

e2e-test-interchainswap-cancel-pool: cleanup-ibc-test-containers
	./scripts/run-e2e.sh TestInterchainswapTestSuite TestCancelPoolMsgPacket
	

e2e-test-interchainswap-single-deposit: cleanup-ibc-test-containers
	./scripts/run-e2e.sh TestInterchainswapTestSuite TestSingleDepositStatus

e2e-test-interchainswap-multi-deposit: cleanup-ibc-test-containers
	./scripts/run-e2e.sh TestInterchainswapTestSuite TestMultiDepositStatus

e2e-test-interchainswap-multi-deposit-cancel: cleanup-ibc-test-containers
	./scripts/run-e2e.sh TestInterchainswapTestSuite TestCancelMultiDepositStatus
	

e2e-test-interchainswap-single-withdraw: cleanup-ibc-test-containers
	./scripts/run-e2e.sh TestInterchainswapTestSuite TestSingleWithdrawStatus

e2e-test-interchainswap-multi-withdraw: cleanup-ibc-test-containers
	./scripts/run-e2e.sh TestInterchainswapTestSuite TestMultiWithdrawStatus

e2e-test-interchainswap-timeout: cleanup-ibc-test-containers
	./scripts/run-e2e.sh TestInterchainswapTestSuite TestRefundMsgOnTimeoutPacket

e2e-test-interchainswap-overflow: cleanup-ibc-test-containers
	./scripts/run-e2e.sh TestInterchainswapTestSuite TestOverflow

e2e-test-interchainswap-swap: cleanup-ibc-test-containers
	./scripts/run-e2e.sh TestInterchainswapTestSuite TestSwapStatus
	
e2e-test-atomic-swap-base: cleanup-ibc-test-containers
	./scripts/run-e2e.sh TestAtomicSwapTestSuite TestAtomicSwap_HappyPath


e2e-test-atomic-swap: cleanup-ibc-test-containers
	./scripts/run-e2e-atomic-swap.sh

e2e-test-atomic-swap-fifo: cleanup-ibc-test-containers
	./scripts/run-e2e.sh TestAtomicSwapTestSuite TestAtomicSwapFiFoPath

compatibility-tests:
	./scripts/run-compatibility-tests.sh $(release_branch)

.PHONY: cleanup-ibc-test-containers e2e-test compatibility-tests
