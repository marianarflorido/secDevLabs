FROM mongo:4.0.3

ADD deployments/mongo-init.sh /docker-entrypoint-initdb.d/
