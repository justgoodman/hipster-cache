First in this folder you need do this commands:
git clone git@github.com:justgoodman/hipster-cache.git 
git clone git@github.com:justgoodman/hipster-cache-proxy.git
git clone git@github.com:justgoodman/hipster-cache-client.git
If you don't install glide install it: "go get github.com/Masterminds/glide" 
Do operatios
cd hipster-cache && glide install && cd ../
cd hipster-cache-proxy && glide install && cd../
cd hipster-cache-client && glide install && cd../
After that you can run all needed enviroment using:
docker-compose up -d

For run test, you need get docker id for client 
docker ps | grep client
and run command:
docker exec -it 38d87ddbb0bf bash -c "cd /go/src/hipster-cache-client && go test ./test/..."

