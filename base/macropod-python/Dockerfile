FROM python:3.13-alpine
WORKDIR /app
COPY . .
RUN pip install --root-user-action=ignore grpcio protobuf requests
CMD ["python", "-u", "app.py"]
