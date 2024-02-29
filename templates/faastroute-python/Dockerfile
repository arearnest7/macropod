FROM python:3.10
WORKDIR /app
COPY . .
RUN pip install flask flask_restful grpcio protobuf redis requests
CMD ["python", "-u", "app.py"]
