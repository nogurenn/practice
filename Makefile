runSimpleEventWorker:
	@echo "Running SimpleEventWorker"
	@(cd simple-event-worker && make tidy run)
.PHONY: runSimpleEventWorker
