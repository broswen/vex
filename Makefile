
.PHONY: compose build publish

compose:
	docker-compose up --build

build:
	docker build . -f Dockerfile -t broswen/vex:latest
	docker build . -f Dockerfile.provisioner -t broswen/vex-provisioner:latest

publish: build
	docker push broswen/vex:latest
	docker push broswen/vex-provisioner:latest