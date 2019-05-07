# ServeBox


![GitHub](https://img.shields.io/github/license/ritwik310/servebox.svg)
![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/ritwik310/servebox.svg)
![Travis (.com)](https://img.shields.io/travis/com/ritwik310/servebox.svg)


**ServeBox** serves files quickly on the web directly from your local machine, what you can access via your IP-address from anywhere. Here's a quick demo of serving files,

<img src="https://gitlab.com/ritwik310/project-documents/raw/master/ServeBox/ServeBox-Demo-GIF-0.gif" alt="demo-gif"/>

> NOTE: Your connection will **not be not private**. So the whole game of authentication is pointless in the open fields of the Internet. If anyone wants to use TLS for communication, open an [Issue](https://github.com/ritwik310/servebox/issues/new), please.

# Instalation

Okay, this might piss you off 😛

```shell
go get github.com/ritwik310/servebox # You need Golang installed on your machine to run this
```

If you want me to host the **Binary** on the web. Again, just throw an [Issue](https://github.com/ritwik310/servebox/issues/new).

# How It Works

This application contains 2 main parts, a **CLI-Client** for managing files; and a **Server** that reads files from the **staging area (file-system)** and serves those from your local machine.

When you add files to ServeBox, for example `servebox add ./example.txt`. That ServeBox copies that file to a **staging area** (generally a directory), and also creates a **Password** file that contains the path to the copied file. When you start the **ServeBox Server** (`servebox start`), anyone can access the files in the staging area with the corresponding **Filepath** and **Password**.

<img src="https://gitlab.com/ritwik310/project-documents/raw/master/ServeBox/ServeBox-Application-Parts-0.png" />

> NOTE: Your connection will **not be not private**, so the communication between client and server is not encrypted. So anyone with the correct skills can read the file and even the password. **Do not share critical information. NOT SAFE!**

# Docs

Here's a list of all the commands that ServeBox supports.

To add a file to the staging area,
```shell
servebox add ./example.txt
```

To list out all the currently tracked files (all the files on the staging area),
```shell
servebox ls
```

To change password associated with a file,
```shell
servebox change-password ./example.txt
```

To remove a file from the staging area,
```shell
servebox remove ./example.txt
```

To remove all from the staging area,
```shell
servebox remove-all
```

To start the **Server** and serve files on the web,
```shell
servebox start
```

And that's it. Yep, just 6 commands (for v1). **Happy Hacking**
