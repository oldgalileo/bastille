# Bastille

## About
Bastille is an Iterated Prisoner's Dilemma tournament manager. The aim of the project is to remove many of the limitations that afflict other IPD tournament managers and to support innovative, interesting strategies.

The core feature of Bastille is the containerization of submitted strategies. By isolating the strategies, we enable ourselves to run binaries, meaning any compiled language that can target Ubuntu is fair game. For beginner's sake, we have also included Python 3 as working with other languages and setting up compilation targets may be daunting for beginners.

Another benefit of being able to run binaries is the drastic increase in the scale of each match. Now, within the tournament, there are hundreds of matches per-strategy. We ensure that every strategy has played every other strategy 200 times (100 times as Player A, another 100 as Player B), and each of those matches can be anywhere from 100 rounds to 10,000+. Because the uploaded strategies communicate with a supervisor within the container using `stdin` and `stdout`, a 10,000 round match can happen in about a second. 

## Getting Started (macOS)

These directions are for those writing their strategy solely for Bastille. If you are coming over from Ethan's PPCG, check out the [examples](examples/). These directions are also only applicable for Python. Please feel free to submit a PR with "Getting Started" instructions for another language.

### Prerequisites

- Homebrew
- Python 3

We will be using python3 as the updated `stdin` and `stdout` api which will make writing our strategy easier. 

### Installing

To install brew, run:
```
ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```
or check https://brew.sh/ for official instructions.

To check if you have python3 installed, run:
```
python3 --version
```
Otherwise, to install python3 run: (this will take a bit)
```
brew install python3
```



### Scoring
The score is on a scale from 0–5, and is calculated by taking the average of the average score per match. This is one step away from Axelrod's method of scoring, where your score is the percentage of the way to a perfect 5.0 average score. We chose to go our route as it gives a tiny bit of insight into the moves of the strategy (and the numbers are bigger, of course).
