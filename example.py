#!/usr/local/bin/python3
import sys

def run():
	numRounds=1 # number of rounds starts at 1 because yall don't know how to program
	opponentMoves=[]
	myMoves=[]

	while True:

		myMove = main(numRounds, opponentMoves, myMoves)
		if myMove:
			sys.stdout.buffer.write(b"\x01")
			myMoves.append(True)
		else:
			sys.stdout.buffer.write(b"\x00")
			myMoves.append(False)
			# possible outputs are \x00 and \x01. anything else is invalid.
		

		theirMove = sys.stdin.buffer.read(1)

		if theirMove==b"":
			return # game is over

		if theirMove==b"\x01":
			opponentMoves.append(True)
		else:
			opponentMoves.append(False)

		numRounds=numRounds+1



#your code below this line

def main(rounds, opMoves, myMoves): # example tit for tat
	if rounds == 1:
		return True
	if opMoves[-1] == 1:
		return True
	if opMoves[-1] == 0:
		return False
	# a way to check if your code is working: printf "\x00\x00\x01\x01" | python program.py | hexdump
	# it will print out what it would move, in the same format of 00 = defect, 01 = coop



# put this at the end
run()