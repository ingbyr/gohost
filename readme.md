## gohost

Gohost is a simple host switcher tool supporting Linux, macOS and Windows.

**To modify the system host file gohost need running in admin mode, such as `sudo gohost` in Linux and macOS
or `Admin CMD` in Windows.**

## Usage

### And new host file

`gohost new HOST_NAME GROUP_NAME_1 GROUP_NAME_2 ... `

For example

```shell
# create a devHost host file which belongs to dev1 and dev2 group
> gohost new devHost dev1 dev2
```

### Edit existed host file
> Make sure `vim` is in system PATH

`gohost edit HOST_NAME`


### Delete existed host file

`gohost rm HOST_NAME`


### List host

`gohost ls [-a]`

For example

```shell
# List host group
> gohost ls
+-------+---------+
| GROUP |  HOSTS  |
+-------+---------+
| dev1  | devHost |
| dev2  | devHost |
+-------+---------+

# List host
> gohost ls -a
+---------+------------+
|  HOST   |   GROUPS   |
+---------+------------+
| devHost | dev1, dev2 |
+---------+------------+

```


### Add group to existed host

`gohost cg HOST_NAME [ -a | -d ] GROUP_NAME_1,GROUP_NAME_2...`

For example

```shell
# add group dev3, dev4 for devHost
> gohost cg devHost -a dev3,dev4
added groups 'dev3, dev4'
+---------+------------------------+
|  HOST   |         GROUPS         |
+---------+------------------------+
| devHost | dev1, dev2, dev3, dev4 |
+---------+------------------------+

# remove group dev1 for devHost
> gohost cg devHost -d dev1
removed groups 'dev1'
+---------+------------------+
|  HOST   |      GROUPS      |
+---------+------------------+
| devHost | dev2, dev3, dev4 |
+---------+------------------+
```


### Use group host as system host

`gohost use GROUP_NAME`


### Display current system host

`gohost sys`
