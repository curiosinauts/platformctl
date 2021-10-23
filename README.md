# platformctl
A command line tool for managing [online coding platform](https://console.curiosityworks.org). The online coding platform is based around `vscode in a browser` aka code server from coder.com 

the online coding platfrom currently provides the following to each user
- individual PostgreSQL schema per user
- online git repo through `gitea`
- python 3
- golang 1.16
- `terminal access in a browser`

# Privacy law compliancy
There are about 4 major privacy laws that govern which information is stored. I decided to be fully compliant by storing no PII. User's information is hashed and anonymized. Users will get funny aliases like `hungry-owl-5231` Just imagine calling your user by such alias on a conference call.

# User provisioning
Users are added to the platform using the following command.
```
platformctl add user user1@example.com
```

# Why CLI?
The coding platform management could have been built as an web application. Perhaps one day I can get there. My vision is to teach young children how to code. I might have a sign up sheet at community centers, barber shops, etc. I just want to pull up my phone and tell my chatops bot to execute `platformctl` command to set up a new user right and then.

# Hardware
I think I have spent about $900. Some of the hardware parts were bought before the pandemic so it might cost a bit more now.
- AMD Ryzen5 2600
- 64 gigs of ram
- 512 gigs SSD

PFSense and VLAN switch are optional.

# Internet speed
I used to have fairly slow internet connection but a user in east coast couldn't tell if the vscode was being served over internet. Gigabit internet is becoming more common these days. My guestmate is that the homelab environment such as mine can handle 10 concurrent users.

# User provisioning flow
![](/assets/vscode_provisioning.png)

# Homelab network topology
All of the servers are virtualized. They are all running in a single workstation.
![](/assets/curiosityworks_network_topology.svg)

# Seeking help from developer community
Small success would be me teaching handful of children how to code in my community. Bigger success would be having the online teaching platform replicated by other developers all over the world. One problem with that is some amount of familiarity with linux operating system, virtualization and automation is required. 

I will be documenting how this platform was put together. Pretty much everything is automated. 
