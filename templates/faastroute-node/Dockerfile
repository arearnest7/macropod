FROM node:latest
WORKDIR /app
COPY . .
RUN npm install -g express @grpc/grpc-js redis moment
CMD [ "node", "app.js" ]
