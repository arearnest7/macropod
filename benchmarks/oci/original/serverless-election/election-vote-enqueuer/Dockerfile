FROM node:current-alpine3.18
WORKDIR /app
COPY . .
RUN npm install redis moment
RUN npm install
CMD [ "node", "app.js" ]
