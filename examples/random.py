#!/usr/bin/python3
import sys
import random
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
		sys.stdout.flush()

		theirMove = sys.stdin.buffer.read(1)

		if theirMove==b"":
			return # game is over

		if theirMove==b"\x01":
			opponentMoves.append(True)
		else:
			opponentMoves.append(False)

		numRounds=numRounds+1
#your code below this line

def main(rounds, opMoves, myMoves): # example random
	return random.random()>0.5
run()