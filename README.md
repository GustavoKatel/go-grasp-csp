
## Greedy implementation

The initial solution is constructed interactively by sorting the most common characters in each position of the set of strings. The characters are then added to the RCL. The size of the RCL, alpha, determines the greediness of the algorithm. The chosen character is them randomly selected from the RCL.

## Local Search

Here is used a simple disturbance of the constructed solution. The method randomly changes characters to others from the alphabet. The total number of permutations can be controlled by `NhbMax`

## The Method

```
func CSP(strings []string, alphabet []string, stringSize int, maxIterations int, alpha int, NhbMax int) (string,int,int)
```

Where:
- `strings` the initial set of strings to work with
- `alphabet` array containing the characters of the strings
- `stringSize` the size of the resulting string
- `maxIterations` number of iterations
- `NhbMax` greediness (a.k.a. the size of the RCL)

Return:
- `string` containing the solution
- `int` lower bound. Min distance to the set
- `int` upper bound. Max distance to the set