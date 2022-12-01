# urlshorten
A golang URL Shortener application with functionalities to shorten and resolve the url.

# Installation using docker file
docker build -t <tag-name> Dockerfile <br />
docker run -itd -p <some-port>:<container-port> --name=<container-name> <tag-name/image-name> /bin/sh.
  
# Installing the binary
cd api/  <br />
go build -t {tag-name} . <br />
./{tag-name}
  
# Usage
curl -X POST -H "Content-Type:application/json" -d '{"longUrl":"https://www.youtube.com/watch?v=F1pWN3Lk7og"}'

Expected output:  <br />
{"id":"e017a9","longUrl":"https://www.youtube.com/watch?v=F1pWN3Lk7og","shortUrl":"http://localhost:9000/e017a9"}

# Use browser or curl to access the short url which will redirect to actual url.

curl -v http://localhost:9000/e017a9
