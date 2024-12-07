FROM sysdevtamu/macropod-python:http
WORKDIR /app
COPY . .
RUN pip install -r requirements.txt
CMD ["python", "-u", "app.py"]
