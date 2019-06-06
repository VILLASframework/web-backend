TITLE := "VILLASweb Backend API Documentation"

DEPLOY_USER ?= deploy
DEPLOY_HOST ?= acs-os-fein-website
DEPLOY_PATH ?= /var/www/villas/api/web/

RSYNC_OPTS ?= --recursive --ignore-missing-args --copy-links --chown $(DEPLOY_USER):$(DEPLOY_USER)

all: index.html

index.html: swagger.json
	redoc-cli bundle --cdn --title $(TITLE) --output $@ swagger.json

deploy: index.html
	rsync $(RSYNC_OPTS) index.html swagger.json $(DEPLOY_USER)@$(DEPLOY_HOST):$(DEPLOY_PATH)

clean:
	rm -rf index.html

.PHONY: docs clean deploy
