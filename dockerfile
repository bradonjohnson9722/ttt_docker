# Use the official MongoDB base image
FROM mongo:latest

# Copy custom MongoDB configuration file
COPY custom-mongo.conf /etc/mongod.conf

# Copy initialization script to be executed at container startup
COPY init-mongo.js /docker-entrypoint-initdb.d/

# Expose MongoDB default port (change to your custom port if needed)
EXPOSE 27017

# Default command to run MongoDB with custom configuration
CMD ["mongod", "--config", "/etc/mongod.conf"]
