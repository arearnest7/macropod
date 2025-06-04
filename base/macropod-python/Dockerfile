FROM python:3.12
WORKDIR /app
COPY . .
RUN pip install --upgrade pip --root-user-action=ignore
RUN pip install --root-user-action=ignore grpcio protobuf requests
CMD ["python", "-u", "app.py"]
