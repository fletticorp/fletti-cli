# FletaloYA Command Line Interface

MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMWXK000000000000000000000000000000000000000000000000000KXWMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMNOocccccccccccc::;;;;;,,,,,,,,,,,,,,,,,,,,,,,,,''''''',,:dXMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMWNXXXXXXXXXXXXXXK0xc;;,,,,,,,,,,,,,,,,,,,,,,cdxxxxxxxxol:,:kNMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMN0OOOOOOOOOOOOOOOOOOO0NW0l;,,,,,,,,,,,,,,,,,,,,,,xMMMMMMMMMMWKo,,oKMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMM0c;;;;;;;;;;;;;;;;;;;:OMNd;,,,,,,,,,,,,,,,,,,,,,,xMMMMMMMMMMMMNk:':kNMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMWK0000000000000000000KNW0c;,,,,,,,,,,,,,,,,,,,,,,xMMMMMMMMMMMMMMKo,,oKWMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMNKKKKKKKKKKKKKKK0Od:,;,,,,,,,,,,,,,,,,,,,,,,xWMMMMMMMMMMMMMMNO:':kNMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMWk:::::::::::::::;;;;,;,,,,,,,,,,,,,,,,,,,,,,xWMMMMMMMMMMMMMMMW0:',xWMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMW0xxxxxxxdol:;;;;;;;;,;,,,,,,,,,,,,,,,,,,,,,,dWMMMMMMMMMMMMMMMMWO:';dXMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMWWWWWWWNWWWWMW0o;;;;;;;,;,,,,,,,,,,,,,,,,,,,,,,;looooooooooooodoodl;''':xO0XWMMMMMMMMMMMMMMMMMMM
MMMMMMMMMKdooooooooooxOXMKl;;;;;,,;,,,,,,,,,,,,,,,,,,,,,,,,'''''''''''''''''''''''''';cdk0XWMMMMMMMMMMMMMM
MMMMMMMMMKdooooooooooxOXMKl;;;;;;;;,,,,,,,,,,,,,,,,,,,,,,,,''''''''''''''''''''''''''.'..';cdONMMMMMMMMMMM
MMMMMMMMMMWWWWWWWWWWWWWN0o;;;;;;;;;,,,,,,,,,,,,,,,,,,,,,,,,,''''''''''''''''''''''''''.......':kNMMMMMMMMM
MMMMMMMMMMMMXkxxxxxxxdol:;;;;;;;;;;,,,,,,,,,,,,,,,,,,,,,,,,,'''''''''''''''''''''..''..........'dNMMMMMMMM
MMMMMMMMMMMM0c;;;;;;;;;;;;:ccc:;,;;,,,,,,,,,,,,,,,,,,,,,,,,,''''''''''''''''''',;;,,'''.........,kMMMMMMMM
MMMMMMMMMMMMOc;;;;;;;:cok0KXXXK0ko:;;,,,,,,,,,,,,,,,,,,,,,,,''''''''''''''';ok0KXXKK0Odc'........xMMMMMMMM
MMMMMMMMMMMMO:;;;;;:dKNNKkddoodk0NXx:,,,,,,,,,,,,,,,,,,,,,,,''''''''''''':xXXOdollllokKXOc'......xMMMMMMMM
MMMMMMMMMMMMk;;;;;:dNWNOc;;;;;;;:xNWkc;,,,,,,,,,,,,,,,,,,,,,'''''''''''';kWKo,'''''..':OWKc.....'xMMMMMMMM
MMMMMMMMMMMMk;;;;;lKWOlc;:lxkxo:;:dXNKd;,,,,,,,,,,,,,,,,,,,,''''''''''',dNKl,';lddol:'';OWO;....'kMMMMMMMM
MMMMMMMMMMMMk;;;;:kWKl;;:xNMMMWk:;;kMM0:,,,,,,,,,,,,,,,,,,,''''''''''''cKWx'':OWMMWWKl'.lXWo....;OMMMMMMMM
MMMMMMMMMMMM0c;;;:OM0c;;:kWMMMM0c;;kMM0c,,,,,,,,,,,,,,,,,,,''''''''''''cKWd''cKMMMMMWd'.cKWd...'oXMMMMMMMM
MMMMMMMMMMMMWKkddd0WNd;;;ck0XKOl;;c0MMKdoooooooooooooooooooooooollloolldKMO;',lkKK0Od,.'dNNklldONMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMX0xc;;:::;;:oKWMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMW0c,'',,'..';xXMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMWXOxdddxkKWMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMN0xolloodOXMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM

## FletaloYA! Desde tu consola!

### Table of Contents
- [Install](#install)
- [Shell](#shell)
- [Commands](#commands)


### Install

Just run `./install.sh`. It just copies a bounch of scripts into `/usr/local/bin`.
You can uninstall them by running `./uninstall.sh`

### Shell

You can use de commands one-by-one, in a independant way, but we also provde a shell for doing it in a easiest way.

```
> fysh
MATÍAS@FletaloYa! () > |
```

### Commands

#### Help

It shows the available commands, and it description:

```
MATÍAS@FletaloYa! (main) > help
login: Authenticate with Google to use FletaloYa API
me: Current logged in user information
requests: Current logged user requests
offers: Current logged user offers
refresh: Refresh token
roles: Current logged user roles
request price: Calculate the price for a request
create request: Starts a reqeuest flow
MATÍAS@FletaloYa! (main) >
```

#### Login

```
> . fylogin
```

Generates FYTKOEN env var, and storage it on .fytoken in your home folder.
Also, storages the refresh_token in .fyrefresh in your home folder, for refreshing the token when necessary.

*Note the dot at the begining, it ensures the export of the environmental variables to the current user session.*

#### Logout

```
> fylogout
```

Removes .fytoken and .fyrefresh from your home folder. It forces you, at next time, when tries to make an API call, to login again.

#### Me

```
> fyme
```

Retrieves the current user information.

#### Roles

```
> fyroles
```

Retrieves the current user roles.

#### Refresh

```
> fyrefresh
```

Refresh token


#### Requests

```
> fyrequests
```

Retrieves the current user requests.


#### Offers

```
> fyoffers
```

Retrieves the current user offers.

#### Create request

```
> fynewrequest "origin address" "destination address" "description" "sender:me|cell" "receiver:me|cell" "vehicle:moto|auto|miniflete"
```

Creates a new request

#### Check request price

```
> fyprice "origin address" "destination address" "vehicle:moto|auto|miniflete"
```

Retrieves the price for the request
