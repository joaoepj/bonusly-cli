# Bonusly CLI

This program enables interacting with the [Bonusly API](https://help.bonus.ly/en/articles/1258685-getting-started-with-the-bonusly-api) from the command line.

It is written in Go using the [Cobra](https://github.com/spf13/cobra) library.

## Quick start guide

1. Create a config.yml file and add your Bonusly API key.
It should look a little something like this:
```yaml
apiToken: 45lknvloi234u8hfa799
```
2. Try a command. Call `bonusly allowance` to get your current remaining allowance for this month.

## Available Commands

|Command|Description|
|-------|----------|
|`allowance`|Lists the amount of Bonuslys you currently have to spend on rewards as well as the amount you can still give away this month.|
|`award`|Use this command to award Bonuslys to another person. Takes arguments `--message`, `--hashtags`, `--recipients`, `--amount`. See below for an example.|
|`makeitrain`|This command spends the entirety of your remaining bonuslys on the specified recipients.|
