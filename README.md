# COVID-19 Data Utility

Downloads updated data from https://covid.ourworldindata.org/data/full_data.csv.

## Installing

Make sure you have Go installed from either [golang.org](https://golang.org) or by running (macOS):

```bash
brew install go
```

Clone this repo and run the install script:

```bash
git clone git@github.com:colinc86/covid-19.git
cd covid-19
./install.sh
```

## Updating data

Run the command

```bash
covid19 update data
```

or add the `-u` flag to commands:

```bash
covid19 -u list data
```

Data is saved to `/usr/local/var/covid_full_data.csv`.

## Listing data

List data by location

```bash
covid19 list data
```

List world data

```bash
covid19 list data -w
```

List data by location and sort by total cases (or new cases/deaths and total deaths)

```bash
covid19 list data --sortBy totalCases
```

List data from a location

```bash
covid19 list data -l [location]
```

## Graphs

Graph world data by total cases

```bash
covid19 graph data
```

Graph location data by total cases
```bash
covid19 graph data -l [location]
```

Graph world data by new deaths
```bash
covid19 graph data --value newDeaths
```
