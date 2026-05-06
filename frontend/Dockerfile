# Build stage
FROM node:24-alpine AS build

WORKDIR /app

# Copy .env from root context
COPY .env ./.env

# Copy package files from frontend directory
COPY frontend/package*.json ./
RUN npm install

# Copy frontend source
COPY frontend/ .

# Run build
RUN npm run build

EXPOSE 3001

# Run Vite preview on port 3001
CMD ["npm", "run", "preview", "--", "--port", "3001", "--host", "0.0.0.0"]
