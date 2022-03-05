package database

// CreateSchema SQL script for creating the schema
var CreateSchema = `
CREATE TABLE platformctl.user (
     id           SERIAL NOT NULL,
     google_id    CHARACTER VARYING(64),
     username     CHARACTER VARYING(50) NOT NULL,
     password     CHARACTER VARYING(100) NOT NULL,
     email        CHARACTER VARYING(50) NOT NULL,
     hashed_email CHARACTER VARYING(64) NOT NULL,
     is_active    BOOLEAN NOT NULL,
     private_key  CHARACTER VARYING(4000),
     public_key   CHARACTER VARYING(1000),
     public_key_id INTEGER, 
     PRIMARY KEY (id),
     UNIQUE (username),
	 UNIQUE (hashed_email)
);


CREATE TABLE ide (
	id SERIAL NOT NULL, 
	name CHARACTER VARYING(30), 
	PRIMARY KEY (id)
);


CREATE TABLE runtime_install (
	id SERIAL NOT NULL, 
	name CHARACTER VARYING(30), 
	script_body CHARACTER VARYING(4000), 
	PRIMARY KEY (id), 
	UNIQUE (name)
);


CREATE TABLE user_repo (
	id SERIAL NOT NULL, 
	uri CHARACTER VARYING(100), 
	user_id INTEGER, 
	PRIMARY KEY (id),
	CONSTRAINT userrepo_fk1 FOREIGN KEY (user_id) REFERENCES "user" ("id")	
);


CREATE TABLE ide_repo (
  id SERIAL NOT NULL, 
  user_ide_id INTEGER, 
  uri CHARACTER VARYING(100), 
  PRIMARY KEY (id),
  CONSTRAINT iderepo_fk1 FOREIGN KEY (user_ide_id) REFERENCES "user_ide" ("id")  
);


CREATE TABLE user_ide (
	id SERIAL NOT NULL, 
	user_id INTEGER, 
	ide_id INTEGER, 
	PRIMARY KEY (id),
	CONSTRAINT useride_fk1 FOREIGN KEY (user_id) REFERENCES "user" ("id"),
	CONSTRAINT useride_fk2 FOREIGN KEY (ide_id)  REFERENCES "ide"  ("id")
);


CREATE TABLE ide_runtime_install (
	id SERIAL NOT NULL, 
	user_ide_id INTEGER, 
	runtime_install_id INTEGER, 
	PRIMARY KEY (id),
	CONSTRAINT ideruntimeinstall_fk1 FOREIGN KEY (user_ide_id)        REFERENCES "user_ide" ("id"),
	CONSTRAINT ideruntimeinstall_fk2 FOREIGN KEY (runtime_install_id) REFERENCES "runtime_install" ("id")
);


-- data goes here
INSERT INTO ide (id, name) VALUES (1, 'vscode');
INSERT INTO ide (id, name) VALUES (2, 'intellij');
INSERT INTO ide (id, name) VALUES (3, 'pycharm');
INSERT INTO ide (id, name) VALUES (4, 'goland');


INSERT INTO runtime_install (id, name, script_body) VALUES (1, 'emberjs', '# ember
 sudo sudo npm install -g ember-cli
');
INSERT INTO runtime_install (id, name, script_body) VALUES (2, 'tmux', '# tmux
sudo apt-get install -y tmux
echo -e "\nalias s=''tmux new -A -s shared''" >> /home/coder/.zshrc
');
INSERT INTO runtime_install (id, name, script_body) VALUES (3, 'github cli', '# github cli gh install
curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo gpg --dearmor -o /usr/share/keyrings/githubcli-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
sudo apt-get update
sudo apt-get install gh
');
INSERT INTO runtime_install (id, name, script_body) VALUES (4, 'poetry', '# poetry
curl -sSL https://raw.githubusercontent.com/python-poetry/poetry/master/get-poetry.py | python3 -
echo -e "\nexport PATH="\$HOME/.poetry/bin:\$PATH"" >> ~/.zshrc
');
INSERT INTO runtime_install (id, name, script_body) VALUES (5, 'postgres', '# vag dependencies
sudo apt-get install -y postgresql
sudo apt-get install -y libpq-dev
');

`

// DropSchema SQL script for dropping the schema
var DropSchema = `
DROP TABLE IF EXISTS ide_runtime_install CASCADE;

DROP TABLE IF EXISTS user_repo;

DROP TABLE IF EXISTS ide_repo CASCADE;

DROP TABLE IF EXISTS user_ide CASCADE;

DROP TABLE IF EXISTS ide CASCADE;

DROP TABLE IF EXISTS runtime_install CASCADE;

DROP TABLE IF EXISTS platformctl.user CASCADE;
`
