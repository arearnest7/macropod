FROM python:3.12
WORKDIR /app
RUN pip install --upgrade pip --root-user-action=ignore
RUN pip install --root-user-action=ignore grpcio protobuf flask flask_restful
COPY . .
CMD ["python", "-u", "app.py"]
