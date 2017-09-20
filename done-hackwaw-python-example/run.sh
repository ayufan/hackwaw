#!/usr/local/bin/dumb-init /bin/sh
python app.py pipe &
python app.py web
