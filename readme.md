## gohost

![test](https://github.com/ingbyr/gohost/actions/workflows/test.yml/badge.svg)
![release](https://github.com/ingbyr/gohost/actions/workflows/release.yml/badge.svg)

Gohost is a simple host switcher tool supporting Windows, Linux and macOS.

**To modify the system host file gohost need running in root mode**

- For Windows user: open console in admin mode.
- For Linux and macOS user: use `sudo gohost` or login as root


## For Windows User

If new hosts not working, you probably need disable `DNS Client Service` by excuting below command in powershell (admin mode) and reboot your compouter.
```powershell
REG add "HKLM\SYSTEM\CurrentControlSet\services\dnscache" /v Start /t REG_DWORD /d 4 /f
``` 

Also you can enbale `DNS Client Service` again by excuting in powershell (admin mode).

```powershell
REG add "HKLM\SYSTEM\CurrentControlSet\services\dnscache" /v Start /t REG_DWORD /d 2 /f
```


## Usage

### Manage Host File

| Description      | Command                                            | Example             | 
|------------------|----------------------------------------------------|---------------------|
| Create host file | `gohost new HOST_NAME GROUP_NAME_1[,GROUP_NAME_2,...]` | `gohost new file1 group1,group2` |   
| Edit host file | `gohost edit HOST_NAME` | `gohost edit file1` |
| Delete host file | `gohost rm HOST_NAME_1[,HOST_NAME_2,...]` | `gohost rm file1,file2` |    
| List host file | `gohost ls -a(-all)` | `gohost ls -a` |
| Rename host file | `gohost mv HOST_NAME NEW_HOST_NAME` | `gohost mv file1 newFile`|

### Manage Group

| Description      | Command                                            | Example             | 
|------------------|----------------------------------------------------|---------------------|
| Add group for host | `gohost cg HOST_NAME -a(--add) GROUP_NAME_1[,GROUP_NAME_2,...]` | `gohost cg file1 -a group3,group4` |
| Remove group for host | `gohost cg HOST_NAME -d(--delete) GROUP_NAME_1[,GROUP_NAME_2,...]` | `gohost cg file1 -d group3,group4` |
| List group | `gohost ls` | `gohost ls` |
| Rename group | `gohost mv -g(--group) GROUP_NAME NEW_GROUP_NAME` | `gohost mv -g group1 newGroup`|

### Apply Group

| Description      | Command                                            | Example             |
|------------------|----------------------------------------------------|---------------------|
| Apply to system host | `gohost use GROUP_NAME` | `gohost use group1`|
| Display group content | `gohost use GROUP_NAME -s(--simulate)` | `gohost use group1 -s`|
| Display system host | `gohost sys` | `gohost sys`|
