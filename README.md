Cache Server
=====================

![](http://dl2.joxi.net/drive/2016/10/13/0011/0529/758289/89/f84998ed49.jpg)

Theory
=====================

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

For this algorithm we have this lemma:

![](http://dl1.joxi.net/drive/2016/10/14/0011/0529/758289/89/6369f8eaaf.jpg)

So that is not an universal family bease for a universal family there shouldn't be any summon L over p the probability of collision shold be at most 1 over m. But we can be very close to universal family becase we can contol "p".We can make P very big and l/p will be very small and the probabolity of collision we be at most 1/m

![](http://dl2.joxi.net/drive/2016/10/14/0011/0529/758289/89/83ce9a16f4.jpg)

For big enought p we will have:
c = O(1 + alpha), where c - lenght of the longest chain,apha - load factor

Computing PolyHash(s) runs in time O(|S|)

If lenght of the names are bounded by constant L, computing h(S) takes O(L) = O(1) time









  


