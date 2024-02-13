FROM node:current-alpine
WORKDIR /app
COPY . .
RUN npm install express @grpc/grpc-js redis moment
CMD [ "node", "app.js" ]
