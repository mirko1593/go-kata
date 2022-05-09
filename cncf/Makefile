build:
	docker build \
		-f Dockerfile \
		-t cncf-service:latest --build-arg BUILD_REF=develop \
		.


reload:
	kubectl delete deployment cncf-service
	minikube image rm cncf-service:latest
	minikube image load cncf-service:latest
	kubectl apply -f k8s/deployment.yaml

