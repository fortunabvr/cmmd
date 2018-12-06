### Table of Contents
1. [About](#About)
2. [Getting Started](#GettingStarted)
    1. [Installation](#Installation)
        1. [Windows](#WindowsInstallation)
        2. [Linux/BSD/MacOSX/POSIX](#PosixInstallation)
    2. [Configuration](#Configuration)
    3. [Controlling and Querying cmmd via cmmctl](#CmmctlConfig)
    4. [Mining](#Mining)
3. [Help](#Help)
    1. [Startup](#Startup)
        1. [Using bootstrap.dat](#BootstrapDat)
    2. [Network Configuration](#NetworkConfig)
    3. [Wallet](#Wallet)
4. [Contact](#Contact)
    1. [IRC](#ContactIRC)
    2. [Mailing Lists](#MailingLists)
5. [Developer Resources](#DeveloperResources)
    1. [Code Contribution Guidelines](#ContributionGuidelines)
    2. [JSON-RPC Reference](#JSONRPCReference)
    3. [The Commercium-related Go Packages](#GoPackages)

<a name="About" />

### 1. About
cmmd is a full node Commercium implementation written in [Go](http://golang.org),
licensed under the [copyfree](http://www.copyfree.org) ISC License.

This project is currently under active development and is in a Beta state. It is
extremely stable and has been in production use since February 2016.

It also properly relays newly mined blocks, maintains a transaction pool, and
relays individual transactions that have not yet made it into a block. It
ensures all individual transactions admitted to the pool follow the rules
required into the block chain and also includes the vast majority of the more
strict checks which filter transactions based on miner requirements ("standard"
transactions).

<a name="GettingStarted" />

### 2. Getting Started

<a name="Installation" />

**2.1 Installation**<br />

The first step is to install cmmd.  See one of the following sections for
details on how to install on the supported operating systems.

<a name="WindowsInstallation" />

**2.1.1 Windows Installation**<br />

* Install the MSI available at: https://github.com/CommerciumBlockchain/cmmd/releases
* Launch cmmd from the Start Menu

<a name="PosixInstallation" />

**2.1.2 Linux/BSD/MacOSX/POSIX Installation**<br />

* Install Go according to the installation instructions here: http://golang.org/doc/install
* Run the following command to ensure your Go version is at least version 1.2: `$ go version`
* Run the following command to obtain cmmd, its dependencies, and install it: `$ go get github.com/CommerciumBlockchain/cmmd/...`<br />
  * To upgrade, run the following command: `$ go get -u github.com/CommerciumBlockchain/cmmd/...`
* Run cmmd: `$ cmmd`

<a name="Configuration" />

**2.2 Configuration**<br />

cmmd has a number of [configuration](http://godoc.org/github.com/CommerciumBlockchain/cmmd)
options, which can be viewed by running: `$ cmmd --help`.

<a name="CmmctlConfig" />

**2.3 Controlling and Querying cmmd via cmmctl**<br />

cmmctl is a command line utility that can be used to both control and query cmmd
via [RPC](http://www.wikipedia.org/wiki/Remote_procedure_call).  cmmd does
**not** enable its RPC server by default;  You must configure at minimum both an
RPC username and password or both an RPC limited username and password:

* cmmd.conf configuration file
```
[Application Options]
rpcuser=myuser
rpcpass=SomeDecentp4ssw0rd
rpclimituser=mylimituser
rpclimitpass=Limitedp4ssw0rd
```
* cmmctl.conf configuration file
```
[Application Options]
rpcuser=myuser
rpcpass=SomeDecentp4ssw0rd
```
OR
```
[Application Options]
rpclimituser=mylimituser
rpclimitpass=Limitedp4ssw0rd
```
For a list of available options, run: `$ cmmctl --help`

<a name="Mining" />

**2.4 Mining**<br />
cmmd supports both the `getwork` and `getblocktemplate` RPCs although the
`getwork` RPC is deprecated and will likely be removed in a future release.
The limited user cannot access these RPCs.<br />

**1. Add the payment addresses with the `miningaddr` option.**<br />

```
[Application Options]
rpcuser=myuser
rpcpass=SomeDecentp4ssw0rd
miningaddr=12c6DSiU4Rq3P4ZxziKxzrL5LmMBrzjrJX
miningaddr=1M83ju3EChKYyysmM2FXtLNftbacagd8FR
```

**2. Add cmmd's RPC TLS certificate to system Certificate Authority list.**<br />

`cgminer` uses [curl](http://curl.haxx.se/) to fetch data from the RPC server.
Since curl validates the certificate by default, we must install the `cmmd` RPC
certificate into the default system Certificate Authority list.

**Ubuntu**<br />

1. Copy rpc.cert to /usr/share/ca-certificates: `# cp /home/user/.cmmd/rpc.cert /usr/share/ca-certificates/cmmd.crt`<br />
2. Add cmmd.crt to /etc/ca-certificates.conf: `# echo cmmd.crt >> /etc/ca-certificates.conf`<br />
3. Update the CA certificate list: `# update-ca-certificates`<br />

**3. Set your mining software url to use https.**<br />

`$ cgminer -o https://127.0.0.1:9109 -u rpcuser -p rpcpassword`

<a name="Help" />

### 3. Help

<a name="Startup" />

**3.1 Startup**<br />

Typically cmmd will run and start downloading the block chain with no extra
configuration necessary, however, there is an optional method to use a
`bootstrap.dat` file that may speed up the initial block chain download process.

<a name="BootstrapDat" />

**3.1.1 bootstrap.dat**<br />
* [Using bootstrap.dat](https://github.com/CommerciumBlockchain/cmmd/tree/master/docs/using_bootstrap_dat.md)

<a name="NetworkConfig" />

**3.1.2 Network Configuration**<br />
* [What Ports Are Used by Default?](https://github.com/CommerciumBlockchain/cmmd/tree/master/docs/default_ports.md)
* [How To Listen on Specific Interfaces](https://github.com/CommerciumBlockchain/cmmd/tree/master/docs/configure_peer_server_listen_interfaces.md)
* [How To Configure RPC Server to Listen on Specific Interfaces](https://github.com/CommerciumBlockchain/cmmd/tree/master/docs/configure_rpc_server_listen_interfaces.md)
* [Configuring cmmd with Tor](https://github.com/CommerciumBlockchain/cmmd/tree/master/docs/configuring_tor.md)

<a name="Wallet" />

**3.1 Wallet**<br />

cmmd was intentionally developed without an integrated wallet for security
reasons.  Please see [cmmwallet](https://github.com/CommerciumBlockchain/cmmwallet) for more
information.

<a name="Contact" />

### 4. Contact

<a name="ContactIRC" />

**4.1 IRC**<br />
* [irc.freenode.net](irc://irc.freenode.net), channel #cmmd

<a name="DeveloperResources" />

### 5. Developer Resources

<a name="ContributionGuidelines" />

* [Code Contribution Guidelines](https://github.com/CommerciumBlockchain/cmmd/tree/master/docs/code_contribution_guidelines.md)
<a name="JSONRPCReference" />

* [JSON-RPC Reference](https://github.com/CommerciumBlockchain/cmmd/tree/master/docs/json_rpc_api.md)
    * [RPC Examples](https://github.com/CommerciumBlockchain/cmmd/tree/master/docs/json_rpc_api.md#ExampleCode)
<a name="GoPackages" />

* The Commercium-related Go Packages:
  * [rpcclient](https://github.com/CommerciumBlockchain/cmmd/tree/master/rpcclient) - Implements a
    robust and easy to use Websocket-enabled Commercium JSON-RPC client
  * [cmmjson](https://github.com/CommerciumBlockchain/cmmd/tree/master/cmmjson) - Provides an extensive API
    for the underlying JSON-RPC command and return values
  * [wire](https://github.com/CommerciumBlockchain/cmmd/tree/master/wire) - Implements the
    Commercium wire protocol
  * [peer](https://github.com/CommerciumBlockchain/cmmd/tree/master/peer) -
    Provides a common base for creating and managing Commercium network peers.
  * [blockchain](https://github.com/CommerciumBlockchain/cmmd/tree/master/blockchain) -
    Implements Commercium block handling and chain selection rules
  * [blockchain/fullblocktests](https://github.com/CommerciumBlockchain/cmmd/tree/master/blockchain/fullblocktests) -
    Provides a set of block tests for testing the consensus validation rules
  * [txscript](https://github.com/CommerciumBlockchain/cmmd/tree/master/txscript) -
    Implements the Commercium transaction scripting language
  * [cmmec](https://github.com/CommerciumBlockchain/cmmd/tree/master/cmmec) - Implements
    support for the elliptic curve cryptographic functions needed for the
    Commercium scripts
  * [database](https://github.com/CommerciumBlockchain/cmmd/tree/master/database) -
    Provides a database interface for the Commercium block chain
  * [mempool](https://github.com/CommerciumBlockchain/cmmd/tree/master/mempool) -
    Package mempool provides a policy-enforced pool of unmined Commercium
    transactions.
  * [cmmutil](https://github.com/CommerciumBlockchain/cmmd/tree/master/cmmutil) - Provides
    Commercium-specific convenience functions and types
  * [chainhash](https://github.com/CommerciumBlockchain/cmmd/tree/master/chaincfg/chainhash) -
    Provides a generic hash type and associated functions that allows the
    specific hash algorithm to be abstracted.
  * [connmgr](https://github.com/CommerciumBlockchain/cmmd/tree/master/connmgr) -
    Package connmgr implements a generic Commercium network connection manager.
