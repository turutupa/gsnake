# gsnake

## How to install 
Execute `install.sh` script 

```bash
$ ./install.sh
```

## How to run

### Locally installed 
After you've installed you can run
```bash 
$ gsnake 
```

or you can run with flag
```bash
$ gsnake --hard
```

The different values for the game difficulty are `--easy`, `--normal`, `--hard` and `--insanity`

### Dev
Go to the root of the project and execute:
```bash
$ go run main.go 
```

## TO DO
- ~~Add bash script to install/run the game~~
- Finish README.md including how to install on machine
- Improve leaderboard 
    - Add some borders
    - Add the option to input username
    - See top scores by (difficulty) mode
- Add main menu 
    - Be able to select difficulty 
- Super fruit! Add a super fruit (maybe every 5 fruit pieces?) that will appear for an X amount of time
    - If you eat it in that time snake will not get longer, 
    - otherwise it simply disappears
- Completely new project would involve to conver this to a SSH application 
