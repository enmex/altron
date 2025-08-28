FROM nikolaik/python-nodejs:python3.9-nodejs20-bullseye

RUN npm install --global curlconverter

ENV LC_ALL=C.UTF-8
ENV LANG=C.UTF-8
WORKDIR /app
RUN mkdir files
COPY ./services/converter/requirements.txt .
RUN pip3 install -r requirements.txt
COPY ./services/converter .

CMD ["python3", "app.py"]