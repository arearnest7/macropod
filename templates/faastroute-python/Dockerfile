FROM python:3
WORKDIR /app
COPY . .
RUN pip install flask flask_restful grpcio protobuf redis requests
CMD ["python","app.py"]
