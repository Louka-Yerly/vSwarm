DOCKER_HUB_ACCOUNT=vhiveease
FUNCTIONS = adservice cartservice checkoutservice currencyservice emailservice paymentservice productcatalogservice recommendationservice shippingservice
ALL_IMAGES = $(addsuffix -image, $(FUNCTIONS))

ROOT = .

all: all-image

all-image: $(ALL_IMAGES)

adservice-image: ./adservice/Dockerfile
	DOCKER_BUILDKIT=1 docker build \
	--tag $(DOCKER_HUB_ACCOUNT)/onlineshop-adservice:latest \
	-f ./adservice/Dockerfile \
	$(ROOT)/adservice/

cartservice-image: ./cartservice/src/Dockerfile
	DOCKER_BUILDKIT=1 docker build \
	--tag $(DOCKER_HUB_ACCOUNT)/onlineshop-cartservice:latest \
	-f ./cartservice/src/Dockerfile \
	$(ROOT)/cartservice/src/

checkoutservice-image: ./checkoutservice/Dockerfile
	DOCKER_BUILDKIT=1 docker build \
	--tag $(DOCKER_HUB_ACCOUNT)/onlineshop-checkoutservice:latest \
	-f ./checkoutservice/Dockerfile \
	$(ROOT)/checkoutservice/

currencyservice-image: ./currencyservice/Dockerfile
	DOCKER_BUILDKIT=1 docker build \
	--tag $(DOCKER_HUB_ACCOUNT)/onlineshop-currencyservice:latest \
	-f ./currencyservice/Dockerfile \
	$(ROOT)/currencyservice/
currencyservice-image-riscv: ./currencyservice/Dockerfile.riscv
	DOCKER_BUILDKIT=1 docker build \
	--tag $(DOCKER_HUB_ACCOUNT)/onlineshop-currencyservice:latest \
	-f ./currencyservice/Dockerfile.riscv \
	$(ROOT)/currencyservice/

emailservice-image: ./emailservice/Dockerfile
	DOCKER_BUILDKIT=1 docker build \
	--tag $(DOCKER_HUB_ACCOUNT)/onlineshop-emailservice:latest \
	-f ./emailservice/Dockerfile \
	$(ROOT)/emailservice/
emailservice-image-riscv: ./emailservice/Dockerfile
	DOCKER_BUILDKIT=1 docker build \
	--tag $(DOCKER_HUB_ACCOUNT)/onlineshop-emailservice:latest \
	-f ./emailservice/Dockerfile.riscv \
	$(ROOT)/emailservice/

paymentservice-image: ./paymentservice/Dockerfile
	DOCKER_BUILDKIT=1 docker build \
	--tag $(DOCKER_HUB_ACCOUNT)/onlineshop-paymentservice:latest \
	-f ./paymentservice/Dockerfile \
	$(ROOT)/paymentservice/

paymentservice-image-riscv: ./paymentservice/Dockerfile.riscv
	DOCKER_BUILDKIT=1 docker build \
	--tag $(DOCKER_HUB_ACCOUNT)/onlineshop-paymentservice:latest \
	-f ./paymentservice/Dockerfile.riscv \
	$(ROOT)/paymentservice/

productcatalogservice-image: ./productcatalogservice/Dockerfile
	DOCKER_BUILDKIT=1 docker build \
	--tag $(DOCKER_HUB_ACCOUNT)/onlineshop-productcatalogservice:latest \
	-f ./productcatalogservice/Dockerfile \
	$(ROOT)/productcatalogservice/

productcatalogservice-image-riscv: ./productcatalogservice/Dockerfile.riscv
	DOCKER_BUILDKIT=1 docker build \
	--tag $(DOCKER_HUB_ACCOUNT)/onlineshop-productcatalogservice:latest \
	-f ./productcatalogservice/Dockerfile.riscv \
	$(ROOT)/productcatalogservice/

recommendationservice-image:  ./recommendationservice/Dockerfile
	DOCKER_BUILDKIT=1 docker build \
	--tag $(DOCKER_HUB_ACCOUNT)/onlineshop-recommendationservice:latest \
	-f ./recommendationservice/Dockerfile \
	$(ROOT)/recommendationservice/

recommendationservice-image-riscv:  ./recommendationservice/Dockerfile.riscv
	DOCKER_BUILDKIT=1 docker build \
	--tag $(DOCKER_HUB_ACCOUNT)/onlineshop-recommendationservice:latest \
	-f ./recommendationservice/Dockerfile.riscv \
	$(ROOT)/recommendationservice/

shippingservice-image: ./shippingservice/Dockerfile
	DOCKER_BUILDKIT=1 docker build \
	--tag $(DOCKER_HUB_ACCOUNT)/onlineshop-shippingservice:latest \
	-f ./shippingservice/Dockerfile \
	$(ROOT)/shippingservice/

shippingservice-image-riscv: ./shippingservice/Dockerfile.riscv
	DOCKER_BUILDKIT=1 docker build \
	--tag $(DOCKER_HUB_ACCOUNT)/onlineshop-shippingservice:latest \
	-f ./shippingservice/Dockerfile.riscv \
	$(ROOT)/shippingservice/

# Push images

push-%: %-image
	docker push docker.io/$(DOCKER_HUB_ACCOUNT)/$(subst push-,onlineshop-,$@):latest

push: $(addprefix push-, $(FUNCTIONS))

# Pull images

pull-%:
	docker pull docker.io/$(DOCKER_HUB_ACCOUNT)/$(subst pull-,onlineshop-,$@):latest

pull: $(addprefix pull-, $(FUNCTIONS))
