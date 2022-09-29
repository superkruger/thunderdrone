# frontend build stage
FROM node:lts-alpine as frontend-builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 3000
RUN chown -R node:node ./node_modules
CMD ["npm", "start"]

