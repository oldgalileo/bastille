#!/usr/bin/python3
import sys
import random

ROCK = b"\x00"
PAPER = b"\x01"
SCISSORS = b"\x02"

def run():
    numRounds=1
    opponentMoves=[]
    myMoves=[]

    while True:
        myMove = main(numRounds, opponentMoves, myMoves)

        sys.stdout.buffer.write(myMove)
        myMoves.append(myMove)
        sys.stdout.flush()

        theirMove = sys.stdin.buffer.read(1)
        if theirMove == b"":
            return
        opponentMoves.append(theirMove)
        numRounds+=1

model = {
    ROCK: 0.33,
    PAPER: 0.33,
    SCISSORS: 0.33,
}

def weighted_choice(choices):
    total = sum(w for c, w in choices)
    r = random.uniform(0, total)
    upto = 0
    for c, w in choices:
        if upto + w >= r:
            return c
        upto += w
    assert False, "Shouldn't get here"

def main(rounds, opMoves, myMoves):
    return weighted_choice([(ROCK, model[ROCK]), (PAPER, model[PAPER]), (SCISSORS, model[SCISSORS])])

run()