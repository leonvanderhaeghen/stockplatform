FROM mongo-express:latest

# Install required packages
RUN apt-get update && apt-get install -y iputils-ping

# Custom entrypoint script
COPY docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["mongo-express"]
