Pre-requisites:
	- PostgreSQL
	- Go
	- .gatorconfig.json

[Installation] - Go
1. Execute the following in the terminal:
	`curl -sS https://webi.sh/golang | sh`


[Installation] - Postgres
1. Execute the following in the terminal:
	`sudo apt update`
	`sudo apt install postgresql postgresql-contrib`

2. (For Linux) update postgres password by executing:
	`sudo passwd postgres`

3. Start Postgres Server in the Background
	`sudo service postgresql start`

(This is for the application proper:)
4. Enter the psql shell by:
	`sudo -u postgres psql`

5. Create Database called Gator
	`CREATE DATABASE gator`

6. Switch to Gator Database:
	`\c gator`

7. (For Linux) update user password:
	`ALTER USER postgres PASSWORD 'postgres';`


[Creation] - Gator Config
	1. Create a JSON file named .gatorconfig.json in your home directory:
	`sudo touch ~/.gatorconfig.json`

	2. Update the content of the config file:
	`sudo vi ~/.gatorconfig.json`
	`i`

3. Paste the following:
	{
	"db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
	"current_user_name": ""
	}


[Installation] - Gator
	1. you can easily install the Gator CLI Application by executing:
	`go install github.com/elitekentoy/blog`


[Commands]
	1. login <name>
		You can login to the existing user
	2. register <name>
		Register to the application
	3. reset
		Reset the database
	4. users
		Displays all the users
	5. agg
		Basically checks all the feeds
	6. addfeed <Title> <URL>
		Adds a feed to the current user
	7. feeds
		Displays all the feeds
	8. follow <URL>
		Set the user to follow the specified feed
	9. following
		Displays all the feeds that the user follows
	10. unfollow <URL>
		Unfollow the specified URL
	11. Browse
		Searches posts for user

