
.PHONY: compose build publish helm-template

compose:
	docker-compose up --build

build:
	docker build . -f Dockerfile -t broswen/vex:latest
	docker build . -f Dockerfile.provisioner -t broswen/vex-provisioner:latest

publish: build
	docker push broswen/vex:latest
	docker push broswen/vex-provisioner:latest

helm-template:
	helm template config k8s/config > k8s/config.yaml
	helm template provisioner k8s/provisioner > k8s/provisioner.yaml