FROM node:24
WORKDIR /app
RUN npm install express @grpc/grpc-js @grpc/proto-loader
COPY . .
CMD [ "node", "app.js" ]
