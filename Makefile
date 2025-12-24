NAMESPACE=social-network

.PHONY: build-base delete-volumes apply-namespace apply-pvc apply-db build-services deploy-users run-migrations logs-users logs-db all reset

# ==== Docker ====

build-base:
	docker build -t social-network/go-base -f backend/docker/go/base2.Dockerfile .

build-services:
	docker build -t social-network/users:dev -f backend/services/users/Dockerfile .
# 	docker build -t social-network/api-gateway:dev -f gateway/Dockerfile .

docker-up:
	$(MAKE) create-network
	$(MAKE) build-base
	docker-compose up --build

delete-volumes:
	docker compose down
	docker volume rm backend_users-db-data backend_posts-db-data backend_chat-db-data backend_notifications-db-data backend_media-db-data

docker-test:
	$(MAKE) create-network
	docker compose -f docker-test.yml up --build

docker-up-test:
	$(MAKE) docker-up
	$(MAKE) docker-test

# this network is used to let the tester connect to the rest of the containers
# since the tester won't be part of the normal docker-compose, it will have it's own so that someone gotta go out of their way to test
create-network:
	@docker network inspect social-network >nul 2>&1 || docker network create social-network

# ==== K8s ====

# 1.
apply-namespace:
	kubectl apply -f k8s/ --recursive --selector stage=namespace

#  2.
deploy-nginx:
	helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
	helm repo update
	helm upgrade --install nginx-ingress ingress-nginx/ingress-nginx \
		-n ingress-nginx --create-namespace

# 3.
apply-pvc:
	kubectl apply -f k8s/ --recursive --selector stage=pvc

# 4.
apply-db:
	kubectl apply -f k8s/ --recursive --selector stage=db

# 5.
run-migrations:
	kubectl apply -f k8s/ --recursive --selector stage=migration

# 6.
apply-configs:
	kubectl apply -R -f k8s --selector stage=config

# 7.
apply-apps:
	kubectl apply -f k8s/ --recursive --selector stage=app

# 8.
apply-ingress:
	kubectl apply -f k8s/nginx/api-gateway-ingress.yaml


build-proto:
	$(MAKE) -f shared/proto/protoMakefile generate

logs-users:
	kubectl logs -l app=users -n users -f

logs-db:
	kubectl logs -l app=users-db -n users -f

reset:
	kubectl delete namespace users --ignore-not-found=true
	kubectl create namespace users


all: build-base build-services apply-namespace deploy-nginx apply-pvc apply-db  deploy-users run-migrations apply-ingress

