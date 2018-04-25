# SSL-CP

![license](https://img.shields.io/github/license/reddec/ssl-cp.svg) ![python](https://img.shields.io/badge/python-3%2B-yellow.svg)

Control panel for organize, manager, sign and revoke certificates.
Project is still under **heavy development** but it's already using for internal projects.

# Docker

Look at `Dockerfile` or pull from dockerhub https://hub.docker.com/r/reddec/ssl-cp/

# Install

Requirements: python3, virtual environment

1. Download or clone archive from github
2. Unpack and `cd` to directory with sources
3. Create virtual environment: `python3 -m venv venv`
4. Enable it: `. ./venv/bin/activate`
5. Install requirements: `pip3 install -r requirements.txt`
6. Run it: `python3 main.py`

# Features

## Web UI + API

Yes, ssl-cp contains full-featured **mobile friendly** flask-based panel.
And automation friendly API!



## Multi-project

System allows you separate certificates by projects. Each project contains own CA (private key + public cert) and list of issued (signed by CA) certificates.
![screencapture-127-0-0-1-5000-2018-04-17-21_28_39](https://user-images.githubusercontent.com/6597086/38889188-97beed0e-4286-11e8-9278-16d05be3ac9e.png)
## Auto-generated CA
One click for generate self-signed CA. In roadmap - upload your own

![screencapture-127-0-0-1-5000-project-1-2018-04-17-21_31_54](https://user-images.githubusercontent.com/6597086/38889235-c43d9240-4286-11e8-87d4-8fd5c3e582c5.png)

After generation you can download and use it as always
![screencapture-127-0-0-1-5000-project-1-2018-04-17-21_32_54](https://user-images.githubusercontent.com/6597086/38889281-e9275b22-4286-11e8-86fd-3bd688ee07c4.png)

## One-click generation of signed certifiactes

Just provide common name (it maybe any label, node name, domain and e.t.c) and press generate.
New certificate will be automatically signed by CA
![screencapture-127-0-0-1-5000-project-1-certificates-2018-04-17-21_35_23](https://user-images.githubusercontent.com/6597086/38889377-40f8d524-4287-11e8-93f0-7d2dc6ab116b.png)

View and manage each certificate:

* Download
* Revoke
* Use cookbooks

![screencapture-127-0-0-1-5000-certificate-2-2018-04-17-21_36_34](https://user-images.githubusercontent.com/6597086/38889448-6eb22f10-4287-11e8-80fb-ace95704a7cf.png)

## Use cook-books

### Stunnel

[stunnel on Wiki](https://en.wikipedia.org/wiki/Stunnel) - SSL it! Even if it was not originally designed for SSL.

Panel can prepare full-featured archive (including SystemD service file and install script) for **client** and **server** configuration of stunnel

![screencapture-127-0-0-1-5000-certificate-2-stunnel-2018-04-17-21_38_16](https://user-images.githubusercontent.com/6597086/38889526-a8538228-4287-11e8-82d2-3f34f9e1adf0.png)
