package queries

const (
	UsersTableQuery = `CREATE TABLE IF NOT EXISTS users(
    	id SERIAL PRIMARY KEY,
        username TEXT NOT NULL UNIQUE,
        password_hash TEXT NOT NULL
	);`
	TasksTableQuery = `CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
        task TEXT NOT NULL ,
        is_completed BOOLEAN DEFAULT FALSE,
        deadline_date DATE,
        user_id INT,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE 
	);`
)
