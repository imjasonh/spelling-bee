# Spelling Bee

Generates NYT Spelling Bee puzzles and answers

## Rules

Spelling Bee puzzles present solvers with seven letters in a circle, with one
letter in the middle. The game is to identify as many words as possible using
only these letters, which can be repeated, and where the letter in the center
must be used.

Each word is worth one point, and words using all seven letters are worth three
points.

Generally, underneath the puzzle there's a scale saying "XX points is good, YY
points is great, ZZ points is genius". The list of answers is printed in next
week's issue.

[Example](https://hackpad-attachments.s3.amazonaws.com/cmsc201f15.hackpad.com_dX1Mr4qvQnX_p.460422_1441302596606_spellingbee.png)

For the purposes of generating puzzles, a puzzle is considered "valid" if it has
at least ten answers, to weed out puzzles with too few answers, and if there is
at least one answer that uses every letter, because finding the three-pointer is
best part.

## Runtime

Without parallelism, generating seven-letter puzzles takes about five and a half
hours, and generates 51,912 unique puzzles.

## Results

The list of all valid puzzles is here:
https://storage.googleapis.com/spelling-bee/ls-7.txt

The name of the puzzle is the letters in the puzzle, sorted A-Z, then rotated so
that the required letter is first.

A puzzle with letters `R`, `T`, `D`, `A`, `S`, `I`, and the required center
letter `E` is named by the letters in the puzzle, sorted (`adeirst`), then
rotated so it starts with `e` (`eirstad`).

So the answers to that puzzle are at:
https://storage.googleapis.com/spelling-bee/eirstad.txt

The last line of the file is the total number of points for that puzzle.

## Findings

### Longest answers

48 puzzles include an answer with 15 letters, all of which are three-point
answers.

```
unconsciousness
superstructures
nonintervention
interconnection
inconveniencing
inconsistencies
consciousnesses
nationalization
```

### Most popular answers

The most popular answer is "deeded", which is found in 3,630 different puzzles.
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

Nearly all of these are three-point answers using all the letters in the puzzle.
Exceptions come from puzzles without enough answers to be valid, or answers like
`chintz` where the only longer answer is `chintzy` which is a three-point
answer.

```
backward
bigotry
blindfold
bloodhound
bobwhite
bullhorn
bumpkin
burdock
chintz
chintzy
crudity
dizzy
dogfight
drawback
equinox
exorcize
forklift
fortify
frigidity
fullback
girlhood
gluttony
gryphon
halfback
helpful
homepage
homophobia
hoodwink
horribly
hunchback
hurtful
imbroglio
invincibly
jackpot
jollity
jovial
jovially
judicial
kazoo
liturgy
luxury
menfolk
mindful
monthly
mortify
mythology
orthodoxy
pettifog
prodigy
public
publicly
quandary
quickly
rhombi
rhomboid
truthful
truthfully
unkempt
unmindful
wakeful
windfall
withhold
```

### Puzzles with the most answers

The puzzle `eirstad` has the most answers, at 456. Interestingly, `deeded` is an
answer in the #1 and #3 puzzles in the list.

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

There are 563 puzzles with only ten answers, among them:

```
abcdkrw
abcfhkl
abdgotu
abehitu
abeiotv
...and 558 more
```

### Puzzles worth the fewest total points:

517 puzzles only have 12 points.

12 points is the fewest possible points for a valid puzzle, since every puzzle
must have ten answers and at least one of them must be worth three points.

```
abcfhkl
abdgotu
abehitu
abeiotv
aceiloz
...and 512 more
```

Puzzles that have only 10 answers worth more than 12 points must include
multiple three-point answers.

### Other puzzle variants

Nothing requires the puzzle to have 7 available letters. If we only provide 6
letters, we get fewer valid puzzles, but it takes less time to compute them.

| # letters | time | # puzzles| 
------------------------------
| 3 | 20s | 0 |
| 4 | 2m20s | 189 |
| 5 | 10m29s | 6139 |
| 6 | 40m39s | 29105 |
| 7 | 180m41s | 51912 |

Times are using `-parallel=100` on a 12-core machine.
