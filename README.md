# Spelling Bee

Generates NYT Spelling Bee puzzles

## Rules

Spelling Bee puzzles present solvers with 7 letters in a circle, with one letter
in the middle. The game is to identify as many words as possible using only
these letters, which can be repeated, and where the letter in the center must be
used.

Each word is worth one point, and words using all seven letters are worth 3
points.

Generally, underneath the puzzle there's a scale saying "XX points is good, YY
points is great, ZZ points is genius". The list of answers is printed in next
week's issue.

For the purposes of generating puzzles, a puzzle is considered "valid" if it has
at least 10 answers, to weed out puzzles with too few answers, and if there is
at least one answer that uses every letter, because finding the 3-pointer is
best part.

## Runtime

Without parallelism, generating 7-letter puzzles takes about 5 and a half hours,
and generates 51912 unique puzzles.

## Findings

### Longest answers

48 puzzles have an answer with 15 letters, all of which are 3-point answers.

```
7 unconsciousness
7 superstructures
7 nonintervention
7 interconnection
7 inconveniencing
7 inconsistencies
7 consciousnesses
6 nationalization
```

### Most popular answers

The most popular answer is "deeded", which is found in 3630 different puzzles.
Basically, any puzzle that has "e" and "d" in it, where either "e" or "d" is the
required letter.

```
3630 deeded
3569 sissies
3569 sises
3516 serer
3516 seers
 ...
```

### Least popular answers

62 different words are only found once in all the puzzles.

Nearly all of these are 3-point answers using all the letters in the puzzle.
Exceptions come from puzzles without enough answers to be valid, or answers like
`chintz` where the only longer answer is `chintzy` which is a 3-point answer.

```
1 backward
1 bigotry
1 blindfold
1 bloodhound
1 bobwhite
1 bullhorn
1 bumpkin
1 burdock
1 chintz
1 chintzy
1 crudity
1 dizzy
1 dogfight
1 drawback
1 equinox
1 exorcize
1 forklift
1 fortify
1 frigidity
1 fullback
1 girlhood
1 gluttony
1 gryphon
1 halfback
1 helpful
1 homepage
1 homophobia
1 hoodwink
1 horribly
1 hunchback
1 hurtful
1 imbroglio
1 invincibly
1 jackpot
1 jollity
1 jovial
1 jovially
1 judicial
1 kazoo
1 liturgy
1 luxury
1 menfolk
1 mindful
1 monthly
1 mortify
1 mythology
1 orthodoxy
1 pettifog
1 prodigy
1 public
1 publicly
1 quandary
1 quickly
1 rhombi
1 rhomboid
1 truthful
1 truthfully
1 unkempt
1 unmindful
1 wakeful
1 windfall
1 withhold
```

### Puzzles with the most answers

The puzzle `eirstad` has the most answers, at 456. Interestingly, `deeded` is an
answer in this puzzle, and #3 in the list.

```
456 eirstad
445 eiprsta
440 eprstad
440 einrsta
439 elprsta
...
```

### Puzzles worth the most total points:

The puzzle `einrsta` (#4 on the list of most answers) has the most total points
available, because it contains a lot 3-point answers.

The puzzle with the most answers, `eirstad`, has the second-most available
points.

```
504 einrsta
484 eirstad
484 staeinr
484 taeinrs
482 einrstd
...
```

### Puzzles with the fewest answers

There are 563 puzzles with only 10 answers, among them:

```
10 abcdkrw
10 abcfhkl
10 abdgotu
10 abehitu
10 abeiotv
...
```

### Puzzles worth the fewest total points:

517 puzzles only have 12 points.

12 points is the fewest possible points for a valid puzzle, since every puzzle
must have 10 answers and at least one of them must be worth 3 points.

```
12 abcfhkl
12 abdgotu
12 abehitu
12 abeiotv
12 aceiloz
...
```

Puzzles that have only 10 answers worth more than 12 points must include
multiple 3-point answers.
