## gohost

Gohost is a simple host switcher tool supporting Linux and macOS.

**To modify the system host file gohost need running in admin mode, such as `sudo gohost`.


## Usage

### Create new host file

`gohost new HOST_NAME GROUP_NAME_1,GROUP_NAME_2,... `

For example

```shell
# create a devHost host file which belongs to dev1 and dev2 group
> gohost new devHost dev1,dev2
```

### Edit existed host file

> Make sure **vim** is installed

`gohost edit HOST_NAME`


### Delete existed host file

`gohost rm HOST_NAME_1,HOST_NAME_2,...`


### Delete existed group

`gohost rm [-g|--group] GROUP_NAME_1,GROUP_NAME_2,...`


### List host group

`gohost ls`


### List host file

`gohost ls [-a|--all]`


### Add group for existed host

`gohost cg HOST_NAME [-a|--add] GROUP_NAME_1,GROUP_NAME_2...`

For example

```shell
# add group dev3, dev4 for devHost
> gohost cg devHost -a dev3,dev4
```


### Remove group for existed host

`gohost cg HOST_NAME [-d|--delete] GROUP_NAME_1,GROUP_NAME_2...`

For example

```shell
# remove group dev1 for devHost
> gohost cg devHost -d dev1
```


### Use group host as system host

`gohost use GROUP_NAME`


### Display group but not apply to system host

`gohost use GROUP_NAME [-s|--simulate]`


### Display current system host

`gohost sys`
