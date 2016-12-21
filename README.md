# pyramixGo
Solution creator for Pyramix puzzle

## Setup
`go get github.com/petar/GoLLRB/llrb`

## Running
`go run goPyramix.go > solutions.txt`

## Understanding the output
The solutions file that it created is a set of lines of the format :
`Level # Before-State Rotation-Type After-State`

The states described is based on the puzzle as shown (the corners are omitted as they are trivial) :

```
     LEFT - RED    FRONT - YELLOW    RIGHT - GREEN
             /\          /\          /\
            /  \        /  \        /  \
           /\16/\      /\ 4/\      /\22/\
          /15\/17\    / 3\/ 5\    /21\/23\
         /\12/\14/\  /\ 0/\ 2/\  /\18/\20/\
        /  \/13\/  \/  \/ 1\/  \/  \/19\/  \
                    \  /\ 7/\  /
                     \/ 6\/ 8\/
                      \ 9/\11/
                       \/10\/
                        \  /
                         \/
                  BOTTOM - BLUE
```

The states consist of a string of 24 character (Y for Yellow, B for Blue, R for Red, G for Green) matching the 24 locations described above. For example, the solved puzzle is 
`YYYYYYBBBBBBRRRRRRGGGGGG`

Rotations are described in the clockwise direction while looking at the front face :

1. A `TOP` rotation results in piece 3-17 moving towards the back, piece 5-21 moving towards the left, and piece 15-23 moving towards the front.
2. A `LEFT` rotation results in piece 3-17 moving forwards, piece 1-7 moving backwards-left, and piece 13-9 moving up-right.
3. A `RIGHT` rotation results in piece 5-21 moving backwards-right, piece 1-7 moving up-right, and piece 19-11 moving forwards-left.
4. A `BACK` rotation results in piece 15-23 moving right, 19-11 moving left and 13-9 moving up.
5. Moves with a single quotation mark are in the opposite direction.  

## Method
1. Position the corner pieces to match their locations in the picture above (don't bother rotating them yet).
2. Find the line with the "before state" field matching your puzzle's state.
3. The "level" field indicates how many moves are remaining to solve the puzzle.
3. Perform the "rotation type" specified on that line.
4. This will result in the "after state" of that line.
5. Repeat from step 2 until the puzzle is solved, except the corners.
6. Solve the corners.

## Example
```
     LEFT - RED    FRONT - YELLOW    RIGHT - GREEN
             /\          /\          /\
            / G\        / R\        / Y\
           /\ G/\      /\ R/\      /\ Y/\
          / G\/ G\    / R\/ B\    / Y\/ Y\
         /\ R/\ R/\  /\ Y/\ B/\  /\ Y/\ G/\
        / R\/ R\/ R\/ Y\/ B\/ G\/ B\/ R\/ G\
                    \ B/\ G/\ Y/
                     \/ B\/ G\/
                      \ B/\ Y/
                       \/ B\/
                        \ B/
                         \/
                  BOTTOM - BLUE
```

1. You will notice that the corners are in their correct positions, although rotated incorrectly.
2. Find `YBBRRBBGGBBYRRRGGGYRGYYY` in the solution :
`Level 2 YBBRRBBGGBBYRRRGGGYRGYY RIGHT' YYYYRRRBBBBBBRRRGGGGGGYYY`
3. Turn the right corner anti-clockwise.
4. Find `YYYYRRRBBBBBBRRRGGGGGGYYY` in the solution :
`Level 1 YYYRRRBBBBBBRRRGGGGGGYYY TOP YYYYYYBBBBBBRRRRRRGGGGGG`
5. Turn the top corner clockwise.
6. Now the problem is solved, with the exception of the bottom-right corner, so just rotate it into position.

