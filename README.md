# Coins

Solution to the coins problem presented [here](https://x.com/littmath/status/1834273354628424080):

> Flip 100 coins, labeled 1 through 100. Alice checks the coins in order (1, 2, 3, …) while Bob checks the odd-labeled coins, then the even-labeled ones (so 1, 3, 5, …, 99, 2, 4, 6, …). Who is more likely to see two heads *first*?

## Install

```bash
go install github.com/nathanjcochran/coins@main
```

## Usage

```
Usage of coins:
  -c int
    	Number of coins to flip (default 100)
  -h int
    	Number of heads for a win (default 2)
  -p format
    	Result printing format (none, short, long, space, heads) (default none)
```
