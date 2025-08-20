-- Drop trigger first (depends on function & table)
DROP TRIGGER IF EXISTS set_timestamp ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_timestamp();

-- Drop table
DROP TABLE IF EXISTS users;