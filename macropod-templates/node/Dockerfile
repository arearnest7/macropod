FROM sysdevtamu/macropod-node:latest
WORKDIR /app
COPY . .
RUN npm install
RUN npm install moment express axios
CMD [ "node", "app.js" ]
