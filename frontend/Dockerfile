# Build stage
FROM node:24-alpine

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

EXPOSE 3001

# Run Vite preview on port 3001
CMD ["npm", "run", "preview", "--", "--port", "3001", "--host", "0.0.0.0"]
