# records
A practice exercise involving parsing and sorting a set of records.

## Problem
### Step 1 - Build a system to parse and sort a set of records
Create a command line app that takes as input a file with a set of
records in one of three formats described below, and outputs (to the
screen) the set of records sorted in one of three ways.

#### Input
A record consists of the following 5 fields: last name, first name,
gender, date of birth and favorite color. The input is 3 files, each
containing records stored in a different format. You may generate
these files yourself, and you can make certain assumptions if it makes
solving your problem easier.

The pipe-delimited file lists each record as follows:
- `LastName | FirstName | Gender | FavoriteColor | DateOfBirth`

The comma-delimited file looks like this:
- `LastName, FirstName, Gender, FavoriteColor, DateOfBirth`

The space-delimited file looks like this:
- `LastName FirstName Gender FavoriteColor DateOfBirth`

You may assume that the delimiters (commas, pipes and spaces) do not
appear anywhere in the data values themselves. Write a program in a
language of your choice to read in records from these files and
combine them into a single set of records.

#### Output
Create and display 3 different views of the data you read in:
- Output 1 – sorted by gender (females before males) then by last name
  ascending.
- Output 2 – sorted by birth date, ascending.
- Output 3 – sorted by last name, descending.

Display dates in the format M/D/YYYY.

### Step 2 - Build a REST API to access your system
Tests for this section are required as well. Within the same code
base, build a standalone REST API with the following endpoints:
- POST /records - Post a single data line in any of the 3 formats
  supported by your existing code
- GET /records/gender - returns records sorted by gender
- GET /records/birthdate - returns records sorted by birthdate
- GET /records/name - returns records sorted by name

It's your choice how you render the output from these endpoints as
long as it well structured data. These endpoints should return JSON.
To keep it simple, don't worry about using a persistent datastore.

## Modifications I've made to the problem
1. The command line application can read from stdin even though it was
   only stated that it needs to read from filenames passed as
   arguments. I made this change because I think it is valuable to be
   consistent with how other command line tools (cat, grep, sed,
   etc...) function.
2. After making decision (1) above, I then decided that the command
   line application should be able to accept a single file containing
   multiple different formats. My thinking was basically that you can
   do: `cat format1.txt format2.txt | grep regex` so you should also
   be able to do: `cat format1.txt format2.txt | thiscmd`.
3. A consequence of (2) is that each line must NOT contain any of the
   other possible delimiters. So the pipe delimited record example
   `LastName | FirstName | Gender | FavoriteColor | DateOfBirth` is
   invalid because it contains space which is another delimiter. The
   line should look like this instead:
   `LastName|FirstName|Gender|FavoriteColor|DateOfBirth`

It's a valuable skill as a programmer to do the minimum amount of work
that is required to solve a problem (which I am not doing here because
of these modifications) but I also think it's important for developers
to think for themselves and make modifications to a specification if
they think it will produce better software. Obviously if this was a
real work task I would discuss these kinds of changes before going
rogue and just making them.
