FROM node:12.18.3-alpine3.12

LABEL maintainer="Galileo Daras <galileo@getcoffee.io>"

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
EXPOSE 8080

CMD ["npm start"]
