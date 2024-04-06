FROM node:latest
WORKDIR /app
COPY . .
RUN npm install express @grpc/grpc-js @grpc/proto-loader moment google-protobuf
CMD [ "node", "app.js" ]
