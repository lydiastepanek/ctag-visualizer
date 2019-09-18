FROM python:alpine
EXPOSE 8080
ADD . .
RUN pip install Flask
ENTRYPOINT ["python"]
CMD ["app.py"]
