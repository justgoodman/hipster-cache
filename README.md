Cache Server
=====================

![](http://dl2.joxi.net/drive/2016/10/13/0011/0529/758289/89/f84998ed49.jpg)


Documentation
====================

I complete all code, you can see it in this links:<br />
1) https://github.com/justgoodman/hipster-cache this is CacheServer (i created my own implementation of HashTable)<br />
2) https://github.com/justgoodman/hipster-cache-proxy this is ProxyServer ( ProxyServer provides Sharding using VirtualNodes scheme.Proxy Server can works in two modes: 1) As a proxy: your client sends any command to ProxyServer, it find needed CacheServer node, sends command to this node and returns this response to your client 2) As a mediator: your client sends a key to ProxyServer, it founds needed CacheServer node and returns address of this node for Client. After that your client can sends command to specified CacheServer.<br />
3) https://github.com/justgoodman/hipster-cache-client this is Client (this client uses Proxy Server as a mediator, this client sends key to proxy server, proxy server returns CacheServerNode address. After that client sends command to this CacheServerNode.<br />
4) You can run test for my application using this code:
https://github.com/justgoodman/hipster-cache-client/tree/master/test 
<br />
<br />
You can see my kubernetes configs in this folder:<br />
https://github.com/justgoodman/hipster-cache/tree/master/kubernetes <br />


Also you can see how to run application locally using docker-compose: 
https://github.com/justgoodman/hipster-cache/tree/master/dockerJuno 

Scheme of work
===================

Proxy for Anonymus client 
-------------------------
![](http://g.recordit.co/a3JmUuYZeX.gif)

Mediator for Golang client
------------------------- 
![](http://g.recordit.co/wfAKdmnhcz.gif)

How to run locally
====================
1. For that you need to copy folder  
https://github.com/justgoodman/hipster-cache/tree/master/dockerJuno 
to you special docker folder<br>
2. In this folder do "git clone" for all applications: <br/>
2.1 git clone git@github.com:justgoodman/hipster-cache.git<br/>
2.2 git clone git@github.com:justgoodman/hipster-cache-proxy.git<br/>
2.3 git clone git@github.com:justgoodman/hipster-cache-client.git<br/>
3. If you don't install glide, install it: "go get github.com/Masterminds/glide"<br/>
4. Install all dependencies: <br/>
4.1 cd hipster-cache && glide install && cd ../ <br/>
4.2 cd hipster-cache-proxy && glide install && cd ../ <br/>
4.3 cd hipster-cache-client && glide install && cd ../ <br/>
5. After that you can run all needed envirinment using: <br/>
docker-compose up -d <br/>
<br/>
For run test, you need get docker id for hipster_cache_client: <br/>
docker ps | grep client <br/>
And using this id run coommand: <br/>
docker exec -it 38d87ddbb0bf bash -c "cd /go/src/hipster-cache-client && go test ./test/..." <br/>
In my example 38d87ddbb0bf was docker id for hipster_cache_client <br/>

**It really works!)**
![](http://dl1.joxi.net/drive/2016/11/01/0011/0529/758289/89/c72e168adf.jpg)

How to check it locally
====================
1. Add to your hosts file, IP your docker machine for "juno.net" , i set my "docker-machine ip"<br/>
2. For directly working with one of 3 HipsterCacheServer you can run:<br/>
telnet juno.net 4022<br/>
telnet juno.net 4032<br/>
telnet juno.net 4001<br/>
3. For send command to ProxyServer you can run:<br/>
telnet jono.net 4001<br/>
In proxy server you can use command:<br/>
**GET_SHARD** key <br/>
This command returns shard address for specified key<br/>
Locally you can observe:<br/>
Consul: http://juno.net:8500/ui <br/>
Promehteus: http://juno.net:9090<br/>
Grafana: http://juno.net:3001/login<br

Consul
====================

All services registereg on consul, also we have health check for CacheServers<br />
You can see working consul service on this link:<br />
http://104.155.104.83:8500/ui/#/dc1/services <br /><br />
![](http://dl1.joxi.net/drive/2016/11/01/0011/0529/758289/89/f1b19308a1.jpg)


Prometheus
===================

All services send metrics to Prometheus
You can see working prometheus service in this link:<br />
http://104.199.49.154:9090/targets <br /><br />
![](http://dl2.joxi.net/drive/2016/11/01/0011/0529/758289/89/ce746fe175.jpg)

Kubernetes 
==================
All services deployed in K8 and Work) <br/>

**Example of usage**:<br/>
![](https://github.com/justgoodman/test/blob/master/pics/hipster-cache-manial.gif)<br/>

Addresses:<br/>
**Hipster-Cache-Proxy**<br/>
IP: 130.211.82.2 <br/>
Application port: 4001 <br/>
Metrics port: 4002 <br/>
** Hipster-Cache-Server1**<br/>
IP: 104.155.106.91 <br/>
Application port: 4011 <br/>
Metrics port: 4012<br/>
** Hipster-Cache-Server2**<br/>
IP: 104.155.86.180<br/>
Application port: 4011 <br/>
Metrics port: 4012<br/>
** Hipster-Cache-Server3**<br/>
IP: 130.211.69.21<br/>
Application port: 4011 <br/>
Metrics port: 4012<br/>
<br/>
You can see this addressed in Services information:</br>
![](http://dl2.joxi.net/drive/2016/11/08/0011/0529/758289/89/778ee18438.jpg)<br/>
Information about Pods:<br/><br/>
![](http://dl1.joxi.net/drive/2016/11/08/0011/0529/758289/89/d783f91d15.jpg)<br/>
Kubernetes config you can see in this page:<br/>
https://github.com/justgoodman/hipster-cache/tree/master/kubernetes <br />

API
=====================

How to connect
-------
You can connet to ProxyServer: it sends your command to needed CacheServer and returns response of founded CacheServer
telnet 104.155.86.180 4001

Also you have opportunity connect directly to needed CacheServer
(on kubernets i don't fix bug with balances, so i can't give you IP addres of Node) <br />
![](http://dl1.joxi.net/drive/2016/11/01/0011/0529/758289/89/13df8f5ef9.jpg)<br />
(i have bug with additional \n :-) ) <br />

Strings
-------

**GET** key <br/>

Get the value of key. If the key does not exist the special value nil is returned. An error is returned if the value stored at key is not a string, because GET only handles string values.

Exmaples:<br/>
hipster_cache>GET nonexisting<br/>
(nil)</br>
hipster_cache>SET mykey "Hello"<br/>
OK <br/>
hipster_cache>GET mykey<br/>
"Hello"<br/>

**SET** key value [EX seconds] [PX milliseconds] <br/>

Set key to hold the string value.<br/>
Options<br/>
EX seconds -- Set the specified expire time, in seconds.<br/>
PX milliseconds -- Set the specified expire time, in milliseconds.<br/>
<br/>
hipster_cache>SET mykey "Hello" <br/>
OK <br/>
hipster_cache>SET mykey "Hello" EX 10 <br/>
OK <br/>


Lists
------
**LPUSH** key value

Insert specified value at the head of the list stored at key.If key does not exist, it is created as empty list before performing the push operations<br/>
hipster_cache>LPUSH junoList Money<br/>
OK <br/>

**LRANGE** key start stop

Returns the specified elements of the list stored at key. The offsets start and stop are zero-based indexes, with 0 being the first element of the list (the head of the list), 1 being the next element and so on.

hipster_cache>LPUSH mylist "one"<br/>
OK<br/>
hipster_cache>LPUSH mylist "two"<br/>
OK<br/>
hipster_cache>LRANGE mylist 0 1<br/>
"one"<br/>
"two"<br/>

**LSET** key index value <br/>

Sets the list element at index to value

hipster_cache>LPUSH mylist "one"<br/>
OK<br/>
hipster_cache>LPUSH mylist "two"<br/>
OK<br/>
hipster_cache>LSET mylist 0 "three"<br/>
OK<br>
hipster_cache>LRANGE mylist 0 1<br/>
"one"<br/>
"three"<br/>

**LLEN** key

Returns the length of the list stored at key

hipster_cache>LPUSH mylist "one"<br/>
OK<br/>
hipster_cache>LPUSH mylist "two"<br/>
OK<br/>
hipster_cache>LLEN mylist <br/>
"2"<br>

Dictionary
---------

**DSET** key field value

Sets field in the dictionary at key to value<br/>

hipster_cache>DSET item label Juno<br/>
OK<br/>
hipster_cache>DGET item label<br/>
"Juno"<br/>

**DGET** key field<br/>

Returns the value associated with field in the dictionary at key.<br/>

hipster_cache>DSET item label Juno<br/>
OK<br/>
hipster_cache>DGET item label<br/>
"Juno"<br/>

**DGETALL** key

Returns all fields and values of the dictionary at key. In the returned value, every field name is followed by its value, so the length of the reply is twice the size of the dictionary<br/>

hipster_cache>DSET item label Juno<br/>
OK<br/>
hipster_cache>DSET item color Yellow<br/>
OK<br/>
hipster_cache>DGETALL item<br/>
"label"<br/>
"Juno"<br/>
"color"<br/>
"Yellow"<br/>



Theory
====================
Hash tables
--------

For imetentation Hash tables i used chaining scheme, in this scheme, we allocated slice size of "m"(cardinality of a hash function), each element of
this slice is chain that has link to the next element.

![](http://dl2.joxi.net/drive/2016/10/14/0011/0529/758289/89/a16d5387cf.jpg)

Evaluation:
- Memory consumption in O(n+m), where n is number of objects currently stored in the hash table and m is the cardinality of the hash function
- Operations work in time O(c+1), where c is the lenght of the longest chain

Main problem: how to make both m and c small?

Solve "c" problem

For getting c small, we can use as a hash function random function from univerral family

Lemma

If h(hash function) choosen ramdomly from a "universal family", the average lenght of the longest chain c is O(1+alpha), where aplha=n/m is the load factor of the hash table

Corollaly

If h is form universal family, operation with hash table run on  average on time -> O(1+alpha). Alpha is actually contant, so in average operation will run on average in a constant time

Solve "m" problem

For getting m small, at first we set small m and we will increse(double) m iteratively, when loadFactor will be more that 0.9 

![](http://dl2.joxi.net/drive/2016/10/14/0011/0529/758289/89/1316ffb61b.jpg)

Universal family hash functions 

![](http://dl2.joxi.net/drive/2016/10/14/0011/0529/758289/89/3c62b8bcf5.jpg)

![](http://dl2.joxi.net/drive/2016/10/14/0011/0529/758289/89/5992d73d33.jpg)

In string Universal family function we have problem that based on this lemma for small "c",we need to take big "p"(cardinality of the hash function).If wetake big "p" we will consumption to much memory.

What we can do?

We can use Univeral family for integer under the result of string hash function  

Algoritm:
1. Apply random hash function from the polynomial family to the string. We get some integer number module "p"
2. Apply random hash function the universal family for integers less than "p". We get a number between 0 and m-1 

![](http://dl2.joxi.net/drive/2016/10/14/0011/0529/758289/89/20d72e0e7c.jpg)

in this formula, we can set a,b randomaly, p is the big prime number, must be more than x

For this algorithm we have this lemma:

![](http://dl1.joxi.net/drive/2016/10/14/0011/0529/758289/89/6369f8eaaf.jpg)

So that is not an universal family bease for a universal family there shouldn't be any summon L over p the probability of collision shold be at most 1 over m. But we can be very close to universal family becase we can contol "p".We can make P very big and l/p will be very small and the probabolity of collision we be at most 1/m

![](http://dl2.joxi.net/drive/2016/10/14/0011/0529/758289/89/83ce9a16f4.jpg)

For big enought p we will have:
c = O(1 + alpha), where c - lenght of the longest chain,apha - load factor

Computing PolyHash(s) runs in time O(|S|)

If lenght of the names are bounded by constant L, computing h(S) takes O(L) = O(1) time









  


