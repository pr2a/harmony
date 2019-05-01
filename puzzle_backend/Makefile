validate:
	swagger validate ./swagger/swagger.yml

gen: validate
	swagger generate server \
		--target=./swagger \
		--spec=./swagger/swagger.yml \
		--exclude-main \
		--name=harmony-puzzle

.PHONY: gen validate
