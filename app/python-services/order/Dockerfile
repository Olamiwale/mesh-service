FROM python:3.12-slim

WORKDIR /app

# Copy app code and requirements
COPY requirements.txt .
COPY main.py .

# Install all dependencies from requirements.txt
RUN pip install -r requirements.txt

CMD ["python", "main.py"]
