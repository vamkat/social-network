NAMESPACE=social-network

.PHONY: build-base delete-volumes build-cnpg build-all apply-namespace apply-pvc apply-db build-services deploy-users run-migrations logs-users logs-db deploy-all reset

# === Utils ===

build-proto:
	$(MAKE) -f backend/shared/proto/protoMakefile generate

# ~~~~~~~~~~~~~~~~~~~~~~~
# ==== Docker ====


# --- Image Creation ---
build-base:
	docker build -t social-network/go-base -f backend/docker/go/base2.Dockerfile .

build-services:
	docker build -f backend/docker/services/api-gateway.Dockerfile -t social-network/api-gateway:dev .
	docker build -f backend/docker/services/chat.Dockerfile -t social-network/chat:dev .
	docker build -f backend/docker/services/live.Dockerfile -t social-network/live:dev .
	docker build -f backend/docker/services/media.Dockerfile -t social-network/media:dev .
	docker build -f backend/docker/services/notifications.Dockerfile -t social-network/notifications:dev .
	docker build -f backend/docker/services/posts.Dockerfile -t social-network/posts:dev .
	docker build -t social-network/users:dev -f backend/docker/services/users.Dockerfile .
	docker build -t social-network/front:dev -f backend/docker/front/front.Dockerfile .
	$(MAKE) build-cnpg

build-cnpg:
	docker buildx bake -f backend/docker/cnpg/bake.hcl postgres16-cloud-native

# --- deploy from docker ---

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

api:
	docker-compose up api-gateway --build

# this network is used to let the tester connect to the rest of the containers
# since the tester won't be part of the normal docker-compose, it will have it's own so that someone gotta go out of their way to test
create-network:
	@docker network inspect social-network >nul 2>&1 || docker network create social-network


# ~~~~~~~~~~~~~~~~~~
# ==== K8s ====


# --- Preliminary ---
# install cnpg operator
op-manifest:
	kubectl apply --server-side -f \
	https://raw.githubusercontent.com/cloudnative-pg/cloudnative-pg/release-1.28/releases/cnpg-1.28.0.yaml


# --- Deployment Order ---

# 1.
apply-namespace:
	kubectl apply -f backend/k8s/ --recursive --selector stage=namespace

#  2.
apply-configs:
	kubectl apply -R -f backend/k8s/ --recursive --selector stage=config

# 3.
deploy-nginx:
	helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
	helm repo update
	helm upgrade --install nginx-ingress ingress-nginx/ingress-nginx \
		-n ingress-nginx --create-namespace

# 4.
apply-db:
	kubectl apply -f backend/k8s/ --recursive --selector stage=db

# !!! WAIT HERE !!!

# 5.
run-migrations:
	kubectl apply -f backend/k8s/ --recursive --selector stage=migration

# 6.
apply-pvc:
	kubectl apply -f backend/k8s/ --recursive --selector stage=pvc

# 7.
apply-apps:
	kubectl apply -f backend/k8s/ --recursive --selector stage=app

# 8.
apply-ingress:
	kubectl apply -f backend/k8s/nginx/api-gateway-ingress.yaml

# 9.
port-forward:
	kubectl port-forward -n frontend svc/nextjs-frontend 3000:80 

# Builds all nessary docker images
build-all:
	$(MAKE) build-base 
	$(MAKE) op-manifest 
	$(MAKE) build-cnpg 
	$(MAKE) build-services 

# Do not run this as it will probalby fail.
# Run all these in order but check that all pods are complete and 
# running before running migrations
deploy-all: 
	$(MAKE) op-manifest
	$(MAKE) apply-namespace
	$(MAKE) apply-configs
	$(MAKE) apply-db
	$(MAKE) deploy-nginx 
	$(MAKE) apply-pvc
	sleep 60  
	$(MAKE) run-migrations 
	$(MAKE) apply-apps
	$(MAKE) apply-ingress
	$(MAKE) port-forward

# Runs the docker and k8s from top to bottom
first-time:
	$(MAKE) build-all
	$(MAKE) deploy-all
	
