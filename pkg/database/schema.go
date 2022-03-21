package database

// CreateSchema SQL script for creating the schema
var CreateSchema = `
CREATE TABLE users (
     id               SERIAL NOT NULL,
     google_id        CHARACTER VARYING(64),
     username         CHARACTER VARYING(50) NOT NULL,
     password         CHARACTER VARYING(100) NOT NULL,
     email            CHARACTER VARYING(50) NOT NULL,
     hashed_email     CHARACTER VARYING(64) NOT NULL,
     is_active        BOOLEAN NOT NULL,
     private_key      CHARACTER VARYING(4000),
     public_key       CHARACTER VARYING(1000),
     public_key_id    INTEGER, 
	 ide              CHARACTER VARYING(100),  
	 git_repo_uri     CHARACTER VARYING(100), 
	 runtime_installs CHARACTER VARYING(100),  
     PRIMARY KEY (id),
     UNIQUE (username),
	 UNIQUE (hashed_email)
);


CREATE TABLE runtime_install (
	id SERIAL NOT NULL, 
	name CHARACTER VARYING(30), 
	script_body CHARACTER VARYING(4000), 
	PRIMARY KEY (id), 
	UNIQUE (name)
);


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
DROP TABLE IF EXISTS runtime_install CASCADE;

DROP TABLE IF EXISTS users CASCADE;
`
