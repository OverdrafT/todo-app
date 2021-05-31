build-image:
	docker build . -t silverspase/todo

run-app-container: build-image
	docker run --rm -p 8000:8000 silverspase/todo



