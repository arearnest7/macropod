FROM sysdevtamu/macropod-node:latest
WORKDIR /app
COPY . .
RUN npm install
RUN npm install express @grpc/grpc-js @grpc/proto-loader moment google-protobuf
CMD [ "node", "app.js" ]
