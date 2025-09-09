FROM python:3.9
WORKDIR /app

RUN apt-get update && apt-get install -y sqlite3 && rm -rf /var/lib/apt/lists/*
RUN pip install numpy pandas scikit-learn torch basedpyright pytest swifter visidata
