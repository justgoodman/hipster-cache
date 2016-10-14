Cache Server
=====================

![](http://dl2.joxi.net/drive/2016/10/13/0011/0529/758289/89/f84998ed49.jpg)

Theory
=====================

Hash tables
--------

For imetentation Hash tables i used chaining scheme, in this scheme, we allocated slice size of "m"(cardinality of a hash function), each element of
this slice is chain that has link to the next element.

![](http://dl2.joxi.net/drive/2016/10/14/0011/0529/758289/89/7f8c6a752b.jpg)

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

![](http://dl2.joxi.net/drive/2016/10/14/0011/0529/758289/89/679e335dcd.jpg)










  


