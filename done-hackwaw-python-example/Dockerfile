FROM python:3-alpine

RUN ["pip", "install", "honcho"]
RUN ["apk", "add", "-U", "git"]

COPY requirements.txt /usr/src/app/
WORKDIR /usr/src/app
RUN pip install --no-cache-dir -r requirements.txt

COPY . /usr/src/app
CMD ["honcho", "start"]
