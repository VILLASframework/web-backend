TITLE := "VILLASweb Backend API Documentation"

DEPLOY_USER ?= deploy
DEPLOY_HOST ?= acs-os-fein-website
DEPLOY_PATH ?= /var/www/villas/api/web/

RSYNC_OPTS ?= --recursive --ignore-missing-args --copy-links --chown $(DEPLOY_USER):$(DEPLOY_USER)

all: index.html

index.html: api.yaml
	redoc-cli bundle --cdn --title "$(TITLE)" --output $@ api.yaml

deploy: index.html
	rsync $(RSYNC_OPTS) index.html api.yaml $(DEPLOY_USER)@$(DEPLOY_HOST):$(DEPLOY_PATH)

clean:
	rm -rf index.html

.PHONY: docs clean deploy
