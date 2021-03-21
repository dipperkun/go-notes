#!/bin/bash
for i in {1..100}
do
	echo $i
	{
		nc 172.17.0.1 12345
	}&
done
wait
exit 0
