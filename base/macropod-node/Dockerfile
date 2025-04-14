FROM node:current-alpine
WORKDIR /app
COPY . .
RUN npm install @grpc/grpc-js @grpc/proto-loader moment google-protobuf
CMD [ "node", "app.js" ]
