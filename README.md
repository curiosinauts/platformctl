# platformctl
A command line tool for managing CuriosityWorks online coding platform. 

VSCode in a browser aka code-server runs in a docker container. 

Users are added to the platform using the following command.
```
platformctl add user foo@example.com
```

This command will add the user to a database. This will also add the new user to Gitea, kick off Jenkins job for building docker image, call kubectl to deploy the docker image and so on.

The final product is VSCode running in a browser behind a reverse proxy with TLS support. 

Good folks at coder.com graciously open sourced code-server. CuriosityWorks online coding platform would not exist without code-server and the work of coder.com devs.

# Why CLI?
Building a full blown UI to manage and build code server is expensive in terms of development time. platformctl is packaged and deployed in CuriosityWorks slack channel as a bot command. This fits my development time budget.

# Future
I am hoping other developers will pick up my tools that I have built and replicate the online coding platform for teaching children how to code.

I am comfortable in systems integration and automation. I understand not every developers like or get to work on building systems or let alone build an entire platform.

I will be documenting how this platform was put together. Pretty much everything is automated. However, this is much to do in terms of cleaning up the hard coded configs and documentation.

![](/assets/vscode_provisioning.png)