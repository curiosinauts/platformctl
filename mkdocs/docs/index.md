# online coding platform setup guide

This documentation is for those who want to set up their own online coding platform at home. You can do this in the cloud but what's the fun in that? I am cheap. I don't like looking at high AWS or GCP bills. It is so easy to spend couple of hundred dollars a month in the cloud environment. My monthly recurring cost is my internet bill, $70 a month plus DNS hosting fee, 50 cents a month.

## Knowledge acquired
I personally learned a lot by doing all of the devops work. I learned to automate by leveraging CloudFlare, LetsEncrypt, Ansible, and K3S. I have much deeper understanding of how to secure servers. For example, adding TLS to Postgres server was simpler that I thought. TLS is something I would have gotten for free as part of AWS RDS. I give credit to my devops community for explaining what hair pinning is in routers. I learned to use Google OIDC for authentication. On the non profit business risk management side of things I complied with all privacy laws by anonymizing PII. There is much more but I will stop here.

## Target audience
developers, devops persons, or any tech savy users. This would be easy task for someone with descent devops experience to follow this guide. However developers who want to teach coding online would be the users of this platform. If you are a developer with limited devops experience I would suggest you seek professinal help from devops community. Having said that, anyone who is willing to learn can use this guide to complete the set up.

## Hardware
My hardware is a workstation with 64 gigs of ram, running AMD Ryzen 5 2600 CPU, and 512gig SSD. 

## Tools

**Promox** - type 1 hypervisor which is used for virtualizing all servers needed by curiosityworks.org coding platform. You can use any other virtualization enviroment but proxmox is the one I prefer.

LearnLinuxTV has good youtube videos on this topic. Check out the following.
<iframe width="560" height="315" src="https://www.youtube.com/embed/MO4CaHn1EjM" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

**Ansible** - configuration management tool similar to chef or puppet. You don't need to possess a deep understanding of the tool since almost all of the execution is initiated from shell scripts with input prompts.

