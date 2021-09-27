## gohost

![test](https://github.com/ingbyr/gohost/actions/workflows/go.yml/badge.svg)

Gohost is a simple host switcher tool supporting Linux and macOS.

**To modify the system host file gohost need running in root mode, such as `sudo gohost`.**

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
