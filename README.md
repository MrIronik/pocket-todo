# <center>README!</center>

<br>

## What is pocketTODO



`pocketTODO` is an cli app. <br>
It scans your all project source files, and look for comments with `TODO:` keyword. <br>
Then it is stored as a json file in ~/.pocketTODO directory. <br>

Example:

```json
{
    "name":"calculator",
    "project_path": "~/home",
    "total_todos": 1,
    "files": [
        {"name":"main.c",
        "path":"~/home/main.c",
        "number_of_todos": 2,
            "todo": [
                {"line": 5,
                "content": "long message"
                },
                {"line": 69,
                "content": "short message"
                }
            ]
        },
        {"name":"foo.h",
        "path":"~/home/foo.h",
        "number_of_todos": 1,
            "todo": [
                {"line": 76,
                "content": "another thing to do"
                }
            ]

        }
    ]
}
```

<br>
<br>

## How to use it

For now after building application you should run this it in your project directory.

<br>

### Set Up
<br>

build

```zsh
go build -o pockettodo
```

<br>

add to path in zsh or bash

```zsh
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

```bash
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

<br>

add to path by copy bin file

```zsh
sudo cp pockettodo /usr/local/bin/
```

<br>

## Running

For now just simple run follow command in your project dir. <br>
In the future I will add some TUI and stuff.

```zsh
pocketTODO
```