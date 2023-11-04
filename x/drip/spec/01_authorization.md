<!--
order: 1
-->

# Authorization

For security reasons, only specific addresses can distribute tokens to $FURYA stakers. We accept any kind of address: multisig, smart contracts, regular and [https://daodao.zone](DAODAO) DAOs. 

Governance can decide wether to approve or deny a new address to be added to the authorized list.

## Query the allowed addresses

You can query the list of allowed addresses directly from x/drip params

```
% furyad q drip params --output json
{"enable_drip":true,"allowed_addresses":[]}
```

## Governance proposal

To update the authorized address is possible to create a onchain new proposal. You can use the following example `proposal.json` file

```json
{
 "messages": [
  {
   "@type": "/furya.drip.v1.MsgUpdateParams",
   "authority": "furya10d07y265gmmuvt4z0w9aw880jnsr700jvss730",
   "params": {
    "enable_drip": false,
    "allowed_addresses": ["furya1j0a9ymgngasfn3l5me8qpd53l5zlm9wurfdk7r65s5mg6tkxal3qpgf5se"]
   }
  }
 ],
 "metadata": "{\"title\": \"Allow an amazing contract to distribute tokens using drip\", \"authors\": [\"dimi\"], \"summary\": \"If this proposal passes furya1j0a9ymgngasfn3l5me8qpd53l5zlm9wurfdk7r65s5mg6tkxal3qpgf5se will be added to the authorized addresses of the drip module\", \"details\": \"If this proposal passes furya1j0a9ymgngasfn3l5me8qpd53l5zlm9wurfdk7r65s5mg6tkxal3qpgf5se will be added to the authorized addresses of the drip module\", \"proposal_forum_url\": \"https://commonwealth.im/furya/discussion/9697-furya-protocol-level-defi-incentives\", \"vote_option_context\": \"yes\"}",
 "deposit": "1000ufury",
 "title": "Allow an amazing contract to distribute tokens using drip",
 "summary": "If this proposal passes furya1j0a9ymgngasfn3l5me8qpd53l5zlm9wurfdk7r65s5mg6tkxal3qpgf5se will be added to the authorized addresses of the drip module"
}
```

It can be submitted with the standard `furyad tx gov submit-proposal proposal.json --from yourkey` command.
