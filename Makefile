.PHONY: run stop

run:
	cd deployments && docker-compose up -d --build && docker image prune -f ;\

stop:
	cd deployments && docker-compose down ;\