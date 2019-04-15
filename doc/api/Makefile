

apidoc: index.html

index.html: swagger-codegen-cli.jar api.yaml
	java -jar swagger-codegen-cli.jar generate -i api.yaml -l html2 -o ./

swagger-codegen-cli.jar:
	wget http://central.maven.org/maven2/io/swagger/swagger-codegen-cli/3.0.0-rc1/swagger-codegen-cli-3.0.0-rc1.jar -O swagger-codegen-cli.jar

clean:
	rm -rf index.html swagger-codegen-cli.jar

.PHONY: docs index.html clean