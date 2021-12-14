KUBECTL_CMD := kubectl

.PHONY: db
db:
	# kubectl delete deploy -n db --ignore-not-found app
	docker build -t dnakano/microservice-handson/db:latest --file ./platform/db/Dockerfile .
	kind load docker-image dnakano/microservice-handson/db:latest --name kind
	kubectl apply -f ./platform/db/deployment.yaml

.PHONY: item
item:
	# kubectl delete deploy -n item --ignore-not-found app
	docker build -t dnakano/microservice-handson/item:latest --file ./services/item/Dockerfile .
	kind load docker-image dnakano/microservice-handson/item:latest --name kind
	kubectl apply -f ./services/item/deployment.yaml
