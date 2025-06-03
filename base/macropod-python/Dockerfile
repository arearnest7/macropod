FROM python:3.10
WORKDIR /app
COPY . .
RUN pip install grpcio protobuf requests
CMD ["python", "-u", "app.py"]
