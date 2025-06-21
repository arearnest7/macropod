FROM node:24
WORKDIR /app
COPY . .
RUN npm install @grpc/grpc-js @grpc/proto-loader
CMD [ "node", "app.js" ]
