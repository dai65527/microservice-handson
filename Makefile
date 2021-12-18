KUBECTL_CMD := kubectl

.PHONY: all
all: db item catalog customer

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

.PHONY: customer
customer:
	# kubectl delete deploy -n customer --ignore-not-found app
	docker build -t dnakano/microservice-handson/customer:latest --file ./services/customer/Dockerfile .
	kind load docker-image dnakano/microservice-handson/customer:latest --name kind
	kubectl apply -f ./services/customer/deployment.yaml

.PHONY: 
catalog:
	# kubectl delete deploy -n catalog --ignore-not-found app
	docker build -t dnakano/microservice-handson/catalog:latest --file ./services/catalog/Dockerfile .
	kind load docker-image dnakano/microservice-handson/catalog:latest --name kind
	kubectl apply -f ./services/catalog/deployment.yaml
