# Build stage
FROM node:18-alpine as build

WORKDIR /app

# Copy package files first for better layer caching
COPY package*.json ./


# Install dependencies
RUN npm install

# Copy the rest of the application
COPY . .


# Set environment variables for the build
ENV REACT_APP_API_URL=/api

# Build the app
RUN npm run build

# Production stage
FROM nginx:alpine

# Copy built assets from build stage
COPY --from=build /app/build /usr/share/nginx/html

# Copy custom nginx configuration
COPY nginx-app.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
