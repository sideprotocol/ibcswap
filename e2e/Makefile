DOCKER := $(shell which docker)
TEST_CONTAINERS=$(shell docker ps --filter "label=ibc-test" -a -q)

cleanup-ibc-test-containers:
	for id in $(TEST_CONTAINERS) ; do \
		$(DOCKER) stop $$id ; \
		$(DOCKER) rm $$id ; \
	done

e2e-test: cleanup-ibc-test-containers
	./scripts/run-e2e.sh $(suite) $(test)

e2e-test-interchainswap: cleanup-ibc-test-containers
	./scripts/run-e2e.sh TestInterchainswapTestSuite TestMsgCreatePool

e2e-test-fee: cleanup-ibc-test-containers
	./scripts/run-e2e.sh TestFeeMiddlewareTestSuite TestMultiMsg_MsgPayPacketFeeSingleSender

.PHONY: cleanup-ibc-test-containers e2e-test
