# frontend build stage
FROM node:lts-alpine as frontend-builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 3000
CMD ["npm", "start"]

