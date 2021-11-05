API_BIN = api
API_CONFIG = config.toml

SMTP_BIN = smtp
SMTP_CONFIG = config.toml

HOST = getzemail.com

.PHONY: all api_clean api_build api_build_linux api_run smtp_clean smtp_build smtp_build_linux smtp_run

all: api_clean api_build api_build_linux api_run smtp_clean smtp_build smtp_build_linux smtp_run

api_clean: 
		rm -f api/$(API_BIN)
		rm -f api/*.log

api_build:
		cd api && make build && cd ..

api_build_linux:
		cd api && make build_linux && cd ..

api_run:
		api/$(API_BIN) --config=api/$(API_CONFIG)

smtp_clean:
		rm -f smtp/$(SMTP_BIN)
		rm -f smtp/*.log

smtp_build:
		cd smtp && make build && cd ..

smtp_build_linux:
		cd smtp && make build_linux && cd ..

smtp_run:
		smtp/$(SMTP_BIN) --config=smtp/$(SMTP_CONFIG)

web_clean:
		rm -rf build/

web_build:
		cd web && yarn run build && make build_linux && cd ..

build: api_clean api_build_linux smtp_clean smtp_build_linux web_clean web_build

deploy:
		scp -r api/ smtp/ root@getzemail.com:/root/getzemail/
		scp -r web/build/ web/public/ root@getzemail.com:/root/getzemail/web/
		scp web/prd.env web/web root@getzemail.com:/root/getzemail/web/

deploy_services:
		scp api/api.service root@getzemail.com:/etc/systemd/system/
		scp smtp/smtp.service root@getzemail.com:/etc/systemd/system/
		scp web/web.service root@getzemail.com:/etc/systemd/system/
