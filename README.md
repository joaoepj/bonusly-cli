:warning: This repo is WIP. Use it at your own risk! :warning:

# Bonusly CLI

This program enables interacting with the [Bonusly API](https://help.bonus.ly/en/articles/1258685-getting-started-with-the-bonusly-api) from the command line.

It is written in Go using the [Cobra](https://github.com/spf13/cobra) library.

## Quick start guide

1. Configure Bonusly CLI to use your API token. If you don't have a token yet, you can create one [here](https://bonus.ly/api).
Then run `bonusly config --set-token <your API token>`. If no errors are returned you sucessfully added your token and are ready to go!

2. Try a command. Call `bonusly allowance` to get your current remaining allowance for this month.

## Available Commands

|Command|Description|
|-------|----------|
|`allowance`|Lists the amount of Bonuslys you currently have to spend on rewards as well as the amount you can still give away this month.|
|`award`|Use this command to award Bonuslys to another person. Takes arguments `--message`, `--hashtags`, `--recipients`, `--amount`. See below for an example.|
|`makeitrain`|This command spends the entirety of your remaining bonuslys on the specified recipients.|

## Examples

1. Get remaining allowances for this month (for giving away and spending on rewards) 
```bash
> bonusly allowance
> ...
> You still have 179 Bonusly left to give away this month.
> You still have 956 Bonusly left to spend on rewards this month. 
```
2. Send Bonuslys to someone
```bash
> bonusly award -m "Here are some bonusly for you! #team" -r "john.doe" -g "awesome, cliIsCool" -a 20
> ...
> Created bonus successfully! Check it out at bonus.ly/bonuses/<bonusPostId>
```
## How To Spend All Your Remaining Bonuslys At The End Of Each Month

This can be achieved by defining three (yes, three) cronjobs (if you are on a Unix-based system).
They will look like this:
```bash
55 23 30 4,6,9,11       bonusly makeitrain -r "john.doe, jane.doe, peter.parker" -m "You are the best #team"
55 23 31 1,3,5,7,8,10,12 bonusly makeitrain -r "john.doe, jane.doe, peter.parker" -m "You are the best #team"
55 23 28 2 bonusly makeitrain -r "john.doe, jane.doe, peter.parker" -m "You are the best #team"
```
What this will do is run `makeitrain` on the 30th of April, June, September, and November; on the 31st of January, March, May, July, August, October, and December; and on the 28th of February.

We have to define three seperate cronjobs since crontabs don't support an easy way to specify the last day of the month.

### Running On Windows Shutdown

If you don't have a server on which you can run the script, but only have a Windows workmachine, you can set up a script to trigger whenever the PC shuts down.
Personally I haven't tested this yet, but [here](https://superuser.com/a/165176) is some information regarding this.
