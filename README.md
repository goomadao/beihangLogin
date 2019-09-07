# beihangLogin

北航校园网登录客户端

## 编译
```
git clone https://github.com/goomadao/beihangLogin
cd beihangLogin
go build -o beihangLogin main.go
```

## 用法

```
A command line auth tool for BUAA

Usage:
  beihangLogin [command]

Available Commands:
  help        Help about any command
  login       Login using your username and password.
  logout      Logout your account
  status      Get current online info.

Flags:
      --debug   Enable debug mode.
  -h, --help    help for beihangLogin

Use "beihangLogin [command] --help" for more information about a command.
```

### 登录

```
Login using your username and password.

Usage:
  beihangLogin login [flags]

Flags:
  -h, --help              help for login
  -p, --password string   Password of your account. (required)
  -u, --username string   Username of your account. (required)

Global Flags:
      --debug   Enable debug mode.
```

### 注销

```
Logout your account

Usage:
  beihangLogin logout [flags]

Flags:
  -h, --help   help for logout

Global Flags:
      --debug   Enable debug mode.
```

### 状态查看

```
Get current online info.

Usage:
  beihangLogin status [flags]

Flags:
  -h, --help   help for status

Global Flags:
      --debug   Enable debug mode.
```