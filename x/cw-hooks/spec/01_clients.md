# Clients

## Command Line Interface

Find below a list of `furyad` commands added with the `x/cw-hooks` module. You can obtain the full list by using the `furyad -h` command. A CLI command can look like this:

```bash
furyad query cw-hooks params
```

### Queries

| Command            | Subcommand             | Description                              |
| :----------------- | :--------------------- | :--------------------------------------- |
| `query` `cw-hooks` | `params`               | Get module params                        |
| `query` `cw-hooks` | `governance-contracts` | Get registered governance contracts      |
| `query` `cw-hooks` | `staking-contracts`    | Get registered staking contracts         |

### Transactions

| Command         | Subcommand   | Description                           |
| :-------------- | :----------- | :------------------------------------ |
| `tx` `cw-hooks` | `register`   | Register a contract for events        |
| `tx` `cw-hooks` | `unregister` | Unregister a contract from events     |

## gRPC Queries

| Verb   | Method                                            |
| :----- | :------------------------------------------------ |
| `gRPC` | `furya.cwhooks.v1.Query/Params`                    |
| `gRPC` | `furya.cwhooks.v1.Query/StakingContracts`          |
| `gRPC` | `furya.cwhooks.v1.Query/GovernanceContracts`       |
| `GET`  | `/furya/cwhooks/v1/params`                         |
| `GET`  | `/furya/cwhooks/v1/staking_contracts`              |
| `GET`  | `/furya/cwhooks/v1/governance_contracts`           |

### gRPC Transactions

| Verb   | Method                                      |
| :----- | :------------------------------------------ |
| `gRPC` | `furya.cwhooks.v1.Msg/RegisterStaking`       |
| `gRPC` | `furya.cwhooks.v1.Msg/UnregisterStaking`     |
| `gRPC` | `furya.cwhooks.v1.Msg/RegisterGovernance`    |
| `gRPC` | `furya.cwhooks.v1.Msg/UnregisterGovernance`  |
| `POST` | `/furya/cwhooks/v1/tx/register_staking`      |
| `POST` | `/furya/cwhooks/v1/tx/unregister_staking`    |
| `POST` | `/furya/cwhooks/v1/tx/register_governance`   |
| `POST` | `/furya/cwhooks/v1/tx/unregister_governance` |
